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

// @Summary Create partner preference
// @Description Create initial partner preference for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param preference body dto.PartnerPreferenceCreateRequest true "Partner preference payload"
// @Success 200 {object} dto.UpdateUserProfileResponse "Update result"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/preference [post]
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
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to create partner preference", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Partner preference created successfully")
	apiresponse.Success(c, "Partner preference updated successfully", nil)
}
