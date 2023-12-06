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
	ctx := gctx.New()
	run(ctx)
}

func run(ctx context.Context) {
	days := gcmd.GetArg(1, "2").Uint32()
	startHour := gcmd.GetArg(2, "8").Uint32() // 0-24
	g.Log().Info(ctx, "start sync blocks from ", days, "days ago")
	syncAndRecordActivation(ctx, days, startHour)
	g.Log().Info(ctx, "end sync blocks")
}

func syncAndRecordActivation(ctx context.Context, days, startHour uint32) {
	var wg sync.WaitGroup

	// Create channels to receive the results
	elaActivationCh := make(chan *common.Activation)
	escActivationCh := make(chan *common.Activation)

	wg.Add(2)

	go func() {
		defer wg.Done()
		elaActivation := ela.Process(ctx, days, startHour)
		elaActivationCh <- elaActivation
	}()

	go func() {
		defer wg.Done()
		escActivation := esc.Process(ctx, days, startHour)
		escActivationCh <- escActivation
	}()

	go func() {
		wg.Wait()
		close(elaActivationCh)
		close(escActivationCh)
	}()

	elaActivation := <-elaActivationCh
	escActivation := <-escActivationCh

	// Print result and store to file
	printAndStore(ctx, elaActivation, escActivation)
}

func printAndStore(ctx context.Context, elaActivation, escActivation *common.Activation) {
	oneDayTxCountInfo := common.ActivationMapToSortedList(elaActivation.OneDayTransactionsCount)
	dailyTxCountInfo := common.ActivationMapToSortedList(elaActivation.DailyTransactionsCount)
	weeklyTxCountInfo := common.ActivationMapToSortedList(elaActivation.WeeklyTransactionsCount)
	monthlyTxCountInfo := common.ActivationMapToSortedList(elaActivation.MonthlyTransactionsCount)
	oneDayActiveAddressesCountInfo := common.ActivationMapToSortedList(elaActivation.OneDayActiveAddressesCount)
	dailyActiveAddressesCountInfo := common.ActivationMapToSortedList(elaActivation.DailyActiveAddressesCount)
	weeklyActiveAddressesCountInfo := common.ActivationMapToSortedList(elaActivation.WeeklyActiveAddressesCount)
	monthlyActiveAddressesCountInfo := common.ActivationMapToSortedList(elaActivation.MonthlyActiveAddressesCount)

	// print result
	g.Log().Info(ctx, "ela one day transactions count:", oneDayTxCountInfo)
	g.Log().Info(ctx, "ela daily transactions count:", dailyTxCountInfo)
	g.Log().Info(ctx, "ela weekly transactions count:", weeklyTxCountInfo)
	g.Log().Info(ctx, "ela monthly transactions count:", monthlyTxCountInfo)
	g.Log().Info(ctx, "ela one day active addresses count:", oneDayActiveAddressesCountInfo)
	g.Log().Info(ctx, "ela daily active addresses count:", dailyActiveAddressesCountInfo)
	g.Log().Info(ctx, "ela weekly active addresses count:", weeklyActiveAddressesCountInfo)
	g.Log().Info(ctx, "ela monthly active addresses count:", monthlyActiveAddressesCountInfo)

	// open result.txt and save all count map result to file
	resultFile, err := os.OpenFile("result.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		g.Log().Error(ctx, "open result.txt error:", err)
		return
	}
	defer resultFile.Close()

	// range map to write result to file
	resultFile.WriteString("ela one day transactions count:\n")
	for _, data := range oneDayTxCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("ela daily transactions count:\n")
	for _, data := range dailyTxCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("ela weekly transactions count:\n")
	for _, data := range weeklyTxCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("ela monthly transactions count:\n")
	for _, data := range monthlyTxCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("ela one day active addresses count:\n")
	for _, data := range oneDayActiveAddressesCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("ela daily active addresses count:\n")
	for _, data := range dailyActiveAddressesCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("ela weekly active addresses count:\n")
	for _, data := range weeklyActiveAddressesCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("ela monthly active addresses count:\n")
	for _, data := range monthlyActiveAddressesCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}

	// esc
	oneDayTxCountInfo = common.ActivationMapToSortedList(escActivation.OneDayTransactionsCount)
	dailyTxCountInfo = common.ActivationMapToSortedList(escActivation.DailyTransactionsCount)
	weeklyTxCountInfo = common.ActivationMapToSortedList(escActivation.WeeklyTransactionsCount)
	monthlyTxCountInfo = common.ActivationMapToSortedList(escActivation.MonthlyTransactionsCount)
	oneDayActiveAddressesCountInfo = common.ActivationMapToSortedList(escActivation.OneDayActiveAddressesCount)
	dailyActiveAddressesCountInfo = common.ActivationMapToSortedList(escActivation.DailyActiveAddressesCount)
	weeklyActiveAddressesCountInfo = common.ActivationMapToSortedList(escActivation.WeeklyActiveAddressesCount)
	monthlyActiveAddressesCountInfo = common.ActivationMapToSortedList(escActivation.MonthlyActiveAddressesCount)

	// print result
	g.Log().Info(ctx, "esc one day transactions count:", oneDayTxCountInfo)
	g.Log().Info(ctx, "esc daily transactions count:", dailyTxCountInfo)
	g.Log().Info(ctx, "esc weekly transactions count:", weeklyTxCountInfo)
	g.Log().Info(ctx, "esc monthly transactions count:", monthlyTxCountInfo)
	g.Log().Info(ctx, "esc one day active addresses count:", oneDayActiveAddressesCountInfo)
	g.Log().Info(ctx, "esc daily active addresses count:", dailyActiveAddressesCountInfo)
	g.Log().Info(ctx, "esc weekly active addresses count:", weeklyActiveAddressesCountInfo)
	g.Log().Info(ctx, "esc monthly active addresses count:", monthlyActiveAddressesCountInfo)

	// range map to write result to file
	resultFile.WriteString("esc one day transactions count:\n")
	for _, data := range oneDayTxCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("esc daily transactions count:\n")
	for _, data := range dailyTxCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("esc weekly transactions count:\n")
	for _, data := range weeklyTxCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("esc monthly transactions count:\n")
	for _, data := range monthlyTxCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("esc one day active addresses count:\n")
	for _, data := range oneDayActiveAddressesCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("esc daily active addresses count:\n")
	for _, data := range dailyActiveAddressesCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("esc weekly active addresses count:\n")
	for _, data := range weeklyActiveAddressesCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
	resultFile.WriteString("esc monthly active addresses count:\n")
	for _, data := range monthlyActiveAddressesCountInfo {
		resultFile.WriteString(data.Date + ": " + strconv.Itoa(data.Count) + "\n")
	}
}
