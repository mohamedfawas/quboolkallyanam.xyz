package server

import (
	"context"
	"fmt"
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
	chatUsecaseImpl "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase/chat"
	userProjectionUsecaseImpl "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase/user_projection"
	eventHandlers "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/handlers/event"
	v1 "github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/handlers/grpc/v1"
	interceptors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/grpc/interceptors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	pgClient        *postgres.Client
	mongoClient     *mongodb.Client
	messagingClient messageBroker.Client
	logger          *zap.Logger
}

func NewServer(ctx context.Context, config *config.Config, rootLogger *zap.Logger) (*Server, error) {

	///////////////////////// POSTGRES INITIALIZATION /////////////////////////
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
	rootLogger.Info("Connected to PostgreSQL ",
		zap.String("host", config.Postgres.Host),
		zap.Int("port", config.Postgres.Port),
		zap.String("ssl_mode", config.Postgres.SSLMode),
		zap.String("time_zone", config.Postgres.TimeZone),
	)

	///////////////////////// MONGODB INITIALIZATION /////////////////////////
	mongoClient, err := mongodb.NewClient(ctx, mongodb.Config{
		URI:      config.MongoDB.URI,
		Database: config.MongoDB.Database,
		Timeout:  config.MongoDB.Timeout,
	})
	if err != nil {
		pgClient.Close()
		return nil, fmt.Errorf("mongodb connect: %w", err)
	}
	rootLogger.Info("Connected to MongoDB ")

	///////////////////////// MESSAGING CLIENT INITIALIZATION /////////////////////////
	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			mongoClient.Close()
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		rootLogger.Info("Connected to PubSub ")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			// Clean up existing connections before returning error
			mongoClient.Close()
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		rootLogger.Info("Connected to RabbitMQ ")
	}

	///////////////////////// GRPC SERVER INITIALIZATION /////////////////////////
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryErrorInterceptor()),
	)
	rootLogger.Info("gRPC server created")

	///////////////////////// REPOSITORIES INITIALIZATION /////////////////////////
	userProjectionRepo := postgresAdapters.NewUserProjectionRepository(pgClient)
	rootLogger.Info("User Projection Repository Initialized")
	conversationRepo := mongodbAdapters.NewConversationRepository(mongoClient)
	rootLogger.Info("Conversation Repository Initialized")
	messageRepo := mongodbAdapters.NewMessageRepository(mongoClient)
	rootLogger.Info("Message Repository Initialized")

	///////////////////////// USE CASES INITIALIZATION /////////////////////////
	userProjectionUC := userProjectionUsecaseImpl.NewUserProjectionUsecase(userProjectionRepo)
	rootLogger.Info("User Projection Use Case Initialized")
	chatUC := chatUsecaseImpl.NewChatUsecase(conversationRepo, messageRepo, userProjectionRepo)
	rootLogger.Info("Chat Use Case Initialized")

	///////////////////////// EVENT HANDLER INITIALIZATION /////////////////////////
	userEventHandler := eventHandlers.NewUserEventListener(messagingClient, userProjectionUC, rootLogger)
	rootLogger.Info("User Event Handler Initialized")

	///////////////////////// GRPC HANDLER INITIALIZATION /////////////////////////
	chatHandler := v1.NewChatHandler(chatUC, rootLogger)
	chatpbv1.RegisterChatServiceServer(grpcServer, chatHandler)
	rootLogger.Info("gRPC Handler Initialized")

	///////////////////////// EVENT LISTENER INITIALIZATION /////////////////////////
	go func() {
		if err := userEventHandler.StartListening(ctx); err != nil {
			rootLogger.Error("Failed to start user event listener", zap.Error(err))
		}
	}()

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
		mongoClient:     mongoClient,
		messagingClient: messagingClient,
		logger:          rootLogger,
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
			s.logger.Error("failed to close messaging client", zap.Error(err))
		}
	}

	if s.mongoClient != nil {
		if err := s.mongoClient.Close(); err != nil {
			s.logger.Error("failed to close mongo client", zap.Error(err))
		}
	}

	if s.pgClient != nil {
		if err := s.pgClient.Close(); err != nil {
			s.logger.Error("failed to close postgres client", zap.Error(err))
		}
	}
}
