package model

type YearStatistic struct {
	Year    int     `json:"year,omitempty"`
	Income  float64 `json:"income,omitempty"`
	Expense float64 `json:"expense,omitempty"`

	Months []MonthCount `json:"months,omitempty"`
}

type MonthCount struct {
	Month   int     `json:"month,omitempty"`
	Income  float64 `json:"income,omitempty"`
	Expense float64 `json:"expense,omitempty"`
}

// CategoryStatistic 按照分类统计结果
type CategoryStatistic struct {
	IncomeCost  float64                 `json:"income_cost"`
	ExpenseCost float64                 `json:"expense_cost"`
	Incomes     []CategoryStatisticInfo `json:"incomes,omitempty"`
	Expenses    []CategoryStatisticInfo `json:"expenses,omitempty"`
}

// CategoryStatisticInfo 按照分类统计结果
type CategoryStatisticInfo struct {
	CategoryId int     `json:"category_id,omitempty"`
	Count      int     `json:"count,omitempty"`
	Cost       float64 `json:"cost,omitempty"`
	Percent    float64 `json:"percent,omitempty"`
}
