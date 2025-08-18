package userprofile

import (
	"context"
	"fmt"
	"log"
	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/ageutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *userProfileUsecase) GetUserDetailsByProfileID(
	ctx context.Context,
	requesterUserID uuid.UUID,
	targetProfileID int64,
	requestedByAdmin bool,
) (*entity.UserProfileResponse, *entity.PartnerPreference, []string, error) {

	targetProfile, err := u.userProfileRepository.GetUserProfileByID(ctx, targetProfileID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load target profile: %w", err)
	}
	if targetProfile == nil {
		return nil, nil, nil, apperrors.ErrInvalidTargetProfileID
	}

	if !requestedByAdmin {
		requesterProfile, err := u.userProfileRepository.GetProfileByUserID(ctx, requesterUserID)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to load requester profile: %w", err)
		}
		if requesterProfile == nil {
			return nil, nil, nil, apperrors.ErrUserProfileNotFound
		}
		if requesterProfile.IsBride == targetProfile.IsBride {
			return nil, nil, nil, apperrors.ErrSameGenderProfileAccessDenied
		}
	}

	age := ageutil.CalculateAge(targetProfile.DateOfBirth)
	var pictureURL *string
	if targetProfile.ProfileImageKey != "" {
		url, err := u.photoStorage.GetDownloadURL(ctx, targetProfile.ProfileImageKey, u.config.MediaStorage.URLExpiry)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to get profile image URL: %w", err)
		}
		pictureURL = &url
	}
	profileResp := &entity.UserProfileResponse{
		ID:                targetProfile.ID,
		FullName:          targetProfile.FullName,
		ProfilePictureURL: pictureURL,
		Age:               int32(age),
		HeightCm:          int32(targetProfile.HeightCm),
		MaritalStatus:     string(targetProfile.MaritalStatus),
		Profession:        string(targetProfile.Profession),
		HomeDistrict:      string(targetProfile.HomeDistrict),
	}

	prefs, err := u.partnerPreferencesRepository.GetPartnerPreferencesByUserProfileID(ctx, targetProfile.ID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get partner preferences: %w", err)
	}
	log.Println("prefs", prefs)
	// Note: prefs may be nil; this is allowed and not an error.

	photos, err := u.GetAdditionalPhotos(ctx, targetProfile.UserID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get additional photos: %w", err)
	}
	// Note: photos may be empty; this is allowed.

	return profileResp, prefs, photos, nil
}
