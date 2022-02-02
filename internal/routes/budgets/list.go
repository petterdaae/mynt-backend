package budgets

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
		"SELECT id, name, color FROM budgets WHERE user_id = $1 ORDER BY name, id",
		sub,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to query budgets: %w", err))
		return
	}
	defer rows.Close()

	var mainBudgetID *int64
	row, err := database.QueryRow("SELECT main_budget FROM users WHERE user_id = $1", sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to get main budget id: %w", err))
		return
	}
	err = row.Scan(&mainBudgetID)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to scan main budget id: %w", err))
		return
	}

	budgets := []types.Budget{}
	for rows.Next() {
		var budget types.Budget
		err := rows.Scan(
			&budget.ID,
			&budget.Name,
			&budget.Color,
		)
		budget.IsMainBudget = mainBudgetID != nil && budget.ID == *mainBudgetID
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to scan row: %w", err))
			return
		}
		budgets = append(budgets, budget)
	}

	c.JSON(http.StatusOK, budgets)
}
