package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
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
		RedirectURL:  os.Getenv("CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("GOOGLE_AUTH_CLIENT_ID")}) // TODO : put in context as well?

	c.Set("oauth2Config", oauth2Config)
	c.Set("oidcProvider", provider)
	c.Set("oidcIDTokenVerifier", verifier)
	c.Next()
}

func CreateToken(c *gin.Context, sub string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SIGNING_SECRET")
	if secret == "" {
		panic("jwt singning secret is empty")
	}

	return token.SignedString([]byte(secret))
}

func ValidateToken(c *gin.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGNING_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}

	return "", fmt.Errorf("token is missing sub")
}
