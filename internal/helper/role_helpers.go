package helper

import (
	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/model"
)

func DoesRoleExist(roleName string) (bool, *model.Role) {
	var role model.Role
	return database.DbConnection.Where("id=$1", roleName).First(&role).Error != nil, &role
}
