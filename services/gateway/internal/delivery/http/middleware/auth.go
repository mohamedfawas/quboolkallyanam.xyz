package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
)

func AuthMiddleware(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(constants.HeaderAuthorization)
		if authHeader == "" || !strings.HasPrefix(authHeader, constants.BearerTokenPrefix) {
			apiresponse.Error(c, apperrors.ErrUnauthorized, nil)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, constants.BearerTokenPrefix)
		userID, role, err := jwtManager.ExtractUserIDAndRole(token)
		if err != nil {
			apiresponse.Error(c, apperrors.ErrUnauthorized, nil)
			c.Abort()
			return
		}

		c.Set(constants.ContextKeyUserID, userID)
		c.Set(constants.ContextKeyRole, role)
		c.Set(constants.ContextKeyAccessToken, token)

		c.Next()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        val, exists := c.Get(constants.ContextKeyRole)
        if !exists {
            apiresponse.Error(c, apperrors.ErrForbidden, nil)
            c.Abort()
            return
        }

        userRole, ok := val.(string)
        if !ok {
            apiresponse.Error(c, apperrors.ErrForbidden, nil)
            c.Abort()
            return
        }

        if !hasRequiredRole(userRole, requiredRole) {
            apiresponse.Error(c, apperrors.ErrForbidden, nil)
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
