package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/server"
	"google.golang.org/grpc"
)

func main() {
	// logger.InitLogger(constants.ServiceAuth)
	// defer logger.Sync()

	configPath := "./config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Catch SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	srv, err := server.NewServer(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			if err != grpc.ErrServerStopped {
				log.Fatalf("Failed to start server: %v", err)
			}
		}
	}()
	log.Println("Auth Service Server started")

	<-quit
	cancel()

	log.Println("Shutting down server...")
	srv.Stop()
	log.Println("Server stopped")
}
