package types

type Categorization struct {
	ID            int64  `json:"id"`
	TransactionId string `json:"transactionId"`
	Amount        int64  `json:"amount"`
	CategoryId    string `json:"categoryId"`
}
