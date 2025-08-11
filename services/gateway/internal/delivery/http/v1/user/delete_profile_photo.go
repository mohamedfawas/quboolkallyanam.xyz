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

// @Summary Delete profile photo
// @Description Delete current profile photo
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} dto.DeleteProfilePhotoResponse "Delete result"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/profile/profile-photo [delete]
func (h *UserHandler) DeleteProfilePhoto(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	err = h.userUsecase.DeleteProfilePhoto(authCtx.Ctx, dto.DeleteProfilePhotoRequest{})
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to delete profile photo", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Profile photo deleted successfully")
	apiresponse.Success(c, "Profile photo deleted successfully", nil)
}