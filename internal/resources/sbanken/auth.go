package sbanken

import (
	"backend/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (resource *Resource) GetAccessToken(clientID, clientSecret string) (string, error) {
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
