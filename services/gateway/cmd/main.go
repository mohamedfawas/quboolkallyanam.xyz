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
	"context"
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
		log.Fatalf("[GATEWAY] Failed to load config: %v", err)
	}
	initialLogger, err := logger.Init(cfg.Environment != constants.EnvProduction) // false for production, true for development
	if err != nil {
		log.Fatalf("[GATEWAY] Failed to init logger: %v", err)
	}
	defer initialLogger.Sync()
	rootLogger := initialLogger.With(zap.String(constants.Service, constants.ServiceGateway))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv, err := server.NewHTTPServer(ctx, cfg, rootLogger)
	if err != nil {
		rootLogger.Fatal("Failed to create server", zap.Error(err))
	}

	go func() {
		if err := srv.Start(); err != nil {
			if err != http.ErrServerClosed {
				rootLogger.Error("Failed to start HTTP server", zap.Error(err))
				stop() // stops listening for signals, and exits the program
			}
		}
	}()
	rootLogger.Info("HTTP server started", zap.String("port", cfg.HTTP.Port))

	<-ctx.Done() // pausing the main goroutine until a shutdown signal is received

	rootLogger.Info("Shutting down server...")
	srv.Stop()
	rootLogger.Info("Server stopped")
}
