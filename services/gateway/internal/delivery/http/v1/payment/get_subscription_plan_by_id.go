package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary Get subscription plan
// @Description Get subscription plan by ID
// @Tags Payment
// @Produce json
// @Param plan_id query string true "Subscription plan ID"
// @Success 200 {object} dto.SubscriptionPlan "Subscription plan"
// @Failure 400 {object} dto.BadRequestError "Bad request"
// @Failure 404 {object} dto.NotFoundError "Not found"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Router /api/v1/payment/subscription-plan [get]
func (h *PaymentHandler) GetSubscriptionPlanByID(c *gin.Context) {
	reqCtx, err := contextutils.ExtractRequestContext(c)
	if err != nil {
		h.logger.Error("Failed to extract request context", zap.Error(err))
		apiresponse.Error(c, err, nil)
		return
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
	)

	planID := c.Query("plan_id")
	if planID == "" {
		apiresponse.Error(c, apperrors.ErrMissingRequiredFields, map[string]string{"plan_id": "required"})
		return
	}

	plan, err := h.paymentUsecase.GetSubscriptionPlan(reqCtx.Ctx, planID)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get subscription plan", zap.Error(err), zap.String("plan_id", planID))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Subscription plan retrieved successfully", zap.String("plan_id", planID))
	apiresponse.Success(c, "Subscription plan retrieved successfully", plan)
}