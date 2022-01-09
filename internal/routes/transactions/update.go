package transactions

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateBody struct {
	ID         string  `json:"id"`
	CustomDate *string `json:"customDate"`
}

func Update(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var body UpdateBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	err = database.Exec(
		"SELECT update_transaction($1, $2, $3)",
		sub,
		body.ID,
		body.CustomDate,
	)

	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
