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

func Callback(c *gin.Context) {
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

	// TODO : Create user if it doesn't exist
	err = createUserIfNotExists(c, claims.Sub)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	token, err := utils.CreateToken(c, claims.Sub)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	utils.SetCookie(c, "auth_token", token, 60)
	c.Redirect(http.StatusFound, os.Getenv("REDIRECT_TO_FRONTEND"))
}

func createUserIfNotExists(c *gin.Context, sub string) error {
	// Connect to database
	database, _ := c.MustGet("database").(*utils.Database)
	connection, err := database.Connect()
	if err != nil {
		return err
	}
	defer connection.Close()

	// Check if user exists
	userExists := false
	rows, err := connection.Query("SELECT id FROM users WHERE id = $1", sub)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		userExists = true
	}

	// Do nothing if user exists
	if userExists {
		return nil
	}

	// Create user
	_, err = connection.Query("INSERT INTO users (id) VALUES ($1)", sub)
	return err
}
