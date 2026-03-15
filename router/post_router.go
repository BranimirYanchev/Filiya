package router

import (
	"github.com/Marionvd/filia-project-backend/internal/handler"
	"github.com/Marionvd/filia-project-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func handlePostRoutes(r *gin.RouterGroup) {
	postRouter := r.Group("/posts")
	{
		postRouter.GET("", handler.GetPosts)
		postRouter.GET("/:id", handler.GetPost)
		postRouter.GET("/:id/comments", handler.GetPostComments)
	}
	postRouter.Use(middleware.JwtMiddleware)
	postRouter.Use(middleware.AuthenticationNecessary)
	{
		postRouter.POST("", middleware.AuthorizePermissions("user:create"), handler.CreatePost)
		postRouter.POST("/:id/comments", middleware.AuthorizePermissions("user:create"), handler.CreateComment)
		postRouter.PUT("/:id", middleware.AuthorizePermissions("user:update"), handler.UpdatePost)
		postRouter.DELETE("/:id", middleware.AuthorizePermissions("user:delete"), handler.DeletePost)

		postRouter.POST("/:id/like", middleware.AuthorizePermissions("user:update"), handler.LikePost)
	}
}
