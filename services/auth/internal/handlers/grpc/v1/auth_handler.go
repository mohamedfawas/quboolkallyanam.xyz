package v1

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"

	authpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/auth/v1"
	logger "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"
)

type AuthHandler struct {
	authpbv1.UnimplementedAuthServiceServer
	userUsecase                usecase.UserUsecase
	adminUsecase               usecase.AdminUsecase
	pendingRegistrationUsecase usecase.PendingRegistrationUsecase
	config                     *config.Config
}

func NewAuthHandler(userUsecase usecase.UserUsecase,
	adminUsecase usecase.AdminUsecase,
	pendingRegistrationUsecase usecase.PendingRegistrationUsecase,
	config *config.Config) *AuthHandler {
	return &AuthHandler{
		userUsecase:                userUsecase,
		adminUsecase:               adminUsecase,
		pendingRegistrationUsecase: pendingRegistrationUsecase,
		config:                     config,
	}
}

func (h *AuthHandler) UserRegister(ctx context.Context, req *authpbv1.UserRegisterRequest) (*authpbv1.UserRegisterResponse, error) {
	logger.Log.Info("Received UserRegistration request in auth handler", "email : ", req.Email, "phone : ", req.Phone)

	userRegistrationRequest := &entity.UserRegistrationRequest{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	logger.Log.Info("Sending request to pendingRegistrationUsecase from auth handler", "email : ", req.Email, "phone : ", req.Phone)
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

	logger.Log.Info("UserRegistration request successful ", "email : ", req.Email, "phone : ", req.Phone)
	return &authpbv1.UserRegisterResponse{
		Email: req.Email,
		Phone: req.Phone,
	}, nil
}

func (h *AuthHandler) UserVerification(ctx context.Context, req *authpbv1.UserVerificationRequest) (*authpbv1.UserVerificationResponse, error) {
	logger.Log.Info("Received UserVerification request in auth handler", "email : ", req.Email)

	err := h.pendingRegistrationUsecase.VerifyUserRegistration(ctx, req.Email, req.Otp)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidOTP):
			return nil, status.Errorf(codes.InvalidArgument, "Invalid OTP")
		case errors.Is(err, appErrors.ErrPendingRegistrationNotFound):
			return nil, status.Errorf(codes.NotFound, "Registration request not found for this email")
		case errors.Is(err, appErrors.ErrOTPNotFound):
			return nil, status.Errorf(codes.NotFound, "No OTP request found for this email")
		default:
			return nil, status.Errorf(codes.Internal, "Failed to verify user: %v", err)
		}
	}

	return &authpbv1.UserVerificationResponse{
		Success: "User verified successfully",
	}, nil
}

func (h *AuthHandler) UserLogin(ctx context.Context, req *authpbv1.UserLoginRequest) (*authpbv1.UserLoginResponse, error) {
	logger.Log.Info("Received UserLogin request in auth handler", "email : ", req.Email)

	result, err := h.userUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidCredentials):
			return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
		default:
			return nil, status.Errorf(codes.Internal, "Failed to login: %v", err)
		}
	}

	return &authpbv1.UserLoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *AuthHandler) UserLogout(ctx context.Context, req *authpbv1.UserLogoutRequest) (*emptypb.Empty, error) {
	logger.Log.Info("Received UserLogout request in auth handler")

	err := h.userUsecase.Logout(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to logout: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) UserDelete(ctx context.Context, req *authpbv1.UserDeleteRequest) (*emptypb.Empty, error) {
	logger.Log.Info("Received UserDelete request in auth handler")

	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	err = h.userUsecase.UserAccountDelete(ctx, userID, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidPassword):
			return nil, status.Errorf(codes.InvalidArgument, "Invalid password")
		default:
			return nil, status.Errorf(codes.Internal, "Failed to delete account: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func extractUserIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata found in context")
	}

	userIDs := md.Get(constants.ContextKeyUserID)
	if len(userIDs) == 0 {
		return "", fmt.Errorf("user ID not found in metadata")
	}

	return userIDs[0], nil
}

func (h *AuthHandler) AdminLogin(ctx context.Context, req *authpbv1.AdminLoginRequest) (*authpbv1.AdminLoginResponse, error) {
	logger.Log.Info("Received AdminLogin request in auth handler", "email : ", req.Email)

	result, err := h.adminUsecase.AdminLogin(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidCredentials):
			return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
		default:
			return nil, status.Errorf(codes.Internal, "Failed to login: %v", err)
		}
	}

	return &authpbv1.AdminLoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *AuthHandler) AdminLogout(ctx context.Context, req *authpbv1.AdminLogoutRequest) (*emptypb.Empty, error) {
	logger.Log.Info("Received AdminLogout request in auth handler")

	err := h.adminUsecase.AdminLogout(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to logout: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpbv1.RefreshTokenRequest) (*authpbv1.RefreshTokenResponse, error) {
	logger.Log.Info("Received RefreshToken request in auth handler")

	result, err := h.userUsecase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidToken):
			return nil, status.Errorf(codes.Unauthenticated, "Invalid refresh token")
		default:
			return nil, status.Errorf(codes.Internal, "Failed to refresh token: %v", err)
		}
	}

	return &authpbv1.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}
