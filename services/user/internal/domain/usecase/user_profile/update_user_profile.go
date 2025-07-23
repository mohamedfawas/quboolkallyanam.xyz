package userprofile

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *userProfileUsecase) UpdateUserProfile(
	ctx context.Context,
	userID uuid.UUID,
	req entity.UpdateUserProfileRequest) error {

	existingProfile, err := u.userProfileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		if err == appError.ErrUserProfileNotFound {
			return appError.ErrUserProfileNotFound
		}
		log.Printf("failed to get user profile: %v", err)
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	existingProfile.IsBride = *req.IsBride
	existingProfile.FullName = req.FullName
	existingProfile.DateOfBirth = req.DateOfBirth
	existingProfile.HeightCm = req.HeightCm
	existingProfile.PhysicallyChallenged = *req.PhysicallyChallenged
	existingProfile.Community = (*entity.CommunityEnum)(req.Community)
	existingProfile.MaritalStatus = (*entity.MaritalStatusEnum)(req.MaritalStatus)
	existingProfile.Profession = (*entity.ProfessionEnum)(req.Profession)
	existingProfile.ProfessionType = (*entity.ProfessionTypeEnum)(req.ProfessionType)
	existingProfile.HighestEducationLevel = (*entity.EducationLevelEnum)(req.HighestEducationLevel)
	existingProfile.HomeDistrict = (*entity.HomeDistrictEnum)(req.HomeDistrict)

	now := time.Now().UTC()
	existingProfile.UpdatedAt = now
	existingProfile.LastLogin = now

	if err := u.userProfileRepository.UpdateUserProfile(ctx, existingProfile); err != nil {
		log.Printf("failed to update user profile: %v", err)
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	log.Printf("user profile updated successfully")
	return nil
}
