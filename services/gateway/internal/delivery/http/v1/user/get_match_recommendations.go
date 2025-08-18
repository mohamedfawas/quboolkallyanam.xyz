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

// @Summary Get match recommendations
// @Description Get recommended profiles for the authenticated user with pagination
// @Tags User
// @Produce json
// @Param limit query int false "Items per page (1-50)" minimum(1) maximum(50)
// @Param offset query int false "Offset (>= 0)" minimum(0)
// @Success 200 {object} dto.GetMatchRecommendationsResponse "Profiles list"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/recommendations [get]
func (h *UserHandler) GetMatchRecommendations(c *gin.Context) {
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

	req := dto.GetMatchRecommendationsRequest{
		Limit:  int32(limit64),
		Offset: int32(offset64),
	}

	resp, err := h.userUsecase.GetMatchRecommendations(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get match recommendations", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Match recommendations retrieved successfully",
		zap.Int32("limit", req.Limit),
		zap.Int32("offset", req.Offset),
	)
	apiresponse.Success(c, "Match recommendations retrieved successfully", resp)
}