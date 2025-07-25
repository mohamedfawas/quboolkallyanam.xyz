package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

type UserProfileRepository interface {
	CreateUserProfile(ctx context.Context,
		userProfile *entity.UserProfile) error
	ProfileExists(ctx context.Context,
		userID uuid.UUID) (bool, error)
	UpdateLastLogin(ctx context.Context,
		userID uuid.UUID,
		lastLogin time.Time) error
	GetProfileByUserID(ctx context.Context,
		userID uuid.UUID) (*entity.UserProfile, error)
	UpdateUserProfile(ctx context.Context,
		userProfile *entity.UserProfile) error
	GetUserProfileByID(ctx context.Context,
		id uint) (*entity.UserProfile, error)
	GetPotentialProfiles(ctx context.Context,
		userID uuid.UUID,
		excludedIDs []uuid.UUID,
		preferences *entity.PartnerPreference,
		isUserBride bool) ([]*entity.UserProfile, error)
}
