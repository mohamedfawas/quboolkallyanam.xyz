package v1

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
)

func (h *UserHandler) UpdateUserPartnerPreferences(
	ctx context.Context,
	req *userpbv1.UpdateUserPartnerPreferencesRequest) (*userpbv1.UpdateUserPartnerPreferencesResponse, error) {

	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)

	userIDUUID, err := uuid.Parse(contextData.UserID)
	if err != nil {
		log.Error("Failed to parse user ID", zap.Error(err))
		return nil, err
	}

	operationType := ""
	if req.OperationType != nil {
		operationType = req.OperationType.Value
	}

	entityReq := entity.UpdateUserPartnerPreferencesRequest{}

	if req.MinAgeYears != nil {
		minAgeYears := int16(req.MinAgeYears.Value)
		entityReq.MinAgeYears = &minAgeYears
	}
	if req.MaxAgeYears != nil {
		maxAgeYears := int16(req.MaxAgeYears.Value)
		entityReq.MaxAgeYears = &maxAgeYears
	}
	if req.MinHeightCm != nil {
		minHeightCm := int16(req.MinHeightCm.Value)
		entityReq.MinHeightCM = &minHeightCm
	}
	if req.MaxHeightCm != nil {
		maxHeightCm := int16(req.MaxHeightCm.Value)
		entityReq.MaxHeightCM = &maxHeightCm
	}
	if req.AcceptPhysicallyChallenged != nil {
		entityReq.AcceptPhysicallyChallenged = &req.AcceptPhysicallyChallenged.Value
	}

	if len(req.PreferredCommunities) > 0 {
		entityReq.PreferredCommunities = &req.PreferredCommunities
	}
	if len(req.PreferredMaritalStatus) > 0 {
		entityReq.PreferredMaritalStatus = &req.PreferredMaritalStatus
	}
	if len(req.PreferredProfessions) > 0 {
		entityReq.PreferredProfessions = &req.PreferredProfessions
	}
	if len(req.PreferredProfessionTypes) > 0 {
		entityReq.PreferredProfessionTypes = &req.PreferredProfessionTypes
	}
	if len(req.PreferredEducationLevels) > 0 {
		entityReq.PreferredEducationLevels = &req.PreferredEducationLevels
	}
	if len(req.PreferredHomeDistricts) > 0 {
		entityReq.PreferredHomeDistricts = &req.PreferredHomeDistricts
	}

	err = h.userProfileUsecase.UpdateUserPartnerPreferences(ctx, userIDUUID, operationType, entityReq)
	if err != nil {
		if !appErrors.IsAppError(err) {
			log.Error("Failed to update partner preferences", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Partner preferences updated successfully")
	return &userpbv1.UpdateUserPartnerPreferencesResponse{Success: &wrapperspb.BoolValue{Value: true}}, nil
}
