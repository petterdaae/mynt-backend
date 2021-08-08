package sync

import (
	"fmt"
	"mynt/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Post(c *gin.Context) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}

	fmt.Println("user id :", userId)

	c.String(http.StatusOK, "Healthy")
}
