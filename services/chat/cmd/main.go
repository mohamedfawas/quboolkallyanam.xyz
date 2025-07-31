package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "./config/config.yaml"
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("[CHAT] Failed to load config: %v", err)
	}

	initialLogger, err := logger.Init(cfg.Environment != constants.EnvProduction) // false for production, true for development
	if err != nil {
		log.Fatalf("[CHAT] Failed to init logger: %v", err)
	}
	defer initialLogger.Sync()

	rootLogger := initialLogger.With(zap.String(constants.Service, constants.ServiceChat))

	srv, err := server.NewServer(cfg, rootLogger)
	if err != nil {
		rootLogger.Fatal("Failed to create server", zap.Error(err))
	}

	go func() {
		if err := srv.Start(); err != nil {
			if err != grpc.ErrServerStopped {
				rootLogger.Fatal("Failed to start server", zap.Error(err))
			}
		}
	}()
	rootLogger.Info("Chat Service Server started on ", zap.Int(constants.Port, cfg.GRPC.Port))

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	rootLogger.Info("Shutting down server...")
	srv.Stop()
	rootLogger.Info("Server stopped")
}
