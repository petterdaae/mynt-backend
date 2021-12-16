package categories

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateCategoryBody struct {
	Name   string `json:"name"`
	Color  string `json:"color"`
	Ignore bool   `json:"ignore"`
}

func Update(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var category UpdateCategoryBody
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
		c.Param("id"),
		category.Ignore,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("insert failed: %w", err))
		return
	}

	updatedCategory := Category{
		Name:   category.Name,
		Color:  &category.Color,
		Ignore: &category.Ignore,
	}

	row, err := database.QueryRow(
		"SELECT id, parent_id FROM categories WHERE user_id = $1 AND id = $2",
		sub,
		c.Param("id"),
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to query parent_id: %w", err))
		return
	}
	err = row.Scan(&updatedCategory.ID, &updatedCategory.ParentID)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to scan parent_id: %w", err))
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}
