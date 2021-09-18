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
	Amount         int    `json:"amount"`
	Text           string `json:"text"`
}

type RequestBody struct {
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to connect to database: %w", err))
		return
	}
	defer connection.Close()

	params := c.Request.URL.Query()

	rows, err := connection.Query(
		"SELECT id, account_id, accounting_date, interest_date, amount, text FROM transactions WHERE user_id = $1 "+
			"AND accounting_date >= $2 AND accounting_date <= $3",
		sub, params["from_date"][0], params["to_date"][0])
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to query database: %w", err))
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
