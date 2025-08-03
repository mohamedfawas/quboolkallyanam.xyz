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

// @Summary Update user profile
// @Description Update user profile
// @Tags User
// @Accept json
// @Produce json
// @Param user_profile_put_request body dto.UserProfilePutRequest true "User profile put request"
// @Success 200 {object} apiresponse.Response "User profile updated successfully"
// @Failure 400 {object} apiresponse.Response "Bad request"
// @Failure 401 {object} apiresponse.Response "Unauthorized"
// @Failure 404 {object} apiresponse.Response "User profile not found"
// @Failure 422 {object} apiresponse.Response "Validation error"
// @Failure 500 {object} apiresponse.Response "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/profile [put]
func (h *UserHandler) PutUserProfile(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	var req dto.UserProfilePutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	request := dto.UserProfilePatchRequest{
		IsBride:               &req.IsBride,
		FullName:              &req.FullName,
		DateOfBirth:           &req.DateOfBirth,
		HeightCm:              &req.HeightCm,
		PhysicallyChallenged:  &req.PhysicallyChallenged,
		Community:             &req.Community,
		MaritalStatus:         &req.MaritalStatus,
		Profession:            &req.Profession,
		ProfessionType:        &req.ProfessionType,
		HighestEducationLevel: &req.HighestEducationLevel,
		HomeDistrict:          &req.HomeDistrict,
	}

	err = h.userUsecase.UpdateUserProfile(authCtx.Ctx, request)
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
