package server

import (
	"context"
	"fmt"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	messageBroker "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/pubsub"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/messagebroker/rabbitmq"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/notifications/smtp"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/config"

	// Adapter imports
	smtpAdapter "github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/adapters/smtp"

	// Use case imports
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/domain/usecase"

	// Handler imports
	eventHandlers "github.com/mohamedfawas/quboolkallyanam.xyz/services/notification/internal/handlers/event"

	"go.uber.org/zap"
)

type Server struct {
	config          *config.Config
	messagingClient messageBroker.Client
	smtpClient      *smtp.Client
	eventHandler    *eventHandlers.EventHandler
	ctx             context.Context    // internal server context
	cancel          context.CancelFunc // to cancel internal context
	logger          *zap.Logger
}

func NewServer(ctx context.Context, config *config.Config, rootLogger *zap.Logger) (*Server, error) {
	// Derive a cancellable child context from the parent
	// Example: if main() calls cancel(), serverCtx.Done() will be closed
	serverCtx, cancel := context.WithCancel(ctx)

	///////////////////////// SMTP INITIALIZATION /////////////////////////
	smtpClient, err := smtp.NewClient(smtp.Config{
		SMTPHost:     config.Email.SMTPHost,
		SMTPPort:     config.Email.SMTPPort,
		SMTPUsername: config.Email.SMTPUsername,
		SMTPPassword: config.Email.SMTPPassword,
		FromEmail:    config.Email.FromEmail,
		FromName:     config.Email.FromName,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create smtp client: %w", err)
	}
	rootLogger.Info("Connected to SMTP ")

	///////////////////////// MESSAGING CLIENT INITIALIZATION /////////////////////////
	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		// Use PubSub with parent ctx, so cancellation propagates
		messagingClient, err = pubsub.NewClient(ctx, config.PubSub.ProjectID)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		rootLogger.Info("Connected to PubSub ")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		rootLogger.Info("Connected to RabbitMQ ")
	}

	///////////////////////// ADAPTERS INITIALIZATION /////////////////////////
	emailAdapter := smtpAdapter.NewEmailAdapter(smtpClient)

	///////////////////////// USE CASES INITIALIZATION /////////////////////////
	notificationUsecase := usecase.NewNotificationUsecase(emailAdapter)

	///////////////////////// EVENT HANDLER INITIALIZATION /////////////////////////
	eventHandler := eventHandlers.NewEventHandler(messagingClient, notificationUsecase)

	//////////////////////// SERVER INITIALIZATION /////////////////////////
	server := &Server{
		config:          config,
		messagingClient: messagingClient,
		smtpClient:      smtpClient,
		eventHandler:    eventHandler,
		ctx:             serverCtx,
		cancel:          cancel,
		logger:          rootLogger,
	}

	return server, nil
}

func (s *Server) Start() error {
	s.logger.Info("Starting event listeners...")

	if err := s.eventHandler.StartListening(s.ctx); err != nil {
		return fmt.Errorf("failed to start event listeners: %w", err)
	}

	s.logger.Info("Event listeners started successfully")

	// Wait until serverCtx is cancelled (via Stop())
	<-s.ctx.Done()
	return s.ctx.Err() // Return the cancellation reason
}

func (s *Server) Stop() {
	s.logger.Info("Stopping Notification Service...")

	// Cancel the internal context
	// Example: if main() calls Stop(), serverCtx.Done() will be closed
	s.cancel()

	if s.messagingClient != nil {
		if err := s.messagingClient.Close(); err != nil {
			log.Printf("Failed to close messaging client: %v", err)
		}
	}

	log.Println("Notification Service stopped gracefully")
}
