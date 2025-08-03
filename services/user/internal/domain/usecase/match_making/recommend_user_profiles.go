package matchmaking

import (
	"context"

	"github.com/google/uuid"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/ageutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pagination"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *matchMakingUsecase) RecommendUserProfiles(
	ctx context.Context,
	userID uuid.UUID,
	limit, offset int) ([]*entity.UserProfileResponse, *pagination.PaginationData, error) {

	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	userProfile, err := u.userProfileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	if userProfile == nil {
		return nil, nil, appError.ErrUserNotFound
	}

	partnerPreference, err := u.partnerPreferencesRepository.GetPartnerPreferencesByUserProfileID(ctx, userProfile.ID)
	if err != nil {
		return nil, nil, err
	}

	if partnerPreference == nil {
		return nil, nil, appError.ErrPartnerPreferencesNotFound
	}

	excludedIDs, err := u.profileMatchRepository.GetMatchedProfileIDs(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	excludedIDs = append(excludedIDs, userID)

	potentialProfiles, totalCount, err := u.userProfileRepository.GetPotentialProfiles(ctx,
		userID,
		excludedIDs,
		partnerPreference,
		limit,
		offset,
		userProfile.IsBride)
	if err != nil {
		return nil, nil, err
	}

	result := potentialProfiles

	if len(result) == 0 {
		return []*entity.UserProfileResponse{}, &pagination.PaginationData{
			TotalCount: totalCount, 
			Limit:      limit,
			Offset:     offset,
			HasMore:    false,
		}, nil
	}

	recommendedProfiles := make([]*entity.UserProfileResponse, len(result))
	for i, profile := range result {
		age := ageutil.CalculateAge(profile.DateOfBirth)
		recommendedProfiles[i] = &entity.UserProfileResponse{
			ID:                profile.ID,
			FullName:          profile.FullName,
			ProfilePictureURL: &profile.ProfilePictureURL,
			Age:               uint32(age),
			HeightCm:          uint32(profile.HeightCm),
			MaritalStatus:     string(profile.MaritalStatus),
			Profession:        string(profile.Profession),
			HomeDistrict:      string(profile.HomeDistrict),
		}
	}

	return recommendedProfiles, &pagination.PaginationData{
		TotalCount: totalCount, 
		Limit:      limit,
		Offset:     offset,
		HasMore:    int64(offset+limit) < totalCount, 
	}, nil
}
