package types

type Transaction struct {
	ID             string  `json:"id"`
	AccountID      string  `json:"account_id"`
	AccountingDate string  `json:"accounting_date"`
	InterestDate   string  `json:"interest_date"`
	Amount         int64   `json:"amount"`
	Text           string  `json:"text"`
	CategoryID     *int64  `json:"category_id"`
	CustomDate     *string `json:"customDate"`
}
