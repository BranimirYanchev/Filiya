// Package routes handles the configuration and setup of all API routes in the Filia application.
// It defines the routing structure and connects routes to their respective controller functions.
package router

import (
	_ "github.com/Marionvd/filia-project-backend/docs"
	"github.com/Marionvd/filia-project-backend/internal/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetUpRoutes configures all the routes for the Filia API.
// It initializes the database connection, performs auto-migration for the User model,
// and sets up the API endpoints with their corresponding controller functions.
//
// Parameters:
//   - r: The Gin engine instance to which routes will be added.
//
// The function creates a routes router group under the "/api" path and defines
// several endpoints for user-related operations.
func SetUpRoutes(r *gin.Engine) {
	mainRouter := r.Group("/api")
	handleAuthenticationRoutes(mainRouter)
	handleUserRoutes(mainRouter)
	handleRoleRoutes(mainRouter)
	handlePostRoutes(mainRouter)
	handleCategoryRoutes(mainRouter)
	handleCommentRoutes(mainRouter)
	handleSearchRoutes(mainRouter)
	handleNotificationRoutes(mainRouter)

	{
		mainRouter.GET("/", handler.GetHome)
		mainRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
