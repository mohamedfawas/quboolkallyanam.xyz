package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) RecordMatchAction(
	ctx context.Context, 
	req dto.RecordMatchActionRequest) (*dto.RecordMatchActionResponse, error) {
	
	if req.Action == "" {
		return nil, appError.ErrInvalidMatchAction
	}

	if req.Action != string(constants.MatchActionLike) && req.Action != string(constants.MatchActionPass) {
		return nil, appError.ErrInvalidMatchAction
	}

	if req.TargetProfileID == 0 {
		return nil, appError.ErrInvalidTargetProfileID
	}

	return u.userClient.RecordMatchAction(ctx, req)
}