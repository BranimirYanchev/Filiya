package model

type Role struct {
	ID          uint64        `gorm:"primaryKey" json:"id"`
	Name        string        `gorm:"uniqueIndex;size:40" json:"name"`
	Permissions []*Permission `gorm:"many2many:role_permissions" json:"permissions"`
	Users       []*User       `gorm:"many2many:user_roles" json:"users"`
}

type SafeRole struct {
	ID          uint64        `gorm:"primaryKey" json:"id"`
	Name        string        `gorm:"size:40" json:"name"`
	Permissions []*Permission `json:"permissions"`
}

func (role *Role) GenerateSafeRole() SafeRole {
	return SafeRole{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: role.Permissions,
	}
}
