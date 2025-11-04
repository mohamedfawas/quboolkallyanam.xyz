package server

import (
	"context"
	"fmt"
	"net"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/redis"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/grpc/interceptors"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
	messageBrokerAdapter "github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/adapters/messageBroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	eventHandlers "github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/handlers/event"

	// Proto imports
	authpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/auth/v1"

	// Repository imports
	postgresAdapters "github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/adapters/postgres"
	redisAdapters "github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/adapters/redis"

	// Use case imports
	adminUsecase "github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase/admin"
	pendingRegistrationUsecase "github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase/pending_registration"
	userUsecase "github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase/user"

	// Handler imports
	grpcHandlerv1 "github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/handlers/grpc/v1"

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
	redisClient     *redis.Client
	messagingClient messageBroker.Client
	jwtManager      jwt.JWTManager
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

	///////////////////////// REDIS INITIALIZATION /////////////////////////
	redisClient, err := redis.NewClient(redis.Config{
		RedisURL: config.Redis.RedisURL,
		Host:     config.Redis.Host,
		Port:     config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if err != nil {
		// Clean up existing connections before returning error
		pgClient.Close()
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}
	rootLogger.Info("Connected to Redis ",
		zap.Int("port", config.Redis.Port),
		zap.Int("db", config.Redis.DB),
	)

	///////////////////////// MESSAGING CLIENT INITIALIZATION /////////////////////////
	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		messagingClient, err = pubsub.NewClient(serverCtx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			redisClient.Close()
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
			redisClient.Close()
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		rootLogger.Info("Connected to RabbitMQ ")
	}

	///////////////////////// GRPC SERVER INITIALIZATION /////////////////////////
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.UnaryErrorInterceptor(),
	))

	///////////////////////// HEALTH SERVER INITIALIZATION /////////////////////////
	// create and register simple standard gRPC health server
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	rootLogger.Info("grpc health server created and registered")

	///////////////////////// REPOSITORIES INITIALIZATION /////////////////////////
	userRepo := postgresAdapters.NewUserRepository(pgClient)
	adminRepo := postgresAdapters.NewAdminRepository(pgClient)
	pendingRegistrationRepo := postgresAdapters.NewPendingRegistrationRepository(pgClient)
	tokenRepo := redisAdapters.NewTokenRepository(redisClient)
	otpRepo := redisAdapters.NewOTPRepository(redisClient)

	///////////////////////// JWT MANAGER INITIALIZATION /////////////////////////
	jwtManager := jwt.NewJWTManager(jwt.JWTConfig{
		SecretKey:          config.Auth.JWT.SecretKey,
		AccessTokenMinutes: config.Auth.JWT.AccessTokenMinutes,
		RefreshTokenDays:   config.Auth.JWT.RefreshTokenDays,
		Issuer:             config.Auth.JWT.Issuer,
	})

	///////////////////////// EVENT PUBLISHER INITIALIZATION /////////////////////////
	eventPublisher := messageBrokerAdapter.NewEventPublisher(messagingClient, rootLogger)

	///////////////////////// USE CASES INITIALIZATION /////////////////////////
	userUC := userUsecase.NewUserUseCase(
		userRepo,
		*jwtManager,
		tokenRepo,
		config,
		messagingClient,
		eventPublisher,
	)

	adminUC := adminUsecase.NewAdminUsecase(
		adminRepo,
		tokenRepo,
		userRepo,
		*jwtManager,
		eventPublisher,
		config,
	)

	pendingRegistrationUC := pendingRegistrationUsecase.NewPendingRegistrationUsecase(
		pendingRegistrationRepo,
		userRepo,
		otpRepo,
		eventPublisher,
	)

	///////////////////////// EVENT HANDLER INITIALIZATION /////////////////////////
	paymentEventHandler := eventHandlers.NewPaymentEventHandler(messagingClient, userUC, rootLogger)

	///////////////////////// GRPC HANDLER INITIALIZATION /////////////////////////
	authHandler := grpcHandlerv1.NewAuthHandler(userUC, adminUC, pendingRegistrationUC, config, rootLogger)
	authpbv1.RegisterAuthServiceServer(grpcServer, authHandler)

	///////////////////////// DEFAULT ADMIN INITIALIZATION /////////////////////////
	if err := adminUC.InitializeDefaultAdmin(ctx, config.Admin.DefaultAdminEmail, config.Admin.DefaultAdminPassword); err != nil {
		rootLogger.Error("Failed to initialize default admin", zap.Error(err))
	}
	rootLogger.Info("Default Admin Initialized")

	///////////////////////// EVENT LISTENER INITIALIZATION /////////////////////////
	go func() {
		if err := paymentEventHandler.StartListening(serverCtx); err != nil {
			rootLogger.Error("Failed to start payment event handler", zap.Error(err))
		}
	}()

	// mark healthy once all deps initialized successfully
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		healthSrv:       healthSrv,
		pgClient:        pgClient,
		redisClient:     redisClient,
		messagingClient: messagingClient,
		jwtManager:      *jwtManager,
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

	if s.pgClient != nil {
		if err := s.pgClient.Close(); err != nil {
			s.logger.Error("failed to close postgres client", zap.Error(err))
		}
	}

	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			s.logger.Error("failed to close redis client", zap.Error(err))
		}
	}
}
