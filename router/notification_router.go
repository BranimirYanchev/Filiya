package router

import (
	"github.com/Marionvd/filia-project-backend/internal/handler"
	"github.com/Marionvd/filia-project-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func handleNotificationRoutes(r *gin.RouterGroup) {
	notificationRouter := r.Group("/notifications")
	notificationRouter.Use(middleware.JwtMiddleware)
	notificationRouter.Use(middleware.AuthenticationNecessary)
	{
		notificationRouter.GET("/", handler.GetNotificationsHandler)
		notificationRouter.GET("/unread-count", handler.GetUnreadNotificationCountHandler)
		notificationRouter.PUT("/read-all", handler.MarkAllNotificationsAsReadHandler)
		notificationRouter.PUT("/:id/read", handler.MarkNotificationAsReadHandler)
		notificationRouter.DELETE("/:id", handler.DeleteNotificationHandler)
	}
}
