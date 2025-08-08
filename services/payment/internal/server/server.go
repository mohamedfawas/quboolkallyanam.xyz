package server

import (
	"context"
	"fmt"
	"net"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/grpc/interceptors"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/payment/razorpay"
	messageBrokerAdapters "github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/adapters/messageBroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/config"

	// Proto imports
	paymentpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/payment/v1"

	// Repository imports
	postgresAdapters "github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/adapters/postgres"

	// Use case imports
	paymentUsecase "github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/usecase/payments"
	subscriptionUsecase "github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/usecase/subscription"

	// Handler imports
	grpcHandlerv1 "github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/handlers/grpc/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	pgClient        *postgres.Client
	messagingClient messageBroker.Client
	razorpayService *razorpay.Service
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

	///////////////////////// MESSAGING INITIALIZATION /////////////////////////
	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		ctx := context.Background()
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
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.UnaryErrorInterceptor(),
	))
	rootLogger.Info("gRPC server created")

	///////////////////////// RAZORPAY SERVICE INITIALIZATION /////////////////////////
	razorpayService := razorpay.NewService(config.Razorpay.KeyID, config.Razorpay.KeySecret)
	rootLogger.Info("Razorpay service created")

	///////////////////////// REPOSITORIES INITIALIZATION /////////////////////////
	paymentsRepo := postgresAdapters.NewPaymentsRepository(pgClient)
	rootLogger.Info("Payments repository created")
	subscriptionPlansRepo := postgresAdapters.NewSubscriptionPlansRepository(pgClient)
	rootLogger.Info("Subscription plans repository created")
	subscriptionsRepo := postgresAdapters.NewSubscriptionsRepository(pgClient)
	rootLogger.Info("Subscriptions repository created")
	txManager := postgres.NewTransactionManager(pgClient)
	rootLogger.Info("Transaction manager created")

	///////////////////////// EVENT PUBLISHER INITIALIZATION /////////////////////////
	eventPublisher := messageBrokerAdapters.NewEventPublisher(messagingClient)
	rootLogger.Info("Event publisher created")

	///////////////////////// USE CASES INITIALIZATION /////////////////////////
	paymentUC := paymentUsecase.NewPaymentUsecase(
		paymentsRepo,
		subscriptionPlansRepo,
		subscriptionsRepo,
		txManager,
		razorpayService,
		eventPublisher,
	)
	rootLogger.Info("Payment use case created")
	subscriptionUC := subscriptionUsecase.NewSubscriptionUsecase(subscriptionPlansRepo, subscriptionsRepo)
	rootLogger.Info("Subscription use case created")

	///////////////////////// GRPC HANDLER INITIALIZATION /////////////////////////
	paymentHandler := grpcHandlerv1.NewPaymentHandler(paymentUC, subscriptionUC)
	paymentpbv1.RegisterPaymentServiceServer(grpcServer, paymentHandler)
	rootLogger.Info("Payment Service gRPC handlers registered")
	

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
		messagingClient: messagingClient,
		razorpayService: razorpayService,
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
