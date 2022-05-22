package names

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

	var name types.Name
	err := utils.ParseBody(c, &name)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	err = database.Exec(
		"UPDATE names SET name = $3, regex = $4, replace_with = $5 WHERE user_id = $1 AND id = $2",
		sub,
		name.ID,
		name.Name,
		name.Fields.Regex,
		name.Fields.ReplaceWith,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("update names failed: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
