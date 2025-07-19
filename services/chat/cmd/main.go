package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/server"
	"google.golang.org/grpc"
)

func main() {

	configPath := "./config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Println("Failed to load config: ", err)
	}
	log.Println("Chat Service Config loaded")

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Println("Failed to create server: ", err)
	}
	log.Println("Chat Service Server created")

	go func() {
		if err := srv.Start(); err != nil {
			if err != grpc.ErrServerStopped {
				log.Println("Failed to start server: ", err)
			}
		}
	}()
	log.Println("Chat Service Server started")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	srv.Stop()
	log.Println("Server stopped")
}
