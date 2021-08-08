package routes

import (
	middleware "mynt/internal/middleware"
	sync "mynt/internal/routes/sync"
	utils "mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes assigns functions to all the different routes
func SetupRoutes(database *utils.Database) *gin.Engine {
	guard := middleware.AuthGuard(database)

	r := gin.Default()

	// Add cors middleware
	cors := middleware.Cors()
	r.Use(cors)

	// Add database to context
	r.Use(func(c *gin.Context) {
		c.Set("database", database)
		c.Next()
	})

	// Public routes
	r.GET("/health", health)

	// Private routes
	r.GET("/authenticated", guard, authenticated)
	r.POST("/sync", guard, sync.Post)

	return r
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func authenticated(c *gin.Context) {
	c.String(http.StatusOK, "Authenticated")
}
