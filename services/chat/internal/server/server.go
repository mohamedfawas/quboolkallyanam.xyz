package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/nosql"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/nosql/firestore"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/nosql/mongodb"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	eventHandlers "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/event"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/config"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	pgClient        *postgres.Client
	nosqlClient     nosql.Client
	messagingClient messageBroker.Client
}

func NewServer(config *config.Config) (*Server, error) {
	ctx := context.Background()

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
	log.Println("Chat Service Connected to PostgreSQL ")

	// Initialize NoSQL client based on environment
	var nosqlClient nosql.Client

	if config.Environment == constants.EnvProduction {
		nosqlClient, err = firestore.NewClient(ctx, config.Firestore.ProjectID)
		if err != nil {
			log.Println("failed to create firestore client: %w", err)
			return nil, fmt.Errorf("failed to create firestore client: %w", err)
		}
		log.Println("Chat Service Connected to Firestore")
	} else {
		// Use MongoDB for development
		nosqlClient, err = mongodb.NewClient(ctx, mongodb.Config{
			URI:      config.MongoDB.URI,
			Database: config.MongoDB.Database,
		})
		if err != nil {
			log.Println("failed to create mongodb client: %w", err)
			return nil, fmt.Errorf("failed to create mongodb client: %w", err)
		}
		log.Println("Chat Service Connected to MongoDB")
	}

	// Initialize messaging client based on environment
	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			nosqlClient.Close()
			log.Println("failed to create pubsub client: %w", err)
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		log.Println("Chat Service Connected to PubSub")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			// Clean up existing connections before returning error
			nosqlClient.Close()
			log.Println("failed to create rabbitmq client: %w", err)
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		log.Println("Chat Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer()

	// Initialize repositories
userProjectionRepo := postgresAdapters.NewUserProjectionRepository(pgClient)

// Initialize use cases  
userProjectionUC := userProjectionUsecaseImpl.NewUserProjectionUsecase(userProjectionRepo)

// Initialize event handler
userEventHandler := eventHandlers.NewUserEventHandler(messagingClient, userProjectionUC)

go func() {
    if err := userEventHandler.StartListening(context.Background()); err != nil {
        log.Printf("Failed to start user event handler: %v", err)
    }
}()
	
	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
		nosqlClient:     nosqlClient,
		messagingClient: messagingClient,
	}

	return server, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.Port))
	if err != nil {
		log.Println("Chat Service failed to listen: %w", err)
		return err
	}

	log.Println("Chat Service starting on port ", s.config.GRPC.Port)
	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	log.Println("Chat Service shutting down...")
	s.grpcServer.GracefulStop()

	// Close messaging client first
	if s.messagingClient != nil {
		if err := s.messagingClient.Close(); err != nil {
			log.Println("Chat Service failed to close messaging client", "error", err)
		}
	}

	// Close NoSQL client
	if s.nosqlClient != nil {
		if err := s.nosqlClient.Close(); err != nil {
			log.Println("Chat Service failed to close nosql client", "error", err)
		}
	}

	if s.pgClient != nil {
		if err := s.pgClient.Close(); err != nil {
			log.Println("failed to close postgres client: %w", err)
		}
	}

	log.Println("Chat Service stopped gracefully")
}
