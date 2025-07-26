package matchmaking

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *matchMakingUsecase) GetProfilesByMatchAction(ctx context.Context,
	userID uuid.UUID,
	action string,
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
		log.Printf("invalid match action: %s", action)
		return nil, nil, fmt.Errorf("invalid match action: %s", action)
	}

	if err != nil {
		log.Printf("failed to get profiles by match action: %v", err)
		return nil, nil, err
	}

	if len(resultUUIDS) == 0 {
		return []*entity.UserProfileResponse{}, &entity.PaginationData{
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
		return []*entity.UserProfileResponse{}, &entity.PaginationData{
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
			log.Printf("failed to get profile for user %s: %v", userUUID, err)
			continue // Skip missing profiles instead of failing entirely
		}
		if profile != nil {
			profiles = append(profiles, profile)
		}
	}

	userProfileResponses := make([]*entity.UserProfileResponse, len(profiles))
	for i, profile := range profiles {
		var age int
		if profile.DateOfBirth != nil {
			age = calculateAge(profile.DateOfBirth)
		}

		userProfileResponses[i] = &entity.UserProfileResponse{
			ID:                profile.ID,
			FullName:          getStringValue(profile.FullName),
			ProfilePictureURL: profile.ProfilePictureURL,
			Age:               age,
			HeightCm:          getIntValue(profile.HeightCm),
			MaritalStatus:     getEnumString(profile.MaritalStatus),
			Profession:        getEnumString(profile.Profession),
			HomeDistrict:      getEnumString(profile.HomeDistrict),
		}
	}

	return userProfileResponses, &entity.PaginationData{
		TotalCount: int64(totalCount),
		Limit:      limit,
		Offset:     offset,
		HasMore:    end < totalCount,
	}, nil
}

func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func getIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func getEnumString(enum interface{}) string {
    if enum == nil {
        return ""
    }
    
    switch v := enum.(type) {
    case *entity.MaritalStatusEnum:
        if v == nil {
            return ""
        }
        return string(*v)
    case *entity.ProfessionEnum:
        if v == nil {
            return ""
        }
        return string(*v)
    case *entity.HomeDistrictEnum:
        if v == nil {
            return ""
        }
        return string(*v)
    case *entity.CommunityEnum:
        if v == nil {
            return ""
        }
        return string(*v)
    case *entity.ProfessionTypeEnum:
        if v == nil {
            return ""
        }
        return string(*v)
    case *entity.EducationLevelEnum:
        if v == nil {
            return ""
        }
        return string(*v)
    default:
        return ""
    }
}
