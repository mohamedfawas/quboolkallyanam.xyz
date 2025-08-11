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

// @Summary Get profile photo upload URL
// @Description Generate a pre-signed upload URL for the profile photo
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.GetProfilePhotoUploadURLRequest true "Upload request"
// @Success 200 {object} dto.GetProfilePhotoUploadURLResponse "Upload URL details"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/profile/profile-photo [post]
func (h *UserHandler) GetProfilePhotoUploadURL(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	var req dto.GetProfilePhotoUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	resp, err := h.userUsecase.GetProfilePhotoUploadURL(c, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get profile photo upload URL", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("GetProfilePhotoUploadURL request successful")
	apiresponse.Success(c, "Profile photo upload URL generated successfully", resp)
}