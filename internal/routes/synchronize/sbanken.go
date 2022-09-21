package synchronize

import (
	"backend/internal/resources/sbanken"
	"backend/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

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

	sbankenResource := sbanken.Configure(sub, database)

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

	accounts, err := sbankenResource.GetAccounts()
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

		transactions, err := getTransactions(accessToken, account.AccountID, sub, database)
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
		responseBodyBytes, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("unexpected status code: (%v, %v)", response.StatusCode, string(responseBodyBytes))
	}

	responseBodyBytes, err := io.ReadAll(response.Body)
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

func getTransactions(accessToken, accountID, sub string, database *utils.Database) (*sbankenTransactions, error) {
	c := context.TODO()

	startDate, err := getTransactionsStartDateParameter(database, sub)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions start date parameter: %w", err)
	}

	// Build request
	request, err := http.NewRequestWithContext(
		c,
		"GET",
		"https://publicapi.sbanken.no/apibeta/api/v1/transactions/archive/"+accountID+"?startDate="+startDate+"&length=1000",
		http.NoBody,
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
		responseBodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code: (%v, %v)", response.StatusCode, string(responseBodyBytes))
	}

	responseBodyBytes, err := io.ReadAll(response.Body)
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

func getTransactionsStartDateParameter(database *utils.Database, sub string) (string, error) {
	row, err := database.QueryRow(
		"SELECT accounting_date FROM transactions WHERE user_id = $1 ORDER BY accounting_date DESC LIMIT 1",
		sub,
	)
	if err != nil {
		return "", fmt.Errorf("failed to query names: %w", err)
	}

	var mostRecentAccountingDate *string
	err = row.Scan(&mostRecentAccountingDate)
	if err != nil {
		// row.Scan fails if there are no rows in transactions
		return time.Now().Add(-300 * 24 * time.Hour).Format("2006-01-02"), nil
	}

	if mostRecentAccountingDate == nil {
		return "", fmt.Errorf("mostRecentAccountingDate is nil")
	}

	parsedTime, err := time.Parse("2006-01-02T15:04:05", *mostRecentAccountingDate)
	if err != nil {
		return "", fmt.Errorf("failed to parse mostRecentAccountingDate string: %w", err)
	}

	parsedTime = parsedTime.Add(-time.Hour * 24 * 7)

	return parsedTime.Format("2006-01-02"), nil
}
