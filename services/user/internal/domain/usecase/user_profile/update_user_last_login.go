package userprofile

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *userProfileUsecase) UpdateUserLastLogin(ctx context.Context,
	userID uuid.UUID,
	email, phone string) error {

	exists, err := u.userProfileRepository.ProfileExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check if user profile exists: %w", err)
	}

	if exists {
		if err := u.userProfileRepository.UpdateLastLogin(ctx, userID); err != nil {
			return fmt.Errorf("failed to update user last login: %w", err)
		}

		return nil
	}

	now := time.Now().UTC()
	community := validation.CommunityAny
	maritalStatus := validation.MaritalStatusAny
	profession := validation.ProfessionAny
	professionType := validation.ProfessionTypeAny
	highestEducationLevel := validation.EducationLevelAny
	homeDistrict := validation.HomeDistrictAny

	// create minimal profile with required fields
	profile := &entity.UserProfile{
		UserID:                userID,
		Email:                 email,
		Phone:                 phone,
		IsBride:               false,
		LastLogin:             now,
		CreatedAt:             now,
		UpdatedAt:             now,
		Community:             community,
		MaritalStatus:         maritalStatus,
		Profession:            profession,
		ProfessionType:        professionType,
		HighestEducationLevel: highestEducationLevel,
		HomeDistrict:          homeDistrict,
	}

	if err := u.userProfileRepository.CreateUserProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	return nil
}
