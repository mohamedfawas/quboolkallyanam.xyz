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

func NewServer(config *config.Config, rootLogger *zap.Logger) (*Server, error) {

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

	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		ctx := context.Background()
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
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
	}

	grpcServer := grpc.NewServer()

	// Initialize repositories
	userProfileRepo := postgresAdapters.NewUserProfileRepository(pgClient)
	partnerPreferencesRepo := postgresAdapters.NewPartnerPreferencesRepository(pgClient)
	profileMatchRepo := postgresAdapters.NewProfileMatchRepository(pgClient)
	mutualMatchRepo := postgresAdapters.NewMutualMatchRepository(pgClient)
	transactionManager := postgres.NewTransactionManager(pgClient)
	// Initialize event publisher
	eventPublisher := messageBrokerAdapter.NewEventPublisher(messagingClient, rootLogger)

	// Initialize use cases
	userProfileUC := userProfileUsecaseImpl.NewUserProfileUsecase(userProfileRepo, partnerPreferencesRepo, eventPublisher)
	matchMakingUC := matchmaking.NewMatchMakingUsecase(userProfileRepo, partnerPreferencesRepo, profileMatchRepo, mutualMatchRepo, transactionManager)

	// Initialize and register gRPC handler
	userHandler := grpcHandlerv1.NewUserHandler(userProfileUC, matchMakingUC)
	userpbv1.RegisterUserServiceServer(grpcServer, userHandler)

	// Initialize event handler
	authEventHandler := eventHandlers.NewAuthEventHandler(messagingClient, userProfileUC, rootLogger)

	// Start event listener
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
