package userprofile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *userProfileUsecase) GetUserPartnerPreferences(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.PartnerPreference, error) {

	profile, err := u.userProfileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	if profile == nil {
		return nil, apperrors.ErrUserProfileNotFound
	}

	if !profile.ProfileCompleted {
		return nil, apperrors.ErrUserProfileNotCompleted
	}

	preferences, err := u.partnerPreferencesRepository.GetPartnerPreferencesByUserProfileID(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner preferences: %w", err)
	}
	if preferences == nil {
		return nil, apperrors.ErrPartnerPreferencesNotFound
	}

	return preferences, nil
}