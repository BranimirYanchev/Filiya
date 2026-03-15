package model

import (
	"time"
)

type Comment struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Content   string    `json:"content"`
	AuthorID  uint64    `json:"author_id"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author"`
	PostID    uint64    `json:"post_id"`
	Likes     []*User   `gorm:"many2many:comment_likes" json:"likes"`
	CreatedAt time.Time `gorm:"default:NOW()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:NOW()" json:"updated_at"`
}
