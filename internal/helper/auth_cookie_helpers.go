package helper

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AccessTokenCookieName  = "token"
	RefreshTokenCookieName = "refresh_token"
)

func cookieSecure() bool {
	return strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true")
}

func SetAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetSameSite(http.SameSiteLaxMode)

	if accessToken != "" {
		c.SetCookie(AccessTokenCookieName, accessToken, int(AccessTokenTTL.Seconds()), "/", "", cookieSecure(), true)
	}

	if refreshToken != "" {
		c.SetCookie(RefreshTokenCookieName, refreshToken, int(RefreshTokenTTL.Seconds()), "/", "", cookieSecure(), true)
	}
}

func ClearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AccessTokenCookieName, "", -1, "/", "", cookieSecure(), true)
	c.SetCookie(RefreshTokenCookieName, "", -1, "/", "", cookieSecure(), true)
}

func ExtractAccessToken(c *gin.Context) string {
	headerValue := strings.TrimSpace(c.GetHeader("Authorization"))
	if strings.HasPrefix(strings.ToLower(headerValue), "bearer ") {
		return strings.TrimSpace(headerValue[7:])
	}

	cookieValue, err := c.Cookie(AccessTokenCookieName)
	if err == nil {
		return strings.TrimSpace(cookieValue)
	}

	return ""
}

func ExtractRefreshToken(c *gin.Context) string {
	cookieValue, err := c.Cookie(RefreshTokenCookieName)
	if err == nil && strings.TrimSpace(cookieValue) != "" {
		return strings.TrimSpace(cookieValue)
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err == nil {
		return strings.TrimSpace(body.RefreshToken)
	}

	return ""
}
