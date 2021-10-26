package auth

import (
	"backend/internal/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Demo(c *gin.Context) {
	demoSub := "demo"

	err := createUserIfNotExists(c, demoSub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to create user: %w", err))
		return
	}

	// Create an auth token
	token, err := utils.CreateToken(c, demoSub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to create token: %w", err))
		return
	}

	tokenCookieAge := 60
	expiryCookieAge := 120
	utils.SetCookie(c, "auth_token", token, tokenCookieAge)
	utils.SetUnsafeCookie(c, "auth_expiry", fmt.Sprintf("%v", unixTimeInOneHour()), expiryCookieAge)
	c.Redirect(http.StatusFound, os.Getenv("REDIRECT_TO_FRONTEND"))
}
