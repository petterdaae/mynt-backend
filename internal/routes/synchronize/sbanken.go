package synchronize

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mynt/internal/utils"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type sbankenAccounts struct {
	AvailableItems int
	Items          []sbankenAccount
}

type sbankenAccount struct {
	AccountID     string `json:"AccountId"`
	AccountNumber string
	Name          string
	Available     float64
	Balance       float64
}

type sbankenTransactions struct {
	AvailableItems int
	Items          []sbankenTransaction
}

type sbankenTransaction struct {
	TransactionID  string `json:"TransactionId"`
	AccountingDate string
	InterestDate   string
	Amount         float64
	Text           string
}

func Sbanken(c *gin.Context) {
	sub := c.GetString("sub")
	database, _ := c.MustGet("database").(*utils.Database)

	// Connect to database
	connection, err := database.Connect()
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("database connection failed: %w", err))
		return
	}
	defer connection.Close()

	// Get client id and secret
	var clientID string
	var clientSecret string
	rows, err := connection.Query("SELECT sbanken_client_id, sbanken_client_secret FROM users WHERE id = $1", sub)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("query for sbanken credentials failed: %w", err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&clientID, &clientSecret)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("parsing query result for sbanken credentials failed: %w", err))
			return
		}
	}

	if clientID == "" || clientSecret == "" {
		utils.Unauthorized(c, fmt.Errorf("missing credentials: client id or client secret is blank"))
		return
	}

	// Query sbanken for accesstoken, accounts and transactions
	accessToken, err := getAccessToken(clientID, clientSecret)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to get sbanken access token: %w", err))
		return
	}

	accounts, err := getAccounts(accessToken)
	if err != nil {
		utils.InternalServerError(c, fmt.Errorf("failed to get sbanken accounts: %w", err))
		return
	}

	// Update accounts and transactions
	for _, account := range accounts.Items {
		_, err := connection.Exec(
			"INSERT INTO accounts (id, user_id, external_id, account_number, name, available, balance) VALUES ($1, $2, $3, $4, $5, $6, $7) "+
				"ON CONFLICT (id) DO UPDATE SET name = $5, available = $6, balance = $7",
			"sbanken:"+account.AccountID,
			sub,
			account.AccountID,
			account.AccountNumber,
			account.Name,
			utils.CurrencyToInt(account.Available),
			utils.CurrencyToInt(account.Balance),
		)

		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to insert sbanken account: %w", err))
			return
		}

		transactions, err := getTransactions(accessToken, account.AccountID)
		if err != nil {
			utils.InternalServerError(c, fmt.Errorf("failed to get sbanken transactions: %w", err))
			return
		}

		for _, transaction := range transactions.Items {
			_, err := connection.Exec(
				"INSERT INTO transactions (id, user_id, account_id, external_id, accounting_date, interest_date, amount, text) "+
					"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) "+
					"ON CONFLICT (id) DO NOTHING",
				"sbanken:"+transaction.TransactionID,
				sub,
				"sbanken:"+account.AccountID,
				transaction.TransactionID,
				transaction.AccountingDate,
				transaction.InterestDate,
				utils.CurrencyToInt(transaction.Amount),
				transaction.Text,
			)

			if err != nil {
				utils.InternalServerError(c, fmt.Errorf("failed to insert sbanken transaction: %w", err))
				return
			}
		}
	}

	c.String(http.StatusOK, "Success")
}

func getAccessToken(clientID, clientSecret string) (string, error) {
	c := context.TODO()

	// Build request
	request, err := http.NewRequestWithContext(
		c,
		"POST",
		"https://auth.sbanken.no/identityserver/connect/token",
		bytes.NewBuffer([]byte("grant_type=client_credentials")),
	)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", authHeader(clientID, clientSecret))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	// Check response
	if response.StatusCode != http.StatusOK {
		responseBodyBytes, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("unexpected status code: (%v, %v)", response.StatusCode, string(responseBodyBytes))
	}

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	responseBody := make(map[string]interface{})

	err = json.Unmarshal(responseBodyBytes, &responseBody)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	accessToken, ok := responseBody["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response body")
	}

	return accessToken, nil
}

func authHeader(clientID, clientSecret string) string {
	return "Basic " + utils.Base64Encode(
		url.QueryEscape(clientID)+":"+url.QueryEscape(clientSecret),
	) + "=="
}

func getAccounts(accessToken string) (*sbankenAccounts, error) {
	c := context.TODO()

	// Build request
	request, err := http.NewRequestWithContext(
		c,
		"GET",
		"https://publicapi.sbanken.no/apibeta/api/v1/Accounts",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/json")

	// Send request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	// Check response
	if response.StatusCode != http.StatusOK {
		responseBodyBytes, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code: (%v, %v)", response.StatusCode, string(responseBodyBytes))
	}

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var responseBody sbankenAccounts
	err = json.Unmarshal(responseBodyBytes, &responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &responseBody, nil
}

func getTransactions(accessToken, accountID string) (*sbankenTransactions, error) {
	c := context.TODO()

	// Build request
	request, err := http.NewRequestWithContext(
		c,
		"GET",
		"https://publicapi.sbanken.no/apibeta/api/v1/transactions/archive/"+accountID,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/json")

	// Send request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	// Check response
	if response.StatusCode != http.StatusOK {
		responseBodyBytes, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code: (%v, %v)", response.StatusCode, string(responseBodyBytes))
	}

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var responseBody sbankenTransactions
	err = json.Unmarshal(responseBodyBytes, &responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &responseBody, nil
}
