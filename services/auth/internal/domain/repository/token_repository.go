package repository

import (
	"context"
	"time"
)

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID string, token string, expiry time.Duration) error
	ValidateRefreshToken(ctx context.Context, userID string, token string) (bool, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
	BlacklistToken(ctx context.Context, tokenID string, expiry time.Duration) error
	IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)
}
