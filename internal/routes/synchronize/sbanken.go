package synchronize

import (
	aResource "backend/internal/resources/accounts"
	itResource "backend/internal/resources/incomingtransactions"
	"backend/internal/resources/sbanken"
	tResource "backend/internal/resources/transactions"
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func Sbanken(c *gin.Context) {
	sub := c.GetString("sub")
	database, _ := c.MustGet("database").(*utils.Database)

	sbankenResource := sbanken.Configure(sub, database)
	transactionsResource := tResource.Configure(sub, database)
	accountsResource := aResource.Configure(sub, database)
	incomingTransactionsResource := itResource.Configure(sub, database)

	accounts, err := sbankenResource.GetAccounts()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to get sbanken accounts: %w", err))
		return
	}

	var wg sync.WaitGroup
	errors := make(chan error)

	for i := range accounts.Items {
		wg.Add(1)
		go func(account sbanken.Account) {
			err = synchronizeAccount(
				account,
				sbankenResource,
				transactionsResource,
				accountsResource,
				incomingTransactionsResource,
			)

			if err != nil {
				errors <- err
			}

			wg.Done()
		}(accounts.Items[i])
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		utils.InternalServerError(c, fmt.Errorf("failed to synchronize some accounts, sample error: %w", err))
		return
	}

	c.String(http.StatusOK, "Success")
}

func synchronizeAccount(
	account sbanken.Account,
	sbankenResource sbanken.Resource,
	transactionsResource tResource.Resource,
	accountsResource aResource.Resource,
	incomingTransactionsResource itResource.Resource,
) error {
	err := accountsResource.CreateIfNotExists(&types.Account{
		ID:            "sbanken:" + account.AccountID,
		AccountNumber: account.AccountNumber,
		Name:          account.Name,
		Available:     utils.CurrencyToInt(account.Available),
		Balance:       utils.CurrencyToInt(account.Balance),
	})

	if err != nil {
		return fmt.Errorf("failed to account: %w", err)
	}

	transactions, err := sbankenResource.GetArchievedTransactions(account.AccountID)
	if err != nil {
		return fmt.Errorf("failed to get arhieved transactions from sbanken: %w", err)
	}

	for _, transaction := range transactions.Items {
		err = transactionsResource.CreateIfNotExists(&types.Transaction{
			ID:             "sbanken:" + transaction.TransactionID,
			AccountID:      "sbanken:" + account.AccountID,
			AccountingDate: transaction.AccountingDate,
			InterestDate:   transaction.InterestDate,
			Amount:         int64(utils.CurrencyToInt(transaction.Amount)),
			Text:           transaction.Text,
		})

		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}
	}

	incomingTransactions, err := sbankenResource.GetIncomingTransactions(account.AccountID)
	if err != nil {
		return fmt.Errorf("failed to get incoming transactions from sbanken: %w", err)
	}

	err = incomingTransactionsResource.DeleteAll("sbanken:" + account.AccountID)
	if err != nil {
		return fmt.Errorf("failed to clear incoming transactions table: %w", err)
	}

	for _, transaction := range incomingTransactions {
		_, err := incomingTransactionsResource.Create(&types.DraftIncomingTransaction{
			AccountID:      "sbanken:" + account.AccountID,
			AccountingDate: transaction.AccountingDate,
			InterestDate:   transaction.InterestDate,
			Amount:         int64(utils.CurrencyToInt(transaction.Amount)),
			Text:           transaction.Text,
		})

		if err != nil {
			return fmt.Errorf("failed to create incoming transaction: %w", err)
		}
	}
	return nil
}
