package user

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/timeutil"
)

func (u *userUseCase) UpdateUserPremium(ctx context.Context, userID string, premiumUntil time.Time) error {
	now := timeutil.NowIST() // for updated_at column
	err := u.userRepository.UpdatePremiumUntil(ctx, userID, premiumUntil, now)
	if err != nil {
		return fmt.Errorf("failed to update user premium value : %w", err)
	}
	return nil
}
