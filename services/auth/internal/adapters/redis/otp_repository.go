package redis

import (
	"context"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/redis"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
)

type otpRepository struct {
	redisClient *redis.Client
}

func NewOTPRepository(redisClient *redis.Client) repository.OTPRepository {
	return &otpRepository{redisClient: redisClient}
}

func (r *otpRepository) GetOTP(ctx context.Context, key string) (string, error) {
	return r.redisClient.Get(ctx, key)
}

func (r *otpRepository) StoreOTP(ctx context.Context, key string, otp string, expiry time.Duration) error {
	return r.redisClient.Set(ctx, key, otp, expiry)
}

func (r *otpRepository) DeleteOTP(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key)
}
