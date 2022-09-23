package synchronize

import (
	"backend/internal/resources/sbanken"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Sbanken(c *gin.Context) {
	sub := c.GetString("sub")
	database, _ := c.MustGet("database").(*utils.Database)

	sbankenResource := sbanken.Configure(sub, database)

	// Connect to database
	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("database connection failed: %w", err))
		return
	}
	defer connection.Close()

	// Get client id and secret
	var clientID string
	var clientSecret string
	rows, err := connection.Query("SELECT sbanken_client_id, sbanken_client_secret FROM users WHERE id = $1", sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("query for sbanken credentials failed: %w", err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&clientID, &clientSecret)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("parsing query result for sbanken credentials failed: %w", err))
			return
		}
	}

	if clientID == "" || clientSecret == "" {
		utils.Unauthorized(c, fmt.Errorf("missing credentials: client id or client secret is blank"))
		return
	}

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
			utils.InternalServerError(c, fmt.Errorf("failed to get sbanken transactions: %w", err))
			return
		}

		for _, transaction := range transactions.Items {
			_, err := connection.Exec(
				"INSERT INTO transactions (id, user_id, account_id, external_id, accounting_date, interest_date, amount, text) "+
					"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) "+
					"ON CONFLICT (id) DO NOTHING",
				"sbanken:"+transaction.TransactionID,
				sub,
				"sbanken:"+account.AccountID,
				transaction.TransactionID,
				transaction.AccountingDate,
				transaction.InterestDate,
				utils.CurrencyToInt(transaction.Amount),
				transaction.Text,
			)

			if err != nil {
				utils.InternalServerError(c, fmt.Errorf("failed to insert sbanken transaction: %w", err))
				return
			}
		}
	}

	c.String(http.StatusOK, "Success")
}
