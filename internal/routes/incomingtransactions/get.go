package incomingtransactions

import (
	"backend/internal/resources/incomingtransactions"
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")
	resource := incomingtransactions.Configure(sub, database)

	result, err := resource.ListAll()

	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
