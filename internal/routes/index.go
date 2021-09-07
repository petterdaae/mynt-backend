package routes

import (
	"encoding/json"
	"mynt/internal/middleware"
	"mynt/internal/routes/accounts"
	"mynt/internal/routes/auth"
	"mynt/internal/routes/categories"
	"mynt/internal/routes/synchronize"
	"mynt/internal/routes/transactions"
	"mynt/internal/routes/user"
	"mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes assigns functions to all the different routes
func SetupRoutes(database *utils.Database) *gin.Engine {
	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		delete(param.Keys, "oauth2Config")
		delete(param.Keys, "database")
		delete(param.Keys, "oidcIDTokenVerifier")
		delete(param.Keys, "oidcProvider")
		bytes, _ := json.Marshal(map[string]interface{}{
			"level":   utils.LevelFromStatusCode(param.StatusCode),
			"method":  param.Method,
			"path":    param.Path,
			"status":  param.StatusCode,
			"latency": param.Latency,
			"ip":      param.ClientIP,
			"context": param.Keys,
			"error":   param.ErrorMessage,
		})
		return string(bytes) + "\n"
	}))

	r.Use(gin.Recovery())

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
	r.GET("/transactions", authGuard, transactions.List)
	r.GET("/accounts", authGuard, accounts.List)
	r.DELETE("/synchronize/delete", authGuard, synchronize.Delete)
	r.GET("/categories", authGuard, categories.List)

	return r
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func authenticated(c *gin.Context) {
	c.String(http.StatusOK, "Authenticated")
}
