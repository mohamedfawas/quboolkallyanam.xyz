package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) setupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	s.registerAuthRoutes(v1)

}

func (s *Server) registerAuthRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	{
		userAuth := auth.Group("/user")
		{
			userAuth.POST("/register", s.authHandler.UserRegister)
			userAuth.POST("/verify", s.authHandler.UserVerification)
			userAuth.POST("/login", s.authHandler.UserLogin)
			userAuth.POST("/logout", s.authHandler.UserLogout)
			// TODO: Add more user auth routes as needed
			// userAuth.POST("/refresh", s.authHandler.RefreshToken)
			// userAuth.POST("/forgot-password", s.authHandler.ForgotPassword)
			// userAuth.POST("/reset-password", s.authHandler.ResetPassword)
		}

		adminAuth := auth.Group("/admin")
		{
			adminAuth.POST("/login", s.authHandler.AdminLogin)
			adminAuth.POST("/logout", s.authHandler.AdminLogout)
			// TODO: Add more admin auth routes as needed
			// adminAuth.POST("/refresh", s.authHandler.AdminRefreshToken)
		}
	}
}

func (s *Server) registerUserRoutes(v1 *gin.RouterGroup)    {}
func (s *Server) registerAdminRoutes(v1 *gin.RouterGroup)   {}
func (s *Server) registerPaymentRoutes(v1 *gin.RouterGroup) {}
func (s *Server) registerChatRoutes(v1 *gin.RouterGroup)    {}
