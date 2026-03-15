package handler

import (
	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/service"

	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// SendFriendRequestHandler sends a friend request to another user
// @Summary Send friend request
// @Description Sends a friend request to a user by their email address (requires user:update permission)
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body object true "Friend request data" example({"recipient_email":"user@example.com"})
// @Success 200 {object} map[string]bool "Request sent confirmation"
// @Failure 400 {object} map[string]string "Invalid input or cannot send to yourself"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/users/friends/requests [post]
func SendFriendRequestHandler(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)
	var requestInput struct {
		RecipientEmail string `json:"recipient_email"`
	}

	if err := c.BindJSON(&requestInput); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid data",
		})
		return
	}
	requestInput.RecipientEmail = strings.TrimSpace(requestInput.RecipientEmail)

	if requestInput.RecipientEmail == user.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot send friend request to yourself",
		})
		return
	}

	recipient := model.User{Email: requestInput.RecipientEmail}

	if err := service.SendFriendRequest(&user, &recipient); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

// GetUserPendingFriendRequestsHandler retrieves pending friend requests received by the user
// @Summary Get pending friend requests
// @Description Returns all pending friend requests received by the authenticated user (requires user:update permission)
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.FriendRequest "List of pending friend requests"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/users/friends/requests/pending [get]
func GetUserPendingFriendRequestsHandler(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)

	friendRequests, err := service.GetPendingFriendRequests(&user)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, friendRequests)
}

// GetUserSentFriendRequestsHandler retrieves friend requests sent by the user
// @Summary Get sent friend requests
// @Description Returns all friend requests sent by the authenticated user (requires user:update permission)
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.FriendRequest "List of sent friend requests"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/users/friends/requests/sent [get]
func GetUserSentFriendRequestsHandler(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)

	friendRequests, err := service.GetSentFriendRequests(&user)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, friendRequests)
}

// AcceptFriendRequestHandler accepts a friend request
// @Summary Accept friend request
// @Description Accepts a pending friend request (requires user:update permission)
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Friend Request ID"
// @Success 202 {object} map[string]interface{} "Accepted friend request"
// @Failure 400 {object} map[string]string "Invalid request ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/users/friends/requests/accept/{id} [post]
func AcceptFriendRequestHandler(c *gin.Context) {
	requestID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)

	req := model.FriendRequest{
		ID: requestID,
	}
	if err := service.AcceptFriendRequest(&req, &user); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data": req,
	})
}

// DeclineFriendRequest declines a friend request
// @Summary Decline friend request
// @Description Declines a pending friend request (requires user:update permission)
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Friend Request ID"
// @Success 202 {object} map[string]interface{} "Declined friend request"
// @Failure 400 {object} map[string]string "Invalid request ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/users/friends/requests/cancel/{id} [post]
func DeclineFriendRequest(c *gin.Context) {
	requestID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)

	req := model.FriendRequest{
		ID: requestID,
	}
	if err := service.DeclineFriendRequest(&req, &user); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data": req,
	})
}

// DeleteFriendRequest deletes a friend request
// @Summary Delete friend request
// @Description Deletes a friend request (requires user:update permission)
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Friend Request ID"
// @Success 202 {object} map[string]string "Deletion confirmation"
// @Failure 400 {object} map[string]string "Invalid request ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/users/friends/requests/cancel/{id} [delete]
func DeleteFriendRequest(c *gin.Context) {
	requestID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)

	req := model.FriendRequest{
		ID: requestID,
	}
	if err := service.DeleteFriendRequest(&req, &user); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data": "friend request deleted",
	})
}

// GetUserFriendsHandler retrieves all friends of the authenticated user
// @Summary Get user friends
// @Description Returns all friends of the authenticated user (requires user:update permission)
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.User "List of friends"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/users/friends [get]
func GetUserFriendsHandler(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)

	friends, err := service.GetFriends(&user)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, friends)
}

// UnfriendHandler removes a friend from the user's friends list
// @Summary Unfriend a user
// @Description Removes a friend from the authenticated user's friends list (requires user:update permission)
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Friend User ID"
// @Success 200 {object} map[string]string "Unfriend confirmation"
// @Failure 400 {object} map[string]string "Invalid friend ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Router /api/users/friends/{id} [delete]
func UnfriendHandler(c *gin.Context) {
	friendID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid friend ID",
		})
		return
	}

	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)

	if err := service.RemoveFriend(&user, friendID); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Friend removed successfully",
	})
}
