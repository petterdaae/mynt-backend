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

type CreatedBudgetResponse struct {
	ID int64 `json:"id"`
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

	var userIDOfBudget string
	row, err := database.QueryRow("SELECT user_id FROM budgets WHERE id = $1", body.BudgetID)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed query budget: %w", err))
		return
	}
	err = row.Scan(&userIDOfBudget)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to scan userIdOfBudget: %w", err))
		return
	}

	if userIDOfBudget != sub {
		c.Status(http.StatusUnauthorized)
		return
	}

	var newBudgetItemID int64
	row, err = database.QueryRow(
		`INSERT INTO budget_items 
		(user_id, budget_id, category_id, negative_amount, positive_amount, name) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
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
	err = row.Scan(&newBudgetItemID)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to scan new budget item id: %w", err))
		return
	}

	response := CreatedBudgetResponse{
		ID: newBudgetItemID,
	}

	c.JSON(http.StatusCreated, response)
}
