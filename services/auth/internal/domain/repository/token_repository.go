package repository

import (
	"context"
	"time"
)

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, key string, token string, expiry time.Duration) error
	IsValidRefreshToken(ctx context.Context, key string) (bool, error)
	DeleteRefreshToken(ctx context.Context, key string) error
	BlacklistToken(ctx context.Context, key string, expiry time.Duration) error
	IsTokenBlacklisted(ctx context.Context, key string) (bool, error)
}
