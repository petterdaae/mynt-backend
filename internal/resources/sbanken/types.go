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
