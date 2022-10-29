package nordigen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetRequisitionLinkRequestBody struct {
	Redirect      string `json:"redirect"`
	InstitutionID string `json:"institution_id"`
	Reference     string `json:"reference"`
	Agreement     string `json:"agreement"`
	UserLanguage  string `json:"user_language"`
}

func (resource Resource) GetRequisitionLink() (string, error) {
	c := context.TODO()

	request, err := buildRequisitionRequest(c)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	response, err := sendRequest(request)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	return readRequisitionResponse(response)
}

func buildRequisitionRequest(c context.Context) (*http.Request, error) {
	body, _ := json.Marshal(GetRequisitionLinkRequestBody{
		Redirect:      "https://localhost:8080/nordigen/callback",
		InstitutionID: "SBANKEN_SBAKNOBB",
		// Reference: "", can use this for internal referencing later
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

func readRequisitionResponse(response *http.Response) (string, error) {
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
