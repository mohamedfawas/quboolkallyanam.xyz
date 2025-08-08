package v1

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"

	authpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/auth/v1"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authpbv1.UnimplementedAuthServiceServer
	userUsecase                usecase.UserUsecase
	adminUsecase               usecase.AdminUsecase
	pendingRegistrationUsecase usecase.PendingRegistrationUsecase
	config                     *config.Config
	logger                     *zap.Logger
}

func NewAuthHandler(userUsecase usecase.UserUsecase,
	adminUsecase usecase.AdminUsecase,
	pendingRegistrationUsecase usecase.PendingRegistrationUsecase,
	config *config.Config,
	logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userUsecase:                userUsecase,
		adminUsecase:               adminUsecase,
		pendingRegistrationUsecase: pendingRegistrationUsecase,
		config:                     config,
		logger:                     logger,
	}
}

func (h *AuthHandler) UserRegister(
	ctx context.Context,
	req *authpbv1.UserRegisterRequest) (*authpbv1.UserRegisterResponse, error) {

	contextData, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
	)

	userRegistrationRequest := &entity.UserRegistrationRequest{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	err = h.pendingRegistrationUsecase.RegisterUser(ctx, userRegistrationRequest, h.config)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to register user", zap.Error(err))
		}
		return nil, err
	}

	log.Info("User register request processed successfully", zap.String("email", req.Email), zap.String("phone", req.Phone))
	return &authpbv1.UserRegisterResponse{
		Email: req.Email,
		Phone: req.Phone,
	}, nil
}

func (h *AuthHandler) UserVerification(ctx context.Context, req *authpbv1.UserVerificationRequest) (*authpbv1.UserVerificationResponse, error) {
	contextData, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
	)

	err = h.pendingRegistrationUsecase.VerifyUserRegistration(ctx, req.Email, req.Otp)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to verify user", zap.Error(err))
		}
		return nil, err
	}

	log.Info("User verification request processed successfully", zap.String("email", req.Email))
	return &authpbv1.UserVerificationResponse{
		Success: &wrapperspb.BoolValue{Value: true},
	}, nil
}

func (h *AuthHandler) UserLogin(ctx context.Context, req *authpbv1.UserLoginRequest) (*authpbv1.UserLoginResponse, error) {
	contextData, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
	)

	result, err := h.userUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to login", zap.Error(err))
		}
		return nil, err
	}

	log.Info("User login request processed successfully", zap.String("email", req.Email))
	return &authpbv1.UserLoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *AuthHandler) UserLogout(ctx context.Context, req *authpbv1.UserLogoutRequest) (*authpbv1.UserLogoutResponse, error) {
	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)
	err = h.userUsecase.Logout(ctx, req.AccessToken)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to logout", zap.Error(err))
		}
		return nil, err
	}

	log.Info("User logout request processed successfully")
	return &authpbv1.UserLogoutResponse{
		Success: &wrapperspb.BoolValue{Value: true},
	}, nil
}

func (h *AuthHandler) UserDelete(ctx context.Context, req *authpbv1.UserDeleteRequest) (*authpbv1.UserDeleteResponse, error) {
	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)

	err = h.userUsecase.UserAccountDelete(ctx, contextData.UserID, req.Password)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to delete account", zap.Error(err))
		}
		return nil, err
	}

	log.Info("User delete request processed successfully")
	return &authpbv1.UserDeleteResponse{
		Success: &wrapperspb.BoolValue{Value: true},
	}, nil
}

func (h *AuthHandler) AdminLogin(ctx context.Context, req *authpbv1.AdminLoginRequest) (*authpbv1.AdminLoginResponse, error) {
	contextData, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
	)

	result, err := h.adminUsecase.AdminLogin(ctx, req.Email, req.Password)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to login", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Admin login request processed successfully")
	return &authpbv1.AdminLoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *AuthHandler) AdminLogout(ctx context.Context, req *authpbv1.AdminLogoutRequest) (*authpbv1.AdminLogoutResponse, error) {
	contextData, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
	)
	err = h.adminUsecase.AdminLogout(ctx, req.AccessToken)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to logout", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Admin logout request processed successfully")
	return &authpbv1.AdminLogoutResponse{
		Success: &wrapperspb.BoolValue{Value: true},
	}, nil
}

func (h *AuthHandler) BlockOrUnblockUser(ctx context.Context, req *authpbv1.BlockOrUnblockUserRequest) (*authpbv1.BlockOrUnblockUserResponse, error) {
	contextData, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
	)

	err = h.adminUsecase.BlockOrUnblockUser(ctx, req.Field, req.Value, req.ShouldBlock)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to block or unblock user", zap.Error(err))
		}
		return nil, err
	}

	log.Info("User block or unblock request processed successfully",
		zap.String("field", req.Field),
		zap.String("value", req.Value),
		zap.Bool("should_block", req.ShouldBlock))
	return &authpbv1.BlockOrUnblockUserResponse{
		Success: &wrapperspb.BoolValue{Value: true},
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpbv1.RefreshTokenRequest) (*authpbv1.RefreshTokenResponse, error) {
	contextData, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
	)
	result, err := h.userUsecase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to refresh token", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Refresh token request processed successfully")
	return &authpbv1.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *AuthHandler) GetUsers(ctx context.Context, req *authpbv1.GetUsersRequest) (*authpbv1.GetUsersResponse, error) {
	contextData, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
	)

	users, err := h.adminUsecase.GetUsers(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to get users", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Users fetched successfully")
	response := make([]*authpbv1.GetUserResponse, len(users))
	for i, user := range users {
		var premiumUntilPb *timestamppb.Timestamp
		if user.PremiumUntil != nil {
			premiumUntilPb = timestamppb.New(*user.PremiumUntil)
		}

		var lastLoginAtPb *timestamppb.Timestamp
		if user.LastLoginAt != nil {
			lastLoginAtPb = timestamppb.New(*user.LastLoginAt)
		}

		response[i] = &authpbv1.GetUserResponse{
			Id:            user.ID.String(),
			Email:         user.Email,
			Phone:         user.Phone,
			EmailVerified: user.EmailVerified,
			PremiumUntil:  premiumUntilPb,
			LastLoginAt:   lastLoginAtPb,
			IsBlocked:     user.IsBlocked,
			CreatedAt:     timestamppb.New(user.CreatedAt),
			UpdatedAt:     timestamppb.New(user.UpdatedAt),
		}
	}

	return &authpbv1.GetUsersResponse{
		Users: response,
	}, nil
}

func (h *AuthHandler) GetUserByField(ctx context.Context, req *authpbv1.GetUserByFieldRequest) (*authpbv1.GetUserByFieldResponse, error) {
	contextData, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
	)

	user, err := h.adminUsecase.GetUserByField(ctx, req.Field, req.Value)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to get user by field", zap.Error(err))
		}
		return nil, err
	}

	var premiumUntilPb *timestamppb.Timestamp
	if user.PremiumUntil != nil {
		premiumUntilPb = timestamppb.New(*user.PremiumUntil)
	}

	var lastLoginAtPb *timestamppb.Timestamp
	if user.LastLoginAt != nil {
		lastLoginAtPb = timestamppb.New(*user.LastLoginAt)
	}

	log.Info("User fetched successfully")
	return &authpbv1.GetUserByFieldResponse{
		User: &authpbv1.GetUserResponse{
			Id:            user.ID.String(),
			Email:         user.Email,
			Phone:         user.Phone,
			EmailVerified: user.EmailVerified,
			PremiumUntil:  premiumUntilPb,
			LastLoginAt:   lastLoginAtPb,
			IsBlocked:     user.IsBlocked,
			CreatedAt:     timestamppb.New(user.CreatedAt),
			UpdatedAt:     timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
