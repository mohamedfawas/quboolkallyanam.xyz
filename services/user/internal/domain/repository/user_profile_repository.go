package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

type UserProfileRepository interface {
	CreateUserProfile(ctx context.Context,
		userProfile *entity.UserProfile) error
	ProfileExists(ctx context.Context,
		userID uuid.UUID) (bool, error)
	UpdateLastLogin(ctx context.Context,
		userID uuid.UUID) error
	GetProfileByUserID(ctx context.Context,
		userID uuid.UUID) (*entity.UserProfile, error)
	UpdateUserProfile(ctx context.Context,
		userProfile *entity.UserProfile) error
	GetUserProfileByID(ctx context.Context,
		id int64) (*entity.UserProfile, error)
	GetPotentialProfiles(ctx context.Context,
		userID uuid.UUID,
		excludedIDs []uuid.UUID,
		preferences *entity.PartnerPreference,
		limit int,
		offset int,
		isUserBride bool) ([]*entity.UserProfile, int64, error)
}
