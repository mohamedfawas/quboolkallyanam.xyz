package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Recover from panics
		defer func() {
			if rec := recover(); rec != nil {
				apiresponse.Fail(c, fmt.Errorf("%v", rec))
			}
		}()

		c.Next()

		// If any handler called c.Error(err), catch it here
		if len(c.Errors) > 0 {
			apiresponse.Fail(c, c.Errors.Last().Err)
		}
	}
}
