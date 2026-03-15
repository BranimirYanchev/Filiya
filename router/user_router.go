package router

import (
	"github.com/Marionvd/filia-project-backend/internal/handler"
	"github.com/Marionvd/filia-project-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func handleUserRoutes(r *gin.RouterGroup) {
	userRouter := r.Group("/users")
	//* public endpoint
	{
		userRouter.GET("/", handler.GetUsers)
		userRouter.GET("/public/:id", handler.GetUser)

		userRouter.GET("/:id/posts", handler.GetUserPosts)
	}
	userRouter.Use(middleware.JwtMiddleware)
	userRouter.Use(middleware.AuthenticationNecessary)
	{
		userRouter.GET("/me", handler.GetLoggedUser)
		userRouter.DELETE("/me", middleware.AuthorizePermissions("user:delete"), handler.DeleteAccountHandler)
		userRouter.PUT("/profile/sensitive", middleware.AuthorizePermissions("user:update"), handler.UserUpdateSensitiveHandler)
		userRouter.PUT("/profile/basic", middleware.AuthorizePermissions("user:update"), handler.UserUpdateInsensitiveHandler)
		userRouter.POST("/profile/avatar", middleware.AuthorizePermissions("user:update"), handler.UpdateAvatarHandler)

		userRouter.PUT("/:id/role", middleware.AuthorizePermissions("moderator:update"), handler.ChangeUserRole)
	}
	{
		userRouter.GET("/friends", middleware.AuthorizePermissions("user:update"), handler.GetUserFriendsHandler)
		userRouter.DELETE("/friends/:id", middleware.AuthorizePermissions("user:update"), handler.UnfriendHandler)
		userRouter.POST("/friends/requests", middleware.AuthorizePermissions("user:update"), handler.SendFriendRequestHandler)
		userRouter.GET("/friends/requests/pending", middleware.AuthorizePermissions("user:update"), handler.GetUserPendingFriendRequestsHandler)
		userRouter.GET("/friends/requests/sent", middleware.AuthorizePermissions("user:update"), handler.GetUserSentFriendRequestsHandler)
		userRouter.POST("/friends/requests/accept/:id", middleware.AuthorizePermissions("user:update"), handler.AcceptFriendRequestHandler)
		userRouter.POST("/friends/requests/cancel/:id", middleware.AuthorizePermissions("user:update"), handler.DeclineFriendRequest)
		userRouter.DELETE("/friends/requests/cancel/:id", middleware.AuthorizePermissions("user:update"), handler.DeleteFriendRequest)
	}

}
