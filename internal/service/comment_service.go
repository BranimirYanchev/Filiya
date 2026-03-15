package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"gorm.io/gorm"
)

func GetComment(comment *model.Comment) error {
	if err := database.DbConnection.Preload("Author").Preload("Likes").First(comment, "id=?", comment.ID).Error; err != nil {
		return errors.New("external: unexisting comment")
	}
	return nil
}

func SaveComment(comment *model.Comment) error {
	if strings.TrimSpace(comment.Content) == "" {
		return errors.New("external: invalid data")
	}

	var post model.Post
	if err := database.DbConnection.First(&post, comment.PostID).Error; err != nil {
		return errors.New("external: unexisting post")
	}

	if err := database.DbConnection.Create(&post).Error; err != nil {
		return fmt.Errorf("internal: failed to save post %v", err)
	}

	return nil
}

func UpdateComment(comment, updatedComment *model.Comment) error {
	if strings.TrimSpace(updatedComment.Content) == "" {
		return errors.New("external: invalid data")
	}

	if err := database.DbConnection.First(comment, comment.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: could not find comment")
		}
		return fmt.Errorf("internal: query failed|%w", err)
	}
	if err := database.DbConnection.Model(comment).Updates(*updatedComment).Error; err != nil {
		return fmt.Errorf("internal: could not process comment|%w", err)
	}
	return nil
}

func DeleteComment(comment *model.Comment) error {
	if err := database.DbConnection.First(comment, comment.ID).Error; err != nil {
		return errors.New("external: comment does not exist")
	}
	if database.DbConnection.Model(comment).Association("Likes").Clear() != nil {
		return errors.New("internal: failed to clear related like info")
	}
	if database.DbConnection.Delete(comment).Error != nil {
		return errors.New("internal: failed to delete comment")
	}
	return nil
}

func LikeComment(comment *model.Comment, user *model.User) (*int64, error) {
	if err := database.DbConnection.First(comment, "id=?", comment.ID).Error; err != nil {
		return nil, errors.New("external: unexisting post")
	}

	if err := database.DbConnection.Model(comment).Association("Likes").Append(user); err != nil {
		return nil, errors.New("internal: failed to like post|" + err.Error())
	}

	likes := database.DbConnection.Model(comment).Association("Likes").Count()
	return &likes, nil
}

func GetPostComments(post *model.Post) (*[]model.Comment, error) {
	if err := database.DbConnection.First(post, "id=?", post.ID).Error; err != nil {
		return nil, errors.New("external: unexisting post")
	}

	var comments *[]model.Comment
	if err := database.DbConnection.Where("post_id=?", post.ID).Find(&comments).Error; err != nil {
		return nil, errors.New("external: failed to receive post comments|post does not contain any comments")
	}

	return comments, nil
}
