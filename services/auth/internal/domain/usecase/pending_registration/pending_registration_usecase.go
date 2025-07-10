package pendingregistration

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/notifications/smtp"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"
)

type pendingRegistrationUsecase struct {
	pendingRegistrationRepository repository.PendingRegistrationRepository
	userRepository                repository.UserRepository
	otpRepository                 repository.OTPRepository
	smtpClient                    smtp.Client
}

func NewPendingRegistrationUsecase(
	pendingRegistrationRepository repository.PendingRegistrationRepository,
	userRepository repository.UserRepository,
	otpRepository repository.OTPRepository,
	smtpClient smtp.Client,
) usecase.PendingRegistrationUsecase {
	return &pendingRegistrationUsecase{
		pendingRegistrationRepository: pendingRegistrationRepository,
		userRepository:                userRepository,
		otpRepository:                 otpRepository,
		smtpClient:                    smtpClient,
	}
}
