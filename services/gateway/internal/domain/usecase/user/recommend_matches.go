package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetMatchRecommendations(ctx context.Context, req dto.GetMatchRecommendationsRequest) (*dto.GetMatchRecommendationsResponse, error) {

	if req.Limit <= 0 {
		req.Limit = constants.DefaultPaginationLimit
	}
	if req.Limit > constants.MaxPaginationLimit {
		req.Limit = constants.MaxPaginationLimit
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	response, err := u.userClient.GetMatchRecommendations(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
