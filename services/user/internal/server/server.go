package server

import (
	"context"
	"fmt"
	"net"

	postgresAdapters "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/adapters/postgres"
    eventHandlers "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/handlers/event"
    userProfileUsecaseImpl "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase/user_profile"
    messageBrokerAdapter "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/adapters/messageBroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/config"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	pgClient        *postgres.Client
	messagingClient messageBroker.Client
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
	logger.Log.Info("✅ User Service Connected to PostgreSQL ")

	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		ctx := context.Background()
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		logger.Log.Info("✅ User Service Connected to PubSub")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		logger.Log.Info("✅ User Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer()

	// Initialize repositories
	userProfileRepo := postgresAdapters.NewUserProfileRepository(pgClient)

	// Initialize event publisher
eventPublisher := messageBrokerAdapter.NewEventPublisher(messagingClient)

	// Initialize use cases
userProfileUC := userProfileUsecaseImpl.NewUserProfileUsecase(userProfileRepo, eventPublisher)

// Initialize event handler
authEventHandler := eventHandlers.NewAuthEventHandler(messagingClient, userProfileUC)


// Start event listener
go func() {
    if err := authEventHandler.StartListening(context.Background()); err != nil {
        logger.Log.Error("Failed to start auth event handler", "error", err)
    }
}()
	

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
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
}
