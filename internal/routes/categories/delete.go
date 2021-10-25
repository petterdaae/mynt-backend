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

	err = recursiveDelete(body.ID, sub, database)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("error occurred when deleting category: %w", err))
		return
	}

	c.Status(http.StatusOK)
}

func recursiveDelete(categoryID int64, sub string, database *utils.Database) error {
	rows, err := database.Query("SELECT id FROM categories WHERE parent_id = $1 AND user_id = $2", categoryID, sub)
	if err != nil {
		return fmt.Errorf("failed to find sub category id's: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		err = recursiveDelete(id, sub, database)
		if err != nil {
			return fmt.Errorf("recursive call failed: %w", err)
		}
	}

	err = database.Exec("DELETE FROM categories WHERE id = $1 AND user_id = $2", categoryID, sub)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	err = database.Exec("DELETE from transactions_to_categories WHERE category_id = $1", categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete rows from transactions_to_categories: %w", err)
	}

	return nil
}
