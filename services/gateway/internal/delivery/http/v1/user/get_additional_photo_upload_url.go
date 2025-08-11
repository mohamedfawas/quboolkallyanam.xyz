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

// @Summary Get additional photo upload URL
// @Description Generate a pre-signed upload URL for an additional photo
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.GetAdditionalPhotoUploadURLRequest true "Upload request"
// @Success 200 {object} dto.GetAdditionalPhotoUploadURLResponse "Upload URL details"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/profile/additional-photo [post]
func (h *UserHandler) GetAdditionalPhotoUploadURL(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	var req dto.GetAdditionalPhotoUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	resp, err := h.userUsecase.GetAdditionalPhotoUploadURL(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get additional photo upload URL", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("GetAdditionalPhotoUploadURL request successful")
	apiresponse.Success(c, "Additional photo upload URL generated successfully", resp)
}