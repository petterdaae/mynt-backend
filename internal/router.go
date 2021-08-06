package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Dependencies) health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

// SetupRoutes assigns functions to all the different routes
func SetupRoutes(d *Dependencies) *gin.Engine {
	r := gin.Default()

	r.GET("/health", d.health)

	return r
}
