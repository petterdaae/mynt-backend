package settings

import (
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	row, err := database.QueryRow(
		"SELECT main_budget FROM users WHERE id = $1",
		sub,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to query users: %w", err))
		return
	}

	settings := types.Settings{}
	err = row.Scan(&settings.MainBudgetID)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to scan row: %w", err))
		return
	}

	c.JSON(http.StatusOK, settings)
}
