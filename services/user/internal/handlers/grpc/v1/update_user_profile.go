package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"go.uber.org/zap"
)

func (h *UserHandler) UpdateUserProfile(
	ctx context.Context,
	req *userpbv1.UpdateUserProfileRequest) (*userpbv1.UpdateUserProfileResponse, error) {

	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)

	entityReq := entity.UpdateUserProfileRequest{}
	if req.IsBride != nil {
		entityReq.IsBride = &req.IsBride.Value
	}
	if req.FullName != nil {
		entityReq.FullName = &req.FullName.Value
	}
	if req.DateOfBirth != nil {
		entityReq.DateOfBirth = &req.DateOfBirth.Value
	}
	if req.HeightCm != nil {
		entityReq.HeightCm = &req.HeightCm.Value
	}
	if req.PhysicallyChallenged != nil {
		entityReq.PhysicallyChallenged = &req.PhysicallyChallenged.Value
	}
	if req.Community != nil {
		entityReq.Community = &req.Community.Value
	}
	if req.MaritalStatus != nil {
		entityReq.MaritalStatus = &req.MaritalStatus.Value
	}
	if req.Profession != nil {
		entityReq.Profession = &req.Profession.Value
	}
	if req.ProfessionType != nil {
		entityReq.ProfessionType = &req.ProfessionType.Value
	}
	if req.HighestEducationLevel != nil {
		entityReq.HighestEducationLevel = &req.HighestEducationLevel.Value
	}
	if req.HomeDistrict != nil {
		entityReq.HomeDistrict = &req.HomeDistrict.Value
	}

	userIDUUID, err := uuid.Parse(contextData.UserID)
	if err != nil {
		log.Error("Failed to parse user ID", zap.Error(err))
		return nil, err
	}

	err = h.userProfileUsecase.UpdateUserProfile(ctx, userIDUUID, entityReq)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to update user profile", zap.Error(err))
		}
		return nil, err
	}

	log.Info("User profile updated successfully")
	return &userpbv1.UpdateUserProfileResponse{Success: &wrapperspb.BoolValue{Value: true}}, nil
}
