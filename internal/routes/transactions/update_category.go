package transactions

import (
	"fmt"
	"mynt/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	TransactionID   string           `json:"transaction_id"`
	Categorizations []Categorization `json:"categorizations"`
}

type Categorization struct {
	CategoryID int64 `json:"category_id"`
	Amount     int64 `json:"amount"`
}

func UpdateCategory(c *gin.Context) {
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

	err = RemoveOldCategorization(database, sub, transaction.ID)
	if err != nil {
		utils.BadRequest(c, err)
		return
	}

	err = ValidateCategorizations(body.Categorizations, transaction)
	if err != nil {
		utils.BadRequest(c, err)
		return
	}

	err = CreateCategorizations(database, sub, transaction.ID, body.Categorizations)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func GetTransaction(database *utils.Database, sub, transactionID string) (*Transaction, error) {
	var transaction Transaction
	row, err := database.QueryRow("SELECT id, amount FROM transactions WHERE user_id = $1 AND id = $2", sub, transactionID)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&transaction.ID, &transaction.Amount)
	return &transaction, err
}

func RemoveOldCategorization(databse *utils.Database, sub, transactionID string) error {
	return databse.Exec("DELETE FROM transactions_to_categories WHERE user_id = $1 AND transaction_id = $2", sub, transactionID)
}

func ValidateCategorizations(categorizations []Categorization, transaction *Transaction) error {
	var sum int64
	for _, categorization := range categorizations {
		sum += categorization.Amount
	}

	if sum != transaction.Amount {
		return fmt.Errorf("the sum of categorizations has to match the transaction amount")
	}

	return nil
}

func CreateCategorizations(database *utils.Database, sub, transactionID string, categorizations []Categorization) error {
	for _, categorization := range categorizations {
		err := database.Exec(
			"INSERT INTO transactions_to_categories (user_id, transaction_id, category_id, amount) VALUES ($1, $2, $3, $4)",
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
