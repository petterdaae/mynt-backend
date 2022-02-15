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
			monthly_amount = $1, 
			category_id = $2,
			name = $3,
			kind = $6,
			custom_items = $7
		WHERE user_id = $4 AND id = $5`,
		budget.MonthlyAmount,
		budget.CategoryID,
		budget.Name,
		sub,
		budget.ID,
		budget.Kind,
		budget.CustomItems,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("update budget_items failed: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
