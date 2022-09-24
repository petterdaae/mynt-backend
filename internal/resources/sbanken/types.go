package sbanken

type Accounts struct {
	AvailableItems int
	Items          []Account
}

type Account struct {
	AccountID     string `json:"AccountId"`
	AccountNumber string
	Name          string
	Available     float64
	Balance       float64
}

type ArchievedTransactions struct {
	AvailableItems int
	Items          []ArchievedTransaction
}

type ArchievedTransaction struct {
	TransactionID  string `json:"TransactionId"`
	AccountingDate string
	InterestDate   string
	Amount         float64
	Text           string
}
