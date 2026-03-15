package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"gorm.io/gorm"
)

func assignPostCategoryIDs(post *model.Post) {
	post.CategoryIDs = make([]uint64, 0, len(post.Categories))
	for _, category := range post.Categories {
		post.CategoryIDs = append(post.CategoryIDs, category.ID)
	}
}

func GetPosts(limit, offset int, user *model.User) ([]model.Post, error) {
	var posts []model.Post

	if user != nil && user.ID > 0 {
		// Load user with view history
		if err := database.DbConnection.First(user, "id=?", user.ID).Error; err != nil {
			return nil, errors.New("external: unexisting user")
		}

		// Get posts user has viewed
		var viewedPosts []model.Post
		if err := database.DbConnection.Model(user).Association("PostHistory").Find(&viewedPosts); err != nil {
			return nil, errors.New("internal: failed to load user history")
		}

		var categoryIDs []uint64
		if len(viewedPosts) > 0 {
			// Extract unique category IDs from viewed posts
			categoryMap := make(map[uint64]bool)
			for _, post := range viewedPosts {
				// Load categories for this post
				var postWithCats model.Post
				if err := database.DbConnection.Preload("Categories").First(&postWithCats, "id=?", post.ID).Error; err == nil {
					for _, cat := range postWithCats.Categories {
						categoryMap[cat.ID] = true
					}
				}
			}

			// Convert map to slice
			for catID := range categoryMap {
				categoryIDs = append(categoryIDs, catID)
			}
		}

		// If user has no viewing history or no categories found, get some categories
		if len(categoryIDs) == 0 {
			var randomCategories []model.Category
			if err := database.DbConnection.Limit(5).Find(&randomCategories).Error; err == nil && len(randomCategories) > 0 {
				for _, cat := range randomCategories {
					categoryIDs = append(categoryIDs, cat.ID)
				}
			}
		}

		// If we have preferred categories, prioritize posts in those categories
		if len(categoryIDs) > 0 {
			// Get posts in user's preferred categories
			if err := database.DbConnection.Limit(limit).Offset(offset).
				Preload("Author").Preload("Categories").Preload("Likes").Preload("Attachments").
				Joins("JOIN category_posts ON posts.id = category_posts.post_id").
				Where("category_posts.category_id IN ? AND posts.is_private = ?", categoryIDs, false).
				Group("posts.id").
				Order("posts.created_at DESC").
				Find(&posts).Error; err != nil {
				return nil, errors.New("internal: failed to receive posts")
			}

			// If we got fewer posts than requested, fill with random posts
			if len(posts) < limit {
				var additionalPosts []model.Post
				var excludeIDs []uint64
				for _, p := range posts {
					excludeIDs = append(excludeIDs, p.ID)
				}

				remainingLimit := limit - len(posts)
				query := database.DbConnection.Limit(remainingLimit).Offset(0).
					Preload("Author").Preload("Categories").Preload("Likes").Preload("Attachments").
					Where("is_private = ?", false)

				if len(excludeIDs) > 0 {
					query = query.Where("id NOT IN ?", excludeIDs)
				}

				if err := query.Order("created_at DESC").Find(&additionalPosts).Error; err == nil {
					posts = append(posts, additionalPosts...)
				}
			}
		} else {
			// No categories found, just get all public posts
			if err := database.DbConnection.Limit(limit).Offset(offset).
				Preload("Author").Preload("Categories").Preload("Likes").Preload("Attachments").
				Where("is_private = ?", false).
				Order("created_at DESC").
				Find(&posts).Error; err != nil {
				return nil, errors.New("internal: failed to receive posts")
			}
		}
	} else {
		// Anonymous user - return all public posts
		if err := database.DbConnection.Limit(limit).Offset(offset).
			Preload("Author").Preload("Categories").Preload("Likes").Preload("Attachments").
			Where("is_private = ?", false).
			Order("created_at DESC").
			Find(&posts).Error; err != nil {
			return nil, errors.New("internal: failed to receive posts")
		}
	}

	for index := range posts {
		assignPostCategoryIDs(&posts[index])
	}

	return posts, nil
}

func SavePost(post *model.Post) error {
	if strings.TrimSpace(post.Title) == "" || strings.TrimSpace(post.Content) == "" {
		return errors.New("external: invalid data|post title and content must be provided")
	}
	if len(post.CategoryIDs) == 0 {
		return errors.New("external: post must fall in at least one category")
	}

	var count int64
	if err := database.DbConnection.
		Model(&model.Category{}).
		Where("id IN ?", post.CategoryIDs).
		Count(&count).Error; err != nil {
		return errors.New("external: failed to validate categories|" + err.Error())
	}

	if count != int64(len(post.CategoryIDs)) {
		return errors.New("external: unexisting category ids")
	}

	var categories []model.Category
	if err := database.DbConnection.Where("id IN ?", post.CategoryIDs).Find(&categories).Error; err != nil {
		return errors.New("internal: failed to load categories|" + err.Error())
	}

	if err := database.DbConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}

		if len(categories) > 0 {
			if err := tx.Model(post).Association("Categories").Append(categories); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return errors.New("internal: failed to persist post|" + err.Error())
	}

	return nil
}

func GetPost(post *model.Post) error {
	if err := database.DbConnection.Preload("Author").
		Preload("Likes").Preload("Comments").Preload("Categories").Preload("Attachments").
		First(post, "id=?", post.ID).Error; err != nil {
		return errors.New("external: unexisting post")
	}
	assignPostCategoryIDs(post)
	return nil
}

// TrackPostView records that a user has viewed a post
func TrackPostView(user *model.User, post *model.Post) error {
	if user == nil || user.ID == 0 {
		// Anonymous user - no tracking
		return nil
	}

	if post == nil || post.ID == 0 {
		return errors.New("external: invalid post")
	}

	// Check if user already viewed this post (to avoid duplicates)
	var count int64
	if err := database.DbConnection.Table("post_view_history").
		Where("user_id = ? AND post_id = ?", user.ID, post.ID).Count(&count).Error; err != nil {
		// If error checking, still try to add
	} else if count > 0 {
		// Already viewed, skip
		return nil
	}

	// Add to view history
	if err := database.DbConnection.Model(user).Association("PostHistory").Append(post); err != nil {
		return fmt.Errorf("internal: failed to track post view|%w", err)
	}

	return nil
}

func LikePost(post *model.Post, user *model.User) (*int64, error) {
	if err := database.DbConnection.First(post, "id=?", post.ID).Error; err != nil {
		return nil, errors.New("external: unexisting post")
	}

	if err := database.DbConnection.Model(post).Association("Likes").Append(user); err != nil {
		return nil, errors.New("internal: failed to like post|" + err.Error())
	}

	likes := database.DbConnection.Model(post).Association("Likes").Count()
	return &likes, nil
}

func DeletePost(post *model.Post) error {
	if err := database.DbConnection.First(post, post.ID).Error; err != nil {
		return errors.New("external: post does not exist")
	}
	if database.DbConnection.Model(post).Association("Categories").Clear() != nil {
		return errors.New("internal: failed to clear related categories")
	}
	if database.DbConnection.Model(post).Association("Likes").Clear() != nil {
		return errors.New("internal: failed to clear related like info")
	}
	if database.DbConnection.Delete(post).Error != nil {
		return errors.New("internal: failed to delete related posts")
	}
	return nil
}

func UpdatePost(post *model.Post, updated *model.Post) error {
	if err := database.DbConnection.First(post, post.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: could not find post")
		}
		return fmt.Errorf("internal: query failed|%w", err)
	}
	if err := database.DbConnection.Model(post).Updates(*updated).Error; err != nil {
		return fmt.Errorf("internal: could not process post|%w", err)
	}
	return nil
}
