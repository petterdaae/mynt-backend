package synchronize

import (
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

	// Connect to database
	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("database connection failed: %w", err))
		return
	}
	defer connection.Close()

	accounts, err := sbankenResource.GetAccounts()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to get sbanken accounts: %w", err))
		return
	}

	// Update accounts and transactions
	for _, account := range accounts.Items {
		_, err := connection.Exec(
			"INSERT INTO accounts (id, user_id, external_id, account_number, name, available, balance) VALUES ($1, $2, $3, $4, $5, $6, $7) "+
				"ON CONFLICT (id) DO UPDATE SET name = $5, available = $6, balance = $7",
			"sbanken:"+account.AccountID,
			sub,
			account.AccountID,
			account.AccountNumber,
			account.Name,
			utils.CurrencyToInt(account.Available),
			utils.CurrencyToInt(account.Balance),
		)

		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to insert sbanken account: %w", err))
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
