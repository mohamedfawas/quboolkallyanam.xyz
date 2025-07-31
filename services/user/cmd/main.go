package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/server"
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
		log.Fatalf("[USER] Failed to load config: %v", err)
	}

	rootLogger, err := logger.Init(cfg.Environment != constants.EnvProduction) // false for production, true for development
	if err != nil {
		log.Fatalf("[USER] Failed to init logger: %v", err)
	}
	defer rootLogger.Sync()

	srv, err := server.NewServer(cfg, rootLogger)
	if err != nil {
		rootLogger.Fatal("[USER] Failed to create server", zap.Error(err))
	}

	go func() {
		if err := srv.Start(); err != nil {
			if err != grpc.ErrServerStopped {
				rootLogger.Fatal("[USER] Failed to start server", zap.Error(err))
			}
		}
	}()
	rootLogger.Info("[USER] User Service Server started", zap.String("port", strconv.Itoa(cfg.GRPC.Port)))

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	rootLogger.Info("[USER] Shutting down server...")
	srv.Stop()
	rootLogger.Info("[USER] Server stopped")
}
