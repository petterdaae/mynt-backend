package auth

import (
	"mynt/internal/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Signout(c *gin.Context) {
	utils.RemoveCookie(c, "auth_token")
	c.Redirect(http.StatusFound, os.Getenv("REDIRECT_TO_FRONTEND"))
}
