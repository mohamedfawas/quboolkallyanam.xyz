package server

import (
	"context"
	"fmt"
	"net"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/redis"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/grpc/interceptors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/notifications/smtp"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"

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

	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	pgClient        *postgres.Client
	redisClient     *redis.Client
	messagingClient messageBroker.Client
	jwtManager      jwt.JWTManager
	smtpClient      smtp.Client
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
	logger.Log.Info("✅ Auth Service Connected to PostgreSQL ")

	redisClient, err := redis.NewClient(redis.Config{
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
	logger.Log.Info("✅ Auth Service Connected to Redis")

	var messagingClient messageBroker.Client
	if config.Environment == "production" {
		ctx := context.Background()
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			redisClient.Close()
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		logger.Log.Info("✅ Auth Service Connected to PubSub")
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
		logger.Log.Info("✅ Auth Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.UnaryErrorInterceptor(),
	))

	// Initialize repositories
	userRepo := postgresAdapters.NewUserRepository(pgClient)
	adminRepo := postgresAdapters.NewAdminRepository(pgClient)
	pendingRegistrationRepo := postgresAdapters.NewPendingRegistrationRepository(pgClient)
	tokenRepo := redisAdapters.NewTokenRepository(redisClient)
	otpRepo := redisAdapters.NewOTPRepository(redisClient)

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager(jwt.JWTConfig{
		SecretKey:          config.Auth.JWT.SecretKey,
		AccessTokenMinutes: config.Auth.JWT.AccessTokenMinutes,
		RefreshTokenDays:   config.Auth.JWT.RefreshTokenDays,
		Issuer:             config.Auth.JWT.Issuer,
	})

	// Initialize SMTP client
	smtpClient, err := smtp.NewClient(smtp.Config{
		SMTPHost:     config.Email.SMTPHost,
		SMTPPort:     config.Email.SMTPPort,
		SMTPUsername: config.Email.SMTPUsername,
		SMTPPassword: config.Email.SMTPPassword,
		FromEmail:    config.Email.FromEmail,
		FromName:     config.Email.FromName,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create smtp client: %w", err)
	}

	// Initialize use cases
	userUC := userUsecase.NewUserUseCase(userRepo, *jwtManager, tokenRepo, config, messagingClient)
	adminUC := adminUsecase.NewAdminUsecase(adminRepo, tokenRepo, *jwtManager, config)
	pendingRegistrationUC := pendingRegistrationUsecase.NewPendingRegistrationUsecase(
		pendingRegistrationRepo,
		userRepo,
		otpRepo,
		*smtpClient,
	)

	// Initialize and register gRPC handler
	authHandler := grpcHandlerv1.NewAuthHandler(userUC, adminUC, pendingRegistrationUC, config)
	authpbv1.RegisterAuthServiceServer(grpcServer, authHandler)

	logger.Log.Info("✅ Auth Service gRPC handlers registered")

	// Initialize default admin if needed
	ctx := context.Background()
	if err := adminUC.InitializeDefaultAdmin(ctx, config.Admin.DefaultAdminEmail, config.Admin.DefaultAdminPassword); err != nil {
		logger.Log.Error("Failed to initialize default admin", "error", err)
	}

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
		redisClient:     redisClient,
		messagingClient: messagingClient,
		jwtManager:      *jwtManager,
		smtpClient:      *smtpClient,
	}

	return server, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.Port))
	if err != nil {
		return err
	}

	logger.Log.Info("🚀 Auth Service gRPC server starting", "port", s.config.GRPC.Port)

	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()

	// Close messaging client first
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

	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			logger.Log.Error("failed to close redis client: %w", err)
		}
	}
}
