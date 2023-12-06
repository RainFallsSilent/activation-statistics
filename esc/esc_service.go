package esc

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gogf/gf/v2/frame/g"

	com "github.com/RainFallsSilent/activation-statistics/common"
	"github.com/RainFallsSilent/activation-statistics/rpc"
	"github.com/RainFallsSilent/activation-statistics/rpc/esc"
)

type ESCService struct {
	client       *esc.Client
	totalDays    uint32
	dayStartHour int64

	activation *com.Activation
	addressMap map[string]com.SetCount // date->addressList
}

func New(days uint32, daystartHour int64) (*ESCService, error) {
	c, err := esc.Dial(rpc.EscRpcConfig.HttpUrl)
	if err != nil {
		return nil, err
	}
	s := &ESCService{
		client: c, totalDays: days, dayStartHour: daystartHour,
	}

	s.activation = &com.Activation{
		OneDayTransactionsCount:     make(map[string]int),
		DailyTransactionsCount:      make(map[string]int),
		WeeklyTransactionsCount:     make(map[string]int),
		MonthlyTransactionsCount:    make(map[string]int),
		OneDayActiveAddressesCount:  make(map[string]int),
		DailyActiveAddressesCount:   make(map[string]int),
		WeeklyActiveAddressesCount:  make(map[string]int),
		MonthlyActiveAddressesCount: make(map[string]int),
	}

	s.addressMap = make(map[string]com.SetCount) //date -> addressList
	return s, nil
}

func (s *ESCService) Start() error {
	latestBlock, err := s.client.LatestBlock()
	if err != nil {
		return err
	}
	var block *types.Block
	ctx := context.Background()
	runTime := time.Now()
	before24Hour := runTime.Add(-24 * time.Hour)
	endDay := time.Date(runTime.Year(), runTime.Month(), runTime.Day()-int(s.totalDays), int(s.dayStartHour), 0, 0, 0, time.UTC)
	g.Log("ESC").Info(ctx, "runTime", runTime, "endDay", endDay, "totalDay", s.totalDays)
	var blockInterfaces []interface{}
	txCount := 0
	for i := latestBlock.Uint64(); i > 0; i-- {
		height := big.NewInt(0).SetUint64(i)
		block, err = s.client.BlockByNumber(ctx, height)
		if err != nil {
			i++
			g.Log("ESC").Error(ctx, "get block error", err, "height", height)
			continue
		}
		btime := time.Unix(int64(block.Time()), 0)
		txCount = block.Transactions().Len()
		startHour := time.Duration(-s.dayStartHour)
		dateStr := com.GetDateByTimeStamp(btime.Add(startHour * time.Hour).Unix())
		days := uint32(len(s.activation.DailyTransactionsCount))
		g.Log("ESC").Info(ctx, "detail block ", i, "timestamp ", block.Time(), "date ", dateStr, "days", days)
		if btime.Before(endDay) {

			break
		}
		if txCount == 0 {
			continue
		}

		if runTime.Sub(btime).Hours() < 24 {
			str := before24Hour.Format("2006-01-02 15:04:05") + "~" + runTime.Format("2006-01-02 15:04:05")
			if count, ok := s.activation.OneDayTransactionsCount[str]; ok {
				s.activation.OneDayTransactionsCount[str] = count + txCount
			} else {
				s.activation.OneDayTransactionsCount[str] = txCount
			}
		}
		if err != nil {
			g.Log("ESC").Fatal(ctx, "BlockByNumber failed", err)
		}
		blockInterfaces, err = s.client.TraceBlockByNumber(ctx, height)
		if err != nil {
			g.Log("ESC").Fatal(ctx, "TraceBlockByNumber failed", err)
		}

		if count, ok := s.activation.DailyTransactionsCount[dateStr]; ok {
			s.activation.DailyTransactionsCount[dateStr] = count + txCount
		} else {
			s.activation.DailyTransactionsCount[dateStr] = txCount
		}
		s.processTraceBlockInfo(blockInterfaces, dateStr)
		s.activation.DailyActiveAddressesCount[dateStr] = s.addressMap[dateStr].Size()
	}

	w, m := com.CalculateWeeklyAndMonthlyActivationData(runTime, com.ActivationMapToSortedList(s.activation.DailyTransactionsCount))
	s.activation.WeeklyTransactionsCount = com.ActivationListToMap(w)
	s.activation.MonthlyTransactionsCount = com.ActivationListToMap(m)

	list := make(map[string]map[string]int)
	for key, value := range s.addressMap {
		list[key] = value
	}
	w, m = com.CalculateWeeklyAndMonthlyActiveAddressData(runTime, com.ActiveAddressesMapToSortedList(list))
	s.activation.WeeklyActiveAddressesCount = com.ActivationListToMap(w)
	s.activation.MonthlyActiveAddressesCount = com.ActivationListToMap(m)

	return nil
}

func (s *ESCService) processTraceBlockInfo(infos []interface{}, date string) {
	for i, info := range infos {
		g.Log("ESC").Info(context.Background(), fmt.Sprintf("processTraceBlockInfo %d", i))
		typeInfo := reflect.TypeOf(info)
		if typeInfo.String() != "map[string]interface {}" {
			continue
		}
		txInfo := info.(map[string]interface{})
		var result map[string]interface{}
		var ok bool
		if result, ok = txInfo["result"].(map[string]interface{}); !ok {
			continue
		}
		s.processTxResult(result, date)
	}
}

func (s *ESCService) processTxResult(txInfo map[string]interface{}, date string) {
	var zeroAddress common.Address
	var from string
	if txInfo["from"] == nil || txInfo["from"] == "" {
		return
	}
	from = txInfo["from"].(string)
	if s.addressMap[date] == nil {
		s.addressMap[date] = com.NewSetCount()
	}
	s.addressMap[date].Add(from)

	var to string
	if txInfo["to"] == nil || txInfo["to"] == "" {
		to = zeroAddress.String()
	} else {
		to = txInfo["to"].(string)
		s.addressMap[date].Add(to)
	}
	g.Log("ESC").Info(context.Background(), "processTxResult", "from", from, "to", to, "activeAddressCount", s.addressMap[date].Size())
	if calls, ok := txInfo["calls"]; ok {
		var items = calls.([]interface{})
		for i := 0; i < len(items); i++ {
			s.processTxResult(items[i].(map[string]interface{}), date)
		}
	}
}
