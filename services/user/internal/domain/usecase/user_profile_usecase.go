package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	mediastorage "github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/mediastorage"
)

type UserProfileUsecase interface {
	// USER PROFILE MANAGEMENT
	UpdateUserLastLogin(ctx context.Context,
		userID uuid.UUID,
		email, phone string) error
	HandleUserDeletion(ctx context.Context,
		userID uuid.UUID) error
	UpdateUserProfile(ctx context.Context,
		userID uuid.UUID,
		req entity.UpdateUserProfileRequest) error
	GetUserProfile(ctx context.Context,
		 userID uuid.UUID) (*entity.UserProfileResponse, error)
	GetProfilePhotoUploadURL(ctx context.Context,
		userID uuid.UUID,
		contentType string) (*mediastorage.PhotoUploadURLResponse, error)
	ConfirmProfilePhotoUpload(ctx context.Context,
		userID uuid.UUID,
		objectKey string) (string, error)
	DeleteProfilePhoto(ctx context.Context,
		userID uuid.UUID) error
	// USER PARTNER PREFERENCES MANAGEMENT
	UpdateUserPartnerPreferences(ctx context.Context,
		userID uuid.UUID,
		operationType string,
		req entity.UpdateUserPartnerPreferencesRequest) error
	
	// ADDITIONAL PHOTOS MANAGEMENT
	GetAdditionalPhotoUploadURL(ctx context.Context,
		userID uuid.UUID,
		displayOrder int32,
		contentType string) (*mediastorage.PhotoUploadURLResponse, error)
	ConfirmAdditionalPhotoUpload(ctx context.Context,
		userID uuid.UUID,
		objectKey string) (string, error)
	DeleteAdditionalPhoto(ctx context.Context,
		userID uuid.UUID,
		displayOrder int32) error
	GetAdditionalPhotos(ctx context.Context, 
		userID uuid.UUID) ([]string, error)
}
