package types

type BudgetItem struct {
	ID            int64   `json:"id"`
	BudgetID      int64   `json:"budgetId"`
	Kind          string  `json:"kind"`
	CategoryID    int64   `json:"categoryId"`
	MonthlyAmount *int64  `json:"monthlyAmount"`
	CustomItems   *string `json:"customItems"`
	Name          string  `json:"name"`
}

type BudgetItemCustomItem struct {
	ID           int64  `json:"id"`
	BudgetItemID int64  `json:"budgetItemId"`
	Amount       int64  `json:"amount"`
	Date         string `json:"date"`
}
