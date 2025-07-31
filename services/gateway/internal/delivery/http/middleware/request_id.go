package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
)


func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(constants.HeaderRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set in Gin context
		c.Set(constants.ContextKeyRequestID, requestID)

		// Set in response headers so client can trace it
		c.Writer.Header().Set(constants.HeaderRequestID, requestID)

		c.Next()
	}
}