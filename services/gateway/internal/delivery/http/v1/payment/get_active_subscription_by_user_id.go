package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

// @Summary Get active subscription by user
// @Description Fetch currently active subscription for authenticated user
// @Tags Payment
// @Produce json
// @Success 200 {object} dto.ActiveSubscription "Active subscription"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 404 {object} dto.NotFoundError "Not found"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/payment/subscriptions [get]
func (h *PaymentHandler) GetActiveSubscriptionByUserID(c *gin.Context) {
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

	subscription, err := h.paymentUsecase.GetActiveSubscriptionByUserID(authCtx.Ctx)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to get active subscription", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Active subscription retrieved successfully")
	apiresponse.Success(c, "Active subscription retrieved successfully", subscription)
}