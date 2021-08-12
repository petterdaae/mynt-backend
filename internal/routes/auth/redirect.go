package auth

import (
	"mynt/internal/utils"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func Redirect(c *gin.Context) {
	oauth2Config, _ := c.MustGet("oauth2Config").(*oauth2.Config)

	state, err := utils.RandomString(16)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	nonce, err := utils.RandomString(16)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	utils.SetCookie(c, "state", state, 60)
	utils.SetCookie(c, "nonce", nonce, 60)

	c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)))
}
