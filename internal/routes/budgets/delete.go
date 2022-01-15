package budgets

import (
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteBudgetBody struct {
	ID int64 `json:"id"`
}

func Delete(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var body DeleteBudgetBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	err = database.Exec("DELETE FROM budgets WHERE id = $1 AND user_id = $2", body.ID, sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("delete budgets failed: %w", err))
		return
	}

	err = deleteBudgetItems(body.ID, sub, database)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("error occured while deleting budget items: %w", err))
		return
	}

	c.Status(http.StatusOK)
}

func deleteBudgetItems(budgetId int64, sub string, database *utils.Database) error {
	err := database.Exec("DELETE FROM budget_items WHERE budget_id = $1 AND user_id = $2", budgetId, sub)
	if err != nil {
		return fmt.Errorf("failed to delete budget items: %w", err)
	}

	return nil
}
