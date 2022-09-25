package types

type IncomingTransaction struct {
	ID string `json:"id"`
	DraftIncomingTransaction
}

type DraftIncomingTransaction struct {
	AccountID      string `json:"accountId"`
	AccountingDate string `json:"accountingDate"`
	InterestDate   string `json:"interestDate"`
	Amount         int64  `json:"amount"`
	Text           string `json:"text"`
}
