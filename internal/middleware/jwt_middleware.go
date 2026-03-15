package middleware

import (
	"net/http"

	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/helper"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

func JwtMiddleware(c *gin.Context) {
	tokenString := helper.ExtractAccessToken(c)
	if tokenString == "" {
		c.Set("authenticated", false)
		c.Next()
		return
	}

	token, err := helper.ValidateToken(&tokenString)
	if err != nil {
		log.Error("Error while validating token: ", err.Error())
		c.Set("authenticated", false)
		c.Next()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.Set("authenticated", false)
		c.Next()
		return
	}

	userData, ok := claims["user"].(map[string]interface{})
	if !ok {
		c.Set("authenticated", false)
		c.Next()
		return
	}

	var foundUser model.User

	email, ok := userData["email"].(string)
	if !ok || email == "" {
		c.Set("authenticated", false)
		c.Next()
		return
	}

	if err = database.DbConnection.Where("email=?", email).Preload("Role").First(&foundUser).Error; err != nil {
		log.Error("Invalid token. User not found: ", err.Error())
		c.Set("authenticated", false)
		c.Next()
		return
	}
	if err := database.DbConnection.Model(&foundUser.Role).Preload("Permissions").Find(&foundUser.Role).Error; err != nil {
		log.Error(err)
		c.Set("authenticated", false)
		c.Next()
		return
	}

	c.Set("user", foundUser)
	c.Set("authenticated", true)
	c.Next()
}

func AuthenticationUnnecessary(c *gin.Context) {
	authenticated, exists := c.Get("authenticated")
	if exists && authenticated.(bool) {
		c.JSON(
			200,
			gin.H{
				"authenticated": true,
				"canRedirect":   true,
			},
		)
		c.Abort()
	}
}

func AuthenticationNecessary(c *gin.Context) {
	authenticated, exists := c.Get("authenticated")
	if !exists || !authenticated.(bool) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"error":         "unauthorized",
		})
		c.Abort()
	}
}
