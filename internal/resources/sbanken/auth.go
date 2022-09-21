package sbanken

import (
	"backend/internal/resources/user"
	"backend/internal/types"
	"backend/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (resource *Resource) GetAccessToken() (string, error) {
	userResource := user.Configure(resource.sub, resource.database)
	userInfo, err := userResource.Read()
	if err != nil {
		return "", fmt.Errorf("failed to read user: %w", err)
	}

	c := context.TODO()

	request, err := buildRequest(c, userInfo)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	response, err := sendRequest(request)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	return readResponse(response)
}

func buildRequest(c context.Context, userInfo types.User) (*http.Request, error) {
	request, err := http.NewRequestWithContext(
		c,
		"POST",
		"https://auth.sbanken.no/identityserver/connect/token",
		bytes.NewBuffer([]byte("grant_type=client_credentials")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", authHeader(userInfo.SbankenClientID, userInfo.SbankenClientSecret))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

func authHeader(clientID, clientSecret string) string {
	return "Basic " + utils.Base64Encode(
		url.QueryEscape(clientID)+":"+url.QueryEscape(clientSecret),
	) + "=="
}

func sendRequest(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		responseBodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code: (%v, %v)", response.StatusCode, string(responseBodyBytes))
	}

	return response, nil
}

func readResponse(response *http.Response) (string, error) {
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
