package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/server"
	"google.golang.org/grpc"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("✅ Logger in chat service is initialized")
	defer logger.Sync()

	configPath := "./config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Log.Fatal("Failed to load config: ", err)
	}
	logger.Log.Info("✅ Chat Service Config loaded")

	srv, err := server.NewServer(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to create server: ", err)
	}
	logger.Log.Info("✅ Chat Service Server created")

	go func() {
		if err := srv.Start(); err != nil {
			if err != grpc.ErrServerStopped {
				logger.Log.Fatal("Failed to start server: ", err)
			}
		}
	}()
	logger.Log.Info("✅ Chat Service Server started")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("🛑 Shutting down server...")
	srv.Stop()
	logger.Log.Info("✅ Server stopped")
}
