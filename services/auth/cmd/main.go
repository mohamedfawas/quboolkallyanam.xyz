package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/server"
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
		log.Fatalf("[AUTH] Failed to load config: %v", err)
	}
	initialLogger, err := logger.Init(cfg.Environment != constants.EnvProduction) // false for production, true for development
	if err != nil {
		log.Fatalf("[AUTH] Failed to init logger: %v", err)
	}
	defer initialLogger.Sync()
	rootLogger := initialLogger.With(zap.String(constants.Service, constants.ServiceAuth))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv, err := server.NewServer(ctx, cfg, rootLogger)
	if err != nil {
		rootLogger.Fatal("Failed to create server", zap.Error(err))
	}

	go func() {
		if err := srv.Start(); err != nil {
			if err != grpc.ErrServerStopped {
				rootLogger.Error("Server exited with error", zap.Error(err))
				stop() // stops listening for signals, and exits the program
			}
		}
	}()
	rootLogger.Info("gRPC server started", zap.Int("port", cfg.GRPC.Port))

	<-ctx.Done() // pausing the main goroutine until a shutdown signal is received

	rootLogger.Info("Shutting down server...")
	srv.Stop()
	rootLogger.Info("Server stopped")
}
