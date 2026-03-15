package service

import (
	"errors"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/model"
)

func GetRoles(limit, offset int) ([]model.Role, error) {
	var roles []model.Role
	if err := database.DbConnection.Limit(limit).Offset(offset).Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, errors.New("internal: failed to receive roles")
	}

	return roles, nil
}

func CreateRole(role *model.Role, perms ...string) error {
	if perms == nil {
		perms = []string{"user:read", "user:update", "user:delete"}
	}

	if err := database.DbConnection.First(&model.Role{}, "name=?", role.Name).Error; err == nil {
		return errors.New("external: role with this name already exists")
	}

	for _, perm := range perms {
		var permission model.Permission
		if err := database.DbConnection.First(&permission, "name=?", perm).Error; err == nil {
			return errors.New("external: specified permission doesn't exists")
		}
		role.Permissions = append(role.Permissions, &permission)
	}

	if err := database.DbConnection.Create(role).Error; err != nil {
		return errors.New("internal: couldn't persist role|" + err.Error())
	}

	return nil
}

func DeleteRole(role *model.Role) error {
	if err := database.DbConnection.First(role, "id=?", role.ID).Error; err != nil {
		return errors.New("external: unexisting role")
	}

	if err := database.DbConnection.Delete(role).Error; err != nil {
		return errors.New("internal: failed to delete role|" + err.Error())
	}

	return nil
}
