package categories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateCategoryBody struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
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

	_, err = connection.Exec(
		"INSERT INTO categories (user_id, name, parent_id) VALUES ($1, $2, $3)",
		sub,
		category.Name,
		category.ParentID,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("insert failed: %w", err))
		return
	}

	c.Status(http.StatusCreated)
}
