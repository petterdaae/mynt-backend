package categories

import (
	"backend/internal/resources/categories"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	resource := categories.Configure(sub, database)
	result, err := resource.List()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to list categories: %w", err))
		return
	}

	c.JSON(http.StatusOK, result)
}
