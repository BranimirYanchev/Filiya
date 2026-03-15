package router

import (
	"github.com/Marionvd/filia-project-backend/internal/handler"
	"github.com/Marionvd/filia-project-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func handleRoleRoutes(group *gin.RouterGroup) {
	roleRouter := group.Group("/roles")
	roleRouter.Use(middleware.JwtMiddleware)
	roleRouter.Use(middleware.AuthenticationNecessary)
	{
		roleRouter.GET("", middleware.AuthorizePermissions("moderator:read"), handler.RetrieveRoles)
		roleRouter.POST("", middleware.AuthorizePermissions("moderator:update"), handler.CreateRole)
		roleRouter.GET("/me", middleware.AuthorizePermissions("user:read"), handler.RetrieveUserRole)
		roleRouter.DELETE("/:id", middleware.AuthorizePermissions("admin:access"), handler.DeleteRole)
	}
}
