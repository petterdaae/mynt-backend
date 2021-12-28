package routes

import (
	"backend/internal/middleware"
	"backend/internal/routes/accounts"
	"backend/internal/routes/auth"
	"backend/internal/routes/categories"
	"backend/internal/routes/spendings"
	"backend/internal/routes/synchronize"
	"backend/internal/routes/transactions"
	"backend/internal/routes/user"
	"backend/internal/utils"
	"encoding/json"
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
	r.GET("/auth/signout", auth.Signout)
	r.GET("/auth/demo", auth.Demo)

	// Public
	r.GET("/health", health)
	r.POST("/demo/reset", synchronize.ResetDemoAccount)

	// Private
	authGuard := middleware.Auth(database)
	r.GET("/authenticated", authGuard, authenticated)
	r.PUT("/user/secrets/sbanken", authGuard, user.UpdateSbankenSecrets)
	r.DELETE("/user/delete", authGuard, user.Delete)
	r.POST("/synchronize/sbanken", authGuard, synchronize.Sbanken)
	r.GET("/transactions", authGuard, transactions.List)
	r.GET("/accounts", authGuard, accounts.List)
	r.GET("/categories", authGuard, categories.List)
	r.POST("/categories", authGuard, categories.Create)
	r.DELETE("/categories", authGuard, categories.Delete)
	r.PUT("/categories/:id", authGuard, categories.Update)
	r.PUT("/transactions/update_category", authGuard, transactions.UpdateCategory)
	r.PUT("/transactions/update_custom_date", authGuard, transactions.UpdateCustomDate)
	r.GET("/spendings", authGuard, spendings.List)

	return r
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func authenticated(c *gin.Context) {
	c.String(http.StatusOK, "Authenticated")
}
