package redis

import (
	"context"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/redis"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
)

type tokenRepository struct {
	redisClient *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) repository.TokenRepository {
	return &tokenRepository{redisClient: redisClient}
}

func (r *tokenRepository) StoreRefreshToken(ctx context.Context, key string, token string, expiry time.Duration) error {
	return r.redisClient.Set(ctx, key, token, expiry)
}

func (r *tokenRepository) IsValidRefreshToken(ctx context.Context, key string) (bool, error) {
	return r.redisClient.Exists(ctx, key)
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key)
}

func (r *tokenRepository) BlacklistToken(ctx context.Context, key string, expiry time.Duration) error {
	return r.redisClient.Set(ctx, key, constants.BlacklistedToken, expiry)
}

func (r *tokenRepository) IsTokenBlacklisted(ctx context.Context, key string) (bool, error) {
	return r.redisClient.Exists(ctx, key)
}
