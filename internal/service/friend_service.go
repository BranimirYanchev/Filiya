package service

import (
	"errors"
	"fmt"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SendFriendRequest(sender, recipient *model.User) error {
	if err := database.DbConnection.First(sender, "id=?", sender.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting sender")
		}
		return fmt.Errorf("internal: failed to query user by email | %w", err)
	}
	if err := database.DbConnection.First(recipient, "email=?", recipient.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting recipient")
		}
		return fmt.Errorf("internal: failed to query user by email | %w", err)
	}
	request := model.FriendRequest{
		Recipient: recipient,
		Sender:    sender,
	}
	if err := database.DbConnection.Create(&request).Error; err != nil {
		return errors.New("internal: failed to send friend request")
	}

	// Load sender with full name
	if err := database.DbConnection.First(sender, "id=?", sender.ID).Error; err == nil {
		if sender.FullName != nil {
			requestID := request.ID
			CreateNotification(recipient.ID, model.NotificationTypeFriendRequest,
				"New friend request",
				fmt.Sprintf("%s sent you a friend request", *sender.FullName),
				&requestID, "friend_request")
		}
	}

	return nil
}

func GetPendingFriendRequests(recipient *model.User) ([]model.FriendRequest, error) {
	var friendRequests []model.FriendRequest
	if err := database.DbConnection.First(recipient, "id=?", recipient.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("external: unexisting recipient")
		}
		return nil, fmt.Errorf("internal: failed to query user by id | %w", err)
	}

	if err := database.DbConnection.Where("recipient_id = ? AND pending = ?", recipient.ID, true).
		Preload("Sender").Preload("Recipient").Find(&friendRequests).Error; err != nil {
		log.Error(err)
		return nil, errors.New("internal: failed to get requests")
	}
	return friendRequests, nil
}

func GetSentFriendRequests(sender *model.User) ([]model.FriendRequest, error) {
	var friendRequests []model.FriendRequest
	if err := database.DbConnection.First(sender, "id=?", sender.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("external: unexisting sender")
		}
		return nil, fmt.Errorf("internal: failed to query user by id | %w", err)
	}

	if err := database.DbConnection.Where("sender_id = ? AND pending = ?", sender.ID, true).
		Preload("Sender").Preload("Recipient").Find(&friendRequests).Error; err != nil {
		log.Error(err)
		return nil, errors.New("internal: failed to get requests")
	}
	return friendRequests, nil
}

func AcceptFriendRequest(req *model.FriendRequest, user *model.User) error {
	if err := database.DbConnection.First(user, "id=?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting user")
		}
		return fmt.Errorf("internal: failed to query user by email | %w", err)
	}

	if err := database.DbConnection.First(req, "id=?", req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting friend request")
		}
		return fmt.Errorf("internal: failed to query request by id | %w", err)
	}

	if !req.Pending {
		return errors.New("external: request is not pending, so it cannot be edited")
	}

	if user.ID != req.RecipientID {
		return errors.New("external: unauthorized action|invalid recipient id")
	}

	database.DbConnection.Model(req).Update("Pending", false)
	database.DbConnection.Model(req).Update("Accepted", true)

	// Add both users to each other's friends list
	err := database.DbConnection.Model(user).Association("Friends").Append(&req.Sender)
	if err != nil {
		return errors.New("internal: cannot process friend request")
	}
	err = database.DbConnection.Model(&req.Sender).Association("Friends").Append(user)
	if err != nil {
		return errors.New("internal: cannot process friend request")
	}

	// Load sender to get full name
	var sender model.User
	if err := database.DbConnection.First(&sender, "id=?", req.SenderID).Error; err == nil {
		if sender.FullName != nil && user.FullName != nil {
			requestID := req.ID
			CreateNotification(req.SenderID, model.NotificationTypeFriendAccept,
				"Friend request accepted",
				fmt.Sprintf("%s accepted your friend request", *user.FullName),
				&requestID, "friend_request")
		}
	}

	return nil
}

func DeleteFriendRequest(req *model.FriendRequest, user *model.User) error {
	if err := database.DbConnection.First(user, "id=?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting user")
		}
		return fmt.Errorf("internal: failed to query user by email | %w", err)
	}

	if err := database.DbConnection.First(req, "id=?", req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting friend request")
		}
		return fmt.Errorf("internal: failed to query request by id | %w", err)
	}

	if !req.Pending {
		return errors.New("external: request is not pending, so it cannot be edited")
	}

	if user.ID != req.SenderID {
		return errors.New("external: unauthorized action|invalid sender id")
	}

	if err := database.DbConnection.Delete(req).Error; err != nil {
		return errors.New("internal: failed to delete friend request|" + err.Error())
	}

	return nil
}

func DeclineFriendRequest(req *model.FriendRequest, user *model.User) error {
	if err := database.DbConnection.First(user, "id=?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting user")
		}
		return fmt.Errorf("internal: failed to query user by email | %w", err)
	}

	if err := database.DbConnection.First(req, "id=?", req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting friend request")
		}
		return fmt.Errorf("internal: failed to query request by id | %w", err)
	}

	if !req.Pending {
		return errors.New("external: request is not pending, so it cannot be edited")
	}

	if user.ID != req.RecipientID {
		return errors.New("external: unauthorized action|invalid recipient id")
	}

	database.DbConnection.Model(req).Update("Pending", false)
	database.DbConnection.Model(req).Update("Accepted", false)

	return nil
}

// GetFriends retrieves all friends of a user
func GetFriends(user *model.User) ([]model.User, error) {
	if err := database.DbConnection.First(user, "id=?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("external: unexisting user")
		}
		return nil, fmt.Errorf("internal: failed to query user by id | %w", err)
	}

	var friends []model.User
	if err := database.DbConnection.Model(user).Association("Friends").Find(&friends); err != nil {
		log.Error(err)
		return nil, errors.New("internal: failed to get friends")
	}

	return friends, nil
}

// RemoveFriend removes a friend from user's friends list
func RemoveFriend(user *model.User, friendID uint64) error {
	if err := database.DbConnection.First(user, "id=?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting user")
		}
		return fmt.Errorf("internal: failed to query user by id | %w", err)
	}

	var friend model.User
	if err := database.DbConnection.First(&friend, "id=?", friendID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("external: unexisting friend")
		}
		return fmt.Errorf("internal: failed to query friend by id | %w", err)
	}

	// Remove from both sides of the friendship
	if err := database.DbConnection.Model(user).Association("Friends").Delete(&friend); err != nil {
		log.Error(err)
		return errors.New("internal: failed to remove friend")
	}

	if err := database.DbConnection.Model(&friend).Association("Friends").Delete(user); err != nil {
		log.Error(err)
		return errors.New("internal: failed to remove friend")
	}

	return nil
}
