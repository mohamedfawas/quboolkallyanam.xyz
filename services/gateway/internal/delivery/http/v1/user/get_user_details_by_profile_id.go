package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

// @Summary Get target user details by profile ID (user)
// @Description Returns target user's profile, partner preferences (if any), and additional photos (if any). Enforces opposite-gender access.
// @Tags User
// @Produce json
// @Param profile_id path int true "Target profile ID"
// @Success 200 {object} dto.GetUserDetailsByProfileIDResponse
// @Failure 400 {object} dto.BadRequestError
// @Failure 401 {object} dto.UnauthorizedError
// @Failure 403 {object} dto.ForbiddenError
// @Failure 404 {object} dto.NotFoundError
// @Failure 500 {object} dto.InternalServerError
// @Security BearerAuth
// @Router /api/v1/user/profiles/{profile_id} [get]
func (h *UserHandler) GetUserDetailsByProfileIDForUser(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	profileIDStr := c.Param("profile_id")
	targetID, err := strconv.ParseInt(profileIDStr, 10, 64)
	if err != nil || targetID <= 0 {
		apiresponse.Error(c, apperrors.ErrInvalidTargetProfileID, nil)
		return
	}

	req := dto.GetUserDetailsByProfileIDRequest{
		TargetProfileID:  targetID,
		RequestedByAdmin: false,
	}

	resp, err := h.userUsecase.GetUserDetailsByProfileID(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get user details by profile id (user)", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Fetched user details by profile id (user)",
		zap.Int64("target_profile_id", targetID),
	)
	apiresponse.Success(c, "Fetched user details successfully", resp)
}