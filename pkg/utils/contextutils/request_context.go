package contextutils

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
)

type RequestContextResult struct {
	Ctx context.Context
}

func ExtractRequestContext(c *gin.Context) (*RequestContextResult, error) {
	requestID, exists := c.Get(constants.ContextKeyRequestID)
	if !exists {
		return nil, fmt.Errorf("request ID context missing")
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyRequestID, requestID)

	return &RequestContextResult{
		Ctx: ctx,
	}, nil
}
