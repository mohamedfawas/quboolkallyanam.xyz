package auth

import (
	"context"
	"slices"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) GetUserByField(
	ctx context.Context,
	req dto.GetUserByFieldRequest) (*dto.GetUserByFieldResponse, error) {

	allowedFields := []string{"email", "phone", "id"}
	if !slices.Contains(allowedFields, req.Field) {
		return nil, apperrors.ErrInvalidField
	}

	return u.authClient.GetUserByField(ctx, req)
}