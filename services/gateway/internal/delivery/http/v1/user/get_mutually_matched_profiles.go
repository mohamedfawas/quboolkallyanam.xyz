package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

// @Summary Get mutually matched profiles
// @Description List profiles that mutually matched with the authenticated user
// @Tags User
// @Produce json
// @Param limit query int false "Items per page (1-50)" minimum(1) maximum(50)
// @Param offset query int false "Offset (>= 0)" minimum(0)
// @Success 200 {object} dto.GetProfilesByMatchActionResponse "Profiles list"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/matches/mutual [get]
func (h *UserHandler) GetMutuallyMatchedProfiles(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit64, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit64 < 1 || limit64 > 50 {
		apiresponse.Error(c, apperrors.ErrInvalidPaginationLimit, nil)
		return
	}
	offset64, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil || offset64 < 0 {
		apiresponse.Error(c, apperrors.ErrInvalidPaginationPage, nil)
		return
	}

	req := dto.GetProfilesByMatchActionRequest{
		Action: validation.MatchMakingOptionMutual,
		Limit:  int32(limit64),
		Offset: int32(offset64),
	}

	resp, err := h.userUsecase.GetProfilesByMatchAction(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get mutually matched profiles", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Mutually matched profiles retrieved successfully")
	apiresponse.Success(c, "Mutually matched profiles retrieved successfully", resp)
}
