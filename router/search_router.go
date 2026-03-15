package router

import (
	"github.com/Marionvd/filia-project-backend/internal/handler"
	"github.com/gin-gonic/gin"
)

func handleSearchRoutes(r *gin.RouterGroup) {
	searchRouter := r.Group("/search")
	{
		searchRouter.GET("/users", handler.SearchUsersHandler)
		searchRouter.GET("/posts", handler.SearchPostsHandler)
		searchRouter.GET("/categories", handler.SearchCategoriesHandler)
	}
}
