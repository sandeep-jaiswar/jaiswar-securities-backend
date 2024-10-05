package paytm

import (
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	clientInstance *PaytmMoneyClient
	once           sync.Once
)

type PaytmMoneyClient struct {
	BaseURL    string
	HTTPClient *http.Client
	ApiKey     string
	SecretKey  string
	Logger     *zap.Logger
}

func NewPaytmMoneyClient(apiKey, secretKey string, logger *zap.Logger) *PaytmMoneyClient {
	once.Do(func() {
		clientInstance = &PaytmMoneyClient{
			BaseURL: os.Getenv("PAYTM_LOGIN_BASE_URL"),
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			ApiKey:    os.Getenv("PAYTM_API_KEY"),
			SecretKey: os.Getenv("PAYTM_SECRET_KEY"),
			Logger:    logger,
		}
	})

	return clientInstance
}
