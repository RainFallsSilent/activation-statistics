package common

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestActiveAddressesMapToSortedList(t *testing.T) {
	addressMap := make(map[string]SetCount)
	addressMap["2023-12-04"] = NewSetCount()
	addressMap["2023-12-06"] = NewSetCount()
	addressMap["2023-12-05"] = NewSetCount()
	addressMap["2023-12-03"] = NewSetCount()
	addressMap["2023-12-02"] = NewSetCount()
	addressMap["2023-12-01"] = NewSetCount()
	addressMap["2023-11-30"] = NewSetCount()

	address := [4]string{"0xbeeaab15628329c2c89bc9f403d34b31fbcb3085", "0x9d5641fc60fa00af9406528d1f41f45c86babb72", "0xbeeaab15628329c2c89bc9f403d34b31fbcb3085", "0x024f5e84cd663c3150552ad6087be59a385468f6"}
	for _, item := range addressMap {
		for i := 0; i < len(address); i++ {
			item.Add(address[i])
		}
	}

	list := make(map[string]map[string]int)
	for key, value := range addressMap {
		list[key] = value
	}
	fmt.Println("list", list)
	data := ActiveAddressesMapToSortedList(list)
	fmt.Println("data", data)

	w, m := CalculateWeeklyAndMonthlyActiveAddressData(data)
	fmt.Println(w)
	fmt.Println(m)
}

func TestCalculateWeeklyAndMonthlyActivationData(t *testing.T) {
	var dailyTransactionsCount = make(map[string]int)
	now := time.Now()
	for i := 0; i < 20; i++ {
		data := time.Date(now.Year(), now.Month(), now.Day()-i, 0, 0, 0, 0, time.Local)
		dailyTransactionsCount[data.Format("2006-01-02")] = rand.Int()
	}
	fmt.Println(dailyTransactionsCount)
	w, m := CalculateWeeklyAndMonthlyActivationData(ActivationMapToSortedList(dailyTransactionsCount))
	fmt.Println(w)
	fmt.Println(m)
}
