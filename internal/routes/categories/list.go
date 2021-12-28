package categories

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Category struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	ParentID *int64  `json:"parent_id"`
	Color    *string `json:"color"`
	Ignore   *bool   `json:"ignore"`
}

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to connect to databse: %w", err))
		return
	}
	defer connection.Close()

	rows, err := connection.Query(
		"SELECT id, name, parent_id, color, ignore FROM categories WHERE user_id = $1 AND (deleted != TRUE OR deleted is NULL) ORDER BY name",
		sub,
	)

	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to query categories: %w", err))
		return
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.ParentID,
			&category.Color,
			&category.Ignore,
		)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to scan row: %w", err))
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, categories)
}
