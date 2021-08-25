package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func SetCookie(c *gin.Context, name, value string, minutes int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   minutes * 60,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		Domain:   os.Getenv("COOKIE_DOMAIN"),
	}
	http.SetCookie(c.Writer, cookie)
}

func SetUnsafeCookie(c *gin.Context, name, value string, minutes int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   minutes * 60,
		Secure:   true,
		HttpOnly: false,
		Path:     "/",
		Domain:   os.Getenv("COOKIE_DOMAIN"),
	}
	http.SetCookie(c.Writer, cookie)
}

func Base64Encode(s string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

func CurrencyToInt(currency float64) int {
	currencyScaleInDatabase := 100
	return int(currency * float64(currencyScaleInDatabase))
}

func IntToCurrency(amount int) float64 {
	currencyScaleInDatabase := 100
	return float64(amount) / float64(currencyScaleInDatabase)
}
