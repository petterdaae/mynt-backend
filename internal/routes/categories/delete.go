package categories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteCategoryBody struct {
	ID int64 `json:"id"`
}

func Delete(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to connect to databse: %w", err))
		return
	}
	defer connection.Close()

	rawBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to read body: %w", err))
		return
	}

	var body DeleteCategoryBody
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to unmarshal body: %w", err))
		return
	}

	_, err = connection.Exec("UPDATE categories SET deleted = TRUE WHERE user_id = $1 AND id = $2", sub, body.ID)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to query categories: %w", err))
		return
	}

	c.Status(http.StatusOK)
}
