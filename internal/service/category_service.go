package service

import (
	"errors"
	"fmt"
	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"gorm.io/gorm"
)

func GetCategories(args ...int) (*[]model.Category, error) {
	var categories []model.Category
	if err := database.DbConnection.Limit(args[0]).Offset(args[1]).Preload("Parent").
		Preload("SubCategories").Where("parent_id IS NULL").Find(&categories).Error; err != nil {
		return nil, errors.New("internal: failed to receive categories")
	}

	return &categories, nil
}

func GetCategoryPosts(category *model.Category, limit, offset int) error {
	if err := database.DbConnection.Where("id=?", category.ID).First(&category).Error; err != nil {
		return errors.New("external: could not find category")
	}

	if err := database.DbConnection.Preload("Posts").
		Limit(limit).Offset(offset).Model(category).Error; err != nil {
		return errors.New("external: could not get posts or the specified category does not contain any posts")
	}

	return nil
}

func GetCategoryByID(id int) (*model.Category, error) {
	var category model.Category
	if err := database.DbConnection.Preload("Parent").
		Preload("SubCategories").Where("parent_id IS NULL").Where("id=?", id).Find(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func GetCategoryByName(name string) (*model.Category, error) {
	var category model.Category
	if err := database.DbConnection.Where("name LIKE ?", name).Preload("Parent").Preload("SubCategories").First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func SaveCategory(category *model.Category) error {
	//?Exists?
	var existing model.Category
	var parent model.Category

	if database.DbConnection.Where("name=?", category.Name).First(&existing).Error == nil {
		return errors.New("external: category with this name already exists")
	}

	if database.DbConnection.Where("parent_id=?", category.ParentID).First(&parent).Error != nil {
		return errors.New("external: specified parent does not exist")
	}

	if err := database.DbConnection.Save(&category).Error; err != nil {
		return errors.New("internal: " + err.Error())
	}
	return nil
}

func DeleteCategory(category *model.Category) error {
	if err := database.DbConnection.Preload("Posts").First(category, category.ID).Error; err != nil {
		return errors.New("external: category does not exist")
	}

	if database.DbConnection.Delete(&model.Category{}, "parent_id=?", category.ID).Error != nil {
		return errors.New("internal: failed to delete subcategories")
	}

	if database.DbConnection.Model(category).Association("Posts").Clear() != nil {
		return errors.New("internal: failed to clear related category posts")
	}

	if database.DbConnection.Where("category_id", category.ID).Delete(&model.Post{}).Error != nil {
		return errors.New("internal: failed to delete related posts")
	}

	return nil
}

func UpdateCategory(category, updatedCategory *model.Category) error {
	if err := database.DbConnection.First(category, category.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: could not find category")
		}
		return fmt.Errorf("internal: query failed|%w", err)
	}

	if err := database.DbConnection.Model(category).Updates(*updatedCategory).Error; err != nil {
		return fmt.Errorf("internal: could not process category|%w", err)
	}

	return nil
}
