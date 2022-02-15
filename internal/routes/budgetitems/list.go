package budgetitems

import (
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	rows, err := database.Query(
		"SELECT id, budget_id, category_id, monthly_amount, name, custom_items, kind "+
			"FROM budget_items WHERE user_id = $1",
		sub,
	)

	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to query budget_items: %w", err))
		return
	}
	defer rows.Close()

	budgetItems := []types.BudgetItem{}
	for rows.Next() {
		var budgetItem types.BudgetItem
		err := rows.Scan(
			&budgetItem.ID,
			&budgetItem.BudgetID,
			&budgetItem.CategoryID,
			&budgetItem.MonthlyAmount,
			&budgetItem.Name,
			&budgetItem.CustomItems,
			&budgetItem.Kind,
		)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to scan row: %w", err))
			return
		}
		budgetItems = append(budgetItems, budgetItem)
	}

	c.JSON(http.StatusOK, budgetItems)
}
