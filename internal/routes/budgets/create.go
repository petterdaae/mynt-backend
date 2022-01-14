package budgets

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateBudgetBody struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func Create(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var body CreateBudgetBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	err = database.Exec(
		"INSERT INTO budgets (user_id, name, color) VALUES ($1, $2, $3)",
		sub,
		body.Name,
		body.Color,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to insert new budget: %w", err))
		return
	}

	c.Status(http.StatusCreated)
}
