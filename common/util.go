package common

import (
	"fmt"
	"sort"
	"time"
)

type Activation struct {
	// transactions count
	OneDayTransactionsCount  map[string]int
	DailyTransactionsCount   map[string]int
	WeeklyTransactionsCount  map[string]int
	MonthlyTransactionsCount map[string]int

	// addresses count
	OneDayActiveAddressesCount  map[string]int
	DailyActiveAddressesCount   map[string]int
	WeeklyActiveAddressesCount  map[string]int
	MonthlyActiveAddressesCount map[string]int
}

type ActivationData struct {
	Date  string
	Count int
}

type AddressData struct {
	Date          string
	AddressesData map[string]int
}

func ActivationListToMap(input []ActivationData) map[string]int {
	result := make(map[string]int)
	for _, data := range input {
		result[data.Date] = data.Count
	}
	return result
}

func ActivationMapToSortedList(input map[string]int) []ActivationData {
	result := make([]ActivationData, 0)
	for k, v := range input {
		result = append(result, ActivationData{
			Date:  k,
			Count: v,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})

	return result
}

func ActiveAddressesMapToSortedList(input map[string]map[string]int) []AddressData {
	result := make([]AddressData, 0)
	for k, v := range input {
		result = append(result, AddressData{
			Date:          k,
			AddressesData: v,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})

	return result
}

func CalculateWeeklyAndMonthlyActivationData(currentDate time.Time, dailyTransactionsCountList []ActivationData) (weeklyTransactionsCount, monthlyTransactionsCount []ActivationData) {
	// Initialize weeklyTransactionsCount and monthlyTransactionsCount maps
	weeklyTransactionsCount = make([]ActivationData, 0)
	monthlyTransactionsCount = make([]ActivationData, 0)

	// Calculate weekly transactions count
	// Start from 7 days ago
	weeklyStartDate := currentDate.AddDate(0, 0, -6)
	weeklyEndDate := currentDate
	lastDayDate := currentDate
	weekValue := int(0)
	for _, data := range dailyTransactionsCountList {
		date, err := time.Parse("2006-01-02", data.Date)
		if err != nil {
			// Skip invalid date strings
			fmt.Println("invalid date string:", data.Date, "err:", err)
			continue
		}
		if !(date.After(weeklyStartDate) || date.Equal(weeklyStartDate)) {
			weekKey := weeklyStartDate.Format("2006-01-02") + "~" + weeklyEndDate.Format("2006-01-02")
			weeklyTransactionsCount = append(weeklyTransactionsCount, ActivationData{
				Date:  weekKey,
				Count: weekValue,
			})
			weekValue = int(0)
			weeklyEndDate = weeklyStartDate.AddDate(0, 0, -1)
			weeklyStartDate = weeklyEndDate.AddDate(0, 0, -6)
		}
		weekValue += data.Count
		lastDayDate = date
	}
	weekKey := lastDayDate.Format("2006-01-02") + "~" + weeklyEndDate.Format("2006-01-02")
	weeklyTransactionsCount = append(weeklyTransactionsCount, ActivationData{
		Date:  weekKey,
		Count: weekValue,
	})

	// Calculate monthly transactions count
	// Start from 30 days ago
	monthlyStartDate := currentDate.AddDate(0, -1, 0)
	monthlyEndDate := currentDate
	monthlyValue := int(0)
	for _, data := range dailyTransactionsCountList {
		date, err := time.Parse("2006-01-02", data.Date)
		if err != nil {
			// Skip invalid date strings
			fmt.Println("invalid date string:", data.Date, "err:", err)
			continue
		}
		if !(date.After(monthlyStartDate) || date.Equal(monthlyStartDate)) {
			monthKey := monthlyStartDate.Format("2006-01-02") + "~" + monthlyEndDate.Format("2006-01-02")
			monthlyTransactionsCount = append(monthlyTransactionsCount, ActivationData{
				Date:  monthKey,
				Count: monthlyValue,
			})
			monthlyValue = int(0)
			monthlyEndDate = monthlyStartDate
			monthlyStartDate = monthlyStartDate.AddDate(0, -1, 0)
		}
		monthlyValue += data.Count
	}
	monthKey := lastDayDate.Format("2006-01-02") + "~" + monthlyEndDate.Format("2006-01-02")
	monthlyTransactionsCount = append(monthlyTransactionsCount, ActivationData{
		Date:  monthKey,
		Count: monthlyValue,
	})

	return weeklyTransactionsCount, monthlyTransactionsCount
}

func GetDateByTimeStamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02")
}

func CalculateWeeklyAndMonthlyActiveAddressData(currentDate time.Time, dailyTransactionsCountList []AddressData) (weeklyTransactionsCount, monthlyTransactionsCount []ActivationData) {
	// Initialize weeklyTransactionsCount and monthlyTransactionsCount maps
	weeklyTransactionsCount = make([]ActivationData, 0)
	monthlyTransactionsCount = make([]ActivationData, 0)

	// Calculate weekly transactions count
	// Start from 7 days ago
	weeklyStartDate := currentDate.AddDate(0, 0, -6)
	weeklyEndDate := currentDate
	lastDayDate := currentDate
	weeklyAddresses := make(map[string]int)
	for _, data := range dailyTransactionsCountList {
		date, err := time.Parse("2006-01-02", data.Date)
		if err != nil {
			// Skip invalid date strings
			fmt.Println("invalid date string:", data.Date, "err:", err)
			continue
		}
		if !(date.After(weeklyStartDate) || date.Equal(weeklyStartDate)) {
			weekKey := weeklyStartDate.Format("2006-01-02") + "~" + weeklyEndDate.Format("2006-01-02")
			weeklyTransactionsCount = append(weeklyTransactionsCount, ActivationData{
				Date:  weekKey,
				Count: len(weeklyAddresses),
			})
			weeklyAddresses = make(map[string]int)
			weeklyEndDate = weeklyStartDate.AddDate(0, 0, -1)
			weeklyStartDate = weeklyEndDate.AddDate(0, 0, -6)
		}
		for addr, count := range data.AddressesData {
			weeklyAddresses[addr] += count
		}
		lastDayDate = date
	}
	weekKey := lastDayDate.Format("2006-01-02") + "~" + weeklyEndDate.Format("2006-01-02")
	weeklyTransactionsCount = append(weeklyTransactionsCount, ActivationData{
		Date:  weekKey,
		Count: len(weeklyAddresses),
	})

	// Calculate monthly transactions count
	// Start from 30 days ago
	monthlyStartDate := currentDate.AddDate(0, -1, 0)
	monthlyEndDate := currentDate
	monthlyAddresses := make(map[string]int)
	for _, data := range dailyTransactionsCountList {
		date, err := time.Parse("2006-01-02", data.Date)
		if err != nil {
			// Skip invalid date strings
			fmt.Println("invalid date string:", data.Date, "err:", err)
			continue
		}
		if !(date.After(monthlyStartDate) || date.Equal(monthlyStartDate)) {
			monthKey := monthlyStartDate.Format("2006-01-02") + "~" + monthlyEndDate.Format("2006-01-02")
			monthlyTransactionsCount = append(monthlyTransactionsCount, ActivationData{
				Date:  monthKey,
				Count: len(monthlyAddresses),
			})
			monthlyAddresses = make(map[string]int)
			monthlyEndDate = monthlyStartDate
			monthlyStartDate = monthlyStartDate.AddDate(0, -1, 0)
		}
		for addr, count := range data.AddressesData {
			monthlyAddresses[addr] += count
		}
	}
	monthKey := lastDayDate.Format("2006-01-02") + "~" + monthlyEndDate.Format("2006-01-02")
	monthlyTransactionsCount = append(monthlyTransactionsCount, ActivationData{
		Date:  monthKey,
		Count: len(monthlyAddresses),
	})

	return weeklyTransactionsCount, monthlyTransactionsCount
}
