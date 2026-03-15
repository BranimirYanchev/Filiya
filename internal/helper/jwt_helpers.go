package helper

import (
	"os"
	"time"

	"github.com/Marionvd/filia-project-backend/internal/model"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GenerateJWT(user model.JWTUser) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"type": "access",
		"exp":  time.Now().Add(AccessTokenTTL).Unix(),
	})

	tokenString, err := token.SignedString(getJWTSecret())

	if err != nil {
		log.Error("Error while signing token", err.Error())
		return ""
	}

	return tokenString
}

// GenerateRefreshToken generates a refresh token with longer expiration
func GenerateRefreshToken(user model.JWTUser) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"type": "refresh",
		"exp":  time.Now().Add(RefreshTokenTTL).Unix(),
	})

	tokenString, err := token.SignedString(getJWTSecret())

	if err != nil {
		log.Error("Error while signing refresh token", err.Error())
		return ""
	}

	return tokenString
}

func ValidateToken(tokenString *string) (*jwt.Token, error) {
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return getJWTSecret(), nil
	})

	if err != nil {
		log.Error("Error while parsing token.", err.Error())
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return token, nil
}

// ExtractUserFromToken extracts JWTUser from a token's claims
func ExtractUserFromToken(token *jwt.Token) (model.JWTUser, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.JWTUser{}, jwt.ErrInvalidKey
	}

	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		return model.JWTUser{}, jwt.ErrInvalidKey
	}

	var user model.JWTUser
	if id, ok := userMap["id"].(float64); ok {
		user.Id = uint64(id)
	}
	if email, ok := userMap["email"].(string); ok {
		user.Email = email
	}
	if fullName, ok := userMap["full_name"].(string); ok {
		user.FullName = fullName
	}

	return user, nil
}
