package user

import (
	"context"
	"log"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetMatchRecommendations(ctx context.Context, req dto.GetMatchRecommendationsRequest) (*dto.GetMatchRecommendationsResponse, error) {
	log.Printf("Getting match recommendations for user with limit: %d, offset: %d", req.Limit, req.Offset)

	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 50 {
		req.Limit = 50
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	response, err := u.userClient.GetMatchRecommendations(ctx, req)
	if err != nil {
		log.Printf("Failed to get match recommendations from user service: %v", err)
		return nil, err
	}

	log.Printf("Successfully retrieved %d match recommendations", len(response.Profiles))
	return response, nil
}
