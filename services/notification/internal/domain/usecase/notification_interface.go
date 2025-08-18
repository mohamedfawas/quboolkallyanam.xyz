package usecase

import (
	"context"
)

type NotificationUsecase interface {
	HandleOTPVerification(ctx context.Context, userEmail, otp string, expiryMinutes int) error
	HandleUserAccountDeletion(ctx context.Context, userEmail string) error
	HandleAdminBlockedUser(ctx context.Context, userEmail string, shouldBlock bool) error
	HandleUserInterestSent(ctx context.Context, receiverEmail string, senderProfileID int64, senderName string) error
	HandleMutualMatchCreated(ctx context.Context, 
		user1Email string, user1ProfileID int64, user1FullName string, 
		user2Email string, user2ProfileID int64, user2FullName string) error
}
