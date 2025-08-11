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

// @Summary Confirm additional photo upload
// @Description Confirm and persist the uploaded additional photo
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.ConfirmAdditionalPhotoUploadRequest true "Confirm request"
// @Success 200 {object} dto.ConfirmAdditionalPhotoUploadResponse "Confirmation result"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/profile/additional-photo/confirm [post]
func (h *UserHandler) ConfirmAdditionalPhotoUpload(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	var req dto.ConfirmAdditionalPhotoUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	resp, err := h.userUsecase.ConfirmAdditionalPhotoUpload(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to confirm additional photo upload", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("ConfirmAdditionalPhotoUpload request successful")
	apiresponse.Success(c, "Additional photo uploaded successfully", resp)
}