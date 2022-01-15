package budgetitems

import (
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Update(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var budget types.BudgetItem
	err := utils.ParseBody(c, &budget)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	err = database.Exec(
		`UPDATE budget_items 
		SET 
			negative_amount = $1, 
			positive_amount = $2,
			category_id = $3,
			name = $4
		WHERE user_id = $5 AND id = $6`,
		budget.NegativeAmount,
		budget.PositiveAmount,
		budget.CategoryID,
		budget.Name,
		sub,
		budget.ID,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("update budget_items failed: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
