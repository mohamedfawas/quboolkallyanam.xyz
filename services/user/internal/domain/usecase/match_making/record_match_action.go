package matchmaking

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	userevents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/user"
	"gorm.io/gorm"
)

func (u *matchMakingUsecase) RecordMatchAction(
	ctx context.Context,
	userID uuid.UUID,
	targetProfileID int64,
	action string) (bool, error) {

	senderProfile, err := u.userProfileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("error retrieving sender profile: %w", err)
	}
	if senderProfile == nil {
		return false, apperrors.ErrUserProfileNotFound
	}

	targetProfile, err := u.userProfileRepository.GetUserProfileByID(ctx, targetProfileID)
	if err != nil {
		return false, fmt.Errorf("error retrieving target profile: %w", err)
	}
	if targetProfile == nil {
		return false, apperrors.ErrUserProfileNotFound
	}
	if targetProfile.UserID == userID {
		return false, apperrors.ErrInvalidMatchAction
	}

	var matchAction constants.MatchAction
	switch action {
	case string(constants.MatchActionLike):
		matchAction = constants.MatchActionLike
	case string(constants.MatchActionPass):
		matchAction = constants.MatchActionPass
	default:
		return false, apperrors.ErrInvalidMatchAction
	}

	// Ensure user id on left always less than right, used for mutual match only
	userID1, userID2 := userID, targetProfile.UserID
	if bytes.Compare(userID1[:], userID2[:]) > 0 {
		userID1, userID2 = userID2, userID1
	}

	// Check if active mutual match already exists
	existingMutual, err := u.mutualMatchRepository.GetMutualMatch(ctx, userID1, userID2)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("error retrieving existing mutual match: %w", err)
	}

	// Handle existing active mutual match
	if existingMutual != nil {
		if matchAction == constants.MatchActionLike {
			// Already matched, like is a no-op
			return true, nil
		}
		if matchAction == constants.MatchActionPass {
			// Deactivate mutual match and record pass action
			if txErr := u.transactionManager.WithTransaction(ctx, func(tx *gorm.DB) error {
				// Deactivate the mutual match
				if err := u.mutualMatchRepository.DeactivateMutualMatchTx(ctx, tx, userID1, userID2); err != nil {
					return err
				}
				// Handle current user's action
				return u.upsertProfileMatchAction(ctx, tx, userID, targetProfile.UserID, false)
			}); txErr != nil {
				return false, fmt.Errorf("could not deactivate mutual match: %w", txErr)
			}
			return true, nil
		}
	}

	// Check if reverse match exists and both are likes (potential mutual match)
	existingReverseMatch, err := u.profileMatchRepository.GetExistingMatch(ctx, targetProfile.UserID, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("error retrieving existing reverse match: %w", err)
	}

	if existingReverseMatch != nil && existingReverseMatch.IsLiked && matchAction == constants.MatchActionLike {
		// Create mutual match and handle current user's action
		if txErr := u.transactionManager.WithTransaction(ctx, func(tx *gorm.DB) error {
			// Create or reactivate mutual match
			if err := u.mutualMatchRepository.UpsertMutualMatchTx(ctx, tx, userID1, userID2); err != nil {
				return err
			}
			// Handle current user's action
			return u.upsertProfileMatchAction(ctx, tx, userID, targetProfile.UserID, true)
		}); txErr != nil {
			return false, fmt.Errorf("error creating mutual match: %w", txErr)
		}

		// Publish mutual match event
		if err := u.eventPublisher.PublishMutualMatchCreated(ctx, userevents.MutualMatchCreatedEvent{
			User1Email:     senderProfile.Email,
			User1ProfileID: senderProfile.ID,
			User1FullName:  senderProfile.FullName,
			User2Email:     targetProfile.Email,
			User2ProfileID: targetProfile.ID,
			User2FullName:  targetProfile.FullName,
		}); err != nil {
			// no need to return error here, just log it in event publisher side.
		}

		return true, nil
	}

	// Handle regular action (no mutual match created)
	if err := u.upsertProfileMatchAction(ctx, nil, userID, targetProfile.UserID, matchAction == constants.MatchActionLike); err != nil {
		return false, err
	}

	if matchAction == constants.MatchActionLike {
		// Publish like event
		if err := u.eventPublisher.PublishUserInterestSent(ctx, userevents.UserInterestSentEvent{
			ReceiverEmail:   targetProfile.Email,
			SenderProfileID: senderProfile.ID,
			SenderName:      senderProfile.FullName,
		}); err != nil {
			// no need to return error here, just log it in event publisher side.
		}
	}

	return true, nil
}

// Helper method to upsert profile match action (create new or update existing)
func (u *matchMakingUsecase) upsertProfileMatchAction(ctx context.Context, tx *gorm.DB, userID, targetUserID uuid.UUID, isLiked bool) error {
	existingMatch, err := u.profileMatchRepository.GetExistingMatch(ctx, userID, targetUserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error retrieving existing user action: %w", err)
	}

	if existingMatch != nil {
		// Update existing action
		if tx != nil {
			return u.profileMatchRepository.UpdateMatchActionTx(ctx, tx, userID, targetUserID, isLiked)
		}

		return u.profileMatchRepository.UpdateMatchAction(ctx, userID, targetUserID, isLiked)
	} else {
		// Create new action
		if tx != nil {
			return u.profileMatchRepository.CreateMatchActionTx(ctx, tx, userID, targetUserID, isLiked)
		}
		return u.profileMatchRepository.CreateMatchAction(ctx, userID, targetUserID, isLiked)
	}
}
