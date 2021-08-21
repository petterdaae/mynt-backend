package routes

import (
	"mynt/internal/middleware"
	"mynt/internal/routes/auth"
	"mynt/internal/routes/synchronize"
	"mynt/internal/routes/user"
	"mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes assigns functions to all the different routes
func SetupRoutes(database *utils.Database) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.Cors())

	// Dependencies
	r.Use(func(c *gin.Context) {
		c.Set("database", database)
		c.Next()
	})

	// Authentication
	r.Use(utils.ConfigureOauth2)
	r.GET("/auth/redirect", auth.Redirect)
	r.GET("/auth/callback", auth.Callback)

	// Public
	r.GET("/health", health)

	// Private
	auth := middleware.Auth(database)
	r.GET("/authenticated", auth, authenticated)
	r.PUT("/user/secrets/sbanken", auth, user.UpdateSbankenSecrets)
	r.POST("/synchronize/sbanken", auth, synchronize.Sbanken)

	return r
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func authenticated(c *gin.Context) {
	c.String(http.StatusOK, "Authenticated")
}
