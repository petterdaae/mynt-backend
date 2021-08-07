package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Dependencies) health(c *gin.Context) {
	c.String(http.StatusOK, "Healthy")
}

func (d *Dependencies) authenticated(c *gin.Context) {
	c.String(http.StatusOK, "Authenticated")
}

// SetupRoutes assigns functions to all the different routes
func SetupRoutes(d *Dependencies) *gin.Engine {
	jwtMiddleware := createJwtMiddleware()
	corsMiddleware := createCorsMiddleware()

	r := gin.Default()
	r.Use(corsMiddleware)

	r.GET("/health", d.health)
	r.GET("/authenticated", checkJwt(jwtMiddleware), d.authenticated)

	return r
}
