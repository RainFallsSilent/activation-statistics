package esc

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/frame/g"

	com "github.com/RainFallsSilent/activation-statistics/common"
	"github.com/RainFallsSilent/activation-statistics/rpc"
	"github.com/RainFallsSilent/activation-statistics/rpc/esc"
)

type ESCService struct {
	client       *esc.Client
	totalDays    uint32
	dayStartHour uint32

	activation *com.Activation
	addressMap map[string]com.SetCount // date->addressList
}

func New(days, daystartHour uint32) (*ESCService, error) {
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
	now := time.Now()
	todayStartTime := time.Date(now.Year(), now.Month(), now.Day(), int(s.dayStartHour), 0, 0, 0, time.UTC)
	ondayStartHour := time.Duration(s.dayStartHour) * time.Hour
	g.Log("ESC").Info(ctx, "todayStartTime", todayStartTime, "stamp ", todayStartTime.Unix(), "ondayStartHour ", ondayStartHour.Seconds())
	var blockInterfaces []interface{}
	//latestBlock.Uint64()
	_ = latestBlock
	txCount := 0
	for i := uint64(21986776); i > 0; i-- {
		height := big.NewInt(0).SetUint64(i)
		block, err = s.client.BlockByNumber(ctx, height)
		txCount = block.Transactions().Len()
		dateStr := com.GetDateByTimeStamp(uint64(float64(block.Time()) - ondayStartHour.Seconds()))
		days := uint32(len(s.activation.DailyTransactionsCount))
		g.Log("ESC").Info(ctx, "detail block ", i, "timestamp ", block.Time(), "date ", dateStr, "days", days)
		if days > s.totalDays {
			g.Log("ESC").Info(ctx, "ESC SERVICE COMPLETED")
			break
		}
		btime := time.Unix(int64(block.Time()), 0)

		if now.Sub(btime).Hours() < 24 {
			if count, ok := s.activation.OneDayTransactionsCount[dateStr]; ok {
				s.activation.OneDayTransactionsCount[dateStr] = count + txCount
			} else {
				s.activation.OneDayTransactionsCount[dateStr] = txCount
			}
		}
		fmt.Println(" past hours ==", now.Sub(btime).Hours(), "past seconds ", now.Sub(btime).Seconds())
		if days > 7 {
			s.activation.WeeklyTransactionsCount[dateStr] = s.activation.DailyTransactionsCount[dateStr]
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
	}
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