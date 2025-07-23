package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"

	// Client imports
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
	authGRPC "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client/grpc/auth/v1"
	paymentGRPC "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client/grpc/payment/v1"
	chatGRPC "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client/grpc/chat"
	userGRPC "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client/grpc/user"

	// Usecase imports
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
	authUsecase "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase/auth"
	chatUsecase "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase/chat"
	paymentUsecase "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase/payment"
	userUsecase "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase/user"

	// Handler imports
	authHandler "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/delivery/http/v1/auth"
	paymentHandler "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/delivery/http/v1/payment"
	chatHandler "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/delivery/http/v1/chat"
	userHandler "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/delivery/http/v1/user"
)

type Server struct {
	config     *config.Config
	httpServer *http.Server
	jwtManager *jwt.JWTManager

	// Interface-based clients (for dependency injection)
	authClient    client.AuthClient
	paymentClient client.PaymentClient
	chatClient    client.ChatClient
	userClient    client.UserClient

	// Usecases (interface-based)
	authUsecase    usecase.AuthUsecase
	paymentUsecase usecase.PaymentUsecase
	chatUsecase   usecase.ChatUsecase
	userUsecase   usecase.UserUsecase

	// Handlers
	authHandler    *authHandler.AuthHandler
	paymentHandler *paymentHandler.PaymentHandler
	chatHandler    *chatHandler.ChatHandler
	userHandler    *userHandler.UserHandler
}

func NewHTTPServer(config *config.Config) (*Server, error) {
	if config.Environment == constants.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &Server{
		config: config,
	}

	server.jwtManager = jwt.NewJWTManager(jwt.JWTConfig{
		SecretKey:          config.Auth.JWT.SecretKey,
		AccessTokenMinutes: config.Auth.JWT.AccessTokenMinutes,
		RefreshTokenDays:   config.Auth.JWT.RefreshTokenDays,
		Issuer:             config.Auth.JWT.Issuer,
	})

	if err := server.initClients(); err != nil {
		return nil, fmt.Errorf("failed to initialize clients: %w", err)
	}

	if err := server.initUsecases(); err != nil {
		return nil, fmt.Errorf("failed to initialize usecases: %w", err)
	}

	if err := server.initHandlers(); err != nil {
		return nil, fmt.Errorf("failed to initialize handlers: %w", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("internal/web/templates/*")

	server.setupRoutes(router)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.HTTP.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(config.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.HTTP.IdleTimeout) * time.Second,
	}

	server.httpServer = httpServer

	return server, nil
}

func (s *Server) initClients() error {
	ctx := context.Background()

	// Initialize Auth gRPC Client
	authGRPCClient, err := authGRPC.NewAuthGRPCClient(
		ctx,
		fmt.Sprintf("localhost:%s", s.config.Services.AuthServicePort),
		false, // useTLS - set to true in production
		nil,   // tlsConfig
	)
	if err != nil {
		return fmt.Errorf("failed to create auth gRPC client: %w", err)
	}

	s.authClient = authGRPCClient

	// Initialize Payment gRPC Client
	paymentGRPCClient, err := paymentGRPC.NewPaymentGRPCClient(
		ctx,
		fmt.Sprintf("localhost:%s", s.config.Services.PaymentServicePort),
		false, // useTLS - set to true in production
		nil,   // tlsConfig
	)
	if err != nil {
		return fmt.Errorf("failed to create payment gRPC client: %w", err)
	}
	s.paymentClient = paymentGRPCClient

	// Initialize Chat gRPC Client
	chatGRPCClient, err := chatGRPC.NewChatGRPCClient(
		ctx,
		fmt.Sprintf("localhost:%s", s.config.Services.ChatServicePort),
		false, // useTLS - set to true in production
		nil,   // tlsConfig
	)
	if err != nil {
		return fmt.Errorf("failed to create chat gRPC client: %w", err)
	}
	s.chatClient = chatGRPCClient


	// Initialize User gRPC Client
	userGRPCClient, err := userGRPC.NewUserGRPCClient(
		ctx,
		fmt.Sprintf("localhost:%s", s.config.Services.UserServicePort),
		false, // useTLS - set to true in production
		nil,   // tlsConfig
	)
	if err != nil {
		return fmt.Errorf("failed to create user gRPC client: %w", err)
	}
	s.userClient = userGRPCClient

	// TODO: Add other clients as you implement them

	return nil
}

func (s *Server) initUsecases() error {
	s.authUsecase = authUsecase.NewAuthUsecase(s.authClient)
	s.paymentUsecase = paymentUsecase.NewPaymentUsecase(s.paymentClient, s.config)
	s.chatUsecase = chatUsecase.NewChatUsecase(s.chatClient)
	s.userUsecase = userUsecase.NewUserUsecase(s.userClient)
	// TODO: Add other usecases as you implement them
	// s.adminUsecase = adminUsecase.NewAdminUsecase(s.adminClient)
	// s.userUsecase = userUsecase.NewUserUsecase(s.userClient)
	// s.chatUsecase = chatUsecase.NewChatUsecase(s.chatClient)
	// s.paymentUsecase = paymentUsecase.NewPaymentUsecase(s.paymentClient)

	return nil
}

func (s *Server) initHandlers() error {
	s.authHandler = authHandler.NewAuthHandler(s.authUsecase, *s.config)
	s.paymentHandler = paymentHandler.NewPaymentHandler(s.paymentUsecase)
	s.chatHandler = chatHandler.NewChatHandler(s.chatUsecase)
	s.userHandler = userHandler.NewUserHandler(s.userUsecase)
	// TODO: Add other handlers as you implement them
	// s.adminHandler = adminHandler.NewAdminHandler(s.adminUsecase, *s.config)
	// s.userHandler = userHandler.NewUserHandler(s.userUsecase, *s.config)
	// s.chatHandler = chatHandler.NewChatHandler(s.chatUsecase, *s.config)
	// s.paymentHandler = paymentHandler.NewPaymentHandler(s.paymentUsecase, *s.config)

	return nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
