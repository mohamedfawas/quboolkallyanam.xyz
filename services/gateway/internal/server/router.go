package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/mohamedfawas/quboolkallyanam.xyz/docs/swagger" // Swagger generated files
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/delivery/http/middleware"
)

func (s *Server) setupRoutes(router *gin.Engine) {
	router.Use(middleware.ErrorHandler())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
			userAuth.POST("/logout", middleware.AuthMiddleware(s.jwtManager),
				middleware.RequireRole(constants.RoleUser),
				s.authHandler.UserLogout)
			userAuth.POST("/delete", middleware.AuthMiddleware(s.jwtManager),
				middleware.RequireRole(constants.RoleUser),
				s.authHandler.UserDelete)
			userAuth.POST("/refresh", s.authHandler.RefreshToken)
			// TODO: Add more user auth routes as needed
			// userAuth.POST("/forgot-password", s.authHandler.ForgotPassword)
			// userAuth.POST("/reset-password", s.authHandler.ResetPassword)
		}

		adminAuth := auth.Group("/admin")
		{
			adminAuth.POST("/login", s.authHandler.AdminLogin)
			adminAuth.POST("/logout", middleware.AuthMiddleware(s.jwtManager),
				middleware.RequireRole(constants.RoleAdmin),
				s.authHandler.AdminLogout)
			// TODO: Add more admin auth routes as needed
			// adminAuth.POST("/refresh", s.authHandler.AdminRefreshToken)
		}
		auth.POST("/refresh", middleware.AuthMiddleware(s.jwtManager),
			s.authHandler.RefreshToken)
	}
}

func (s *Server) registerUserRoutes(v1 *gin.RouterGroup)    {}
func (s *Server) registerAdminRoutes(v1 *gin.RouterGroup)   {}
func (s *Server) registerPaymentRoutes(v1 *gin.RouterGroup) {}
func (s *Server) registerChatRoutes(v1 *gin.RouterGroup)    {}
