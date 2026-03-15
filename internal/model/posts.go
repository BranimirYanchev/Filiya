package model

import "time"

type CreatePost struct {
	Title         string   `json:"title" example:"My First Post"`
	Content       string   `json:"content" example:"This is the content of my first post."`
	CategoryIDs   []uint64 `json:"category_ids" example:"1"`
	Tags          []string `json:"tags" example:"tag1,tag2"`
	TaggedUsersID []int64  `json:"tagged_users" example:"2,3"`
	IsPrivate     *bool    `json:"is_private" example:"false"`
}

type PostAttachment struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	PostID    uint64    `gorm:"index;not null" json:"post_id"`
	FileName  string    `gorm:"size:255;not null" json:"file_name"`
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	FileURL   string    `gorm:"size:500;not null" json:"file_url"`
	MimeType  string    `gorm:"size:120" json:"mime_type"`
	Size      int64     `json:"size"`
	Kind      string    `gorm:"size:40" json:"kind"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID            uint64           `gorm:"primaryKey" json:"id" example:"1"`
	Title         string           `gorm:"unique;not null;size:255" json:"title" example:"My First Post"`
	Content       string           `gorm:"not null" json:"content" example:"This is the content of my first post."`
	CategoryIDs   []uint64         `gorm:"-" json:"category_ids"`
	Categories    []Category       `gorm:"constraint:OnDelete:CASCADE;many2many:category_posts" json:"categories"`
	AuthorID      uint64           `gorm:"not null" json:"author_id" example:"1"`
	Author        User             `gorm:"foreignKey:AuthorID" json:"author"`
	Tags          []string         `gorm:"type:text[]" json:"tags" example:"tag1,tag2"`
	TaggedUsersID []int64          `gorm:"type:bigint[]" json:"tagged_users" example:"2,3"`
	Likes         []User           `gorm:"many2many:post_likes" json:"likes"`
	Comments      []Comment        `gorm:"foreignKey:PostID" json:"comments"`
	Attachments   []PostAttachment `gorm:"foreignKey:PostID" json:"attachments"`
	IsPrivate     bool             `gorm:"default:false" json:"is_private"`
	ViewHistory   []User           `gorm:"many2many:post_view_history" json:"post_history"`
	CreatedAt     time.Time        `gorm:"default:NOW()" json:"created_at" example:"2023-08-17T10:00:00Z"`
	UpdatedAt     time.Time        `gorm:"default:NOW()" json:"updated_at" example:"2023-08-17T10:00:00Z"`
	DeletedAt     time.Time        `json:"deleted_at"`
}

type UpdatePost struct {
	Title         string   `example:"My First Post Updated"`
	Content       string   `example:"This is the content of my first post updated."`
	Tags          []string `json:"tags" example:"tag1,tag2"`
	TaggedUsersID []int64  `json:"tagged_users" example:"2,3"`
	CategoryIDs   []uint64 `json:"category_ids" example:"2,3"`
	IsPrivate     *bool    `json:"is_private" example:"false"`
}
