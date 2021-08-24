package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO : Is this actually necessary?
func AbortWithError(c *gin.Context, code int, abortError error) int {
	err := c.AbortWithError(http.StatusInternalServerError, abortError)
	if err != nil {
		// This error is not actually an error, its just
		// c.Error() returning the struct that gin uses
		// to store the error internally.
		return 1
	}
	return 0
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
