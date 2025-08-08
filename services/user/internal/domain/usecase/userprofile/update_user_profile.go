package userprofile

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	appError "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/dateutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
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

	if req.IsBride != nil {
		existingProfile.IsBride = *req.IsBride
	}
	if req.FullName != nil {
		existingProfile.FullName = *req.FullName
	}
	if req.DateOfBirth != nil {
		dateOfBirth, err := dateutil.ParseDate(*req.DateOfBirth)
		if err != nil {
			return fmt.Errorf("failed to parse date of birth: %w", err)
		}
		existingProfile.DateOfBirth = *dateOfBirth
	}
	if req.HeightCm != nil {
		existingProfile.HeightCm = int16(*req.HeightCm)
	}
	if req.PhysicallyChallenged != nil {
		existingProfile.PhysicallyChallenged = *req.PhysicallyChallenged
	}

	if req.Community != nil {
		if !validation.IsValidCommunity(*req.Community) {
			return appError.ErrInvalidCommunity
		}
		existingProfile.Community = validation.Community(*req.Community)
	}

	if req.MaritalStatus != nil {
		if !validation.IsValidMaritalStatus(*req.MaritalStatus) {
			return appError.ErrInvalidMaritalStatus
		}
		existingProfile.MaritalStatus = validation.MaritalStatus(*req.MaritalStatus)
	}

	if req.Profession != nil {
		if !validation.IsValidProfession(*req.Profession) {
			return appError.ErrInvalidProfession
		}
		existingProfile.Profession = validation.Profession(*req.Profession)
	}

	if req.ProfessionType != nil {
		if !validation.IsValidProfessionType(*req.ProfessionType) {
			return appError.ErrInvalidProfessionType
		}
		existingProfile.ProfessionType = validation.ProfessionType(*req.ProfessionType)
	}

	if req.HighestEducationLevel != nil {
		if !validation.IsValidEducationLevel(*req.HighestEducationLevel) {
			return appError.ErrInvalidEducationLevel
		}
		existingProfile.HighestEducationLevel = validation.EducationLevel(*req.HighestEducationLevel)
	}

	if req.HomeDistrict != nil {
		if !validation.IsValidHomeDistrict(*req.HomeDistrict) {
			return appError.ErrInvalidHomeDistrict
		}
		existingProfile.HomeDistrict = validation.HomeDistrict(*req.HomeDistrict)
	}

	existingProfile.ProfileCompleted = true // when the first update happens , the user have self verified the profile
	// because when the first user profile is created it's done based on login event, so user haven't self verified the profile

	now := time.Now().UTC()
	existingProfile.UpdatedAt = now
	existingProfile.LastLogin = now

	if err := u.userProfileRepository.UpdateUserProfile(ctx, existingProfile); err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	userProfileUpdatedEvent := userevents.UserProfileUpdatedEvent{
		UserID:        userID,
		UserProfileID: existingProfile.ID,
		Email:         existingProfile.Email,
		Phone:         existingProfile.Phone,
		FullName:      existingProfile.FullName,
		CreatedAt:     existingProfile.CreatedAt,
		UpdatedAt:     existingProfile.UpdatedAt,
	}

	if err := u.eventPublisher.PublishUserProfileUpdated(ctx, userProfileUpdatedEvent); err != nil {
		// No need to fail, logging is done in event publishing code.
	}

	return nil
}
