package spendings

import (
	"mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Spending struct {
	CategoryID int64  `json:"category_id"`
	Amount     int64  `json:"amount"`
	ParentID   *int64 `json:"parent_id"`
}

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	rows, err := database.Query(
		`SELECT tc.category_id, SUM(tc.amount), c.parent_id
		FROM transactions_to_categories AS tc, transactions AS t, categories AS c
		WHERE tc.transaction_id = t.id
		AND tc.category_id = c.id
		AND t.user_id = $1
		AND tc.user_id = $1
		AND t.accounting_date >= $2
		AND t.accounting_date <= $3
		GROUP BY tc.category_id, c.id`,
		sub,
		c.Query("from_date"),
		c.Query("to_date"),
	)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}
	defer rows.Close()

	spendings := []*Spending{}
	for rows.Next() {
		var spending Spending
		err := rows.Scan(&spending.CategoryID, &spending.Amount, &spending.ParentID)
		if err != nil {
			utils.InternalServerError(c, err)
			return
		}
		spendings = append(spendings, &spending)
	}

	groupSpendings(nil, spendings)

	c.JSON(http.StatusOK, spendings)
}

func groupSpendings(categoryID *int64, spendings []*Spending) int64 {
	var sum int64
	for _, spending := range spendings {
		if spending.ParentID == categoryID {
			sum += spending.Amount
			childSum := groupSpendings(&spending.CategoryID, spendings)
			sum += childSum
		}
	}
	for _, spending := range spendings {
		if &spending.CategoryID == categoryID {
			spending.Amount = sum + spending.Amount
		}
	}
	return sum
}
