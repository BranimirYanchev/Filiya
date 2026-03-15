package handler

import (
	"net/http"
	"strconv"

	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// GetNotificationsHandler retrieves notifications for the authenticated user
// @Summary Get user notifications
// @Description Returns all notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of notifications to return (default: 20, max: 100)" default(20)
// @Param offset query int false "Number of notifications to skip (default: 0)" default(0)
// @Param unread_only query bool false "Return only unread notifications" default(false)
// @Success 200 {object} map[string]interface{} "List of notifications"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/notifications [get]
func GetNotificationsHandler(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, ok := userObj.(model.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	unreadOnly := c.Query("unread_only") == "true"

	notifications, err := service.GetUserNotifications(user.ID, limit, offset, unreadOnly)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notifications,
	})
}

// MarkNotificationAsReadHandler marks a notification as read
// @Summary Mark notification as read
// @Description Marks a specific notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Notification ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid notification ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/notifications/{id}/read [put]
func MarkNotificationAsReadHandler(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid notification ID",
		})
		return
	}

	userObj, _ := c.Get("user")
	user, ok := userObj.(model.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if err := service.MarkNotificationAsRead(notificationID, user.ID); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification marked as read",
	})
}

// MarkAllNotificationsAsReadHandler marks all notifications as read
// @Summary Mark all notifications as read
// @Description Marks all notifications for the authenticated user as read
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Success message"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/notifications/read-all [put]
func MarkAllNotificationsAsReadHandler(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, ok := userObj.(model.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if err := service.MarkAllNotificationsAsRead(user.ID); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All notifications marked as read",
	})
}

// DeleteNotificationHandler deletes a notification
// @Summary Delete notification
// @Description Deletes a specific notification
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Notification ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid notification ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/notifications/{id} [delete]
func DeleteNotificationHandler(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid notification ID",
		})
		return
	}

	userObj, _ := c.Get("user")
	user, ok := userObj.(model.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if err := service.DeleteNotification(notificationID, user.ID); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification deleted",
	})
}

// GetUnreadNotificationCountHandler returns the count of unread notifications
// @Summary Get unread notification count
// @Description Returns the count of unread notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Unread notification count"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/notifications/unread-count [get]
func GetUnreadNotificationCountHandler(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, ok := userObj.(model.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	count, err := service.GetUnreadNotificationCount(user.ID)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

