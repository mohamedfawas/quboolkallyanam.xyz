package payment

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

// @Summary Get completed payment details (Admin)
// @Description List all completed payments with pagination
// @Tags Payment
// @Accept json
// @Produce json
// @Param page query int true "Page number (>=1)" default(1)
// @Param limit query int true "Page size (1-50)" default(10)
// @Success 200 {object} dto.GetCompletedPaymentDetailsResponse
// @Failure 400 {object} dto.BadRequestError
// @Failure 401 {object} dto.UnauthorizedError
// @Failure 403 {object} dto.ForbiddenError
// @Failure 500 {object} dto.InternalServerError
// @Security BearerAuth
// @Router /api/v1/payment/admin/completed-payments [get]
func (h *PaymentHandler) GetCompletedPaymentDetails(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page64, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page64 < 1 {
		apiresponse.Error(c, apperrors.ErrInvalidPaginationPage, nil)
		return
	}
	limit64, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit64 < 1 || limit64 > 50 {
		apiresponse.Error(c, apperrors.ErrInvalidPaginationLimit, nil)
		return
	}

	req := dto.GetCompletedPaymentDetailsRequest{
		Page:  int32(page64),
		Limit: int32(limit64),
	}

	resp, err := h.paymentUsecase.GetCompletedPaymentDetails(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get completed payment details", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Completed payment details fetched")
	apiresponse.Success(c, "Completed payment details fetched", resp)
}