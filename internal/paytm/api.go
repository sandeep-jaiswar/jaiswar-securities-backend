package paytm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"go.uber.org/zap"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *PaytmMoneyClient) Login(stateKey string) (map[string]interface{}, error) {
	c.Logger.Info("Logging in to Paytm Money")

	// Construct the login URL
	loginURL := fmt.Sprintf("https://login.paytmmoney.com/merchant-login?apiKey=%s&state=%s", c.ApiKey, stateKey)

	// Create a new request
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		c.Logger.Error("Error creating login request", zap.Error(err))
		return nil, err
	}

	// Make the HTTP request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Error("Error making login request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		c.Logger.Error("Failed to login", zap.String("status", resp.Status))
		return nil, fmt.Errorf("failed to login: %s", resp.Status)
	}

	// Decode the response body
	var loginResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		c.Logger.Error("Error decoding login response", zap.Error(err))
		return nil, err
	}

	c.Logger.Info("Login successful", zap.Any("response", loginResponse))
	return loginResponse, nil
}
