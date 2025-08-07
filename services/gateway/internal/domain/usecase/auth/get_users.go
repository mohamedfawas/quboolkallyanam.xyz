package auth

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *authUsecase) GetUsers(
	ctx context.Context,
	req dto.GetUsersRequest) (*dto.GetUsersResponse, error) {

	if req.Page < 1 {
		return nil, apperrors.ErrInvalidPaginationPage
	}
	if req.Limit < 1 || req.Limit > 100 {
		return nil, apperrors.ErrInvalidPaginationLimit
	}

	return u.authClient.GetUsers(ctx, req)
}