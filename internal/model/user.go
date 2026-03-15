// Package model defines the data structures used throughout the Filia application.
// It contains the database model and input/output data transfer objects (DTOs).
package model

import (
	"time"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

// User represents a user in the Filia application.
//
// @Description Full user model stored in the database.
// @Description Includes personal info, role, social connections, and metadata.
type User struct {
	ID             uint64     `gorm:"primaryKey" json:"id" example:"1"`
	Email          string     `gorm:"unique;not null;size:60" json:"email" example:"user@example.com"`
	Username       *string    `gorm:"uniqueIndex;size:120" json:"username" swaggertype:"string" example:"john.doe"`
	FullName       *string    `gorm:"size:254" json:"full_name" swaggertype:"string" example:"John Doe"`
	Bio            *string    `gorm:"size:500;null" json:"bio" swaggertype:"string" example:"Loves hiking and coding."`
	RoleID         uint64     `gorm:"default:1" json:"role_id" example:"2"`
	Role           *Role      `gorm:"foreignKey:RoleID" json:"role"`
	ProfilePicture *string    `gorm:"size:500" json:"profile_picture" swaggertype:"string" example:"https://example.com/avatar.jpg"`
	GoogleID       *string    `gorm:"unique;null" json:"google_id" swaggertype:"string" example:"google-oauth-id"`
	Password       *string    `gorm:"size:61" json:"password" swaggertype:"string" example:"passDAdsa@"`
	Friends        []*User    `gorm:"many2many:user_friends" json:"friends"`
	LikedPosts     []*Post    `gorm:"many2many:post_likes" json:"liked_posts"`
	PostHistory    []Post     `gorm:"many2many:post_view_history" json:"post_history"`
	LikedComments  []*Comment `gorm:"many2many:comment_likes" json:"likes"`
	CreatedAt      time.Time  `json:"created_at" example:"2023-08-17T10:00:00Z"`
	UpdatedAt      time.Time  `json:"updated_at" example:"2023-08-17T10:00:00Z"`
}

// SafeUser is a sanitized version of the User model.
//
// @Description Used for returning user data without sensitive fields like password or GoogleID.
type SafeUser struct {
	ID             uint64    `json:"id" example:"1"`
	Email          string    `json:"email" example:"user@example.com"`
	Username       *string   `json:"username" swaggertype:"string" example:"john.doe"`
	FullName       *string   `json:"full_name" swaggertype:"string" example:"John Doe"`
	Bio            *string   `json:"bio" swaggertype:"string" example:"Loves hiking and coding."`
	Role           *Role     `json:"role"`
	ProfilePicture *string   `json:"profile_picture" swaggertype:"string" example:"https://example.com/avatar.jpg"`
	Friends        []*User   `json:"friends"`
	CreatedAt      time.Time `json:"created_at" example:"2023-08-17T10:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (user *User) GenerateSafeUser() SafeUser {
	return SafeUser{
		ID:             user.ID,
		Email:          user.Email,
		Username:       user.Username,
		FullName:       user.FullName,
		Bio:            user.Bio,
		Role:           user.Role,
		ProfilePicture: user.ProfilePicture,
		Friends:        user.Friends,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

// PasswordRegisterInput is used for registering a new user.
//
// @Description DTO for user registration with password confirmation.
type PasswordRegisterInput struct {
	Email            string `json:"email" example:"user@example.com"`
	FullName         string `json:"full_name" example:"John Doe"`
	Password         string `json:"password" example:"securePassword123"`
	RepeatedPassword string `json:"repeated_password" example:"securePassword123"`
}

// PasswordLoginInput is used for logging in with email and password.
//
// @Description DTO for user login.
type PasswordLoginInput struct {
	Email    string `json:"email" example:"standard_user@mail.com"`
	Password string `json:"password" example:"m123#@S"`
}

// JWTUser represents the minimal user info encoded in JWT.
//
// @Description Used for generating and decoding JWT tokens.
type JWTUser struct {
	Id       uint64 `json:"id" example:"1"`
	FullName string `json:"full_name" example:"John Doe"`
	Email    string `json:"email" example:"user@example.com"`
}

func (user *User) GenerateJwtUser() JWTUser {
	return JWTUser{
		Id:       user.ID,
		FullName: *user.FullName,
		Email:    user.Email,
	}
}

type UserUpdateSensitiveInput struct {
	Email       OptionalString `json:"email"`
	Password    OptionalString `json:"password"`
	OldPassword string         `json:"old_password"`
}

type UserUpdateSensitiveBody struct {
	Email       *string `json:"email"`
	FullName    *string `json:"full_name"`
	Password    *string `json:"password"`
	OldPassword string  `json:"old_password"`
}

// UserUpdateInsensitiveInput is used for updating non-sensitive user fields.
//
// @Description DTO for updating user bio or other public info.
type UserUpdateInsensitiveInput struct {
	Bio OptionalString `json:"bio"`
}
