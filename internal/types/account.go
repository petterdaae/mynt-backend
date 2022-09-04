package types

type Account struct {
	ID            string `json:"id"`
	AccountNumber string `json:"accountNumber"`
	Name          string `json:"name"`
	Available     int    `json:"available"`
	Balance       int    `json:"balance"`
}
