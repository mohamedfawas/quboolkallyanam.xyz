package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

type UserProfileRepository interface {
	CreateUserProfile(ctx context.Context, userProfile *entity.UserProfile) error
	ProfileExists(ctx context.Context, userID uuid.UUID) (bool, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLogin time.Time) error
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	UpdateUserProfile(ctx context.Context, userProfile *entity.UserProfile) error
	GetUserProfileByID(ctx context.Context, id uint) (*entity.UserProfile, error)
	/*
		GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
		GetProfileByID(ctx context.Context, id uint) (*models.UserProfile, error)
		GetBasicProfileByUUID(ctx context.Context, userUUID uuid.UUID) (*models.UserProfile, error)
		GetUserUUIDByProfileID(ctx context.Context, profileID uint64) (uuid.UUID, error)
		ProfileExists(ctx context.Context, userID uuid.UUID) (bool, error)
		UpdateProfile(ctx context.Context, profile *models.UserProfile) error
		PatchProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error
		UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLogin time.Time) error
		UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error
		UpdateProfilePhoto(ctx context.Context, userID uuid.UUID, photoURL string) error
		RemoveProfilePhoto(ctx context.Context, userID uuid.UUID) error
		SoftDeleteUserProfile(ctx context.Context, userID uuid.UUID) error
	*/
}
