package names

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateNameBody struct {
	Name   string `json:"name"`
	Fields struct {
		Regex       string `json:"regex"`
		ReplaceWith string `json:"replaceWith"`
	} `json:"fields"`
}

type CreatedNameResponse struct {
	ID int64 `json:"id"`
}

func Create(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var body CreateNameBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	var newNameID int64
	row, err := database.QueryRow(
		"INSERT INTO names (user_id, name, regex, replace_with) VALUES ($1, $2, $3, $4) RETURNING id",
		sub,
		body.Name,
		body.Fields.Regex,
		body.Fields.ReplaceWith,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to insert new name: %w", err))
		return
	}
	err = row.Scan(&newNameID)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to scan new name id: %w", err))
		return
	}

	response := CreatedNameResponse{
		ID: newNameID,
	}

	c.JSON(http.StatusCreated, response)
}
