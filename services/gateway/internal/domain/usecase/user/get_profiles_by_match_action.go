package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetProfilesByMatchAction(ctx context.Context, req dto.GetProfilesByMatchActionRequest) (*dto.GetProfilesByMatchActionResponse, error) {
	if req.Action == "" || !validation.IsValidMatchMakingOption(req.Action) {
		return nil, apperrors.ErrInvalidMatchAction
	}

	if req.Limit <= 0 {
		req.Limit = constants.DefaultPaginationLimit
	}
	if req.Limit > constants.MaxPaginationLimit {
		req.Limit = constants.MaxPaginationLimit
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	resp, err := u.userClient.GetProfilesByMatchAction(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}