package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

type UserProfileUsecase interface {
	UpdateUserLastLogin(ctx context.Context,
		userID uuid.UUID,
		email, phone string) error
	UpdateUserProfile(ctx context.Context,
		userID uuid.UUID,
		req entity.UpdateUserProfileRequest) error
	UpdateUserPartnerPreferences(ctx context.Context,
		userID uuid.UUID,
		operationType string,
		req entity.UpdateUserPartnerPreferencesRequest) error
	// UpdateUserDeletion
	// PatchUserProfile
	// GetUserProfile
	// PostPartnerPreference
	// PatchPartnerPreference
	// DeletePartnerPreference
	// GetPartnerPreference
	// PostUserProfilePhoto
	// DeleteUserProfilePhoto
	// PostUserAdditionalPhotos
	// DeleteUserAdditionalPhotos
	// GetUserAdditionalPhotos
	// PostUserVideo
	// DeleteUserVideo
}
