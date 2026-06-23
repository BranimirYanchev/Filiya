package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Marionvd/filia-project-backend/internal/helper"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/internal/service"
	"github.com/Marionvd/filia-project-backend/internal/username"

	"github.com/Marionvd/filia-project-backend/config"
	"github.com/Marionvd/filia-project-backend/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	googleapi "google.golang.org/api/oauth2/v2"
)

func redirectToFrontend(c *gin.Context, baseURL string, params map[string]string) bool {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return false
	}

	redirectURL, err := url.Parse(baseURL)
	if err != nil {
		log.WithError(err).Warn("failed to parse frontend redirect url")
		return false
	}

	query := redirectURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	redirectURL.RawQuery = query.Encode()

	c.Redirect(http.StatusSeeOther, redirectURL.String())
	return true
}

// AuthenticationHome returns the authentication home message
// @Summary Get authentication home
// @Description Returns a welcome message for the authentication endpoint
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Welcome message"
// @Router /api/auth/ [get]
func AuthenticationHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "Authentication home",
	})
}

// Register handles user registration via email and password.
// @Summary Register a new user with email and password
// @Description Validates input, hashes password, assigns default role, and creates user in the database
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.PasswordRegisterInput true "User registration input"
// @Success 201 {object} map[string]interface{} "Validation flags"
// @Failure 400 {object} map[string]interface{} "Invalid input or hashing error"
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
	var userJSON model.PasswordRegisterInput

	if err := c.BindJSON(&userJSON); err != nil {
		log.Error("error while binding json", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|invalid data",
		})
		return
	}

	if userJSON.Password != userJSON.RepeatedPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request|password do not match",
		})
		return
	}

	user := model.User{
		Email:    userJSON.Email,
		FullName: &userJSON.FullName,
		Password: &userJSON.Password,
		GoogleID: nil,
	}

	if err := service.SaveRegisterUser(&user); err != nil {
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

	token := helper.GenerateJWT(user.GenerateJwtUser())
	refreshToken := helper.GenerateRefreshToken(user.GenerateJwtUser())

	helper.SetAuthCookies(c, token, refreshToken)

	c.JSON(http.StatusCreated, gin.H{
		"error":         "",
		"data":          user,
		"token":         token,
		"refresh_token": refreshToken,
	})
}

// EmailLogin authenticates a user using email and password.
//
// @Summary Login with email and password
// @Description Validates credentials and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body model.PasswordLoginInput true "Login credentials"
// @Success 200 {object} map[string]interface{} "Validation flags"
// @Failure 400 {object} map[string]interface{} "Invalid input or login failure"
// @Router /api/auth/login [post]
func EmailLogin(c *gin.Context) {
	var loginInput model.PasswordLoginInput

	if err := c.BindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	user, err := service.Login(loginInput)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()[10:],
		})
		return
	}

	token := helper.GenerateJWT(user.GenerateJwtUser())
	refreshToken := helper.GenerateRefreshToken(user.GenerateJwtUser())

	helper.SetAuthCookies(c, token, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"error":         "",
		"data":          user,
		"token":         token,
		"refresh_token": refreshToken,
	})
}

// GoogleLogin initiates Google OAuth login flow
// @Summary Initiate Google OAuth login
// @Description Redirects user to Google OAuth consent screen
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 307 "Redirect to Google OAuth"
// @Router /api/auth/google/login [get]
func GoogleLogin(c *gin.Context) {
	url := config.GoogleConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallbackHandler handles the OAuth callback from Google
// @Summary Handle Google OAuth callback
// @Description Processes the OAuth callback, creates or authenticates user, and returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code query string true "OAuth authorization code from Google"
// @Success 200 {object} map[string]string "User information and token"
// @Failure 400 {object} map[string]string "Invalid or missing code"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/google/callback [get]
func GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in google authentication request."})
		return
	}

	token, err := config.GoogleConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := config.GoogleConfig.Client(c.Request.Context(), token)
	service, err := googleapi.New(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Google API client"})
		return
	}

	userinfo, err := service.Userinfo.Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	var user model.User
	result := database.DbConnection.Where("email = ?", userinfo.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, sql.ErrNoRows) || result.RowsAffected == 0 {
			usernameValue, err := username.Generate(userinfo.Name, func(candidate string) (bool, error) {
				var count int64
				if err := database.DbConnection.Model(&model.User{}).Where("username = ?", candidate).Count(&count).Error; err != nil {
					return false, err
				}
				return count > 0, nil
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate username"})
				return
			}

			user = model.User{
				Email:          userinfo.Email,
				Username:       database.Strptr(usernameValue),
				FullName:       &userinfo.Name,
				GoogleID:       &userinfo.Id,
				RoleID:         1,
				ProfilePicture: &userinfo.Picture,
			}
			if err := database.DbConnection.Save(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	jwtToken := helper.GenerateJWT(user.GenerateJwtUser())
	refreshToken := helper.GenerateRefreshToken(user.GenerateJwtUser())
	helper.SetAuthCookies(c, jwtToken, refreshToken)

	if redirectToFrontend(c, config.FrontendAuthSuccessURL(), map[string]string{
		"auth":     "success",
		"provider": "google",
	}) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":         userinfo.Email,
		"name":          userinfo.Name,
		"token":         jwtToken,
		"refresh_token": refreshToken,
	})
}

// Logout handles user logout by clearing the authentication cookie
// @Summary Logout user
// @Description Clears the authentication cookie to log out the user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Logout confirmation"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/auth/logout [post]
func Logout(c *gin.Context) {
	helper.ClearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// RefreshToken handles refresh token request to get a new access token
// @Summary Refresh access token
// @Description Uses refresh token to generate a new access token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body object true "Refresh token" example({"refresh_token":"string"})
// @Success 200 {object} map[string]interface{} "New access token"
// @Failure 400 {object} map[string]string "Invalid refresh token"
// @Router /api/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	refreshToken := helper.ExtractRefreshToken(c)
	if refreshToken == "" {
		helper.ClearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token required",
		})
		return
	}

	token, err := helper.ValidateToken(&refreshToken)
	if err != nil {
		helper.ClearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	// Verify it's a refresh token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helper.ClearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token format",
		})
		return
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		helper.ClearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token type",
		})
		return
	}

	// Extract user from token
	user, err := helper.ExtractUserFromToken(token)
	if err != nil {
		helper.ClearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "failed to extract user from token",
		})
		return
	}

	var foundUser model.User
	if err := database.DbConnection.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
		helper.ClearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
		return
	}

	// Generate new access token
	jwtUser := foundUser.GenerateJwtUser()
	newToken := helper.GenerateJWT(jwtUser)
	newRefreshToken := helper.GenerateRefreshToken(jwtUser)

	helper.SetAuthCookies(c, newToken, newRefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"token":         newToken,
		"refresh_token": newRefreshToken,
	})
}

// PasswordResetRequest handles password reset request
// @Summary Request password reset
// @Description Sends a password reset token to the user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param email body object true "User email" example({"email":"user@example.com"})
// @Success 200 {object} map[string]string "Reset token sent confirmation"
// @Failure 400 {object} map[string]string "Invalid email or user not found"
// @Router /api/auth/reset-password-request [post]
func PasswordResetRequest(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	if strings.TrimSpace(input.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email is required",
		})
		return
	}

	// Check if user exists
	var user model.User
	if err := database.DbConnection.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Don't reveal if user exists or not for security
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent",
		})
		return
	}

	// Check if user has a password set (not OAuth only)
	if user.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user has not set a password",
		})
		return
	}

	// Generate reset token (valid for 1 hour)
	resetToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
		"type":  "password_reset",
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := resetToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Error("Error generating reset token: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate reset token",
		})
		return
	}

	// TODO: Send email with reset token
	// For now, return token in response (in production, send via email only)
	log.Info("Password reset token for ", input.Email, ": ", tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message":     "If the email exists, a password reset link has been sent",
		"reset_token": tokenString, // Remove in production - only for development
	})
}

// PasswordReset handles password reset with token
// @Summary Reset password
// @Description Resets user password using a valid reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body object true "Reset data" example({"token":"reset_token","new_password":"newPassword123"})
// @Success 200 {object} map[string]string "Password reset confirmation"
// @Failure 400 {object} map[string]string "Invalid token or password"
// @Router /api/auth/reset-password [post]
func PasswordReset(c *gin.Context) {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	if input.Token == "" || input.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token and new_password are required",
		})
		return
	}

	// Validate token
	token, err := helper.ValidateToken(&input.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid or expired token",
		})
		return
	}

	// Verify it's a password reset token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token format",
		})
		return
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "password_reset" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token type",
		})
		return
	}

	// Get user ID from token
	userID, ok := claims["id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token data",
		})
		return
	}

	// Find user
	var user model.User
	if err := database.DbConnection.First(&user, "id = ?", uint64(userID)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not found",
		})
		return
	}

	// Validate new password
	if pwdErr := helper.ValidatePassword(input.NewPassword); pwdErr != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password" + pwdErr,
		})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error hashing password: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process password",
		})
		return
	}

	// Update password
	hashedPasswordStr := string(hashedPassword)
	if err := database.DbConnection.Model(&user).Update("password", hashedPasswordStr).Error; err != nil {
		log.Error("Error updating password: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}
