package auth

import (
	"backend/internal/utils"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func Redirect(c *gin.Context) {
	oauth2Config, _ := c.MustGet("oauth2Config").(*oauth2.Config)

	stateAndNonceLength := 16

	state, err := utils.RandomString(stateAndNonceLength)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	nonce, err := utils.RandomString(stateAndNonceLength)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	cookieMaxAgeInMinutes := 60
	utils.SetCookieWithoutDomain(c, "state", state, cookieMaxAgeInMinutes)
	utils.SetCookieWithoutDomain(c, "nonce", nonce, cookieMaxAgeInMinutes)

	c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)))
}
