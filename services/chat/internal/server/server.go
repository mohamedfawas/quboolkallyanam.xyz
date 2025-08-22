package server

import (
	"context"
	"fmt"
	"net"

	chatpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/chat/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/mongodb"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	interceptors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/grpc/interceptors"
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
	"go.uber.org/zap"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	healthSrv       *health.Server
	pgClient        *postgres.Client
	mongoClient     *mongodb.Client
	messagingClient messageBroker.Client
	logger          *zap.Logger
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewServer(ctx context.Context, config *config.Config, rootLogger *zap.Logger) (*Server, error) {
	// Create child context with cancellation
	serverCtx, cancel := context.WithCancel(ctx)
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
	mongoClient, err := mongodb.NewClient(serverCtx, mongodb.Config{
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
		messagingClient, err = pubsub.NewClient(serverCtx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
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
			pgClient.Close()
			mongoClient.Close()
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		rootLogger.Info("Connected to RabbitMQ ")
	}

	///////////////////////// GRPC SERVER INITIALIZATION /////////////////////////
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryErrorInterceptor()),
	)

	///////////////////////// HEALTH SERVER INITIALIZATION /////////////////////////
	// create and register simple standard gRPC health server
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	rootLogger.Info("grpc health server created and registered")

	///////////////////////// REPOSITORIES INITIALIZATION /////////////////////////
	userProjectionRepo := postgresAdapters.NewUserProjectionRepository(pgClient)
	conversationRepo := mongodbAdapters.NewConversationRepository(mongoClient)
	messageRepo := mongodbAdapters.NewMessageRepository(mongoClient)

	///////////////////////// USE CASES INITIALIZATION /////////////////////////
	userProjectionUC := userProjectionUsecaseImpl.NewUserProjectionUsecase(userProjectionRepo)
	chatUC := chatUsecaseImpl.NewChatUsecase(conversationRepo, messageRepo, userProjectionRepo)

	///////////////////////// EVENT HANDLER INITIALIZATION /////////////////////////
	userEventHandler := eventHandlers.NewUserEventListener(messagingClient, userProjectionUC, rootLogger)

	///////////////////////// GRPC HANDLER INITIALIZATION /////////////////////////
	chatHandler := v1.NewChatHandler(chatUC, rootLogger)
	chatpbv1.RegisterChatServiceServer(grpcServer, chatHandler)

	///////////////////////// EVENT LISTENER INITIALIZATION /////////////////////////
	go func() {
		if err := userEventHandler.StartListening(serverCtx); err != nil {
			rootLogger.Error("Failed to start user event listener", zap.Error(err))
		}
	}()

	// mark healthy once all deps initialized successfully
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		healthSrv:       healthSrv,
		pgClient:        pgClient,
		mongoClient:     mongoClient,
		messagingClient: messagingClient,
		logger:          rootLogger,
		ctx:             serverCtx,
		cancel:          cancel,
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
	s.cancel()
	// mark not-serving so readiness probe fails quickly
	if s.healthSrv != nil {
		s.healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	}

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
