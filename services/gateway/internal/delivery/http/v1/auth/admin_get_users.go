package auth

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

func (h *AuthHandler) AdminGetUsers(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	// Extract query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		apiresponse.Error(c, apperrors.ErrInvalidPaginationPage, nil)
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit < 1 || limit > 100 {
		apiresponse.Error(c, apperrors.ErrInvalidPaginationLimit, nil)
		return
	}

	req := dto.GetUsersRequest{
		Page:  int32(page),
		Limit: int32(limit),
	}

	resp, err := h.authUsecase.GetUsers(reqCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get users", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Users retrieved successfully", zap.Int32("page", req.Page), zap.Int32("limit", req.Limit))
	apiresponse.Success(c, "Users retrieved successfully", resp)
}