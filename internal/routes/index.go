package routes

import (
	"backend/internal/middleware"
	"backend/internal/routes/accounts"
	"backend/internal/routes/auth"
	"backend/internal/routes/budgetitems"
	"backend/internal/routes/budgets"
	"backend/internal/routes/categories"
	"backend/internal/routes/categorizations"
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

	// Unauthenticated
	r.GET("/health", health)
	r.POST("/demo/reset", synchronize.ResetDemoAccount)

	// Authenticated
	authGuard := middleware.Auth(database)
	r.GET("/authenticated", authGuard, authenticated)
	r.PUT("/user/secrets/sbanken", authGuard, user.UpdateSbankenSecrets)
	r.POST("/synchronize/sbanken", authGuard, synchronize.Sbanken)

	r.GET("/transactions", authGuard, transactions.List)
	r.PUT("/transactions", authGuard, transactions.Update)

	r.GET("/categorizations", authGuard, categorizations.List)
	r.PUT("/categorizations", authGuard, categorizations.UpdateCategorizationsForTransaction)

	r.GET("/categories", authGuard, categories.List)
	r.PUT("/categories", authGuard, categories.Update)
	r.POST("/categories", authGuard, categories.Create)
	r.DELETE("/categories", authGuard, categories.Delete)

	r.GET("/accounts", authGuard, accounts.List)

	r.GET("/budgets", authGuard, budgets.List)
	r.POST("/budgets", authGuard, budgets.Create)
	r.PUT("/budgets", authGuard, budgets.Update)
	r.DELETE("/budgets", authGuard, budgets.Delete)

	r.GET("budget_items", authGuard, budgetitems.List)
	r.POST("budget_items", authGuard, budgetitems.Create)
	r.PUT("budget_items", authGuard, budgetitems.Update)
	r.DELETE("budget_items", authGuard, budgetitems.Delete)

	return r
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func authenticated(c *gin.Context) {
	c.String(http.StatusOK, "Authenticated")
}
