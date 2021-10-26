package middleware

import (
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(database *utils.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract auth token from cookie
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			utils.Unauthorized(c, err)
			return
		}

		// Validate auth token
		sub, err := utils.ValidateToken(c, cookie)
		if err != nil {
			utils.Unauthorized(c, err)
			return
		}

		// Check if used exists
		connection, err := database.Connect()
		if err != nil {
			utils.InternalServerError(c, err)
			return
		}
		defer connection.Close()

		exists := false
		rows, err := connection.Query("SELECT * FROM users WHERE id = $1", sub)
		if err != nil {
			utils.InternalServerError(c, err)
			return
		}
		if rows.Err() != nil {
			utils.InternalServerError(c, rows.Err())
			return
		}

		for rows.Next() {
			exists = true
		}

		// Abort if user doesn't exist
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Set user id in context if the token is valid and the user exists
		c.Set("sub", sub)
		c.Next()
	}
}
