package categorizations

import (
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	TransactionID   string `json:"transactionId"`
	Categorizations []struct {
		CategoryID int64 `json:"categoryId"`
		Amount     int64 `json:"amount"`
	} `json:"categorizations"`
}

func UpdateCategorizationsForTransaction(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	var body RequestBody
	err := utils.ParseBody(c, &body)
	if err != nil {
		utils.BadRequest(c, err)
		return
	}

	transaction, err := GetTransaction(database, sub, body.TransactionID)
	if err != nil {
		utils.BadRequest(c, err)
		return
	}

	err = RemoveOldCategorizations(database, sub, transaction.ID)
	if err != nil {
		utils.BadRequest(c, err)
		return
	}

	err = ValidateCategorizations(body, transaction)
	if err != nil {
		utils.BadRequest(c, err)
		return
	}

	err = CreateCategorizations(database, sub, transaction.ID, body)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func GetTransaction(database *utils.Database, sub, transactionID string) (*types.Transaction, error) {
	var transaction types.Transaction
	row, err := database.QueryRow("SELECT id, amount FROM transactions WHERE user_id = $1 AND id = $2", sub, transactionID)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&transaction.ID, &transaction.Amount)
	return &transaction, err
}

func RemoveOldCategorizations(databse *utils.Database, sub, transactionID string) error {
	return databse.Exec("DELETE FROM transactions_to_categories WHERE user_id = $1 AND transaction_id = $2", sub, transactionID)
}

func ValidateCategorizations(body RequestBody, transaction *types.Transaction) error {
	var sum int64
	for _, categorization := range body.Categorizations {
		sum += categorization.Amount
	}

	if sum != transaction.Amount {
		return fmt.Errorf("the sum of categorizations has to match the transaction amount")
	}

	return nil
}

func CreateCategorizations(database *utils.Database, sub, transactionID string, body RequestBody) error {
	for _, categorization := range body.Categorizations {
		err := database.Exec(
			"INSERT INTO categorizations (user_id, transaction_id, category_id, amount) VALUES ($1, $2, $3, $4)",
			sub,
			transactionID,
			categorization.CategoryID,
			categorization.Amount,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
