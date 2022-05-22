package names

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
		"SELECT id, name, regex, replace_with FROM names WHERE user_id = $1 ORDER BY name, id",
		sub,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to query names: %w", err))
		return
	}
	defer rows.Close()

	names := []types.Name{}
	for rows.Next() {
		var name types.Name
		err := rows.Scan(
			&name.ID,
			&name.Name,
			&name.Fields.Regex,
			&name.Fields.ReplaceWith,
		)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to scan row: %w", err))
			return
		}
		names = append(names, name)
	}

	c.JSON(http.StatusOK, names)
}
