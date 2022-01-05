package categorizations

import (
	"backend/internal/types"
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	rows, err := database.Query(
		`SELECT * FROM list_categorizations($1, $2, $3)`,
		sub,
		c.Query("from_date"),
		c.Query("to_date"),
	)

	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	defer rows.Close()

	categorizations := []types.Categorization{}
	for rows.Next() {
		var categorization types.Categorization
		err := rows.Scan(
			&categorization.ID,
			&categorization.TransactionID,
			&categorization.Amount,
			&categorization.CategoryID,
		)
		if err != nil {
			utils.InternalServerError(c, err)
			return
		}
		categorizations = append(categorizations, categorization)
	}

	c.JSON(http.StatusOK, categorizations)
}
