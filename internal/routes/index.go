package routes

import (
	middleware "mynt/internal/middleware"
	auth "mynt/internal/routes/auth"
	sync "mynt/internal/routes/sync"
	utils "mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes assigns functions to all the different routes
func SetupRoutes(database *utils.Database) *gin.Engine {
	r := gin.Default()

	// Add cors middleware
	cors := middleware.Cors()
	r.Use(cors)

	// Add database to context
	r.Use(func(c *gin.Context) {
		c.Set("database", database)
		c.Next()
	})

	// Configure oauth2
	r.Use(auth.ConfigureOauth2)

	// Oauth2 routes
	r.GET("/auth/redirect", auth.HandleRedirect)
	r.GET("/auth/callback", auth.HandleOauth2Callback)

	// Public routes
	r.GET("/health", health)

	// Private routes
	r.GET("/authenticated", authenticated)
	r.POST("/sync", sync.Post)

	return r
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func authenticated(c *gin.Context) {
	c.String(http.StatusOK, "Authenticated")
}
