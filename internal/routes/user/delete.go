package user

import (
	"fmt"
	"mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	sub := c.GetString("sub")
	database, _ := c.MustGet("database").(*utils.Database)

	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to connect to database: %w", err))
		return
	}
	defer connection.Close()

	_, err = connection.Exec(`DELETE FROM transactions WHERE user_id = $1`, sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to delete transactions: %w", err))
		return
	}

	_, err = connection.Exec(`DELETE FROM accounts WHERE user_id = $1`, sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to delete accounts: %w", err))
		return
	}

	_, err = connection.Exec(`DELETE FROM categories WHERE user_id = $1`, sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to delete categories: %w", err))
		return
	}

	c.String(http.StatusOK, "Success")
}
