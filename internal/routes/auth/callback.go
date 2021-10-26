package auth

import (
	"backend/internal/utils"
	"fmt"
	"net/http"
	"os"
	"time"

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
		utils.BadRequest(c, fmt.Errorf("state not found"))
		return
	}

	if c.Query("state") != state {
		utils.BadRequest(c, fmt.Errorf("state did not match"))
		return
	}

	// Exchange code for token
	oauth2Token, err := oauth2Config.Exchange(c, c.Request.URL.Query().Get("code"))
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to exchange code for token: %w", err))
		return
	}

	// Extract the ID Token from OAuth2 Token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		utils.InternalServerError(c, fmt.Errorf("no id_token field in oauth2 token"))
		return
	}

	// Parse and verify ID Token payload
	idToken, err := verifier.Verify(c, rawIDToken)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to verify id token: %w", err))
		return
	}

	// Verify nonce
	nonce, err := c.Cookie("nonce")
	if err != nil {
		utils.BadRequest(c, fmt.Errorf("nonce not found"))
		return
	}

	if idToken.Nonce != nonce {
		utils.InternalServerError(c, fmt.Errorf("nonce did not match"))
		return
	}

	// Extract identity
	var claims struct {
		Sub string `json:"sub"`
	}

	if err = idToken.Claims(&claims); err != nil {
		return
	}

	// Create user if not exists
	err = createUserIfNotExists(c, claims.Sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to create user: %w", err))
		return
	}

	// Create an auth token
	token, err := utils.CreateToken(c, claims.Sub)
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
	_, err = connection.Exec("INSERT INTO users (id) VALUES ($1)", sub)
	return err
}

func unixTimeInOneHour() int64 {
	return time.Now().Add(time.Hour).Unix()
}
