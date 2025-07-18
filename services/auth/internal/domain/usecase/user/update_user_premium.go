package user

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (u *userUseCase) UpdateUserPremium(ctx context.Context, userID string, premiumUntil time.Time) error {
	now := time.Now().UTC()
	err := u.userRepository.UpdatePremiumUntil(ctx, userID, premiumUntil, now)
	if err != nil {
		log.Printf("failed to update user premium value : %v", err)
		return fmt.Errorf("failed to update user premium value : %w", err)
	}
	return nil
}
