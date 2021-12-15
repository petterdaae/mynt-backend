package transactions

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateCustomDateBody struct {
	Id         string  `json:"id"`
	CustomDate *string `json:"customDate"`
}

func UpdateCustomDate(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var body UpdateCustomDateBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to unmarshal body: %w", err))
		return
	}

	err = database.Exec(
		"UPDATE transactions SET custom_date = $3 WHERE user_id = $1 AND id = $2",
		sub,
		body.Id,
		body.CustomDate,
	)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("insert failed: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
