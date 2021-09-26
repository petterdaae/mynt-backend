package transactions

import (
	"database/sql"
	"errors"
	"fmt"
	"mynt/internal/utils"
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
		`SELECT id, account_id, accounting_date, interest_date, amount, text 
		 FROM transactions 
		 WHERE user_id = $1 
	     AND accounting_date >= $2 
		 AND accounting_date <= $3`,
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
		)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to scan row: %w", err))
			return
		}

		row, err := database.QueryRow(
			"SELECT category_id FROM transactions_to_categories WHERE user_id = $1 AND transaction_id = $2 LIMIT 1",
			sub,
			transaction.ID,
		)
		if err != nil {
			utils.InternalServerError(c, err)
			return
		}

		err = row.Scan(&transaction.CategoryID)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.InternalServerError(c, fmt.Errorf("failed to scan row: %w", err))
			return
		}

		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, transactions)
}
