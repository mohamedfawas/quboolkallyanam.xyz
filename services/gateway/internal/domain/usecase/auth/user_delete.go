package auth

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) UserDelete(
	ctx context.Context,
	req dto.UserDeleteRequest) error {

	if !validation.IsValidPassword(req.Password, validation.DefaultPasswordRequirements()) {
		return apperrors.ErrInvalidPassword
	}

	return u.authClient.UserDelete(ctx, req)
}
