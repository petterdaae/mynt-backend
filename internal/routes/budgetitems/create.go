package budgetitems

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateBudgetItemBody struct {
	BudgetID       int64  `json:"budgetId"`
	CategoryID     int64  `json:"categoryId"`
	NegativeAmount int64  `json:"negativeAmount"`
	PositiveAmount int64  `json:"positiveAmount"`
	Name           string `json:"name"`
}

func Create(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var body CreateBudgetItemBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	var userIdOfBudget string
	row, err := database.QueryRow("SELECT user_id FROM budgets WHERE id = $1", body.BudgetID)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed query budget: %w", err))
		return
	}
	row.Scan(&userIdOfBudget)

	if userIdOfBudget != sub {
		c.Status(http.StatusUnauthorized)
		return
	}

	err = database.Exec(
		`INSERT INTO budget_items 
		(user_id, budget_id, category_id, negative_amount, positive_amount, name) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		sub,
		body.BudgetID,
		body.CategoryID,
		body.NegativeAmount,
		body.PositiveAmount,
		body.Name,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to insert new budget_item: %w", err))
		return
	}

	c.Status(http.StatusCreated)
}
