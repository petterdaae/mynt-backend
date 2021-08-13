package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func SetCookie(c *gin.Context, name string, value string, minutes int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   minutes * 60,
		Secure:   c.Request.TLS != nil,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(c.Writer, cookie)
}
