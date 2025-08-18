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

// @Summary Update subscription plan
// @Description Update fields of an existing subscription plan (admin only)
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body dto.UpdateSubscriptionPlanRequest true "Update subscription plan"
// @Success 200 {object} dto.CreateOrUpdateSubscriptionPlanResponse "Plan updated"
// @Failure 400 {object} dto.BadRequestError "Bad request"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 403 {object} dto.ForbiddenError "Forbidden"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/payment/subscription-plan [patch]
func (h *PaymentHandler) UpdateSubscriptionPlan(c *gin.Context) {
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

	var req dto.UpdateSubscriptionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	resp, err := h.paymentUsecase.CreateOrUpdateSubscriptionPlan(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to update subscription plan", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Subscription plan updated successfully", zap.String("plan_id", req.ID))
	apiresponse.Success(c, "Subscription plan updated successfully", resp)
}