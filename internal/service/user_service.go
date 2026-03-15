package service

import (
	"errors"
	"strings"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/helper"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/username"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func assignBasicRole(user *model.User) {
	var roleUser model.Role

	if err := database.DbConnection.Where("name=?", "RoleUser").First(&roleUser).Error; err != nil {
		log.Error("failed to retrieve basic role")
		return
	}

	user.RoleID = roleUser.ID
}

func SaveRegisterUser(user *model.User) error {
	assignBasicRole(user)

	if strings.TrimSpace(user.Email) == "" || user.FullName == nil || user.Password == nil {
		return errors.New("external: unprovided fields")
	}

	if err := database.DbConnection.Where("email=?", user.Email).First(&model.User{}).Error; err == nil {
		log.Error(err)
		return errors.New("external: failed to process user|user with this email already exists")
	}

	if !helper.ValidateEmailRegex(user.Email) {
		return errors.New("external: invalid email")
	}

	fnError := helper.ValidateFullName(*user.FullName)
	if fnError != "" {
		return errors.New("external: invalid full name|" + fnError)
	}

	usernameValue, err := username.Generate(*user.FullName, func(candidate string) (bool, error) {
		var count int64
		if err := database.DbConnection.Model(&model.User{}).Where("username = ?", candidate).Count(&count).Error; err != nil {
			return false, err
		}
		return count > 0, nil
	})
	if err != nil {
		log.Error(err)
		return errors.New("internal: could not generate username")
	}
	user.Username = database.Strptr(usernameValue)

	pwdErr := helper.ValidatePassword(*user.Password)
	if pwdErr != "" {
		return errors.New("external: invalid password" + pwdErr)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Debug(err)
		return errors.New("internal: hash error")
	}
	user.Password = database.Strptr(string(hashedPassword))

	if err := database.DbConnection.Save(user).Error; err != nil {
		log.Debugf(err.Error())
		return errors.New("internal: could not persist user")
	}
	return nil
}

func Login(body model.PasswordLoginInput) (*model.User, error) {
	if strings.TrimSpace(body.Email) == "" || strings.TrimSpace(body.Password) == "" {
		return nil, errors.New("external: invalid data")
	}

	var targetUser model.User
	if err := database.DbConnection.Where("email=?", body.Email).First(&targetUser).Error; err != nil {
		return nil, errors.New("external: unexisting user")
	}

	if targetUser.Password == nil {
		return nil, errors.New("external: incorrect password|user has not set a password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*targetUser.Password), []byte(body.Password)); err != nil {
		return nil, errors.New("external: incorrect password")
	}

	return &targetUser, nil
}

func UpdateUserSensitiveInformation(prev, after *model.User) error {
	if prev.Password == nil {
		return errors.New("external: invalid password|cannot change user information|user has not set a password")
	}

	if strings.TrimSpace(after.Email) != "" {
		if prev.Email != after.Email {
			if !helper.ValidateEmailRegex(after.Email) {
				return errors.New("external: invalid email")
			}

			if err := database.DbConnection.Model(prev).Update("Email", after.Email).Error; err != nil {
				log.Error(err)
				return errors.New("internal: could not process email")
			}
		}
	}

	if after.FullName != nil && *after.FullName != "" {
		if fnError := helper.ValidateFullName(*after.FullName); fnError != "" {
			return errors.New("external: invalid name" + fnError)
		}

		if err := database.DbConnection.Model(prev).Update("FullName", after.FullName).Error; err != nil {
			log.Error(err)
			return errors.New("internal: could not process email")
		}
	}

	if after.Password != nil {
		if err := helper.ValidatePassword(*after.Password); err != "" {
			return errors.New("external: invalid password" + err)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(*prev.Password), []byte(*after.Password)); err == nil {
			return errors.New("external: invalid password|new password and previous password must not match")
		}
	}

	if err := database.DbConnection.Model(prev).Update("UpdatedAt", after.UpdatedAt).Error; err != nil {
		log.Error(err)
		return errors.New("internal: could not process email")
	}
	return nil
}

func UpdateUserSensitiveInsensitiveInformation() {}

func PublicGetUserByID(usr *model.User) error {
	if err := database.DbConnection.First(usr, "id=?", usr.ID).Error; err != nil {
		return errors.New("external: unexisting user")
	}

	return nil
}

func PublicGetUserPosts(usr *model.User) ([]model.Post, error) {
	if err := database.DbConnection.Preload("Categories").Preload("Attachments").First(usr, "id=?", usr.ID).Error; err != nil {
		return nil, errors.New("external: unexisting user")
	}

	var usrPosts []model.Post
	if err := database.DbConnection.Preload("Categories").Preload("Likes").Preload("Attachments").Where("author_id=?", usr.ID).Find(&usrPosts).Error; err != nil {
		return nil, errors.New("external: couldn't receive user posts")
	}

	for index := range usrPosts {
		usrPosts[index].CategoryIDs = make([]uint64, 0, len(usrPosts[index].Categories))
		for _, category := range usrPosts[index].Categories {
			usrPosts[index].CategoryIDs = append(usrPosts[index].CategoryIDs, category.ID)
		}
	}

	return usrPosts, nil
}

// DeleteUserAccount soft deletes or anonymizes a user account
// It handles related data appropriately (anonymizes posts, comments, etc.)
func DeleteUserAccount(user *model.User) error {
	if err := database.DbConnection.First(user, "id=?", user.ID).Error; err != nil {
		return errors.New("external: unexisting user")
	}

	// Anonymize user data instead of hard deleting to preserve data integrity
	anonymousName := "Deleted User"
	emailPrefix := ""
	if user.Email != "" {
		if idx := strings.Index(user.Email, "@"); idx > 0 {
			emailPrefix = user.Email[:idx]
		}
	}
	anonymousEmail := "deleted_" + emailPrefix + "@deleted.com"

	// Update user to anonymized state
	anonymousNamePtr := database.Strptr(anonymousName)
	if err := database.DbConnection.Model(user).Updates(map[string]interface{}{
		"email":           anonymousEmail,
		"full_name":       anonymousNamePtr,
		"bio":             nil,
		"password":        nil,
		"google_id":       nil,
		"profile_picture": nil,
	}).Error; err != nil {
		log.Error("Failed to anonymize user: ", err)
		return errors.New("internal: failed to delete user account")
	}

	// Anonymize user's posts
	if err := database.DbConnection.Model(&model.Post{}).
		Where("author_id = ?", user.ID).
		Update("title", "[Deleted]").Error; err != nil {
		log.Error("Failed to anonymize user posts: ", err)
		// Don't fail, just log
	}

	// Clear friend relationships
	if err := database.DbConnection.Model(user).Association("Friends").Clear(); err != nil {
		log.Error("Failed to clear friendships: ", err)
		// Don't fail, just log
	}

	// Clear friend requests
	if err := database.DbConnection.Where("sender_id = ? OR recipient_id = ?", user.ID, user.ID).
		Delete(&model.FriendRequest{}).Error; err != nil {
		log.Error("Failed to delete friend requests: ", err)
		// Don't fail, just log
	}

	// Clear view history
	if err := database.DbConnection.Model(user).Association("PostHistory").Clear(); err != nil {
		log.Error("Failed to clear view history: ", err)
		// Don't fail, just log
	}

	// Clear likes
	if err := database.DbConnection.Model(user).Association("LikedPosts").Clear(); err != nil {
		log.Error("Failed to clear liked posts: ", err)
		// Don't fail, just log
	}

	if err := database.DbConnection.Model(user).Association("LikedComments").Clear(); err != nil {
		log.Error("Failed to clear liked comments: ", err)
		// Don't fail, just log
	}

	return nil
}
