package synchronize

import (
	"backend/internal/resources/sbanken"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Sbanken(c *gin.Context) {
	sub := c.GetString("sub")
	database, _ := c.MustGet("database").(*utils.Database)

	sbankenResource := sbanken.Configure(sub, database)

	err := sbankenResource.Synchronize()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("sbanken synch failed: %w", err))
		return
	}

	c.String(http.StatusOK, "Success")
}
