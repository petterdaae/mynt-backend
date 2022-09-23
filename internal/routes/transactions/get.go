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

	result, err := resource.List(c.Query("from"), c.Query("to"))

	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
