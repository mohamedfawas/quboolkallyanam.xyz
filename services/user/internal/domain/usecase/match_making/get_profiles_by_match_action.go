package matchmaking

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/ageutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pagination"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *matchMakingUsecase) GetProfilesByMatchAction(ctx context.Context,
	userID uuid.UUID,
	action string,
	limit, offset int) ([]*entity.UserProfileResponse, *pagination.PaginationData, error) {

	if limit <= 0 {
		limit = constants.DefaultPaginationLimit
	}
	if limit > constants.MaxPaginationLimit {
		limit = constants.MaxPaginationLimit
	}
	if offset < 0 {
		offset = 0
	}

	var resultUUIDS []uuid.UUID
	var err error

	switch action {
	case constants.MatchMakingOptionMutual:
		resultUUIDS, err = u.mutualMatchRepository.GetMutualMatchedUserIDs(ctx, userID)
	case constants.MatchMakingOptionLiked:
		resultUUIDS, err = u.profileMatchRepository.GetLikedUserIDs(ctx, userID)
	case constants.MatchMakingOptionPassed:
		resultUUIDS, err = u.profileMatchRepository.GetPassedUserIDs(ctx, userID)
	default:
		return nil, nil, fmt.Errorf("invalid match action: %s", action)
	}

	if err != nil {
		return nil, nil, err
	}

	if len(resultUUIDS) == 0 {
		return []*entity.UserProfileResponse{}, &pagination.PaginationData{
			TotalCount: 0,
			Limit:      limit,
			Offset:     offset,
			HasMore:    false,
		}, nil
	}

	// Apply pagination to UUIDs first
	totalCount := len(resultUUIDS)
	end := offset + limit
	if end > totalCount {
		end = totalCount
	}

	var paginatedUUIDs []uuid.UUID
	if offset < totalCount {
		paginatedUUIDs = resultUUIDS[offset:end]
	} else {
		paginatedUUIDs = []uuid.UUID{}
	}

	if len(paginatedUUIDs) == 0 {
		return []*entity.UserProfileResponse{}, &pagination.PaginationData{
			TotalCount: int64(totalCount),
			Limit:      limit,
			Offset:     offset,
			HasMore:    false,
		}, nil
	}

	var profiles []*entity.UserProfile
	for _, userUUID := range paginatedUUIDs {
		profile, err := u.userProfileRepository.GetProfileByUserID(ctx, userUUID)
		if err != nil {
			continue // Skip missing profiles instead of failing entirely
		}
		if profile != nil {
			profiles = append(profiles, profile)
		}
	}

	userProfileResponses := make([]*entity.UserProfileResponse, len(profiles))
	for i, profile := range profiles {
		age := ageutil.CalculateAge(profile.DateOfBirth)
		profilePictureURL, err := u.photoStorage.GetDownloadURL(ctx, profile.ProfileImageKey, u.config.MediaStorage.URLExpiry)
		if err != nil {
			return nil, nil, err
		}
		userProfileResponses[i] = &entity.UserProfileResponse{
			ID:                profile.ID,
			FullName:          profile.FullName,
			ProfilePictureURL: &profilePictureURL,
			Age:               int32(age),
			HeightCm:          int32(profile.HeightCm),
			MaritalStatus:     string(profile.MaritalStatus),
			Profession:        string(profile.Profession),
			HomeDistrict:      string(profile.HomeDistrict),
		}
	}

	return userProfileResponses, &pagination.PaginationData{
		TotalCount: int64(totalCount),
		Limit:      limit,
		Offset:     offset,
		HasMore:    end < totalCount,
	}, nil
}
