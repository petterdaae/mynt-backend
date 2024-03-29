package categories

import (
	"backend/internal/resources/categories"
	"backend/internal/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteCategoryBody struct {
	ID int64 `json:"id"`
}

func Delete(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	resource := categories.Configure(sub, database)

	categoryIDString := c.Param("category_id")
	categoryID, err := strconv.Atoi(categoryIDString)
	if err != nil {
		utils.BadRequest(c, fmt.Errorf("failed to parse category_id: %w", err))
		return
	}

	err = resource.Delete(int64(categoryID))

	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to delete category: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
