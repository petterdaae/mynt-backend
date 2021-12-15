package transactions

import (
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	transactionType := c.Query("type")
	queryString :=
		`SELECT t.id, t.account_id, t.accounting_date, t.interest_date, t.amount, t.text, tc.category_id, t.custom_date
		FROM transactions AS t LEFT JOIN transactions_to_categories AS tc ON t.id = tc.transaction_id
		WHERE t.user_id = $1`

	if transactionType == "uncategorized" {
		queryString += ` AND tc.category_id IS NULL `
	}

	queryString +=
		`AND accounting_date >= $2
		AND accounting_date <= $3
		ORDER BY t.accounting_date DESC, t.id`

	rows, err := database.Query(
		queryString,
		sub,
		c.Query("from_date")+"T00:00:00",
		c.Query("to_date")+"T00:00:00",
	)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}
	defer rows.Close()

	transactions := []types.Transaction{}
	for rows.Next() {
		var transaction types.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.AccountingDate,
			&transaction.InterestDate,
			&transaction.Amount,
			&transaction.Text,
			&transaction.CategoryID,
			&transaction.CustomDate,
		)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to scan row: %w", err))
			return
		}

		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, transactions)
}
