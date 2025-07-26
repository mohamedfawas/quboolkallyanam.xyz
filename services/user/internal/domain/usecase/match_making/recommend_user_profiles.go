package matchmaking

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *matchMakingUsecase) RecommendUserProfiles(
	ctx context.Context,
	userID uuid.UUID,
	limit, offset int) ([]*entity.UserProfileResponse, *entity.PaginationData, error) {

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
		log.Printf("failed to get user profile: %v", err)
		return nil, nil, err
	}

	if userProfile == nil {
		return nil, nil, appError.ErrUserNotFound
	}

	partnerPreference, err := u.partnerPreferencesRepository.GetPartnerPreferencesByUserProfileID(ctx, uint(userProfile.ID))
	if err != nil {
		log.Printf("failed to get partner preference: %v", err)
		return nil, nil, err
	}

	if partnerPreference == nil {
		return nil, nil, appError.ErrPartnerPreferencesNotFound
	}

	excludedIDs, err := u.profileMatchRepository.GetMatchedProfileIDs(ctx, userID)
	if err != nil {
		log.Printf("failed to get matched profile IDs: %v", err)
		return nil, nil, err
	}
	excludedIDs = append(excludedIDs, userID)

	potentialProfiles, err := u.userProfileRepository.GetPotentialProfiles(ctx, 
		userID, 
		excludedIDs, 
		partnerPreference, 
		userProfile.IsBride)
	if err != nil {
		log.Printf("failed to get potential profiles: %v", err)
		return nil, nil, err
	}

	totalCount := len(potentialProfiles)
	end := offset + limit
	if end > totalCount {
		end = totalCount
	}

	var result []*entity.UserProfile
	if offset < totalCount {
		result = potentialProfiles[offset:end]
	} else {
		result = []*entity.UserProfile{}
	}

	if len(result) == 0 {
		return []*entity.UserProfileResponse{}, &entity.PaginationData{
			TotalCount: 0,
			Limit:      limit,
			Offset:     offset,
			HasMore:    false,
		}, nil
	}

	recommendedProfiles := make([]*entity.UserProfileResponse, len(result))
	for i, profile := range result {
		var age int
		if profile.DateOfBirth != nil {
			age = calculateAge(profile.DateOfBirth)
		}
		recommendedProfiles[i] = &entity.UserProfileResponse{
			ID:                profile.ID,
			FullName:          *profile.FullName,
			ProfilePictureURL: profile.ProfilePictureURL,
			Age:               age,
			HeightCm:          *profile.HeightCm,
			MaritalStatus:     string(*profile.MaritalStatus),
			Profession:        string(*profile.Profession),
			HomeDistrict:      string(*profile.HomeDistrict),
		}
	}

	return recommendedProfiles, &entity.PaginationData{
		TotalCount: int64(totalCount),
		Limit:      limit,
		Offset:     offset,
		HasMore:    end < totalCount,
	}, nil
}

func calculateAge(dateOfBirth *time.Time) int {
	if dateOfBirth == nil {
		return 0
	}
	now := time.Now().UTC()
	age := now.Year() - dateOfBirth.Year()
	if now.Month() < dateOfBirth.Month() || (now.Month() == dateOfBirth.Month() && now.Day() < dateOfBirth.Day()) {
		age--
	}
	return age
}
