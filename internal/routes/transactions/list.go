package transactions

import (
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
		}
		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, transactions)
}
