package sync

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Post(c *gin.Context) {
	c.String(http.StatusOK, "Good")
}
