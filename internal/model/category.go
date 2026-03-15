package model

import "time"

type Category struct {
	ID            uint64      `gorm:"primaryKey" json:"id"`
	Name          string      `gorm:"size:40;uniqueIndex;not null" json:"name"`
	ParentID      *uint64     `json:"parent_id"`
	Parent        *Category   `gorm:"foreignKey:ParentID;references:ID" json:"parent"`
	SubCategories []*Category `gorm:"foreignKey:ParentID;references:ID" json:"sub_categories"`
	Posts         []*Post     `gorm:"many2many:category_posts" json:"posts"`
	DeletedAt     time.Time   `json:"deleted_at"`
}

type CreateCategoryInput struct {
	Name     string  `json:"name"`
	ParentID *uint64 `json:"parent_id"`
}

type UpdateCategoryInput struct {
	Name     string  `json:"name"`
	ParentID *uint64 `json:"parent_id"`
}
