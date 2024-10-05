package main

import (
	"log"

	"github.com/sandeep-jaiswar/jaiswar-securities/internal/config"
	"github.com/sandeep-jaiswar/jaiswar-securities/internal/paytm"
	"github.com/sandeep-jaiswar/jaiswar-securities/internal/server"
	"github.com/sandeep-jaiswar/jaiswar-securities/pkg"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.LoadConfig()
	logger, err := pkg.InitializeLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting jaiswar-securities...")

	if err := run(); err != nil {
		logger.Fatal("Application encountered an error", zap.Error(err))
	}

	paytmClient := paytm.NewPaytmMoneyClient(appConfig.PaytmApiKey, appConfig.PaytmSecretKey, logger)

	srv := server.NewServer(logger, appConfig.Port, paytmClient)
	srv.Start()
}

func run() error {
	zap.L().Info("Application is running")
	return nil
}