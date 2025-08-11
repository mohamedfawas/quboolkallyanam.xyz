package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/errutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary Get user profile
// @Description Fetch the authenticated user's profile
// @Tags User
// @Produce json
// @Success 200 {object} dto.UserProfileRecommendation "User profile"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/profile [get]
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	profile, err := h.userUsecase.GetUserProfile(authCtx.Ctx)
	if err != nil {
		if apperrors.ShouldLogError(err) && !errutil.IsGRPCError(err) {
			log.Error("Failed to get user profile", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("User profile fetched successfully")
	apiresponse.Success(c, "User profile fetched successfully", profile)
}