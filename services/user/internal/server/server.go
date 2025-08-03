package server

import (
	"context"
	"fmt"
	"net"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/rabbitmq"
	messageBrokerAdapter "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/adapters/messageBroker"
	postgresAdapters "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/adapters/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/config"
	matchmaking "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase/match_making"
	userProfileUsecaseImpl "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase/user_profile"
	eventHandlers "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/handlers/event"
	grpcHandlerv1 "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/handlers/grpc/v1"
	interceptors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/grpc/interceptors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	pgClient        *postgres.Client
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

	///////////////////////// MESSAGING CLIENT INITIALIZATION /////////////////////////	
	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
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
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		rootLogger.Info("Connected to RabbitMQ ")
	}

	///////////////////////// GRPC SERVER INITIALIZATION /////////////////////////
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryErrorInterceptor()),
	)

	///////////////////////// REPOSITORIES INITIALIZATION /////////////////////////
	userProfileRepo := postgresAdapters.NewUserProfileRepository(pgClient)
	rootLogger.Info("User Profile Repository Initialized")
	partnerPreferencesRepo := postgresAdapters.NewPartnerPreferencesRepository(pgClient)
	rootLogger.Info("Partner Preferences Repository Initialized")
	profileMatchRepo := postgresAdapters.NewProfileMatchRepository(pgClient)
	rootLogger.Info("Profile Match Repository Initialized")
	mutualMatchRepo := postgresAdapters.NewMutualMatchRepository(pgClient)
	rootLogger.Info("Mutual Match Repository Initialized")
	transactionManager := postgres.NewTransactionManager(pgClient)
	rootLogger.Info("Transaction Manager Initialized")

	///////////////////////// EVENT PUBLISHER INITIALIZATION /////////////////////////
	eventPublisher := messageBrokerAdapter.NewEventPublisher(messagingClient, rootLogger)
	rootLogger.Info("Event Publisher Initialized")

	///////////////////////// USE CASES INITIALIZATION /////////////////////////
	userProfileUC := userProfileUsecaseImpl.NewUserProfileUsecase(userProfileRepo, partnerPreferencesRepo, eventPublisher)
	rootLogger.Info("User Profile Use Case Initialized")
	matchMakingUC := matchmaking.NewMatchMakingUsecase(userProfileRepo, partnerPreferencesRepo, profileMatchRepo, mutualMatchRepo, transactionManager)
	rootLogger.Info("Match Making Use Case Initialized")


	///////////////////////// EVENT HANDLER INITIALIZATION /////////////////////////
	authEventHandler := eventHandlers.NewAuthEventHandler(messagingClient, userProfileUC, rootLogger)
	rootLogger.Info("Auth Event Handler Initialized")

	///////////////////////// GRPC HANDLER INITIALIZATION /////////////////////////
	userHandler := grpcHandlerv1.NewUserHandler(userProfileUC, matchMakingUC, rootLogger)
	userpbv1.RegisterUserServiceServer(grpcServer, userHandler)
	rootLogger.Info("gRPC Handler Initialized")

	///////////////////////// EVENT LISTENER INITIALIZATION /////////////////////////
	go func() {
		if err := authEventHandler.StartListening(context.Background()); err != nil {
			rootLogger.Error("failed to start auth event listener", zap.Error(err))
		}
	}()

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
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

	if s.pgClient != nil {
		if err := s.pgClient.Close(); err != nil {
			s.logger.Error("failed to close postgres client", zap.Error(err))
		}
	}
}
