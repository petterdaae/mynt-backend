package utils

import (
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
