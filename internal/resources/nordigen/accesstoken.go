package nordigen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetAccessTokenRequestBody struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

func (resource Resource) GetAccessToken() (string, error) {
	c := context.TODO()

	request, err := buildRequest(c)
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

func buildRequest(c context.Context) (*http.Request, error) {
	body, _ := json.Marshal(GetAccessTokenRequestBody{
		SecretID:  "",
		SecretKey: "",
	})

	request, err := http.NewRequestWithContext(
		c,
		"POST",
		"https://ob.nordigen.com/api/v2/token/new",
		bytes.NewBuffer(body),
	)

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	return request, err
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

	accessToken, ok := responseBody["access"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response body")
	}

	return accessToken, nil
}
