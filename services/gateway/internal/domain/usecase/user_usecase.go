package usecase

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type UserUsecase interface {
	///////// USER PROFILE MANAGEMENT //////////
	UpdateUserProfile(ctx context.Context, req dto.UserProfilePatchRequest) error
	GetProfilePhotoUploadURL(ctx context.Context, req dto.GetProfilePhotoUploadURLRequest) (*dto.GetProfilePhotoUploadURLResponse, error)
	ConfirmProfilePhotoUpload(ctx context.Context, req dto.ConfirmProfilePhotoUploadRequest) (*dto.ConfirmProfilePhotoUploadResponse, error)
	
	
	///////// PARTNER PREFERENCES MANAGEMENT //////////
	UpdateUserPartnerPreferences(ctx context.Context, operationType string, req dto.PartnerPreferencePatchRequest) error
	

	///////// MATCH MAKING //////////
	RecordMatchAction(ctx context.Context, req dto.RecordMatchActionRequest) (*dto.RecordMatchActionResponse, error)
	GetMatchRecommendations(ctx context.Context, req dto.GetMatchRecommendationsRequest) (*dto.GetMatchRecommendationsResponse, error)
	GetProfilesByMatchAction(ctx context.Context, req dto.GetProfilesByMatchActionRequest) (*dto.GetProfilesByMatchActionResponse, error)
}
