package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

func (h *UserHandler) PostPartnerPreference(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)
	
	var req dto.PartnerPreferenceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	request := dto.PartnerPreferencePatchRequest{
		MinAgeYears:                &req.MinAgeYears,
		MaxAgeYears:                &req.MaxAgeYears,
		MinHeightCM:                &req.MinHeightCM,
		MaxHeightCM:                &req.MaxHeightCM,
		AcceptPhysicallyChallenged: &req.AcceptPhysicallyChallenged,
		PreferredCommunities:       &req.PreferredCommunities,
		PreferredMaritalStatus:     &req.PreferredMaritalStatus,
		PreferredProfessions:       &req.PreferredProfessions,
		PreferredProfessionTypes:   &req.PreferredProfessionTypes,
		PreferredEducationLevels:   &req.PreferredEducationLevels,
		PreferredHomeDistricts:     &req.PreferredHomeDistricts,
	}
	err = h.userUsecase.UpdateUserPartnerPreferences(authCtx.Ctx, constants.CreateOperationType, request)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("failed to update partner preference", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("partner preference updated successfully")
	apiresponse.Success(c, "Partner preference updated successfully", nil)
}
