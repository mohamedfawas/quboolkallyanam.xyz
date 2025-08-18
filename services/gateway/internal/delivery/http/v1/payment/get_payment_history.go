package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary Get payment history
// @Description Fetch payments for authenticated user
// @Tags Payment
// @Produce json
// @Success 200 {array} dto.GetPaymentHistoryResponse "Payment history"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/payment/payments-history [get]
func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		h.logger.Error("Failed to extract auth context", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	history, err := h.paymentUsecase.GetPaymentHistory(authCtx.Ctx)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get payment history", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Payment history retrieved successfully")
	apiresponse.Success(c, "Payment history retrieved successfully", history)
}