package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type UserUsecase interface {
	///////// USER PROFILE MANAGEMENT //////////
	UpdateUserProfile(ctx context.Context, req dto.UserProfilePatchRequest) error
	GetUserProfile(ctx context.Context) (*dto.UserProfileRecommendation, error)
	GetProfilePhotoUploadURL(ctx context.Context, req dto.GetProfilePhotoUploadURLRequest) (*dto.GetProfilePhotoUploadURLResponse, error)
	ConfirmProfilePhotoUpload(ctx context.Context, req dto.ConfirmProfilePhotoUploadRequest) (*dto.ConfirmProfilePhotoUploadResponse, error)
	DeleteProfilePhoto(ctx context.Context, req dto.DeleteProfilePhotoRequest) error
	
	///////// USER ADDITIONAL PHOTO MANAGEMENT //////////
	GetAdditionalPhotoUploadURL(ctx context.Context, req dto.GetAdditionalPhotoUploadURLRequest) (*dto.GetAdditionalPhotoUploadURLResponse, error)
	ConfirmAdditionalPhotoUpload(ctx context.Context, req dto.ConfirmAdditionalPhotoUploadRequest) (*dto.ConfirmAdditionalPhotoUploadResponse, error)
	DeleteAdditionalPhoto(ctx context.Context, req dto.DeleteAdditionalPhotoRequest) (*dto.DeleteAdditionalPhotoResponse, error)
	
	///////// PARTNER PREFERENCES MANAGEMENT //////////
	UpdateUserPartnerPreferences(ctx context.Context, operationType string, req dto.PartnerPreferencePatchRequest) error
	

	///////// MATCH MAKING //////////
	RecordMatchAction(ctx context.Context, req dto.RecordMatchActionRequest) (*dto.RecordMatchActionResponse, error)
	GetMatchRecommendations(ctx context.Context, req dto.GetMatchRecommendationsRequest) (*dto.GetMatchRecommendationsResponse, error)
	GetProfilesByMatchAction(ctx context.Context, req dto.GetProfilesByMatchActionRequest) (*dto.GetProfilesByMatchActionResponse, error)
}
