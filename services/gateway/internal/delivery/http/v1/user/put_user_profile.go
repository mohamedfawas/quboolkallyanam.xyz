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


// @Summary Replace user profile
// @Description Replace full user profile
// @Tags User
// @Accept json
// @Produce json
// @Param user_profile body dto.UserProfilePutRequest true "Complete user profile"
// @Success 200 {object} dto.UpdateUserProfileResponse "Update result"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
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
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to update user profile", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("User profile updated successfully")
	apiresponse.Success(c, "User profile updated successfully", nil)
}
