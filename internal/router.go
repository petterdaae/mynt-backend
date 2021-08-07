package internal

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (d *Dependencies) health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// SetupRoutes assigns functions to all the different routes
func SetupRoutes(d *Dependencies) *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())

	r.GET("/health", d.health)

	return r
}
