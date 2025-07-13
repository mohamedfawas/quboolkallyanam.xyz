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

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/server"
)

func main() {
	// logger.InitLogger(constants.ServiceGateway)
	// defer logger.Sync()

	configPath := "./config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv, err := server.NewHTTPServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("Failed to start HTTP server: %v", err)
			}
		}
	}()
	log.Println("HTTP server started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server")
	if err := srv.Stop(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed")
	}
}
