package repository

import (
	"context"
	"time"
)

type OTPRepository interface {
	GetOTP(ctx context.Context, key string) (string, error)
	StoreOTP(ctx context.Context, key string, otp string, expiry time.Duration) error
	DeleteOTP(ctx context.Context, key string) error
}
