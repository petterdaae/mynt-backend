package accounts

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateBody struct {
	ID       string `json:"id"`
	Favorite bool   `json:"favorite"`
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
		"UPDATE accounts SET favorite = $3 WHERE user_id = $1 AND id = $2",
		sub,
		body.ID,
		body.Favorite,
	)

	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
