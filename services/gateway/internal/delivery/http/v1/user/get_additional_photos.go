package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary Get additional photos
// @Description Fetch all additional photos for the authenticated user
// @Tags User
// @Produce json
// @Success 200 {object} dto.GetAdditionalPhotosResponse "Additional photos"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/profile/additional-photos [get]
func (h *UserHandler) GetAdditionalPhotos(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	resp, err := h.userUsecase.GetAdditionalPhotos(authCtx.Ctx)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get additional photos", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Additional photos fetched successfully")
	apiresponse.Success(c, "Additional photos fetched successfully", resp)
}
