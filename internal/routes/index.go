package routes

import (
	middleware "mynt/internal/middleware"
	auth "mynt/internal/routes/auth"
	utils "mynt/internal/utils"
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
	r.GET("/authenticated", authenticated)

	return r
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func authenticated(c *gin.Context) {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		c.String(http.StatusUnauthorized, "Not authenticated")
		return
	}

	_, err = utils.ValidateToken(c, cookie)
	if err != nil {
		c.String(http.StatusUnauthorized, "Not authenticated")
		return
	}

	c.String(http.StatusOK, "Authenticated")
}
