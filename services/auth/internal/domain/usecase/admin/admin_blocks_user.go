package admin

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
)


func (u *adminUsecase) BlockUser(ctx context.Context, field string, value string) error {
	user, err := u.userRepository.GetUser(ctx, field, value)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return apperrors.ErrUserNotFound
	}

	if user.IsBlocked {
		return apperrors.ErrUserAlreadyBlocked
	}
	
	err = u.userRepository.BlockUser(ctx, field, value)
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}

	return nil
}