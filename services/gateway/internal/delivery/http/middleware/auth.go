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
		c.Set(constants.ContextKeyAccessToken, token)

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
		userRole := val.(string)
		if !hasRequiredRole(userRole, role) {
			apiresponse.Fail(c, status.Errorf(codes.PermissionDenied, "insufficient role"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func hasRequiredRole(userRole, requiredRole string) bool {
	switch requiredRole {
	case constants.RoleUser:
		return userRole == constants.RoleUser || userRole == constants.RolePremiumUser
	case constants.RolePremiumUser:
		return userRole == constants.RolePremiumUser
	case constants.RoleAdmin:
		return userRole == constants.RoleAdmin
	default:
		return false
	}
}
