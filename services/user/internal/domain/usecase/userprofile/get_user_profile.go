package userprofile

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/ageutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *userProfileUsecase) GetUserProfile(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.UserProfileResponse, error) {

	profile, err := u.userProfileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, apperrors.ErrUserProfileNotFound
	}

	age := ageutil.CalculateAge(profile.DateOfBirth)

	// Generate picture URL only if an image exists
	var pictureURL *string
	if profile.ProfileImageKey != "" {
		url, err := u.photoStorage.GetDownloadURL(ctx, profile.ProfileImageKey, u.config.MediaStorage.URLExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to get profile image URL for user %s: %w", userID, err)
		}
		pictureURL = &url
	}

	return &entity.UserProfileResponse{
		ID:                profile.ID,
		FullName:          profile.FullName,
		ProfilePictureURL: pictureURL,
		Age:               int32(age),
		HeightCm:          int32(profile.HeightCm),
		MaritalStatus:     string(profile.MaritalStatus),
		Profession:        string(profile.Profession),
		HomeDistrict:      string(profile.HomeDistrict),
	}, nil
}