package userprofile

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *userProfileUsecase) UpdateUserProfile(
	ctx context.Context,
	userID uuid.UUID,
	req entity.UpdateUserProfileRequest) error {

	existingProfile, err := u.userProfileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	if existingProfile == nil {
		return appError.ErrUserProfileNotFound
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
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	userProfileUpdatedEvent := userevents.UserProfileUpdatedEvent{
		UserID:        userID,
		UserProfileID: int64(existingProfile.ID),
		Email:         *existingProfile.Email,
		Phone:         *existingProfile.Phone,
		FullName:      *existingProfile.FullName,
		CreatedAt:     existingProfile.CreatedAt,
		UpdatedAt:     existingProfile.UpdatedAt,
	}

	if err := u.eventPublisher.PublishUserProfileUpdated(ctx, userProfileUpdatedEvent); err != nil {
		// No need to fail, logging is done in event publishing code.
	}

	return nil
}
