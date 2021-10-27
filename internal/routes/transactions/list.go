package transactions

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Transaction struct {
	ID             string `json:"id"`
	AccountID      string `json:"account_id"`
	AccountingDate string `json:"accounting_date"`
	InterestDate   string `json:"interest_date"`
	Amount         int64  `json:"amount"`
	Text           string `json:"text"`
	CategoryID     *int64 `json:"category_id"`
}

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	rows, err := database.Query(
		`SELECT t.id, t.account_id, t.accounting_date, t.interest_date, t.amount, t.text, tc.category_id
		FROM transactions AS t LEFT JOIN transactions_to_categories AS tc ON t.id = tc.transaction_id
		WHERE t.user_id = $1
		AND accounting_date >= $2
		AND accounting_date <= $3
		ORDER BY t.accounting_date, t.id DESC`,
		sub,
		c.Query("from_date"),
		c.Query("to_date"),
	)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}
	defer rows.Close()

	transactions := []Transaction{}
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.AccountingDate,
			&transaction.InterestDate,
			&transaction.Amount,
			&transaction.Text,
			&transaction.CategoryID,
		)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to scan row: %w", err))
			return
		}

		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, transactions)
}
