package types

type Categorization struct {
	ID            int64  `json:"id"`
	TransactionID string `json:"transactionId"`
	Amount        int64  `json:"amount"`
	CategoryID    string `json:"categoryId"`
}
