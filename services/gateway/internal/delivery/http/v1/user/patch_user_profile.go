package user

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
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
	requestID, _ := c.Get(constants.ContextKeyRequestID)

	var req dto.UserProfilePatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		apiresponse.Error(c, apperrors.ErrUserContextMissing, nil)
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)

	err := h.userUsecase.UpdateUserProfile(ctx, req)
	if err != nil {
		// Log only internal server errors
		if !apperrors.IsAppError(err) {
			h.logger.Error("Failed to update user profile",
				zap.String(constants.RequestID, requestID.(string)),
				zap.String(constants.UserIDS, userID.(string)),
				zap.Error(err))
		}

		apiresponse.Error(c, err, nil)
		return
	}

	h.logger.Info("User profile updated successfully",
		zap.String(constants.RequestID, requestID.(string)),
		zap.String(constants.UserIDS, userID.(string)))

	apiresponse.Success(c, "User profile updated successfully", nil)
}
