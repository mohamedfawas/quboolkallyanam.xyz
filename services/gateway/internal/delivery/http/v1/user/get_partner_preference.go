package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary Get user partner preferences
// @Description Fetch the authenticated user's partner preferences
// @Tags User
// @Produce json
// @Success 200 {object} dto.GetUserPartnerPreferencesResponse "Partner preferences"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 404 {object} dto.NotFoundError "Not found"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/preference [get]
func (h *UserHandler) GetPartnerPreference(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	resp, err := h.userUsecase.GetUserPartnerPreferences(authCtx.Ctx)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get partner preferences", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Partner preferences fetched successfully")
	apiresponse.Success(c, "Partner preferences fetched successfully", resp)
}