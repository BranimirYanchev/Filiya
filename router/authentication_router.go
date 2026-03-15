package router

import (
	"github.com/Marionvd/filia-project-backend/internal/handler"
	"github.com/Marionvd/filia-project-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

	func handleAuthenticationRoutes(g *gin.RouterGroup) {
		router := g.Group("/auth")
		{
			router.GET("/", middleware.JwtMiddleware, middleware.AuthenticationUnnecessary, handler.AuthenticationHome)
			router.POST("/login", middleware.JwtMiddleware, middleware.AuthenticationUnnecessary, handler.EmailLogin)
			router.POST("/register", middleware.JwtMiddleware, middleware.AuthenticationUnnecessary, handler.Register)
			router.POST("/logout", middleware.JwtMiddleware, handler.Logout)
			router.POST("/refresh", handler.RefreshToken)
			router.POST("/reset-password-request", handler.PasswordResetRequest)
			router.POST("/reset-password", handler.PasswordReset)
			router.GET("/google/login", middleware.JwtMiddleware, middleware.AuthenticationUnnecessary, handler.GoogleLogin)
			router.GET("/google/callback", middleware.JwtMiddleware, middleware.AuthenticationUnnecessary, handler.GoogleCallbackHandler)
		}
	}
