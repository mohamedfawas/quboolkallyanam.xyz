package v1

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"

	authpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/auth/v1"
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
			log.Printf("Failed to register user: %v", err)
			return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
		}
	}

	return &authpbv1.UserRegisterResponse{
		Email: req.Email,
		Phone: req.Phone,
	}, nil
}

func (h *AuthHandler) UserVerification(ctx context.Context, req *authpbv1.UserVerificationRequest) (*authpbv1.UserVerificationResponse, error) {
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
			log.Printf("Failed to verify user: %v", err)
			return nil, status.Errorf(codes.Internal, "Failed to verify user: %v", err)
		}
	}

	return &authpbv1.UserVerificationResponse{
		Success: "User verified successfully",
	}, nil
}

func (h *AuthHandler) UserLogin(ctx context.Context, req *authpbv1.UserLoginRequest) (*authpbv1.UserLoginResponse, error) {
	result, err := h.userUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidCredentials):
			return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
		case errors.Is(err, appErrors.ErrUserNotFound):
			return nil, status.Errorf(codes.NotFound, "User not found")
		case errors.Is(err, appErrors.ErrAccountDisabled):
			return nil, status.Errorf(codes.PermissionDenied, "Account disabled")
		case errors.Is(err, appErrors.ErrAccountBlocked):
			return nil, status.Errorf(codes.PermissionDenied, "Account blocked")
		default:
			log.Printf("Failed to login: %v", err)
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
	err := h.userUsecase.Logout(ctx, req.AccessToken)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidToken):
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
		case errors.Is(err, appErrors.ErrExpiredToken):
			return nil, status.Errorf(codes.Unauthenticated, "Token expired")
		case errors.Is(err, appErrors.ErrTokenNotActive):
			return nil, status.Errorf(codes.Unauthenticated, "Token not active")
		case errors.Is(err, appErrors.ErrUnauthorized):
			return nil, status.Errorf(codes.PermissionDenied, "Unauthorized")
		default:
			log.Printf("Failed to logout: %v", err)
			return nil, status.Errorf(codes.Internal, "Failed to logout: %v", err)
		}
	}
	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) UserDelete(ctx context.Context, req *authpbv1.UserDeleteRequest) (*emptypb.Empty, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	err = h.userUsecase.UserAccountDelete(ctx, userID, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrUserNotFound):
			return nil, status.Errorf(codes.NotFound, "User not found")
		case errors.Is(err, appErrors.ErrInvalidCredentials):
			return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
		default:
			log.Printf("Failed to delete account: %v", err)
			return nil, status.Errorf(codes.Internal, "Failed to delete account: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) AdminLogin(ctx context.Context, req *authpbv1.AdminLoginRequest) (*authpbv1.AdminLoginResponse, error) {
	result, err := h.adminUsecase.AdminLogin(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidPassword):
			return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
		case errors.Is(err, appErrors.ErrAdminNotFound):
			return nil, status.Errorf(codes.NotFound, "Admin not found")
		case errors.Is(err, appErrors.ErrAdminAccountDisabled):
			return nil, status.Errorf(codes.PermissionDenied, "Admin account disabled")
		default:
			log.Printf("Failed to login: %v", err)
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
	err := h.adminUsecase.AdminLogout(ctx, req.AccessToken)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidToken):
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
		case errors.Is(err, appErrors.ErrExpiredToken):
			return nil, status.Errorf(codes.Unauthenticated, "Token expired")
		case errors.Is(err, appErrors.ErrTokenNotActive):
			return nil, status.Errorf(codes.Unauthenticated, "Token not active")
		case errors.Is(err, appErrors.ErrUnauthorized):
			return nil, status.Errorf(codes.PermissionDenied, "Unauthorized")
		default:
			log.Printf("Failed to logout: %v", err)
			return nil, status.Errorf(codes.Internal, "Failed to logout: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpbv1.RefreshTokenRequest) (*authpbv1.RefreshTokenResponse, error) {
	result, err := h.userUsecase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
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
