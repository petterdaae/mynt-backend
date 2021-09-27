package spendings

import (
	"mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Spending struct {
	CategoryID int64 `json:"category_id"`
	Amount     int64 `json:"amount"`
}

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	rows, err := database.Query(
		`SELECT category_id, SUM(tc.amount)
		FROM transactions_to_categories AS tc, transactions AS t
		WHERE tc.transaction_id = t.id
		AND t.user_id = $1
		AND tc.user_id = $1
		AND t.accounting_date >= $2
		AND t.accounting_date <= $3
		GROUP BY category_id`,
		sub,
		c.Query("from_date"),
		c.Query("to_date"),
	)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}
	defer rows.Close()

	spendings := []Spending{}
	for rows.Next() {
		var spending Spending
		err := rows.Scan(&spending.CategoryID, &spending.Amount)
		if err != nil {
			utils.InternalServerError(c, err)
			return
		}
		spendings = append(spendings, spending)
	}

	c.JSON(http.StatusOK, spendings)
}
