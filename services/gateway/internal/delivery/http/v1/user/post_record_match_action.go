package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

// @Summary Record match action (POST)
// @Description Record a new match action for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param match_action body dto.RecordMatchActionRequest true "Match action payload"
// @Success 200 {object} dto.RecordMatchActionResponse "Result"
// @Failure 400 {object} dto.BadRequestError "Bad request - validation errors"
// @Failure 401 {object} dto.UnauthorizedError "Unauthorized"
// @Failure 500 {object} dto.InternalServerError "Internal server error"
// @Security BearerAuth
// @Router /api/v1/user/match-action [post]
func (h *UserHandler) PostRecordMatchAction(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	var req dto.RecordMatchActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperrors.ErrBindingJSON, nil)
		return
	}

	resp, err := h.userUsecase.RecordMatchAction(authCtx.Ctx, req)
	if err != nil {
		if apperrors.ShouldLogError(err) {
			log.Error("Failed to record match action (POST)", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("Match action recorded successfully (POST)")
	apiresponse.Success(c, "Match action recorded successfully", resp)
}
