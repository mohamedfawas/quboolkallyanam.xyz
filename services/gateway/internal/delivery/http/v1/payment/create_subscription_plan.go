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

// @Summary Create subscription plan
// @Description Create a new subscription plan (admin only)
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body dto.CreateSubscriptionPlanRequest true "Create subscription plan"
// @Success 200 {object} dto.CreateOrUpdateSubscriptionPlanResponse "Plan created"
// @Failure 400 {object} dto.BadRequestError "Bad request"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 403 {object} dto.ForbiddenError "Forbidden"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/payment/subscription-plan [post]
func (h *PaymentHandler) CreateSubscriptionPlan(c *gin.Context) {
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

	var req dto.CreateSubscriptionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	updateReq := dto.UpdateSubscriptionPlanRequest{
		ID:           req.ID,
		DurationDays: &req.DurationDays,
		Amount:       &req.Amount,
		Currency:     &req.Currency,
		Description:  &req.Description,
		IsActive:     &req.IsActive,
	}

	resp, err := h.paymentUsecase.CreateOrUpdateSubscriptionPlan(authCtx.Ctx, updateReq)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to create subscription plan", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Subscription plan created successfully", zap.String("plan_id", req.ID))
	apiresponse.Success(c, "Subscription plan created successfully", resp)
}