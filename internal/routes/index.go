package routes

import (
	"mynt/internal/middleware"
	"mynt/internal/routes/accounts"
	"mynt/internal/routes/auth"
	"mynt/internal/routes/synchronize"
	"mynt/internal/routes/transactions"
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
	authGuard := middleware.Auth(database)
	r.GET("/authenticated", authGuard, authenticated)
	r.PUT("/user/secrets/sbanken", authGuard, user.UpdateSbankenSecrets)
	r.POST("/synchronize/sbanken", authGuard, synchronize.Sbanken)
	r.GET("/transactions", authGuard, transactions.Get)
	r.GET("/accounts", authGuard, accounts.Get)
	r.DELETE("/synchronize/delete", authGuard, synchronize.Delete)

	return r
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func authenticated(c *gin.Context) {
	c.String(http.StatusOK, "Authenticated")
}
