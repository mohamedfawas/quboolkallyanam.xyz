package server

import (
	"context"
	"fmt"
	"log"
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
		log.Println("failed to create postgres client", err)
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}
	log.Println("User Service Connected to PostgreSQL ")

	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		ctx := context.Background()
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			log.Println("failed to create pubsub client", err)
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		log.Println("User Service Connected to PubSub")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			log.Println("failed to create rabbitmq client", err)
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		log.Println("User Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer()

	// Initialize repositories
	userProfileRepo := postgresAdapters.NewUserProfileRepository(pgClient)
	partnerPreferencesRepo := postgresAdapters.NewPartnerPreferencesRepository(pgClient)
	profileMatchRepo := postgresAdapters.NewProfileMatchRepository(pgClient)
	mutualMatchRepo := postgresAdapters.NewMutualMatchRepository(pgClient)
	transactionManager := postgres.NewTransactionManager(pgClient)
	// Initialize event publisher
	eventPublisher := messageBrokerAdapter.NewEventPublisher(messagingClient)

	// Initialize use cases
	userProfileUC := userProfileUsecaseImpl.NewUserProfileUsecase(userProfileRepo, partnerPreferencesRepo, eventPublisher)
	matchMakingUC := matchmaking.NewMatchMakingUsecase(userProfileRepo, partnerPreferencesRepo, profileMatchRepo, mutualMatchRepo, transactionManager)

	// Initialize and register gRPC handler
	userHandler := grpcHandlerv1.NewUserHandler(userProfileUC, matchMakingUC)
	userpbv1.RegisterUserServiceServer(grpcServer, userHandler)

	// Initialize event handler
	authEventHandler := eventHandlers.NewAuthEventHandler(messagingClient, userProfileUC)

	// Start event listener
	go func() {
		if err := authEventHandler.StartListening(context.Background()); err != nil {
			log.Println("Failed to start auth event handler", err)
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

	log.Println("User Service starting on port ", s.config.GRPC.Port)
	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()

	if s.messagingClient != nil {
		if err := s.messagingClient.Close(); err != nil {
			log.Println("failed to close messaging client", err)
		}
	}

	if s.pgClient != nil {
		if err := s.pgClient.Close(); err != nil {
			log.Println("failed to close postgres client", err)
		}
	}
}
