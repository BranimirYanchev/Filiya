package model

import "time"

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeFriendRequest NotificationType = "friend_request"
	NotificationTypeFriendAccept  NotificationType = "friend_accept"
	NotificationTypePostLike      NotificationType = "post_like"
	NotificationTypeComment       NotificationType = "comment"
	NotificationTypeCommentLike   NotificationType = "comment_like"
	NotificationTypePostMention   NotificationType = "post_mention"
)

// Notification represents a notification in the system
type Notification struct {
	ID         uint64           `gorm:"primaryKey" json:"id"`
	UserID     uint64           `gorm:"not null;index" json:"user_id"`
	User       *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type       NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	Title      string           `gorm:"size:255;not null" json:"title"`
	Message    string           `gorm:"size:500" json:"message"`
	Read       bool             `gorm:"default:false;index" json:"read"`
	EntityID   *uint64          `gorm:"index" json:"entity_id,omitempty"`     // ID of the related entity (post, comment, friend request, etc.)
	EntityType string           `gorm:"size:50" json:"entity_type,omitempty"` // Type of entity (post, comment, friend_request)
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
