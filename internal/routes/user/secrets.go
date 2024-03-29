package user

import (
	"backend/internal/utils"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sbankenSecrets struct {
	ClientID     string `json:"sbanken_client_id"`
	ClientSecret string `json:"sbanken_client_secret"`
}

func UpdateSbankenSecrets(c *gin.Context) {
	sub, _ := c.MustGet("sub").(string)
	database, _ := c.MustGet("database").(*utils.Database)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	secrets := &sbankenSecrets{}
	err = json.Unmarshal(body, secrets)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}
	defer connection.Close()

	_, err = connection.Exec(
		"UPDATE users SET sbanken_client_id = $1, sbanken_client_secret = $2 WHERE id = $3",
		secrets.ClientID,
		secrets.ClientSecret,
		sub,
	)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
