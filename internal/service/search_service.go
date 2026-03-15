package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/model"
)

// SearchUsers searches for users by email, full name, or bio
func SearchUsers(query string, limit, offset int) ([]model.User, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("external: search query cannot be empty")
	}

	var users []model.User
	searchPattern := "%" + strings.ToLower(query) + "%"

	if err := database.DbConnection.
		Limit(limit).Offset(offset).
		Where("LOWER(email) LIKE ? OR LOWER(full_name) LIKE ? OR LOWER(COALESCE(bio, '')) LIKE ?",
			searchPattern, searchPattern, searchPattern).
		Preload("Role").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("internal: failed to search users|%w", err)
	}

	return users, nil
}

// SearchPosts searches for posts by title or content
func SearchPosts(query string, limit, offset int, user *model.User) ([]model.Post, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("external: search query cannot be empty")
	}

	var posts []model.Post
	searchPattern := "%" + strings.ToLower(query) + "%"

	queryBuilder := database.DbConnection.
		Limit(limit).Offset(offset).
		Where("(LOWER(title) LIKE ? OR LOWER(content) LIKE ?) AND is_private = ?",
			searchPattern, searchPattern, false).
		Preload("Author").Preload("Categories").Preload("Likes").
		Order("created_at DESC")

	// If user is authenticated, also show their private posts
	if user != nil && user.ID > 0 {
		queryBuilder = database.DbConnection.
			Limit(limit).Offset(offset).
			Where("(LOWER(title) LIKE ? OR LOWER(content) LIKE ?) AND (is_private = ? OR author_id = ?)",
				searchPattern, searchPattern, false, user.ID).
			Preload("Author").Preload("Categories").Preload("Likes").
			Order("created_at DESC")
	}

	if err := queryBuilder.Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("internal: failed to search posts|%w", err)
	}

	return posts, nil
}

// SearchCategories searches for categories by name
func SearchCategories(query string, limit, offset int) ([]model.Category, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("external: search query cannot be empty")
	}

	var categories []model.Category
	searchPattern := "%" + strings.ToLower(query) + "%"

	if err := database.DbConnection.
		Limit(limit).Offset(offset).
		Where("LOWER(name) LIKE ?", searchPattern).
		Preload("Parent").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("internal: failed to search categories|%w", err)
	}

	return categories, nil
}
