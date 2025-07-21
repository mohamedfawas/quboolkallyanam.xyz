package server

import (
	"context"
	"fmt"
	"log"
	"net"

	chatpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/chat/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/mongodb"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/rabbitmq"
	mongodbAdapters "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/adapters/mongodb"
	postgresAdapters "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/adapters/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/config"
	eventHandlers "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/event"
	chatUsecaseImpl "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase/chat"
	userProjectionUsecaseImpl "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase/user_projection"
	v1 "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/handlers/grpc/v1"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	pgClient        *postgres.Client
	mongoClient     *mongodb.Client
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

	mongoClient, err := mongodb.NewClient(ctx, mongodb.Config{
		URI:      config.MongoDB.URI,
		Database: config.MongoDB.Database,
		Timeout:  config.MongoDB.Timeout,
	})
	if err != nil {
		pgClient.Close()
		return nil, fmt.Errorf("mongodb connect: %w", err)
	}
	log.Printf("Connected to MongoDB (%s)\n", config.Environment)

	// Initialize messaging client based on environment
	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			mongoClient.Close()
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
			mongoClient.Close()
			log.Println("failed to create rabbitmq client: %w", err)
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		log.Println("Chat Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer()

	// Initialize repositories
	userProjectionRepo := postgresAdapters.NewUserProjectionRepository(pgClient)
	conversationRepo := mongodbAdapters.NewConversationRepository(mongoClient)
	//messageRepo := mongodbAdapters.NewMessageRepository(mongoClient)

	// Initialize use cases
	userProjectionUC := userProjectionUsecaseImpl.NewUserProjectionUsecase(userProjectionRepo)
	chatUC := chatUsecaseImpl.NewChatUsecase(conversationRepo, userProjectionRepo)

	// Initialize event handler
	userEventHandler := eventHandlers.NewUserEventHandler(messagingClient, userProjectionUC)

	// Initialize handlers
	chatHandler := v1.NewChatHandler(chatUC)

	chatpbv1.RegisterChatServiceServer(grpcServer, chatHandler)
	log.Println("Chat Service gRPC handlers registered")

	go func() {
		if err := userEventHandler.StartListening(context.Background()); err != nil {
			log.Printf("Failed to start user event handler: %v", err)
		}
	}()

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
		mongoClient:     mongoClient,
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
	if s.mongoClient != nil {
		if err := s.mongoClient.Close(); err != nil {
			log.Println("Chat Service failed to close mongo client", "error", err)
		}
	}

	if s.pgClient != nil {
		if err := s.pgClient.Close(); err != nil {
			log.Println("failed to close postgres client: %w", err)
		}
	}

	log.Println("Chat Service stopped gracefully")
}
