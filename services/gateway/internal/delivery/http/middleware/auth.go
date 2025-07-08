package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/auth/jwt"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
)

func AuthMiddleware(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(constants.AuthHeaderKey)
		if authHeader == "" || !strings.HasPrefix(authHeader, constants.BearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, constants.BearerPrefix)
		userID, role, err := jwtManager.ExtractUserIDAndRole(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(constants.UserIDKey, userID)
		c.Set(constants.RoleKey, role)

		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(constants.RoleKey)
		if !exists || val != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
			return
		}
		c.Next()
	}
}
