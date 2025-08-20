package server

import (
	"context"
	"fmt"
	"net"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	interceptors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/grpc/interceptors"
	gcsstore "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/mediastorage/gcs"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/rabbitmq"
	gcs "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/adapters/gcs"
	messageBrokerAdapter "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/adapters/messageBroker"
	postgresAdapters "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/adapters/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/config"
	matchmaking "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase/match_making"
	userProfileUsecaseImpl "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase/userprofile"
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
	gcsStore        *gcsstore.GCSStore
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
	///////////////////////// GCS STORE INITIALIZATION /////////////////////////
	gcsConfig := gcsstore.MediaStorageConfig{
		Bucket:      config.MediaStorage.Bucket,
		URLExpiry:   config.MediaStorage.URLExpiry,
		SignerEmail: config.MediaStorage.SignerEmail, 
	}

	gcsStore, err := gcsstore.NewGCSStore(ctx, gcsConfig)
	if err != nil {
		pgClient.Close()
		if messagingClient != nil {
			messagingClient.Close()
		}
		return nil, fmt.Errorf("failed to create GCS store: %w", err)
	}
	rootLogger.Info("Connected to GCS")
	///////////////////////// GRPC SERVER INITIALIZATION /////////////////////////
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryErrorInterceptor()),
	)

	///////////////////////// REPOSITORIES INITIALIZATION /////////////////////////
	userProfileRepo := postgresAdapters.NewUserProfileRepository(pgClient)
	userImageRepo := postgresAdapters.NewUserImageRepository(pgClient)
	partnerPreferencesRepo := postgresAdapters.NewPartnerPreferencesRepository(pgClient)
	profileMatchRepo := postgresAdapters.NewProfileMatchRepository(pgClient)
	mutualMatchRepo := postgresAdapters.NewMutualMatchRepository(pgClient)
	transactionManager := postgres.NewTransactionManager(pgClient)

	///////////////////////// EVENT PUBLISHER INITIALIZATION /////////////////////////
	eventPublisher := messageBrokerAdapter.NewEventPublisher(messagingClient, rootLogger)

	///////////////////////// MEDIA STORAGE INITIALIZATION /////////////////////////
	photoStorage := gcs.NewPhotoStorageAdapter(gcsStore, config.MediaStorage.Bucket)

	///////////////////////// USE CASES INITIALIZATION /////////////////////////
	userProfileUC := userProfileUsecaseImpl.NewUserProfileUsecase(userProfileRepo, userImageRepo, partnerPreferencesRepo, eventPublisher, photoStorage, config)
	matchMakingUC := matchmaking.NewMatchMakingUsecase(userProfileRepo, partnerPreferencesRepo, profileMatchRepo, mutualMatchRepo, transactionManager, photoStorage, config, eventPublisher)

	///////////////////////// EVENT HANDLER INITIALIZATION /////////////////////////
	authEventHandler := eventHandlers.NewAuthEventHandler(messagingClient, userProfileUC, rootLogger)

	///////////////////////// GRPC HANDLER INITIALIZATION /////////////////////////
	userHandler := grpcHandlerv1.NewUserHandler(userProfileUC, matchMakingUC, rootLogger)
	userpbv1.RegisterUserServiceServer(grpcServer, userHandler)

	///////////////////////// EVENT LISTENER INITIALIZATION /////////////////////////
	go func() {
		if err := authEventHandler.StartListening(serverCtx); err != nil {
			rootLogger.Error("failed to start auth event listener", zap.Error(err))
		}
	}()

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
		messagingClient: messagingClient,
		gcsStore:        gcsStore,
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
	s.grpcServer.GracefulStop()

	// Close GCS store
	if s.gcsStore != nil {
		if err := s.gcsStore.Close(); err != nil {
			s.logger.Error("failed to close GCS store", zap.Error(err))
		}
	}

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
