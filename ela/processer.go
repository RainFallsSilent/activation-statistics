package ela

import (
	"context"

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
	weeklyTransactionsCount := make(map[string]int)
	monthlyTransactionsCount := make(map[string]int)

	oneDayAddressesCount := make(map[string]int)
	dailyActiveAddressesCount := make(map[string]int)
	weeklyActiveAddressesCount := make(map[string]int)
	monthlyActiveAddressesCount := make(map[string]int)

	// add sample data for test
	{
		oneDayTransactionsCount["2021-10-01"] = 100
		dailyTransactionsCount["2021-10-01"] = 100
		dailyTransactionsCount["2021-10-02"] = 200
		dailyTransactionsCount["2021-10-03"] = 300
		weeklyTransactionsCount["2021-10-01"] = 100
		weeklyTransactionsCount["2021-10-08"] = 200
		monthlyTransactionsCount["2021-10-01"] = 100

		oneDayAddressesCount["2021-10-01"] = 100
		dailyActiveAddressesCount["2021-10-01"] = 100
		dailyActiveAddressesCount["2021-10-02"] = 200
		dailyActiveAddressesCount["2021-10-03"] = 300
		weeklyActiveAddressesCount["2021-10-01"] = 100
		weeklyActiveAddressesCount["2021-10-08"] = 200
		monthlyActiveAddressesCount["2021-10-01"] = 100
	}

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
