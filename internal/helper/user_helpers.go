package helper

import (
	"bytes"
	"errors"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

func ValidateFullName(fullname string) string {
	names := strings.Fields(strings.TrimSpace(fullname))
	if len(names) == 0 {
		return "|invalid name"
	}

	for _, name := range names {
		nameRe := regexp.MustCompile(`^[A-Z][a-z]+(-[A-Z][a-z]+)?$`)
		if !nameRe.MatchString(name) {

			return "|invalid name characters"
		}
	}

	return ""
}
func ValidatePassword(password string) string {
	flags := strings.Builder{}
	// Check password length and allowed characters
	passwordRe := regexp.MustCompile(`^[a-zA-Z0-9_!@#$%^&*]+$`)
	if !passwordRe.MatchString(password) {
		flags.WriteString("|invalid password characters")
	}
	// Check for at least one lowercase letter
	lowerRe := regexp.MustCompile(`[a-z]`)
	if !lowerRe.MatchString(password) {
		flags.WriteString("|no lowercase letters provided")
	}

	// Check for at least one uppercase letter
	upperRe := regexp.MustCompile(`[A-Z]`)
	if !upperRe.MatchString(password) {
		flags.WriteString("|no uppercase letters provided")
	}

	// Check for at least one digit
	digitRe := regexp.MustCompile(`\d`)
	if !digitRe.MatchString(password) {
		flags.WriteString("|no digits provided")
	}

	// Check for at least one special character
	specialRe := regexp.MustCompile(`[_!@#$%^&*]`)
	if !specialRe.MatchString(password) {
		flags.WriteString("|no special characters provided")
	}
	return flags.String()
}
func ValidateEmailRegex(email string) bool {
	emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRe.MatchString(email)
}

func ValidateUserUpdateInsensitive(user *model.UserUpdateInsensitiveInput) map[string]any {
	var flags = make(map[string]any)
	flags["necessaryFieldsProvided"] = user.Bio.Valid
	return flags
}

func SanitizeBio(bio string) (string, error) {
	var buffer bytes.Buffer
	if err := goldmark.Convert([]byte(bio), &buffer); err != nil {
		return "", errors.New("failed to parse markdown, error: " + err.Error())
	}

	sanitizedHtml := bluemonday.UGCPolicy().SanitizeBytes(buffer.Bytes())

	return string(sanitizedHtml), nil
}
