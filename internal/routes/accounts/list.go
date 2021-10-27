package accounts

import (
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("database connection failed: %w", err))
		return
	}
	defer connection.Close()

	rows, err := connection.Query("SELECT id, account_number, name, available, balance FROM accounts WHERE user_id = $1", sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("database query failed: %w", err))
		return
	}
	defer rows.Close()

	accounts := []types.Account{}

	for rows.Next() {
		var account types.Account
		err := rows.Scan(&account.ID, &account.AccountNumber, &account.Name, &account.Available, &account.Balance)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("database scan failed: %w", err))
			return
		}
		accounts = append(accounts, account)
	}

	c.JSON(http.StatusOK, accounts)
}
