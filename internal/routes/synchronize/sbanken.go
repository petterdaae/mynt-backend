package synchronize

import (
	"bytes"
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
	AccountId     string
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
	TransactionId  string
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
		c.AbortWithError(500, fmt.Errorf("database connection failed: %w", err))
		return
	}
	defer connection.Close()

	// Get client id and secret
	var clientId string
	var clientSecret string
	rows, err := connection.Query("SELECT sbanken_client_id, sbanken_client_secret FROM users WHERE id = $1", sub)
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("query for sbanken credentials failed"))
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&clientId, &clientSecret)
		if err != nil {
			c.AbortWithError(500, fmt.Errorf("parsing query result for sbanken credentials failed: %w", err))
			return
		}
	}

	if clientId == "" || clientSecret == "" {
		c.AbortWithError(400, fmt.Errorf("missing credentials: client id or client secret is blank"))
		return
	}

	// Query sbanken for accesstoken, accounts and transactions
	accessToken, err := getAccessToken(clientId, clientSecret)
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("failed to get sbanken access token: %w", err))
		return
	}

	accounts, err := getAccounts(accessToken)
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("failed to get sbanken accounts: %w", err))
		return
	}

	transactions, err := getTransactions(accessToken, accounts.Items[0].AccountId)
	if err != nil {
		c.AbortWithError(500, fmt.Errorf("failed to get sbanken transactions: %w", err))
		return
	}

	// Update accounts and transactions
	c.JSON(200, transactions)
}

func getAccessToken(clientId, clientSecret string) (string, error) {
	// Build request
	request, err := http.NewRequest(
		"POST",
		"https://auth.sbanken.no/identityserver/connect/token",
		bytes.NewBuffer([]byte("grant_type=client_credentials")),
	)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", authHeader(clientId, clientSecret))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	// Check response
	if response.StatusCode != 200 {
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

func authHeader(clientId string, clientSecret string) string {
	return "Basic " + utils.Base64Encode(
		url.QueryEscape(clientId)+":"+url.QueryEscape(clientSecret),
	) + "=="
}

func getAccounts(accessToken string) (*sbankenAccounts, error) {
	// Build request
	request, err := http.NewRequest(
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
	if response.StatusCode != 200 {
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

func getTransactions(accessToken string, accountId string) (*sbankenTransactions, error) {
	// Build request
	request, err := http.NewRequest(
		"GET",
		"https://publicapi.sbanken.no/apibeta/api/v1/transactions/archive/"+accountId,
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
	if response.StatusCode != 200 {
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
