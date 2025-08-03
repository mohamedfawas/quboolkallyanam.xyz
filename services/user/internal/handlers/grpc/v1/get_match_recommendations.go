package v1

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pointerutil"
)



func (h *UserHandler) GetMatchRecommendations(
	ctx context.Context, 
	req *userpbv1.GetMatchRecommendationsRequest) (*userpbv1.GetMatchRecommendationsResponse, error) {
	
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

	profiles, pagination, err := h.matchMakingUsecase.RecommendUserProfiles(ctx, userIDUUID, limit, offset)
	if err != nil {
		if !appErrors.IsAppError(err) {
			log.Error("Failed to get match recommendations", zap.Error(err))
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
		Limit:      uint32(pagination.Limit),
		Offset:     uint32(pagination.Offset),
		HasMore:    pagination.HasMore,
	}

	log.Info("Successfully fetched match recommendations")
	return &userpbv1.GetMatchRecommendationsResponse{
		Profiles:   protoProfiles,
		Pagination: protoPagination,
	}, nil
}