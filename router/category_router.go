package router

import (
	"github.com/Marionvd/filia-project-backend/internal/handler"
	"github.com/Marionvd/filia-project-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func handleCategoryRoutes(group *gin.RouterGroup) {
	areaRouter := group.Group("/categories")
	{
		areaRouter.GET("", handler.GetCategories)
		areaRouter.GET("/:id", handler.GetCategory)
		areaRouter.GET("/:id/posts", handler.GetCategoryPosts)
	}
	areaRouter.Use(middleware.JwtMiddleware)
	areaRouter.Use(middleware.AuthenticationNecessary)
	{
		areaRouter.GET("/name", handler.FindCategoryByName)
		areaRouter.POST("", middleware.AuthorizePermissions("moderator:create"), handler.CreateCategory)
		areaRouter.PUT("/:id", middleware.AuthorizePermissions("moderator:update"), handler.UpdateCategoryHandler)
		areaRouter.DELETE("/:id", middleware.AuthorizePermissions("moderator:delete"), handler.DeleteCategory)
	}

}
