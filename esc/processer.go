package esc

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
	// 2023-10-02  200
	// 2023-10-03  300
	oneDayTransactionsCount := make(map[string]int)
	dailyTransactionsCount := make(map[string]int)
	weeklyTransactionsCount := make(map[string]int)
	monthlyTransactionsCount := make(map[string]int)

	oneDayActiveAddressesCount := make(map[string]int)
	dailyActiveAddressesCount := make(map[string]int)
	weeklyActiveAddressesCount := make(map[string]int)
	monthlyActiveAddressesCount := make(map[string]int)

	// todo get esc height

	// todo get esc block and transactions

	return &common.Activation{
		OneDayTransactionsCount:     oneDayTransactionsCount,
		DailyTransactionsCount:      dailyTransactionsCount,
		WeeklyTransactionsCount:     weeklyTransactionsCount,
		MonthlyTransactionsCount:    monthlyTransactionsCount,
		OneDayActiveAddressesCount:  oneDayActiveAddressesCount,
		DailyActiveAddressesCount:   dailyActiveAddressesCount,
		WeeklyActiveAddressesCount:  weeklyActiveAddressesCount,
		MonthlyActiveAddressesCount: monthlyActiveAddressesCount,
	}
}
