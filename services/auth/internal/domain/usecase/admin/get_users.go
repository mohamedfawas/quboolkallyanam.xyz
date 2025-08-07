package admin

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *adminUsecase) GetUsers(ctx context.Context, page, limit int) ([]*entity.GetUserResponse, error) {
	users, err := u.userRepository.GetUsers(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	response := make([]*entity.GetUserResponse, len(users))
	for i, user := range users {
		response[i] = &entity.GetUserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Phone:         user.Phone,
			EmailVerified: user.EmailVerified,
			PremiumUntil:  user.PremiumUntil,
			LastLoginAt:   user.LastLoginAt,
			IsBlocked:     user.IsBlocked,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		}
	}

	return response, nil
}