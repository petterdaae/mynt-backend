package budgets

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetMainBudgetBody struct {
	ID int64 `json:"id"`
}

func SetMainBudget(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var body SetMainBudgetBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	err = database.Exec(
		"UPDATE users SET main_budget = $1 WHERE user_id = $1",
		sub,
		body.ID, // TODO: Check that budget exists for user
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to update field main_budget in users: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
