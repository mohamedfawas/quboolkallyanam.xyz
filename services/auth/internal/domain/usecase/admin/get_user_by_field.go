package admin

import (
	"context"
	"fmt"
	"slices"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

func (u *adminUsecase) GetUserByField(ctx context.Context, field, value string) (*entity.GetUserResponse, error) {

	allowedFields := []string{"email", "phone", "id"}
	if !slices.Contains(allowedFields, field) {
		return nil, apperrors.ErrInvalidField
	}

	user, err := u.userRepository.GetUser(ctx, field, value)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by field: %w", err)
	}

	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	return &entity.GetUserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Phone:         user.Phone,
		EmailVerified: user.EmailVerified,
		PremiumUntil:  user.PremiumUntil,
		LastLoginAt:   user.LastLoginAt,
		IsBlocked:     user.IsBlocked,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}
