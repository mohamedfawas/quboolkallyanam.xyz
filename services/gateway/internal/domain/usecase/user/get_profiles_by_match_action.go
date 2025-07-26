package user

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetProfilesByMatchAction(ctx context.Context, req dto.GetProfilesByMatchActionRequest) (*dto.GetProfilesByMatchActionResponse, error) {
	log.Printf("Getting profiles by match action: %s with limit: %d, offset: %d", req.Action, req.Limit, req.Offset)

	// Validate action
	if req.Action == "" {
		return nil, appError.ErrInvalidMatchAction
	}

	if !validation.IsValidMatchMakingOption(req.Action) {
		return nil, appError.ErrInvalidMatchAction
	}

	// Set default values for pagination
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 50 {
		req.Limit = 50
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	response, err := u.userClient.GetProfilesByMatchAction(ctx, req)
	if err != nil {
		log.Printf("Failed to get profiles by match action from user service: %v", err)
		return nil, err
	}

	log.Printf("Successfully retrieved %d profiles for action: %s", len(response.Profiles), req.Action)
	return response, nil
}