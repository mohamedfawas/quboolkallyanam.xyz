package middleware

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				apiresponse.Fail(c, fmt.Errorf("internal server error"))
			}
		}()

		c.Next()

		// Handle errors from handlers
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Log the error with context
			log.Printf("request error: %v", err)

			// Only send error response if not already sent
			if !c.Writer.Written() {
				apiresponse.Fail(c, err)
			}
		}
	}
}
