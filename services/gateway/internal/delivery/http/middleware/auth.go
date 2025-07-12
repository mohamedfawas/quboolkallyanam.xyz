package middleware

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
)

func AuthMiddleware(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(constants.HeaderAuthorization)
		if authHeader == "" || !strings.HasPrefix(authHeader, constants.BearerTokenPrefix) {
			apiresponse.Fail(c, status.Errorf(codes.Unauthenticated, "missing or invalid Authorization header"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, constants.BearerTokenPrefix)
		userID, role, err := jwtManager.ExtractUserIDAndRole(token)
		if err != nil {
			apiresponse.Fail(c, status.Errorf(codes.Unauthenticated, "invalid token: %v", err))
			c.Abort()
			return
		}

		c.Set(constants.ContextKeyUserID, userID)
		c.Set(constants.ContextKeyRole, role)

		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(constants.ContextKeyRole)
		if !exists || val != role {
			apiresponse.Fail(c, status.Errorf(codes.PermissionDenied, "insufficient role"))
			c.Abort()
			return
		}
		c.Next()
	}
}
