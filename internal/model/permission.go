package model

type Permission struct {
	ID          uint64  `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"unique;size:40" json:"name"`
	Description string  `json:"description"`
	Roles       []*Role `gorm:"many2many:role_permissions" json:"roles"`
}
