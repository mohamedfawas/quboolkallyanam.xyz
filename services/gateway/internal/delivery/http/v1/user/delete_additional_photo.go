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

func (h *UserHandler) DeleteAdditionalPhoto(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	// Get display_order from path parameter
	displayOrderStr := c.Param("display_order")
	displayOrderInt, err := strconv.Atoi(displayOrderStr)
	if err != nil {
		apiresponse.Error(c, apperrors.ErrInvalidInput, nil)
		return
	}

	req := dto.DeleteAdditionalPhotoRequest{
		DisplayOrder: int32(displayOrderInt),
	}

	resp, err := h.userUsecase.DeleteAdditionalPhoto(authCtx.Ctx, req)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to delete additional photo", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Additional photo deleted successfully")
	apiresponse.Success(c, "Additional photo deleted successfully", resp)
}