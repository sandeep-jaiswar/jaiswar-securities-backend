package main

import (
	"log"

	"github.com/sandeep-jaiswar/jaiswar-securities/pkg"
	"go.uber.org/zap"
)

func main() {
	logger, err := pkg.InitializeLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting jaiswar-securities...")

	if err := run(); err != nil {
		logger.Fatal("Application encountered an error", zap.Error(err))
	}
}

func run() error {
	zap.L().Info("Application is running")
	return nil
}