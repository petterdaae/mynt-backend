package transactions

import (
	"backend/internal/resources/transactions"
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")
	resource := transactions.Configure(sub, database)

	from := c.Query("from_date")
	to := c.Query("to_date")
	result, err := resource.List(from, to)

	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
