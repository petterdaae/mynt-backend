package transactions

import (
	"mynt/internal/utils"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	ID int64 `json:"id"`
}

func UpdateCategory(c *gin.Context) {
	// "-q-3rdatabase, _ := c.MustGet("database").(*utils.Database)
	// sub := c.GetString("sub")

	var body RequestBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}
}
