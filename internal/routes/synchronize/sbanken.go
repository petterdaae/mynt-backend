package synchronize

import (
	aResource "backend/internal/resources/accounts"
	"backend/internal/resources/sbanken"
	tResource "backend/internal/resources/transactions"
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Sbanken(c *gin.Context) {
	sub := c.GetString("sub")
	database, _ := c.MustGet("database").(*utils.Database)

	sbankenResource := sbanken.Configure(sub, database)
	transactionsResource := tResource.Configure(sub, database)
	accountsResource := aResource.Configure(sub, database)

	accounts, err := sbankenResource.GetAccounts()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to get sbanken accounts: %w", err))
		return
	}

	for _, account := range accounts.Items {
		err = accountsResource.CreateIfNotExists(&types.Account{
			ID:            "sbanken:" + account.AccountID,
			AccountNumber: account.AccountNumber,
			Name:          account.Name,
			Available:     utils.CurrencyToInt(account.Available),
			Balance:       utils.CurrencyToInt(account.Balance),
		})

		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to account: %w", err))
			return
		}

		transactions, err := sbankenResource.GetArchievedTransactions(account.AccountID)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to get arhieved transactions from sbanken: %w", err))
			return
		}

		for _, transaction := range transactions.Items {
			err := transactionsResource.CreateIfNotExists(&types.Transaction{
				ID:             "sbanken:" + transaction.TransactionID,
				AccountID:      "sbanken:" + account.AccountID,
				AccountingDate: transaction.AccountingDate,
				InterestDate:   transaction.InterestDate,
				Amount:         int64(utils.CurrencyToInt(transaction.Amount)),
				Text:           transaction.Text,
			})

			if err != nil {
				utils.InternalServerError(c, fmt.Errorf("failed to create transaction: %w", err))
				return
			}
		}
	}

	c.String(http.StatusOK, "Success")
}
