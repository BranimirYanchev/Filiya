package service

import (
	"errors"
	"fmt"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"gorm.io/gorm"
)

// CreateNotification creates a new notification for a user
func CreateNotification(userID uint64, notifType model.NotificationType, title, message string, entityID *uint64, entityType string) error {
	notification := model.Notification{
		UserID:     userID,
		Type:       notifType,
		Title:      title,
		Message:    message,
		EntityID:   entityID,
		EntityType: entityType,
		Read:       false,
	}

	if err := database.DbConnection.Create(&notification).Error; err != nil {
		return fmt.Errorf("internal: failed to create notification|%w", err)
	}

	return nil
}

// GetUserNotifications retrieves all notifications for a user
func GetUserNotifications(userID uint64, limit, offset int, unreadOnly bool) ([]model.Notification, error) {
	var notifications []model.Notification

	query := database.DbConnection.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset)

	if unreadOnly {
		query = query.Where("read = ?", false)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("internal: failed to get notifications|%w", err)
	}

	return notifications, nil
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(notificationID, userID uint64) error {
	var notification model.Notification
	if err := database.DbConnection.First(&notification, "id = ? AND user_id = ?", notificationID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: notification not found")
		}
		return fmt.Errorf("internal: failed to find notification|%w", err)
	}

	if err := database.DbConnection.Model(&notification).Update("read", true).Error; err != nil {
		return fmt.Errorf("internal: failed to mark notification as read|%w", err)
	}

	return nil
}

// MarkAllNotificationsAsRead marks all notifications for a user as read
func MarkAllNotificationsAsRead(userID uint64) error {
	if err := database.DbConnection.Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true).Error; err != nil {
		return fmt.Errorf("internal: failed to mark all notifications as read|%w", err)
	}

	return nil
}

// DeleteNotification deletes a notification
func DeleteNotification(notificationID, userID uint64) error {
	var notification model.Notification
	if err := database.DbConnection.First(&notification, "id = ? AND user_id = ?", notificationID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: notification not found")
		}
		return fmt.Errorf("internal: failed to find notification|%w", err)
	}

	if err := database.DbConnection.Delete(&notification).Error; err != nil {
		return fmt.Errorf("internal: failed to delete notification|%w", err)
	}

	return nil
}

// GetUnreadNotificationCount returns the count of unread notifications for a user
func GetUnreadNotificationCount(userID uint64) (int64, error) {
	var count int64
	if err := database.DbConnection.Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("internal: failed to count notifications|%w", err)
	}

	return count, nil
}
