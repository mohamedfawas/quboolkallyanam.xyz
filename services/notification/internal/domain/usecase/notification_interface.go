package usecase

import (
	"context"
)

type NotificationUsecase interface {
	HandleOTPVerification(ctx context.Context, userEmail, otp string, expiryMinutes int) error
	HandleUserAccountDeletion(ctx context.Context, userEmail string) error
	HandleAdminBlockedUser(ctx context.Context, userEmail string, shouldBlock bool) error
}
