package main

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

	logger.Log.Info("Loading config from: ", configPath)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Log.Fatal("Failed to load config: ", err)
	}
	logger.Log.Info("Config loaded successfully")

	srv, err := server.NewHTTPServer(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to create server: ", err)
	}

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
