package ela

import (
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/RainFallsSilent/activation-statistics/common"
	"github.com/RainFallsSilent/activation-statistics/rpc"
	"github.com/RainFallsSilent/activation-statistics/rpc/ela"
)

func Process(ctx context.Context, days, startHour uint32) *common.Activation {
	g.Log("ELA").Info(ctx, "start sync ela blocks")
	currentELAHeight, err := rpc.ELAGetCurrentBlockHeight()
	if err != nil {
		g.Log("ELA").Error(ctx, "get current ela height error:", err)
		return nil
	}
	g.Log("ELA").Info(ctx, "current ela height:", currentELAHeight)

	// todo get ela block and transactions
	oneDayTransactionsCount := make(map[string]int)
	dailyTransactionsCount := make(map[string]int)

	oneDayAddressesCount := make(map[string]int)
	dailyActiveAddressesCount := make(map[string]int)

	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(startHour), 0, 0, 0, time.UTC)
	if currentTime.Hour() < int(startHour) {
		startTime = startTime.Add(-24 * time.Hour)
	}
	dailyAddressesMap := make(map[string]map[string]int)
	for i := currentELAHeight - 1; i > 0; i-- {
		block, err := rpc.ELAGetBlockbyheight(strconv.Itoa(int(i)))
		if err != nil {
			g.Log("ELA").Error(ctx, "get block by height error:", err)
			continue
		}
		utcTimestamp := int64(block.Time)
		blockTime := time.Unix(utcTimestamp, 0)
		g.Log("ELA").Info(ctx, "main chain height:", block.Height, "time:", blockTime.Format("2006-01-02 15:04:05"), "tx count:", len(block.Tx))

		// get active addresses
		addressesMap := make(map[string]int)
		for i, tx := range block.Tx {
			var res ela.TransactionContextInfo
			if err = ela.Unmarshal(tx, &res); err != nil {
				g.Log("ELA").Error(ctx, "Error parsing JSON:", err)
				continue
			}
			for _, output := range res.Outputs {
				addressesMap[output.Address] += 1
			}

			if i != 0 {
				g.Log().Info(ctx, "inputs count:", len(res.Inputs))
				for i, input := range res.Inputs {
					itx, err := rpc.ELAGetRawTransaction(input.TxID)
					if err != nil {
						g.Log("ELA").Error(ctx, "get raw transaction error:", err)
						continue
					}
					addressesMap[itx.Outputs[input.VOut].Address] += 1
					if i > 10 {
						// we think the inputs from same address
						break
					}
				}
			}
		}

		// record daily transactions and active addresses count
		if !currentTime.After(blockTime.Add(24 * time.Hour)) {
			oneDayKey := currentTime.Format("2006-01-02 15:04:05") + "~" + currentTime.Add(-24*time.Hour).Format("2006-01-02 15:04:05")
			oneDayTransactionsCount[oneDayKey] += len(block.Tx)
			oneDayAddressesCount[oneDayKey] += len(addressesMap)
		}

		// record daily transactions and active addreses count
		if blockTime.Hour() >= int(startHour) {
			dailyTransactionsCount[blockTime.Format("2006-01-02")] += len(block.Tx)

			// record detailed daily addresses information
			if _, ok := dailyAddressesMap[blockTime.Format("2006-01-02")]; !ok {
				dailyAddressesMap[blockTime.Format("2006-01-02")] = make(map[string]int)
			}
			for k, v := range addressesMap {
				dailyAddressesMap[blockTime.Format("2006-01-02")][k] += v
			}
		} else {
			dailyTransactionsCount[blockTime.Add(-24*time.Hour).Format("2006-01-02")] += len(block.Tx)

			// record detailed daily addresses information
			if _, ok := dailyAddressesMap[blockTime.Add(-24*time.Hour).Format("2006-01-02")]; !ok {
				dailyAddressesMap[blockTime.Add(-24*time.Hour).Format("2006-01-02")] = make(map[string]int)
			}
			for k, v := range addressesMap {
				dailyAddressesMap[blockTime.Add(-24*time.Hour).Format("2006-01-02")][k] += v
			}
		}

		if startTime.After(blockTime.Add(time.Duration(days) * 24 * time.Hour)) {
			break
		}
	}
	for k, v := range dailyAddressesMap {
		dailyActiveAddressesCount[k] = len(v)
	}

	// calculate weekly and monthly transactions count
	wtc, mtc := common.CalculateWeeklyAndMonthlyActivationData(currentTime, common.ActivationMapToSortedList(dailyTransactionsCount))
	weeklyTransactionsCount := common.ActivationListToMap(wtc)
	monthlyTransactionsCount := common.ActivationListToMap(mtc)

	// calculate weekly and monthly active addresses count
	wac, mac := common.CalculateWeeklyAndMonthlyActiveAddressData(currentTime, common.ActiveAddressesMapToSortedList(dailyAddressesMap))
	weeklyActiveAddressesCount := common.ActivationListToMap(wac)
	monthlyActiveAddressesCount := common.ActivationListToMap(mac)

	return &common.Activation{
		OneDayTransactionsCount:     oneDayTransactionsCount,
		DailyTransactionsCount:      dailyTransactionsCount,
		WeeklyTransactionsCount:     weeklyTransactionsCount,
		MonthlyTransactionsCount:    monthlyTransactionsCount,
		OneDayActiveAddressesCount:  oneDayAddressesCount,
		DailyActiveAddressesCount:   dailyActiveAddressesCount,
		WeeklyActiveAddressesCount:  weeklyActiveAddressesCount,
		MonthlyActiveAddressesCount: monthlyActiveAddressesCount,
	}
}
