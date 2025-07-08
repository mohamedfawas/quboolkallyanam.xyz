package main

// @title Qubool Kallyanam Gateway API
// @version 1.0
// @description REST API Gateway for Qubool Kallyanam microservices
// @termsOfService http://quboolkallyanam.xyz/terms/
// @contact.name Mohamed Fawas
// @contact.email adminquboolkallyanam@gmail.com
// @contact.url https://quboolkallyanam.xyz
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host quboolkallyanam.xyz
// @BasePath /api/v1

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/server"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("Logger in gateway service is initialized")
	defer logger.Sync()

	configPath := "./config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Log.Fatal("Failed to load config: ", err)
	}

	srv, err := server.NewHTTPServer(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to create server: ", err)
	}
	logger.Log.Info("Gateway server created successfully")

	go func() {
		if err := srv.Start(); err != nil {
			if err != http.ErrServerClosed {
				logger.Log.Fatal("Failed to start server: ", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")
	if err := srv.Stop(); err != nil {
		logger.Log.Error("Server forced to shutdown: ", err)
	}

	logger.Log.Info("Server stopped")
}
