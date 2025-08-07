package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

func (h *AuthHandler) AdminGetUserByField(c *gin.Context) {
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
	field := c.Query("field")
	value := c.Query("value")

	if field == "" || value == "" {
		apiresponse.Error(c, apperrors.ErrMissingRequiredFields , nil)
		return
	}

	req := dto.GetUserByFieldRequest{
		Field: field,
		Value: value,
	}

	resp, err := h.authUsecase.GetUserByField(reqCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get user by field", zap.Error(err), zap.String("field", field), zap.String("value", value))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("User retrieved successfully", zap.String("field", field), zap.String("value", value))
	apiresponse.Success(c, "User retrieved successfully", resp)
}