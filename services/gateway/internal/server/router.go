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
	s.registerPaymentRoutes(v1, router)
	s.registerChatRoutes(v1)
	s.registerUserRoutes(v1)

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

func (s *Server) registerPaymentRoutes(v1 *gin.RouterGroup, router *gin.Engine) {
	// ============= PAYMENT API ROUTES  =============
	apiPayment := v1.Group("/payment")
	{
		apiPayment.GET("/subscription-plans", s.paymentHandler.GetActiveSubscriptionPlans)
		apiPayment.GET("/subscription-plan", s.paymentHandler.GetSubscriptionPlanByID)
		apiPayment.POST("/subscription-plan",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleAdmin),
			s.paymentHandler.CreateSubscriptionPlan)
		apiPayment.PATCH("/subscription-plan",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleAdmin),
			s.paymentHandler.UpdateSubscriptionPlan)
		apiPayment.POST("/order",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.paymentHandler.CreatePaymentOrder)
	}

	// ============= PAYMENT WEB PAGES (Public) =============
	payment := router.Group("/payment")
	{

		payment.GET("/checkout", s.paymentHandler.ShowPaymentPage)
		payment.GET("/verify", s.paymentHandler.VerifyPayment)
		payment.GET("/success", s.paymentHandler.ShowSuccessPage)
		payment.GET("/failed", s.paymentHandler.ShowFailurePage)
	}
}

func (s *Server) registerUserRoutes(v1 *gin.RouterGroup) {
	user := v1.Group("/user")
	{
		user.PATCH("/profile",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.PatchUserProfile)
		user.PUT("/profile",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.PutUserProfile)
		user.POST("/preference",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.PostPartnerPreference)
		user.PATCH("/preference",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.PatchPartnerPreference)
		user.GET("/recommendations",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.GetMatchRecommendations)
		user.POST("/match-action",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.PostRecordMatchAction)
		user.PUT("/match-action",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.PutRecordMatchAction)
		user.GET("/matches/liked",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.GetLikedProfiles)
		user.GET("/matches/passed",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.GetPassedProfiles)
		user.GET("/matches/mutual",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RoleUser),
			s.userHandler.GetMutuallyMatchedProfiles)
	}
}

func (s *Server) registerChatRoutes(v1 *gin.RouterGroup) {
	chat := v1.Group("/chat")
	{
		chat.POST("/conversation",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RolePremiumUser),
			s.chatHandler.CreateConversation)
		chat.GET("/conversations",
			middleware.AuthMiddleware(s.jwtManager),
			middleware.RequireRole(constants.RolePremiumUser),
			s.chatHandler.GetUserConversations)
		chat.GET("/ws", s.chatHandler.HandleWebSocket)
	}
}
