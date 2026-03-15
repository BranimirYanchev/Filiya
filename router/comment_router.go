package router

import (
	"github.com/Marionvd/filia-project-backend/internal/handler"
	"github.com/Marionvd/filia-project-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func handleCommentRoutes(r *gin.RouterGroup) {
	cr := r.Group("/comments")
	{
		cr.GET("/:id", handler.GetComment)
	}
	cr.Use(middleware.JwtMiddleware)
	cr.Use(middleware.AuthenticationNecessary)
	//? Which permissions to authorize
	cr.Use(middleware.AuthorizePermissions("user:update"))
	{
		cr.PUT("/:id", handler.UpdateComment)
		cr.DELETE("/:id", handler.DeleteComment)
		cr.POST("/:id/like", middleware.AuthorizePermissions("user:update"), handler.LikeComment)
	}
}
