package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) GetUserDetailsByProfileID(ctx context.Context, req dto.GetUserDetailsByProfileIDRequest) (*dto.GetUserDetailsByProfileIDResponse, error) {
	return u.userClient.GetUserDetailsByProfileID(ctx, req)
}