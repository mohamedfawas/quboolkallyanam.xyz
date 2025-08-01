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

// @Summary Partial update user profile
// @Description Partial update user profile
// @Tags User
// @Accept json
// @Produce json
// @Param user_profile_patch_request body dto.UserProfilePatchRequest true "User profile patch request"
// @Success 200 {object} apiresponse.Response "User profile updated successfully"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Router /api/v1/user/profile [patch]

func (h *UserHandler) PatchUserProfile(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	var req dto.UserProfilePatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	err = h.userUsecase.UpdateUserProfile(authCtx.Ctx, req)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to update user profile", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("User profile updated successfully")

	apiresponse.Success(c, "User profile updated successfully", nil)
}
