package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func ParseBody(c *gin.Context, parsed interface{}) error {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	err = json.Unmarshal(rawBody, parsed)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return nil
}
