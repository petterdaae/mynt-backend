package types

type BudgetItem struct {
	ID            int64  `json:"id"`
	BudgetID      int64  `json:"budgetId"`
	CategoryID    int64  `json:"categoryId"`
	MonthlyAmount *int64 `json:"monthlyAmount"`
	Name          string `json:"name"`
}
