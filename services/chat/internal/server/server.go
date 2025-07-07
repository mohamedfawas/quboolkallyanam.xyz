package server

import (
	"context"
	"fmt"
	"net"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/nosql"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/nosql/firestore"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/nosql/mongodb"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/config"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	nosqlClient     nosql.Client
	messagingClient messaging.Client
}

func NewServer(config *config.Config) (*Server, error) {
	ctx := context.Background()

	// Initialize NoSQL client based on environment
	var nosqlClient nosql.Client
	var err error

	if config.Environment == "production" {
		nosqlClient, err = firestore.NewClient(ctx, config.Firestore.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("failed to create firestore client: %w", err)
		}
		logger.Log.Info("✅ Chat Service Connected to Firestore")
	} else {
		// Use MongoDB for development
		nosqlClient, err = mongodb.NewClient(ctx, mongodb.Config{
			URI:      config.MongoDB.URI,
			Database: config.MongoDB.Database,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create mongodb client: %w", err)
		}
		logger.Log.Info("✅ Chat Service Connected to MongoDB")
	}

	// Initialize messaging client based on environment
	var messagingClient messaging.Client
	if config.Environment == "production" {
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			nosqlClient.Close()
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		logger.Log.Info("✅ Chat Service Connected to PubSub")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			// Clean up existing connections before returning error
			nosqlClient.Close()
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		logger.Log.Info("✅ Chat Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer()

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		nosqlClient:     nosqlClient,
		messagingClient: messagingClient,
	}

	return server, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.Port))
	if err != nil {
		return err
	}

	logger.Log.Info("🚀 Chat Service starting on port ", s.config.GRPC.Port)
	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	logger.Log.Info("🛑 Chat Service shutting down...")
	s.grpcServer.GracefulStop()

	// Close messaging client first
	if s.messagingClient != nil {
		if err := s.messagingClient.Close(); err != nil {
			logger.Log.Error("failed to close messaging client", "error", err)
		}
	}

	// Close NoSQL client
	if s.nosqlClient != nil {
		if err := s.nosqlClient.Close(); err != nil {
			logger.Log.Error("failed to close nosql client", "error", err)
		}
	}

	logger.Log.Info("✅ Chat Service stopped gracefully")
}
