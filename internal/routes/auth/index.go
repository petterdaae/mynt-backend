/*

Based on documentation from:
- https://github.com/coreos/go-oidc
- https://developers.google.com/identity/protocols/oauth2/openid-connect

*/

package auth

import (
	"fmt"
	"mynt/internal/utils"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func ConfigureOauth2(c *gin.Context) {
	provider, err := oidc.NewProvider(c, "https://accounts.google.com")
	if err != nil {
		panic(err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_AUTH_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("GOOGLE_AUTH_CLIENT_ID")}) // TODO : put in context as well?

	c.Set("oauth2Config", oauth2Config)
	c.Set("oidcProvider", provider)
	c.Set("oidcIDTokenVerifier", verifier)
	c.Next()
}

func HandleRedirect(c *gin.Context) {
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

func HandleOauth2Callback(c *gin.Context) {
	oauth2Config, _ := c.MustGet("oauth2Config").(*oauth2.Config)
	verifier, _ := c.MustGet("oidcIDTokenVerifier").(*oidc.IDTokenVerifier)

	// Verify state
	state, err := c.Cookie("state")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("state not found"))
		return
	}

	if c.Query("state") != state {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("state did not match"))
		return
	}

	// Exchange code for token
	oauth2Token, err := oauth2Config.Exchange(c, c.Request.URL.Query().Get("code"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to exchange code for token: %w", err))
		return
	}

	// Extract the ID Token from OAuth2 Token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("no id_token field in oauth2 token"))
		return
	}

	// Parse and verify ID Token payload
	idToken, err := verifier.Verify(c, rawIDToken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to verify id token: %w", err))
		return
	}

	// Verify nonce
	nonce, err := c.Cookie("nonce")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("nonce not found"))
		return
	}

	if idToken.Nonce != nonce {
		c.AbortWithError(500, err)
		return
	}

	// Extract identity
	var claims struct {
		Sub string `json:"sub"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return
	}

	fmt.Println("sub:", claims.Sub)

	// TODO : Create jwt token and put it in cookie
	// TODO : Redirect to web app
	c.Redirect(http.StatusFound, "http://localhost:3000/authenticated")
}
