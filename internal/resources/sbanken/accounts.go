package sbanken

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (resource *Resource) GetAccounts() (*Accounts, error) {
	accessToken, err := resource.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	c := context.TODO()

	// Build request
	request, err := http.NewRequestWithContext(
		c,
		"GET",
		"https://publicapi.sbanken.no/apibeta/api/v1/Accounts",
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

	// Read response
	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var responseBody Accounts
	err = json.Unmarshal(responseBodyBytes, &responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &responseBody, nil
}
