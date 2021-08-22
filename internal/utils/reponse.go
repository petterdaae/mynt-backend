package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO : Is this actually necessary?
func AbortWithError(c *gin.Context, code int, abortError error) {
	err := c.AbortWithError(http.StatusInternalServerError, abortError)
	if err != nil {
		panic(err)
	}
}

func InternalServerError(c *gin.Context, err error) {
	AbortWithError(c, http.StatusInternalServerError, err)
}

func Unauthorized(c *gin.Context, err error) {
	AbortWithError(c, http.StatusUnauthorized, err)
}

func BadRequest(c *gin.Context, err error) {
	AbortWithError(c, http.StatusBadRequest, err)
}
