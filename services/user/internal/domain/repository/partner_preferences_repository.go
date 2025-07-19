package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

type PartnerPreferencesRepository interface {
	GetPartnerPreferencesByUserProfileID(ctx context.Context, userProfileID uint) (*entity.PartnerPreference, error)
	PatchPartnerPreferences(ctx context.Context, userProfileID uint, patch map[string]interface{}) error
}
