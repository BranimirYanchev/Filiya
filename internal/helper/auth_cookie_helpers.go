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

func authCookiesEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("AUTH_COOKIES_ENABLED")), "true")
}

func cookieSecure() bool {
	if value := strings.TrimSpace(os.Getenv("COOKIE_SECURE")); value != "" {
		return strings.EqualFold(value, "true")
	}

	frontendURL := strings.ToLower(strings.TrimSpace(os.Getenv("FRONTEND_URL")))
	if strings.HasPrefix(frontendURL, "https://") &&
		!strings.Contains(frontendURL, "localhost") &&
		!strings.Contains(frontendURL, "127.0.0.1") {
		return true
	}

	return false
}

func cookieSameSite() http.SameSite {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("COOKIE_SAME_SITE"))) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	case "lax":
		return http.SameSiteLaxMode
	}

	if cookieSecure() {
		return http.SameSiteNoneMode
	}

	return http.SameSiteLaxMode
}

func SetAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	if !authCookiesEnabled() {
		return
	}

	c.SetSameSite(cookieSameSite())

	if accessToken != "" {
		c.SetCookie(AccessTokenCookieName, accessToken, int(AccessTokenTTL.Seconds()), "/", "", cookieSecure(), true)
	}

	if refreshToken != "" {
		c.SetCookie(RefreshTokenCookieName, refreshToken, int(RefreshTokenTTL.Seconds()), "/", "", cookieSecure(), true)
	}
}

func ClearAuthCookies(c *gin.Context) {
	if !authCookiesEnabled() {
		return
	}

	c.SetSameSite(cookieSameSite())
	c.SetCookie(AccessTokenCookieName, "", -1, "/", "", cookieSecure(), true)
	c.SetCookie(RefreshTokenCookieName, "", -1, "/", "", cookieSecure(), true)
}

func ExtractAccessToken(c *gin.Context) string {
	headerValue := strings.TrimSpace(c.GetHeader("Authorization"))
	if strings.HasPrefix(strings.ToLower(headerValue), "bearer ") {
		return strings.TrimSpace(headerValue[7:])
	}

	if !authCookiesEnabled() {
		return ""
	}

	cookieValue, err := c.Cookie(AccessTokenCookieName)
	if err == nil {
		return strings.TrimSpace(cookieValue)
	}

	return ""
}

func ExtractRefreshToken(c *gin.Context) string {
	if authCookiesEnabled() {
		cookieValue, err := c.Cookie(RefreshTokenCookieName)
		if err == nil && strings.TrimSpace(cookieValue) != "" {
			return strings.TrimSpace(cookieValue)
		}
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err == nil {
		return strings.TrimSpace(body.RefreshToken)
	}

	return ""
}
