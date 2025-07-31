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

// @Summary Update user profile
// @Description Update user profile
// @Tags User
// @Accept json
// @Produce json
// @Param user_profile_put_request body dto.UserProfilePutRequest true "User profile put request"
// @Success 200 {object} apiresponse.Response "User profile updated successfully"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/profile [put]

func (h *UserHandler) PutUserProfile(c *gin.Context) {
	requestID, exists := c.Get(constants.ContextKeyRequestID)
	if !exists {
		h.logger.Error("failed to get request ID from context")
	}

	var req dto.UserProfilePutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		h.logger.Error("failed to get user ID from context")
		apiresponse.Error(c, apperrors.ErrUserContextMissing, nil)
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, requestID.(string)),
		zap.String(constants.ContextKeyUserID, userID.(string)),
	)

	request := dto.UserProfilePatchRequest{
		IsBride:               req.IsBride,
		FullName:              req.FullName,
		DateOfBirth:           req.DateOfBirth,
		HeightCm:              req.HeightCm,
		PhysicallyChallenged:  req.PhysicallyChallenged,
		Community:             req.Community,
		MaritalStatus:         req.MaritalStatus,
		Profession:            req.Profession,
		ProfessionType:        req.ProfessionType,
		HighestEducationLevel: req.HighestEducationLevel,
		HomeDistrict:          req.HomeDistrict,
	}

	err := h.userUsecase.UpdateUserProfile(ctx, request)
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
