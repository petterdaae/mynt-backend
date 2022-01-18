package types

type BudgetItem struct {
	ID             int64  `json:"id"`
	BudgetID       int64  `json:"budgetId"`
	CategoryID     int64  `json:"categoryId"`
	NegativeAmount int64  `json:"negativeAmount"`
	PositiveAmount int64  `json:"positiveAmount"`
	Name           string `json:"name"`
}
