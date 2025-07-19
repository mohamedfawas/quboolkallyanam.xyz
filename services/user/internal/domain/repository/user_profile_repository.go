package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

type UserProfileRepository interface {
	CreateUserProfile(ctx context.Context, userProfile *entity.UserProfile) error
	// GetUserProfileByID(ctx context.Context, id uint) (*entity.UserProfile, error)
	// GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	ProfileExists(ctx context.Context, userID uuid.UUID) (bool, error)
	// PatchUserProfile(ctx context.Context, userID uuid.UUID, patch map[string]interface{}) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLogin time.Time) error
	// UpdateProfilePicture(ctx context.Context, userID uuid.UUID, profilePictureURL string) error
	// RemoveProfilePicture(ctx context.Context, userID uuid.UUID) error
	// SoftDeleteUserProfile(ctx context.Context, userID uuid.UUID) error
}
