package budgets

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateBudgetBody struct {
	Name         string `json:"name"`
	Color        string `json:"color"`
	IsMainBudget bool   `json:"isMainBudget"`
}

type CreatedBudgetResponse struct {
	ID int64 `json:"id"`
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

	var newBudgetID int64
	row, err := database.QueryRow(
		"INSERT INTO budgets (user_id, name, color) VALUES ($1, $2, $3) RETURNING id",
		sub,
		body.Name,
		body.Color,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to insert new budget: %w", err))
		return
	}
	err = row.Scan(&newBudgetID)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to scan new budget id: %w", err))
		return
	}

	response := CreatedBudgetResponse{
		ID: newBudgetID,
	}

	c.JSON(http.StatusCreated, response)
}
