package types

type Account struct {
	ID            string `json:"id"`
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	Available     int    `json:"available"`
	Balance       int    `json:"balance"`
}
