package auth

import (
	"context"

	errors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	logger "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/logger"

	validation "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/vaidation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) UserRegister(ctx context.Context, req dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	logger.Log.Info("🔑 UserRegister request received in usecase : ", "email : ", req.Email, "phone : ", req.Phone)
	if !validation.IsValidEmail(req.Email) {
		return nil, errors.ErrInvalidEmail
	}

	if !validation.IsValidPhoneNumber(req.Phone) {
		return nil, errors.ErrInvalidPhoneNumber
	}

	if !validation.IsValidPassword(req.Password, validation.DefaultPasswordRequirements()) {
		return nil, errors.ErrInvalidPassword
	}
	logger.Log.Info("🔑 UserRegister request validated in usecase : ", "email : ", req.Email, "phone : ", req.Phone)

	logger.Log.Info("🔑 UserRegister request sent to client", "email", req.Email, "phone", req.Phone)
	return u.authClient.UserRegister(ctx, req)
}
