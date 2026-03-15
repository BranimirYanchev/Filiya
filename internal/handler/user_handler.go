package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Marionvd/filia-project-backend/database"
	helpers "github.com/Marionvd/filia-project-backend/internal/helper"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers retrieves all users
// @Summary Get all users
// @Description Returns a list of all users in the system (public endpoint)
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} model.User "List of users"
// @Router /api/users/ [get]
func GetUsers(c *gin.Context) {
	var users []model.User
	database.DbConnection.Find(&users)
	c.JSON(200, users)
}

// GetHome returns the API home message
// @Summary Get API home
// @Description Returns a welcome message for the API
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Welcome message"
// @Router /api/ [get]
func GetHome(c *gin.Context) {
	fmt.Println("Home endpoint hit.")

	c.JSON(200, gin.H{
		"msg": "Hello.",
	})
}

// GetLoggedUser retrieves the currently authenticated user
// @Summary Get current user
// @Description Returns the user information for the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Current user information"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/users/me [get]
func GetLoggedUser(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(200, gin.H{
		"user": user,
	})
}

// GetUser retrieves a specific user by ID
// @Summary Get user by ID
// @Description Returns public information about a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User information"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/users/public/{id} [get]
func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid path param",
		})
		return
	}

	user := model.User{
		ID: id,
	}
	if err := service.PublicGetUserByID(&user); err != nil {
		code, msg := respondError(err)

		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// UserUpdateSensitiveHandler updates sensitive user information (email, password, full name)
// @Summary Update sensitive user information
// @Description Updates email, password, or full name. Requires old password confirmation.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body model.UserUpdateSensitiveBody true "Updated user information"
// @Success 200 {object} map[string]string "Update confirmation"
// @Failure 400 {object} map[string]string "Invalid input or password mismatch"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/users/profile/sensitive [put]
// * Requires confirmation (via confirm_password)
func UserUpdateSensitiveHandler(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, ok := userObj.(model.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unauthenticated",
		})
		return
	}

	var updateModel model.UserUpdateSensitiveBody
	if err := c.BindJSON(&updateModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid data",
		})
		return
	}

	updatedUser := model.User{
		FullName:  updateModel.FullName,
		Password:  updateModel.Password,
		UpdatedAt: time.Now(),
	}

	if updateModel.Email != nil {
		updatedUser.Email = *updateModel.Email
	}
	if user.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password|cannot change user information|user has not set a password",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(updateModel.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password|cannot change user information|invalid confirmation password",
		})
		return
	}

	err := service.UpdateUserSensitiveInformation(&user, &updatedUser)

	if err != nil {
		if strings.Contains(err.Error(), "external") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()[10:],
			})
			return
		} else if strings.Contains(err.Error(), "internal") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()[10:],
			})
			return
		}
	}

	if updateModel.FullName != nil && *updateModel.FullName != "" {
		user.FullName = updateModel.FullName
	}
	if updateModel.Email != nil && *updateModel.Email != "" {
		user.Email = *updateModel.Email
	}

	jwtUser := user.GenerateJwtUser()
	jwt := helpers.GenerateJWT(jwtUser)
	refreshToken := helpers.GenerateRefreshToken(jwtUser)
	helpers.SetAuthCookies(c, jwt, refreshToken)
	c.JSON(http.StatusOK, gin.H{
		"data": "updated user",
	})
}

// UserUpdateInsensitiveHandler updates non-sensitive user information (bio)
// @Summary Update user bio
// @Description Updates the user's bio information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body model.UserUpdateInsensitiveInput true "Updated bio information"
// @Success 200 {object} map[string]interface{} "Updated user information"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/users/profile/basic [put]
func UserUpdateInsensitiveHandler(c *gin.Context) {
	var updatedUser model.UserUpdateInsensitiveInput
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(400, gin.H{
			"validInput": false,
			"info":       err.Error(),
		})
		return
	}

	var userObj, _ = c.Get("user")
	var user, _ = userObj.(model.User)

	if !updatedUser.Bio.Valid {
		sanitizedBio, err := helpers.SanitizeBio(updatedUser.Bio.Value)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not sanitize bio|invalid bio",
			})
			return
		}

		database.DbConnection.Model(&user).Update("bio", sql.NullString{Valid: true, String: sanitizedBio})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// GetUserPosts retrieves all posts by a specific user
// @Summary Get user posts
// @Description Returns all posts created by the specified user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} model.Post "List of user posts"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/users/{id}/posts [get]
func GetUserPosts(c *gin.Context) {
	var userObj, _ = c.Get("user")
	var user, _ = userObj.(model.User)

	posts, err := service.PublicGetUserPosts(&user)
	if err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// UpdateAvatarHandler uploads and updates user profile picture
// @Summary Update user avatar
// @Description Uploads a new profile picture for the authenticated user
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param profile_picture formData file true "Profile picture image file"
// @Success 200 {object} map[string]bool "Upload confirmation"
// @Failure 400 {object} map[string]interface{} "Invalid file"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/profile/avatar [post]
func UpdateAvatarHandler(c *gin.Context) {
	file, err := c.FormFile("profile_picture")
	if err != nil {
		c.JSON(400, gin.H{"validFile": false, "passed": false})
		return
	}

	userObj, _ := c.Get("user")
	user, _ := userObj.(model.User)

	saveDir := filepath.Join("..", "profiles", fmt.Sprintf("%d", user.ID))
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"passed": false})
		log.Error("Error creating directory: ", err.Error())
		return
	}

	savePath := filepath.Join(saveDir, file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(500, gin.H{
			"error": "failed to save file",
		})
		log.Error("Error saving file: ", err)
		return
	}

	database.DbConnection.Model(&user).Update("profile_picture", savePath)

	c.JSON(200, gin.H{"passed": true})
}

// DeleteAccountHandler deletes the authenticated user's account
// @Summary Delete user account
// @Description Anonymizes and deletes the authenticated user's account (requires user:delete permission)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Account deletion confirmation"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/me [delete]
func DeleteAccountHandler(c *gin.Context) {
	userObj, _ := c.Get("user")
	user, ok := userObj.(model.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unauthenticated",
		})
		return
	}

	if err := service.DeleteUserAccount(&user); err != nil {
		code, msg := respondError(err)
		c.JSON(code, gin.H{
			"error": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}
