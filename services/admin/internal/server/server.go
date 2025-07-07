package server

import (
	"context"
	"fmt"
	"net"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messaging/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/admin/internal/config"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	messagingClient messaging.Client
}

func NewServer(config *config.Config) (*Server, error) {
	var messagingClient messaging.Client
	var err error

	if config.Environment == "production" {
		ctx := context.Background()
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		logger.Log.Info("✅ Admin Service Connected to PubSub")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		logger.Log.Info("✅ Admin Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer()

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		messagingClient: messagingClient,
	}

	return server, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.config.GRPC.Port, err)
	}

	logger.Log.Info("✅ Admin Service's GRPC Server started", "port", s.config.GRPC.Port)
	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()

	if s.messagingClient != nil {
		if err := s.messagingClient.Close(); err != nil {
			logger.Log.Error("failed to close messaging client", "error", err)
		}
	}
	logger.Log.Info("✅ Admin Service stopped")
}
