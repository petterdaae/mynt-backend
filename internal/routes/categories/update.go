package categories

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

	var category types.Category
	err := utils.ParseBody(c, &category)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to unmarshal body: %w", err))
		return
	}

	err = database.Exec(
		"UPDATE categories SET name = $2, color = $3, ignore = $5 WHERE user_id = $1 AND id = $4",
		sub,
		category.Name,
		category.Color,
		category.ID,
		category.Ignore,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("insert failed: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
