package repository

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

type PartnerPreferencesRepository interface {
	CreatePartnerPreferences(ctx context.Context, preferences *entity.PartnerPreference) error
	GetPartnerPreferencesByUserProfileID(ctx context.Context, userProfileID int64) (*entity.PartnerPreference, error)
	PatchPartnerPreferences(ctx context.Context, userProfileID int64, patch map[string]interface{}) error
}
