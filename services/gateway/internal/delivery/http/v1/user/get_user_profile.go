package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/errutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"go.uber.org/zap"
)

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	authCtx, err := contextutils.ExtractAuthContext(c)
	if err != nil {
		apiresponse.Error(c, err, nil)
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, authCtx.Ctx.Value(constants.ContextKeyRequestID).(string)),
		zap.String(constants.ContextKeyUserID, authCtx.Ctx.Value(constants.ContextKeyUserID).(string)),
	)

	profile, err := h.userUsecase.GetUserProfile(authCtx.Ctx)
	if err != nil {
		if apperrors.ShouldLogError(err) && !errutil.IsGRPCError(err) {
			log.Error("Failed to get user profile", zap.Error(err))
		}
		apiresponse.Error(c, err, nil)
		return
	}

	log.Info("User profile fetched successfully")
	apiresponse.Success(c, "User profile fetched successfully", profile)
}