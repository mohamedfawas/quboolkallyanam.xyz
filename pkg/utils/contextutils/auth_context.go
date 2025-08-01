package contextutils

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
)

type AuthContextResult struct {
	Ctx    context.Context
}

func ExtractAuthContext(c *gin.Context) (*AuthContextResult, error) {
	requestID, exists := c.Get(constants.ContextKeyRequestID)
	if !exists {
		return nil, fmt.Errorf("request ID context missing")
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		return nil, fmt.Errorf("user ID context missing")
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)
	ctx = context.WithValue(ctx, constants.ContextKeyRequestID, requestID)

	return &AuthContextResult{
		Ctx: ctx,
	}, nil
}