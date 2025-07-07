package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/config"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	httpServer      *http.Server
	pgClient        *postgres.Client
	messagingClient messaging.Client
}

func NewServer(config *config.Config) (*Server, error) {
	pgClient, err := postgres.NewClient(postgres.Config{
		Host:     config.Postgres.Host,
		Port:     config.Postgres.Port,
		User:     config.Postgres.User,
		Password: config.Postgres.Password,
		DBName:   config.Postgres.DBName,
		SSLMode:  config.Postgres.SSLMode,
		TimeZone: config.Postgres.TimeZone,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}
	logger.Log.Info("✅ Payment Service Connected to PostgreSQL ")

	var messagingClient messaging.Client
	if config.Environment == "production" {
		ctx := context.Background()
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		logger.Log.Info("✅ Payment Service Connected to PubSub")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		logger.Log.Info("✅ Payment Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer()

	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.HTTP.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(config.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.HTTP.IdleTimeout) * time.Second,
	}

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
		messagingClient: messagingClient,
		httpServer:      httpServer,
	}

	return server, nil
}

func (s *Server) Start() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("failed to start http server", "error", err)
		}
	}()
	logger.Log.Info("✅ Payment Service's HTTP Server started")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.Port))
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Log.Error("HTTP server forced to shutdown", "error", err)
	}

	if s.pgClient != nil {
		if err := s.pgClient.Close(); err != nil {
			logger.Log.Error("failed to close postgres client: %w", err)
		}
	}

	if s.messagingClient != nil {
		if err := s.messagingClient.Close(); err != nil {
			logger.Log.Error("failed to close messaging client", "error", err)
		}
	}
	logger.Log.Info("✅ Payment Service stopped")
}
