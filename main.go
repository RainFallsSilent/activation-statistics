package main

import (
	"context"
	"os"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	_ "github.com/lib/pq"

	"github.com/RainFallsSilent/activation-statistics/rpc"
)

func main() {
	ctx := gctx.New()
	run(ctx)
}

func run(ctx context.Context) {
	days := gcmd.GetArg(1, "100").Uint32()
	g.Log().Info(ctx, "start sync blocks from ", days, "days ago")
	syncAndRecordActivation(ctx, days)
	g.Log().Info(ctx, "end sync blocks")
}

func syncAndRecordActivation(ctx context.Context, days uint32) {
	currentELAHeight, err := rpc.ELAGetCurrentBlockHeight()
	if err != nil {
		g.Log().Error(ctx, "get current ela height error:", err)
		return
	}

	g.Log().Info(ctx, "current ela height:", currentELAHeight)

	// todo get ela block and transactions
	elaDailyTransactionsCount := make(map[string]int)
	elaWeeklyTransactionsCount := make(map[string]int)
	elaMonthlyTransactionsCount := make(map[string]int)
	elaDailyActiveAddressesCount := make(map[string]int)
	elaWeeklyActiveAddressesCount := make(map[string]int)
	elaMonthlyActiveAddressesCount := make(map[string]int)

	// calculate daily\weekly\mothly ela active addresses and transactions

	// todo get esc height

	// todo get esc block and transactions
	escDailyTransactionsCount := make(map[string]int)
	escWeeklyTransactionsCount := make(map[string]int)
	escMonthlyTransactionsCount := make(map[string]int)
	escDailyActiveAddressesCount := make(map[string]int)
	escWeeklyActiveAddressesCount := make(map[string]int)
	escMonthlyActiveAddressesCount := make(map[string]int)

	// calculate daily\weekly\mothly esc active addresses and transactions

	// print result
	g.Log().Info(ctx, "ela daily transactions count:", elaDailyTransactionsCount)
	g.Log().Info(ctx, "ela weekly transactions count:", elaWeeklyTransactionsCount)
	g.Log().Info(ctx, "ela monthly transactions count:", elaMonthlyTransactionsCount)
	g.Log().Info(ctx, "ela daily active addresses count:", elaDailyActiveAddressesCount)
	g.Log().Info(ctx, "ela weekly active addresses count:", elaWeeklyActiveAddressesCount)
	g.Log().Info(ctx, "ela monthly active addresses count:", elaMonthlyActiveAddressesCount)

	g.Log().Info(ctx, "esc daily transactions count:", escDailyTransactionsCount)
	g.Log().Info(ctx, "esc weekly transactions count:", escWeeklyTransactionsCount)
	g.Log().Info(ctx, "esc monthly transactions count:", escMonthlyTransactionsCount)
	g.Log().Info(ctx, "esc daily active addresses count:", escDailyActiveAddressesCount)
	g.Log().Info(ctx, "esc weekly active addresses count:", escWeeklyActiveAddressesCount)
	g.Log().Info(ctx, "esc monthly active addresses count:", escMonthlyActiveAddressesCount)

	// open result.txt and save all count map result to file
	resultFile, err := os.OpenFile("result.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		g.Log().Error(ctx, "open result.txt error:", err)
		return
	}
	defer resultFile.Close()

	// range map to write result to file
	resultFile.WriteString("ela daily transactions count:\n")
	for k, v := range elaDailyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("ela weekly transactions count:\n")
	for k, v := range elaWeeklyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("ela monthly transactions count:\n")
	for k, v := range elaMonthlyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("ela daily active addresses count:\n")
	for k, v := range elaDailyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("ela weekly active addresses count:\n")
	for k, v := range elaWeeklyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("ela monthly active addresses count:\n")
	for k, v := range elaMonthlyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}

	resultFile.WriteString("esc daily transactions count:\n")
	for k, v := range escDailyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("esc weekly transactions count:\n")
	for k, v := range escWeeklyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("esc monthly transactions count:\n")
	for k, v := range escMonthlyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("esc daily active addresses count:\n")
	for k, v := range escDailyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("esc weekly active addresses count:\n")
	for k, v := range escWeeklyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}
	resultFile.WriteString("esc monthly active addresses count:\n")
	for k, v := range escMonthlyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v))
	}

}
