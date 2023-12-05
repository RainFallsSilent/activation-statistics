package ela

import (
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/RainFallsSilent/activation-statistics/common"
	"github.com/RainFallsSilent/activation-statistics/rpc"
)

func Process(ctx context.Context, days, startHour uint32) *common.Activation {
	currentELAHeight, err := rpc.ELAGetCurrentBlockHeight()
	if err != nil {
		g.Log().Error(ctx, "get current ela height error:", err)
		return nil
	}

	g.Log().Info(ctx, "current ela height:", currentELAHeight)

	// todo get ela block and transactions
	// 2023-10-01  100
	oneDayTransactionsCount := make(map[string]int)
	dailyTransactionsCount := make(map[string]int)

	oneDayAddressesCount := make(map[string]int)
	dailyActiveAddressesCount := make(map[string]int)
	weeklyActiveAddressesCount := make(map[string]int)
	monthlyActiveAddressesCount := make(map[string]int)

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
		g.Log().Info(ctx, "block height:", block.Height, "block time:", block.Time)
		utcTimestamp := int64(block.Time)
		blockTime := time.Unix(utcTimestamp, 0)
		if !currentTime.After(blockTime.Add(24 * time.Hour)) {
			oneDayTransactionsCount[currentTime.Format("2006-01-02 12:03:04")+"~"+currentTime.Add(-24*time.Hour).Format("2006-01-02 12:03:04")] += len(block.Tx)
		} else {
			g.Log().Info(ctx, "current time after block time add 1 day")
		}

		// record daily transactions count
		if blockTime.Hour() >= int(startHour) {
			dailyTransactionsCount[blockTime.Format("2006-01-02")] += len(block.Tx)
		} else {
			dailyTransactionsCount[blockTime.Add(-24*time.Hour).Format("2006-01-02")] += len(block.Tx)
		}

		if startTime.After(blockTime.Add(time.Duration(days) * 24 * time.Hour)) {
			break
		}

		// g.Log().Print(ctx, "block:", block)
		// break
	}

	// tempTime := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	// for !tempTime.After(currentTime) {
	// 	tempTime = tempTime.Add(1 * time.Hour)
	// 	dailyTransactionsCount[tempTime.Format("2006-01-02")] += int(rand.Intn(10))
	// }

	w, m := common.CalculateWeeklyAndMonthlyActivationData(common.ActivationMapToSortedList(dailyTransactionsCount))
	weeklyTransactionsCount := common.ActivationListToMap(w)
	monthlyTransactionsCount := common.ActivationListToMap(m)

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
