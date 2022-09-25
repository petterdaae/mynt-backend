package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (resource *Resource) GetIncomingTransactions(accountID string) ([]IncomingTransaction, error) {
	accessToken, err := resource.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	startDate := time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02")

	c := context.TODO()

	// Build request
	request, err := http.NewRequestWithContext(
		c,
		"GET",
		"https://publicapi.sbanken.no/apibeta/api/v1/transactions/"+accountID+"?startDate="+startDate+"&length=1000",
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

	var lastWeeksTransansactions IncomingTransactions
	err = json.Unmarshal(responseBodyBytes, &lastWeeksTransansactions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	var reservedTransactions []IncomingTransaction

	for _, transaction := range lastWeeksTransansactions.Items {
		if transaction.IsReservation {
			reservedTransactions = append(reservedTransactions, transaction)
		}
	}

	return reservedTransactions, nil
}
