package paytm

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

type PaytmMoneyClient struct {
	BaseURL    string
	HTTPClient *http.Client
	ApiKey     string
	SecretKey  string
	Logger     *zap.Logger
}

// NewPaytmMoneyClient initializes a new PaytmMoneyClient instance.
func NewPaytmMoneyClient(apiKey, secretKey string, logger *zap.Logger) *PaytmMoneyClient {
	return &PaytmMoneyClient{
		BaseURL: os.Getenv("PaytmLoginBaseUrl"),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ApiKey:    apiKey,
		SecretKey: secretKey,
		Logger:    logger,
	}
}
