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
		`SELECT t.id, t.account_id, split_part(t.accounting_date, 'T', 1) as accounting_date, split_part(t.interest_date, 'T', 1) as interest_date, t.amount, t.text, tc.category_id, t.custom_date
		FROM transactions AS t LEFT JOIN transactions_to_categories AS tc ON t.id = tc.transaction_id
		WHERE t.user_id = $1`

	if transactionType == "uncategorized" {
		queryString += ` AND tc.category_id IS NULL `
	}

	queryString +=
		`AND (
			(t.custom_date IS NOT NULL AND t.custom_date >= $2 AND t.custom_date <= $3)
			OR
			(t.custom_date IS NULL AND accounting_date >= $2 AND accounting_date <= $3)
		)
		ORDER BY COALESCE(t.custom_date, accounting_date) DESC, t.id`

	rows, err := database.Query(
		queryString,
		sub,
		c.Query("from_date"),
		c.Query("to_date"),
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
