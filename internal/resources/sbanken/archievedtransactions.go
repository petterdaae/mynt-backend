package sbanken

import (
	"backend/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (resource *Resource) GetArchievedTransactions(accountID string) (*ArchievedTransactions, error) {
	accessToken, err := resource.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	startDate, err := getTransactionsStartDateParameter(nil, resource.sub)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions start date parameter: %w", err)
	}

	c := context.TODO()

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

	var responseBody ArchievedTransactions
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
