package paytm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (c *PaytmMoneyClient) Login(stateKey string) (string, error) {
	c.Logger.Info("Logging in to Paytm Money")

	// Construct the login URL
	loginURL := fmt.Sprintf("https://login.paytmmoney.com/merchant-login?apiKey=%s&state=%s", c.ApiKey, stateKey)

	// Redirect to the login URL
	return loginURL, nil
}

func (c *PaytmMoneyClient) GenerateAccessToken(requestToken string) (*http.Response, error) {
	c.Logger.Info("GenerateAccessToken in to Paytm Money")

	// Construct the access token URL
	accessTokenURL := "https://developer.paytmmoney.com/accounts/v2/gettoken"

	// Define the request body structure
	type requestBodyStruct struct {
		ApiKey       string `json:"api_key"`
		ApiSecretKey string `json:"api_secret_key"`
		RequestToken string `json:"request_token"`
	}

	// Prepare the request body
	requestBodyJSON := requestBodyStruct{
		ApiKey:       c.ApiKey,
		ApiSecretKey: c.SecretKey,
		RequestToken: requestToken,
	}

	// Marshal the JSON object to a JSON string
	requestBodyJSONStr, err := json.Marshal(requestBodyJSON)
	if err != nil {
		c.Logger.Error("Failed to marshal request body to JSON", zap.Error(err))
		return nil, err
	}

	// Create and execute the HTTP request
	req, err := http.NewRequest("POST", accessTokenURL, bytes.NewBuffer(requestBodyJSONStr))
	if err != nil {
		c.Logger.Error("Failed to create request for access token", zap.Error(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Error("Failed to generate access token", zap.Error(err))
		return nil, err
	}

	return response, nil
}