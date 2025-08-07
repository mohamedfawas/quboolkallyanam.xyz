package user

import (
	"context"
	"fmt"
	"time"
)

func (u *userUseCase) UpdateUserPremium(ctx context.Context, userID string, premiumUntil time.Time) error {
	now := time.Now().UTC()
	err := u.userRepository.UpdatePremiumUntil(ctx, userID, premiumUntil, now)
	if err != nil {
		return fmt.Errorf("failed to update user premium value : %w", err)
	}
	return nil
}
