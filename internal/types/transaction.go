package types

type Transaction struct {
	ID             string  `json:"id"`
	AccountID      string  `json:"accountId"`
	AccountingDate string  `json:"accountingDate"`
	InterestDate   string  `json:"interestDate"`
	CustomDate     *string `json:"customDate"`
	Amount         int64   `json:"amount"`
	Text           string  `json:"text"`
}
