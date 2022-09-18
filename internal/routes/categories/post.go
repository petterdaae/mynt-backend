package categories

import (
	"backend/internal/resources/categories"
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreatedCategoryResponse struct {
	ID int64 `json:"id"`
}

func Post(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var draftCategory types.DraftCategory
	err := utils.ParseBody(c, &draftCategory)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to parse body: %w", err))
		return
	}

	resource := categories.Configure(sub, database)
	id, err := resource.Create(draftCategory)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to create category: %w", err))
		return
	}

	response := CreatedCategoryResponse{
		ID: id,
	}

	c.JSON(http.StatusCreated, response)
}
