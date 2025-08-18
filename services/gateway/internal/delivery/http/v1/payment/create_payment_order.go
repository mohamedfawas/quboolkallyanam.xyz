package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

// @Summary Create payment order
// @Description Create a Razorpay order for the selected subscription plan
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body dto.PaymentOrderRequest true "Create payment order request"
// @Success 200 {object} dto.PaymentOrderResponse "Payment order created"
// @Failure 400 {object} dto.BadRequestError "Bad request"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/payment/order [post]
func (h *PaymentHandler) CreatePaymentOrder(c *gin.Context) {
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

	var req dto.PaymentOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	resp, err := h.paymentUsecase.CreatePaymentOrder(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err){
			log.Error("Failed to create payment order", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Payment order created successfully")
	apiresponse.Success(c, "Payment order created successfully", resp)
}