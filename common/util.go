package common

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
