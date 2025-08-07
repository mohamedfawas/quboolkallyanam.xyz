package pendingregistration

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/event"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/auth/internal/domain/usecase"
)

type pendingRegistrationUsecase struct {
	pendingRegistrationRepository repository.PendingRegistrationRepository
	userRepository                repository.UserRepository
	otpRepository                 repository.OTPRepository
	eventPublisher                event.EventPublisher
}

func NewPendingRegistrationUsecase(
	pendingRegistrationRepository repository.PendingRegistrationRepository,
	userRepository repository.UserRepository,
	otpRepository repository.OTPRepository,
	eventPublisher event.EventPublisher,
) usecase.PendingRegistrationUsecase {
	return &pendingRegistrationUsecase{
		pendingRegistrationRepository: pendingRegistrationRepository,
		userRepository:                userRepository,
		otpRepository:                 otpRepository,
		eventPublisher:                eventPublisher,
	}
}
