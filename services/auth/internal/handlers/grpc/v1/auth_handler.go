package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"

	authpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/auth/v1"
	logger "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
)

type authHandler struct {
	authUsecase                usecase.UserUsecase
	adminUsecase               usecase.AdminUsecase
	pendingRegistrationUsecase usecase.PendingRegistrationUsecase
	config                     *config.Config
}

func NewAuthHandler(authUsecase usecase.UserUsecase,
	adminUsecase usecase.AdminUsecase,
	pendingRegistrationUsecase usecase.PendingRegistrationUsecase,
	config *config.Config) *authHandler {
	return &authHandler{
		authUsecase:                authUsecase,
		adminUsecase:               adminUsecase,
		pendingRegistrationUsecase: pendingRegistrationUsecase,
	}
}

func (h *authHandler) UserRegister(ctx context.Context, req *authpbv1.UserRegisterRequest) (*authpbv1.UserRegisterResponse, error) {
	logger.Log.Info("Received UserRegistration request ", "email ", req.Email, "phone ", req.Phone)

	userRegistrationRequest := &entity.UserRegistrationRequest{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	err := h.pendingRegistrationUsecase.RegisterUser(ctx, userRegistrationRequest, h.config)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrEmailAlreadyExists):
			return nil, status.Errorf(codes.AlreadyExists, "An account with this email already exists")
		case errors.Is(err, appErrors.ErrPhoneAlreadyExists):
			return nil, status.Errorf(codes.AlreadyExists, "An account with this phone number already exists")
		default:
			return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
		}
	}

	logger.Log.Info("UserRegistration request successful ", "email ", req.Email, "phone ", req.Phone)
	return &authpbv1.UserRegisterResponse{
		Email: req.Email,
		Phone: req.Phone,
	}, nil
}
