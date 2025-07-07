package server

import (
	"context"
	"fmt"
	"net"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/redis"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	pgClient        *postgres.Client
	redisClient     *redis.Client
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
	logger.Log.Info("✅ Auth Service Connected to PostgreSQL ")

	redisClient, err := redis.NewClient(redis.Config{
		Host:     config.Redis.Host,
		Port:     config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if err != nil {
		// Clean up existing connections before returning error
		pgClient.Close()
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}
	logger.Log.Info("✅ Auth Service Connected to Redis")

	var messagingClient messaging.Client
	if config.Environment == "production" {
		ctx := context.Background()
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			redisClient.Close()
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		logger.Log.Info("✅ Auth Service Connected to PubSub")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			redisClient.Close()
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		logger.Log.Info("✅ Auth Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer()

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
		redisClient:     redisClient,
		messagingClient: messagingClient,
	}

	return server, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.Port))
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()

	// Close messaging client first
	if s.messagingClient != nil {
		if err := s.messagingClient.Close(); err != nil {
			logger.Log.Error("failed to close messaging client", "error", err)
		}
	}

	if s.pgClient != nil {
		if err := s.pgClient.Close(); err != nil {
			logger.Log.Error("failed to close postgres client: %w", err)
		}
	}

	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			logger.Log.Error("failed to close redis client: %w", err)
		}
	}
}
