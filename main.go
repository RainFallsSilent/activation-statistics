package main

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	_ "github.com/lib/pq"

	"github.com/RainFallsSilent/activation-statistics/common"
	"github.com/RainFallsSilent/activation-statistics/ela"
	"github.com/RainFallsSilent/activation-statistics/esc"
)

func main() {
	var wg sync.WaitGroup
	ctx := gctx.New()
	wg.Add(1)
	run(ctx)
	wg.Wait()
}

func run(ctx context.Context) {
	days := gcmd.GetArg(1, "2").Uint32()
	startHour := gcmd.GetArg(2, "8").Uint32() // 0-24
	g.Log().Info(ctx, "start sync blocks from ", days, "days ago")
	syncAndRecordActivation(ctx, days, startHour)
	g.Log().Info(ctx, "end sync blocks")
}

func syncAndRecordActivation(ctx context.Context, days, startHour uint32) {

	// calculate daily\weekly\mothly ela active addresses and transactions
	elaActivation := ela.Process(ctx, days, startHour)

	// calculate daily\weekly\mothly esc active addresses and transactions
	escActivation := esc.Process(ctx, days, startHour)

	// print result and store to file
	printAndStore(ctx, elaActivation, escActivation)
}

func printAndStore(ctx context.Context, elaActivation, escActivation *common.Activation) {
	// print result
	g.Log().Info(ctx, "ela one day transactions count:", elaActivation.OneDayTransactionsCount)
	g.Log().Info(ctx, "ela daily transactions count:", elaActivation.DailyTransactionsCount)
	g.Log().Info(ctx, "ela weekly transactions count:", elaActivation.WeeklyTransactionsCount)
	g.Log().Info(ctx, "ela monthly transactions count:", elaActivation.MonthlyTransactionsCount)
	g.Log().Info(ctx, "ela one day active addresses count:", elaActivation.OneDayActiveAddressesCount)
	g.Log().Info(ctx, "ela daily active addresses count:", elaActivation.DailyActiveAddressesCount)
	g.Log().Info(ctx, "ela weekly active addresses count:", elaActivation.WeeklyActiveAddressesCount)
	g.Log().Info(ctx, "ela monthly active addresses count:", elaActivation.MonthlyActiveAddressesCount)

	g.Log().Info(ctx, "esc one day transactions count:", escActivation.OneDayTransactionsCount)
	g.Log().Info(ctx, "esc daily transactions count:", escActivation.DailyTransactionsCount)
	g.Log().Info(ctx, "esc weekly transactions count:", escActivation.WeeklyTransactionsCount)
	g.Log().Info(ctx, "esc monthly transactions count:", escActivation.MonthlyTransactionsCount)
	g.Log().Info(ctx, "esc one day active addresses count:", escActivation.OneDayActiveAddressesCount)
	g.Log().Info(ctx, "esc daily active addresses count:", escActivation.DailyActiveAddressesCount)
	g.Log().Info(ctx, "esc weekly active addresses count:", escActivation.WeeklyActiveAddressesCount)
	g.Log().Info(ctx, "esc monthly active addresses count:", escActivation.MonthlyActiveAddressesCount)

	// open result.txt and save all count map result to file
	resultFile, err := os.OpenFile("result.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		g.Log().Error(ctx, "open result.txt error:", err)
		return
	}
	defer resultFile.Close()

	// range map to write result to file
	resultFile.WriteString("ela daily transactions count:\n")
	for k, v := range elaActivation.DailyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("ela weekly transactions count:\n")
	for k, v := range elaActivation.WeeklyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("ela monthly transactions count:\n")
	for k, v := range elaActivation.MonthlyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("ela daily active addresses count:\n")
	for k, v := range elaActivation.DailyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("ela weekly active addresses count:\n")
	for k, v := range elaActivation.WeeklyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("ela monthly active addresses count:\n")
	for k, v := range elaActivation.MonthlyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}

	resultFile.WriteString("esc daily transactions count:\n")
	for k, v := range escActivation.DailyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("esc weekly transactions count:\n")
	for k, v := range escActivation.WeeklyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("esc monthly transactions count:\n")
	for k, v := range escActivation.MonthlyTransactionsCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("esc daily active addresses count:\n")
	for k, v := range escActivation.DailyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("esc weekly active addresses count:\n")
	for k, v := range escActivation.WeeklyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	resultFile.WriteString("esc monthly active addresses count:\n")
	for k, v := range escActivation.MonthlyActiveAddressesCount {
		resultFile.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}

}
