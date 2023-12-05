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
	g.Log().Info(ctx, "start sync ela blocks")
	currentELAHeight, err := rpc.ELAGetCurrentBlockHeight()
	if err != nil {
		g.Log().Error(ctx, "get current ela height error:", err)
		return nil
	}
	g.Log().Info(ctx, "current ela height:", currentELAHeight)

	// todo get ela block and transactions
	oneDayTransactionsCount := make(map[string]int)
	dailyTransactionsCount := make(map[string]int)

	oneDayAddressesCount := make(map[string]int)
	dailyActiveAddressesCount := make(map[string]int)

	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), int(startHour), 0, 0, 0, time.UTC)
	if currentTime.Hour() < int(startHour) {
		startTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-1, int(startHour), 0, 0, 0, time.UTC)
	}
	for i := currentELAHeight - 1; i > 0; i-- {
		block, err := rpc.ELAGetBlockbyheight(strconv.Itoa(int(i)))
		if err != nil {
			g.Log().Error(ctx, "get block by height error:", err)
			return nil
		}
		utcTimestamp := int64(block.Time)
		blockTime := time.Unix(utcTimestamp, 0)
		g.Log().Info(ctx, "main chain height:", block.Height, "time:", blockTime.Format("2006-01-02 15:04:05"), "tx count:", len(block.Tx))

		// get active addresses
		addressesMap := make(map[string]int)
		for i, tx := range block.Tx {
			var res ela.TransactionContextInfo
			if err = ela.Unmarshal(tx, &res); err != nil {
				g.Log().Error(ctx, "Error parsing JSON:", err)
				continue
			}
			for _, output := range res.Outputs {
				addressesMap[output.Address] += 1
			}

			if i != 0 {
				for _, input := range res.Inputs {
					itx, err := rpc.ELAGetRawTransaction(input.TxID)
					if err != nil {
						g.Log().Error(ctx, "get raw transaction error:", err)
						continue
					}
					for _, output := range itx.Outputs {
						addressesMap[output.Address] += 1
					}
				}
			}
		}

		// record daily transactions and active addresses count
		if !currentTime.After(blockTime.Add(24 * time.Hour)) {
			oneDayKey := currentTime.Format("2006-01-02 12:03:04") + "~" + currentTime.Add(-24*time.Hour).Format("2006-01-02 12:03:04")
			oneDayTransactionsCount[oneDayKey] += len(block.Tx)
			oneDayAddressesCount[oneDayKey] += len(addressesMap)
		} else {
			g.Log().Info(ctx, "current time after block time add 1 day")
		}

		// record daily transactions and active addreses count
		if blockTime.Hour() >= int(startHour) {
			dailyTransactionsCount[blockTime.Format("2006-01-02")] += len(block.Tx)
			dailyActiveAddressesCount[blockTime.Format("2006-01-02")] += len(addressesMap)
		} else {
			dailyTransactionsCount[blockTime.Add(-24*time.Hour).Format("2006-01-02")] += len(block.Tx)
			dailyActiveAddressesCount[blockTime.Add(-24*time.Hour).Format("2006-01-02")] += len(addressesMap)
		}

		if startTime.After(blockTime.Add(time.Duration(days) * 24 * time.Hour)) {
			break
		}
	}

	// tempTime := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	// for !tempTime.After(currentTime) {
	// 	tempTime = tempTime.Add(1 * time.Hour)
	// 	dailyTransactionsCount[tempTime.Format("2006-01-02")] += int(rand.Intn(10))
	// 	dailyActiveAddressesCount[tempTime.Format("2006-01-02")] += int(rand.Intn(10))
	// }

	// calculate weekly and monthly transactions count
	wtc, mtc := common.CalculateWeeklyAndMonthlyActivationData(common.ActivationMapToSortedList(dailyTransactionsCount))
	weeklyTransactionsCount := common.ActivationListToMap(wtc)
	monthlyTransactionsCount := common.ActivationListToMap(mtc)

	// calculate weekly and monthly active addresses count
	wac, mac := common.CalculateWeeklyAndMonthlyActivationData(common.ActivationMapToSortedList(dailyActiveAddressesCount))
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
