package server

import (
	"context"
	"fmt"
	"log"
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

	"google.golang.org/grpc"
)

type Server struct {
	config          *config.Config
	grpcServer      *grpc.Server
	pgClient        *postgres.Client
	messagingClient messageBroker.Client
	razorpayService *razorpay.Service
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
	log.Println("Payment Service Connected to PostgreSQL")

	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		ctx := context.Background()
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			// Clean up existing connections before returning error
			pgClient.Close()
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		log.Println("Payment Service Connected to PubSub")
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
		log.Println("Payment Service Connected to RabbitMQ")
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.UnaryErrorInterceptor(),
	))

	// Initialize Razorpay service
	razorpayService := razorpay.NewService(config.Razorpay.KeyID, config.Razorpay.KeySecret)
	log.Println("Payment Service Connected to Razorpay")

	// Initialize repositories
	paymentsRepo := postgresAdapters.NewPaymentsRepository(pgClient)
	subscriptionPlansRepo := postgresAdapters.NewSubscriptionPlansRepository(pgClient)
	subscriptionsRepo := postgresAdapters.NewSubscriptionsRepository(pgClient)
	txManager := postgresAdapters.NewTxManager(pgClient)

	eventPublisher := messageBrokerAdapters.NewEventPublisher(messagingClient)

	// Initialize use cases
	paymentUC := paymentUsecase.NewPaymentUsecase(
		paymentsRepo,
		subscriptionPlansRepo,
		subscriptionsRepo,
		txManager,
		razorpayService,
		eventPublisher,
	)
	subscriptionUC := subscriptionUsecase.NewSubscriptionUsecase(subscriptionPlansRepo, subscriptionsRepo)

	// Initialize and register gRPC handler
	paymentHandler := grpcHandlerv1.NewPaymentHandler(paymentUC, subscriptionUC)
	paymentpbv1.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	log.Println("Payment Service gRPC handlers registered")

	server := &Server{
		config:          config,
		grpcServer:      grpcServer,
		pgClient:        pgClient,
		messagingClient: messagingClient,
		razorpayService: razorpayService,
	}

	return server, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.Port))
	if err != nil {
		return err
	}

	log.Printf("Payment Service gRPC server starting on port %d", s.config.GRPC.Port)

	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()

	// Close messaging client first
	if s.messagingClient != nil {
		if err := s.messagingClient.Close(); err != nil {
			log.Printf("failed to close messaging client: %v", err)
		}
	}

	if s.pgClient != nil {
		if err := s.pgClient.Close(); err != nil {
			log.Printf("failed to close postgres client: %v", err)
		}
	}

	log.Println("Payment Service stopped")
}
