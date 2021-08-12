package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InternalServerError(err error, c *gin.Context) {
	log.WithField("err", err).Error("internal server error")
	c.String(http.StatusInternalServerError, "Internal server error")
}

func NotFound(err error, c *gin.Context) {
	log.WithField("err", err).Warn("not found")
	c.String(http.StatusNotFound, "Not found")
}

func RandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func SetCookie(c *gin.Context, name string, value string, seconds int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   seconds,
		Secure:   c.Request.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, cookie)
}
