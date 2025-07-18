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
)

type Server struct {
	config          *config.Config
	messagingClient messageBroker.Client
	smtpClient      *smtp.Client
	eventHandler    *eventHandlers.EventHandler
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewServer(ctx context.Context, config *config.Config) (*Server, error) {
	serverCtx, cancel := context.WithCancel(ctx)

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
	log.Println("Notification Service Connected to SMTP")

	var messagingClient messageBroker.Client
	if config.Environment == constants.EnvProduction {
		messagingClient, err = pubsub.NewClient(serverCtx, config.PubSub.ProjectID)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create pubsub client: %w", err)
		}
		log.Println("Notification Service Connected to PubSub")
	} else {
		messagingClient, err = rabbitmq.NewClient(rabbitmq.Config{
			DSN:          config.RabbitMQ.DSN,
			ExchangeName: config.RabbitMQ.ExchangeName,
		})
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create rabbitmq client: %w", err)
		}
		log.Println("Notification Service Connected to RabbitMQ")
	}

	// Initialize adapters
	emailAdapter := smtpAdapter.NewEmailAdapter(smtpClient)

	// Initialize use cases
	notificationUsecase := usecase.NewNotificationUsecase(emailAdapter)

	// Initialize event handler
	eventHandler := eventHandlers.NewEventHandler(messagingClient, notificationUsecase)

	server := &Server{
		config:          config,
		messagingClient: messagingClient,
		smtpClient:      smtpClient,
		eventHandler:    eventHandler,
		ctx:             serverCtx,
		cancel:          cancel,
	}

	return server, nil
}

func (s *Server) Start() error {
	log.Println("Notification Service starting event listeners...")

	if err := s.eventHandler.StartListening(s.ctx); err != nil {
		return fmt.Errorf("failed to start event listeners: %w", err)
	}

	log.Println("Notification Service event listeners started successfully")

	// Block until context is cancelled
	<-s.ctx.Done()
	return nil
}

func (s *Server) Stop() {
	log.Println("Stopping Notification Service...")

	s.cancel()

	if s.messagingClient != nil {
		if err := s.messagingClient.Close(); err != nil {
			log.Printf("Failed to close messaging client: %v", err)
		}
	}

	log.Println("Notification Service stopped gracefully")
}
