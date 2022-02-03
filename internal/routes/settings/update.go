package settings

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

	var body types.Settings
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	err = database.Exec(
		"UPDATE users SET main_budget = $2 WHERE id = $1",
		sub,
		body.MainBudgetID,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to update users: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
