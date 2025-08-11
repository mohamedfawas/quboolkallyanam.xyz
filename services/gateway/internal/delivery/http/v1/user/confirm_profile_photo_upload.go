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

// @Summary Confirm profile photo upload
// @Description Confirm and persist the uploaded profile photo
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.ConfirmProfilePhotoUploadRequest true "Confirm request"
// @Success 200 {object} dto.ConfirmProfilePhotoUploadResponse "Confirmation result"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/profile/profile-photo/confirm [post]
func (h *UserHandler) ConfirmProfilePhotoUpload(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	var req dto.ConfirmProfilePhotoUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	resp, err := h.userUsecase.ConfirmProfilePhotoUpload(c, req)
	if err != nil {
		if apperrors.ShouldLogError(err)  {
			log.Error("Failed to confirm profile photo upload", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("ConfirmProfilePhotoUpload request successful")
	apiresponse.Success(c, "Profile photo uploaded successfully", resp)
}