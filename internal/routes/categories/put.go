package categories

import (
	"backend/internal/resources/categories"
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Put(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var draftCategory types.DraftCategory
	err := utils.ParseBody(c, &draftCategory)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	resource := categories.Configure(sub, database)

	categoryIDString := c.Param("category_id")
	categoryID, err := strconv.Atoi(categoryIDString)
	if err != nil {
		utils.BadRequest(c, fmt.Errorf("failed to parse category_id: %w", err))
		return
	}

	err = resource.Update(int64(categoryID), draftCategory)

	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to insert category: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
