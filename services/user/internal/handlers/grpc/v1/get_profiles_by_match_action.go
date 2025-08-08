package v1

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pointerutil"
)

func (h *UserHandler) GetProfilesByMatchAction(
	ctx context.Context, 
	req *userpbv1.GetProfilesByMatchActionRequest) (*userpbv1.GetProfilesByMatchActionResponse, error) {
	
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

	limit := int(req.Limit)
	offset := int(req.Offset)

	profiles, pagination, err := h.matchMakingUsecase.GetProfilesByMatchAction(ctx, userIDUUID, req.Action, limit, offset)
	if err != nil {
		if !appErrors.IsAppError(err) {
			log.Error("Failed to get profiles by match action", zap.Error(err))
		}
		return nil, err
	}

	protoProfiles := make([]*userpbv1.UserProfileRecommendation, len(profiles))
	for i, profile := range profiles {
		protoProfiles[i] = &userpbv1.UserProfileRecommendation{
			Id:                profile.ID,
			FullName:          profile.FullName,
			ProfilePictureUrl: pointerutil.GetStringValue(profile.ProfilePictureURL),
			Age:               profile.Age,
			HeightCm:          profile.HeightCm,
			MaritalStatus:     profile.MaritalStatus,
			Profession:        profile.Profession,
			HomeDistrict:      profile.HomeDistrict,
		}
	}

	protoPagination := &userpbv1.PaginationInfo{
		TotalCount: pagination.TotalCount,
		Limit:      int32(pagination.Limit),
		Offset:     int32(pagination.Offset),
		HasMore:    pagination.HasMore,
	}

	log.Info("Successfully fetched profiles by match action")
	return &userpbv1.GetProfilesByMatchActionResponse{
		Profiles:   protoProfiles,
		Pagination: protoPagination,
	}, nil
}








