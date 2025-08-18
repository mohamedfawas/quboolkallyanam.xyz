package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary List active subscription plans
// @Description Get all active subscription plans
// @Tags Payment
// @Produce json
// @Success 200 {array} dto.SubscriptionPlan "Active plans"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Router /api/v1/payment/subscription-plans [get]
func (h *PaymentHandler) GetActiveSubscriptionPlans(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract request context", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	plans, err := h.paymentUsecase.GetActiveSubscriptionPlans(reqCtx.Ctx)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get active subscription plans", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Active subscription plans retrieved successfully")
	apiresponse.Success(c, "Active subscription plans retrieved successfully", plans)
}