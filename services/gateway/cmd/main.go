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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "./config/config.yaml"
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	initLogger, err := logger.Init(cfg.Environment != constants.EnvProduction) // false for production, true for development
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer initLogger.Sync()

	rootLogger := initLogger.With(zap.String(constants.Service, constants.ServiceGateway))

	srv, err := server.NewHTTPServer(cfg, rootLogger)
	if err != nil {
		rootLogger.Fatal("Failed to create server", zap.Error(err))
	}

	go func() {
		if err := srv.Start(); err != nil {
			if err != http.ErrServerClosed {
				rootLogger.Fatal("Failed to start HTTP server", zap.Error(err))
			}
		}
	}()
	rootLogger.Info("HTTP server started successfully", zap.String("port", cfg.HTTP.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	rootLogger.Info("Shutting down server")
	if err := srv.Stop(); err != nil {
		rootLogger.Error("Error during server shutdown", zap.Error(err))
	} else {
		rootLogger.Info("Server shutdown completed")
	}
}
