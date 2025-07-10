package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/entity"
)

type PendingRegistrationUsecase interface {
	RegisterUser(ctx context.Context, req *entity.UserRegistrationRequest, config *config.Config) error
	VerifyUserRegistration(ctx context.Context, email, otp string) error
}
