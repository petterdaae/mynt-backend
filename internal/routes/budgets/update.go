package budgets

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

	var budget types.Budget
	err := utils.ParseBody(c, &budget)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	err = database.Exec(
		"UPDATE budgets SET name = $2, color = $3 WHERE user_id = $1 AND id = $4",
		sub,
		budget.Name,
		budget.Color,
		budget.ID,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("update budgets failed: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
