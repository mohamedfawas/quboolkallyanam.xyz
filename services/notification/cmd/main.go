package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/server"

	"go.uber.org/zap"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "./config/config.yaml"
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("[NOTIFICATION] Failed to load config: %v", err)
	}
	initialLogger, err := logger.Init(cfg.Environment != constants.EnvProduction) // false for production, true for development
	if err != nil {
		log.Fatalf("[NOTIFICATION] Failed to init logger: %v", err)
	}
	defer initialLogger.Sync()
	rootLogger := initialLogger.With(zap.String(constants.Service, constants.ServiceNotification))

	// Create a context that will be cancelled on SIGINT or SIGTERM
	// Example: if the process receives Ctrl+C, ctx.Done() will unblock
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv, err := server.NewServer(ctx, cfg, rootLogger)
	if err != nil {
		rootLogger.Fatal("Failed to create server", zap.Error(err))
	}

	go func() {
		if err := srv.Start(); err != nil {
			rootLogger.Error("Failed to start notification service", zap.Error(err))
		}
	}()
	rootLogger.Info("notification service started")

	<-ctx.Done() // pausing the main goroutine until a shutdown signal is received

	rootLogger.Info("Shutting down server...")
	srv.Stop()
	rootLogger.Info("Server stopped")
}
