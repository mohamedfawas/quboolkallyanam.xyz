package usecase

import (
	"context"
)

type NotificationUsecase interface {
	HandleOTPVerification(ctx context.Context, userEmail, otp string, expiryMinutes int) error
	HandleUserAccountDeletion(ctx context.Context, userEmail string) error
}
