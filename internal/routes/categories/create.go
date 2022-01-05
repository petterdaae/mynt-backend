package categories

import (
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateCategoryBody struct {
	Name     string `json:"name"`
	Color    string `json:"color"`
	ParentID *int64 `json:"parent_id"`
	Ignore   bool   `json:"ignore"`
}

func Create(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to connect to databse: %w", err))
		return
	}
	defer connection.Close()

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to read body: %w", err))
		return
	}

	var category CreateCategoryBody
	err = json.Unmarshal(body, &category)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to unmarshal body: %w", err))
		return
	}

	var id int64
	err = connection.QueryRow(
		"INSERT INTO categories (user_id, name, parent_id, color, ignore) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		sub,
		category.Name,
		category.ParentID,
		category.Color,
		category.Ignore,
	).Scan(&id)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("insert failed: %w", err))
		return
	}

	c.Status(http.StatusCreated)
}
