package userprofile

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (u *userProfileUsecase) UpdateUserLastLogin(ctx context.Context,
	userID uuid.UUID,
	email, phone string) error {

	exists, err := u.userProfileRepository.ProfileExists(ctx, userID)
	if err != nil {
		log.Printf("failed to check if user profile exists: %v", err)
		return fmt.Errorf("failed to check if user profile exists: %w", err)
	}

	if exists {
		lastLogin := time.Now().UTC()
		if err := u.userProfileRepository.UpdateLastLogin(ctx, userID, lastLogin); err != nil {
			log.Printf("failed to update user last login: %v", err)
			return fmt.Errorf("failed to update user last login: %w", err)
		}

		log.Printf("user last login updated successfully for user: %s", userID)
		return nil
	}

	now := time.Now().UTC()
	community := entity.CommunityNotMentioned
	maritalStatus := entity.MaritalStatusNotMentioned
	profession := entity.ProfessionNotMentioned
	professionType := entity.ProfessionTypeNotMentioned
	highestEducationLevel := entity.EducationLevelNotMentioned
	homeDistrict := entity.HomeDistrictNotMentioned

	// create minimal profile with required fields
	profile := &entity.UserProfile{
		UserID:                userID,
		Email:                 &email,
		Phone:                 &phone,
		IsBride:               false,
		LastLogin:             now,
		CreatedAt:             now,
		UpdatedAt:             now,
		Community:             &community,
		MaritalStatus:         &maritalStatus,
		Profession:            &profession,
		ProfessionType:        &professionType,
		HighestEducationLevel: &highestEducationLevel,
		HomeDistrict:          &homeDistrict,
	}

	if err := u.userProfileRepository.CreateUserProfile(ctx, profile); err != nil {
		log.Printf("failed to create user profile: %v", err)
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	if u.eventPublisher != nil {
		userProfileCreatedEvent := userevents.UserProfileCreatedEvent{
			UserID:        userID,
			Email:         email,
			Phone:         phone,
			CreatedAt:     now,
			UpdatedAt:     now,
			UserProfileID: profile.ID,
		}

		if err := u.eventPublisher.PublishUserProfileCreated(ctx, userProfileCreatedEvent); err != nil {
			log.Printf("failed to publish user profile created event: %v", err)
		}
	}

	log.Printf("user profile created successfully for user: %s", userID)
	return nil
}
