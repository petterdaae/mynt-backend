package budgetitems

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteBudgetItemBody struct {
	ID int64 `json:"id"`
}

func Delete(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var body DeleteBudgetItemBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	err = database.Exec("DELETE FROM budget_items WHERE id = $1 AND user_id = $2", body.ID, sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to delete budget_items: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
