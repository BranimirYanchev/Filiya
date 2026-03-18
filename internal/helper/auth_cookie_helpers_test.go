package helper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCookieSameSiteDefaultsToLaxWithoutSecureCookies(t *testing.T) {
	t.Setenv("COOKIE_SECURE", "false")
	t.Setenv("COOKIE_SAME_SITE", "")

	if got := cookieSameSite(); got != http.SameSiteLaxMode {
		t.Fatalf("cookieSameSite() = %v, want %v", got, http.SameSiteLaxMode)
	}
}

func TestCookieSameSiteDefaultsToNoneForSecureCookies(t *testing.T) {
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("COOKIE_SAME_SITE", "")

	if got := cookieSameSite(); got != http.SameSiteNoneMode {
		t.Fatalf("cookieSameSite() = %v, want %v", got, http.SameSiteNoneMode)
	}
}

func TestCookieSecureDefaultsToTrueForHttpsFrontend(t *testing.T) {
	t.Setenv("COOKIE_SECURE", "")
	t.Setenv("FRONTEND_URL", "https://filiya-frontend.onrender.com")

	if !cookieSecure() {
		t.Fatal("expected cookieSecure() to default to true for HTTPS frontend URL")
	}
}

func TestCookieSecureStaysFalseForLocalFrontend(t *testing.T) {
	t.Setenv("COOKIE_SECURE", "")
	t.Setenv("FRONTEND_URL", "http://localhost:3000")

	if cookieSecure() {
		t.Fatal("expected cookieSecure() to stay false for local frontend URL")
	}
}

func TestSetAuthCookiesUsesSameSiteNoneForSecureCookies(t *testing.T) {
	t.Setenv("AUTH_COOKIES_ENABLED", "true")
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("COOKIE_SAME_SITE", "")

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	SetAuthCookies(c, "access-token", "refresh-token")

	cookies := recorder.Result().Cookies()
	if len(cookies) != 2 {
		t.Fatalf("expected 2 cookies, got %d", len(cookies))
	}

	for _, cookie := range cookies {
		if !cookie.Secure {
			t.Fatalf("expected cookie %q to be Secure", cookie.Name)
		}
		if !cookie.HttpOnly {
			t.Fatalf("expected cookie %q to be HttpOnly", cookie.Name)
		}
		if cookie.SameSite != http.SameSiteNoneMode {
			t.Fatalf("expected cookie %q SameSite=None, got %v", cookie.Name, cookie.SameSite)
		}
	}

	setCookieHeader := strings.Join(recorder.Header().Values("Set-Cookie"), "\n")
	if !strings.Contains(setCookieHeader, "SameSite=None") {
		t.Fatalf("expected Set-Cookie header to contain SameSite=None, got %q", setCookieHeader)
	}
}

func TestSetAuthCookiesDisabledByDefault(t *testing.T) {
	t.Setenv("AUTH_COOKIES_ENABLED", "")

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	SetAuthCookies(c, "access-token", "refresh-token")

	if cookies := recorder.Result().Cookies(); len(cookies) != 0 {
		t.Fatalf("expected no cookies when auth cookies are disabled, got %d", len(cookies))
	}
}

func TestCookieSameSiteHonorsExplicitOverride(t *testing.T) {
	t.Setenv("AUTH_COOKIES_ENABLED", "true")
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("COOKIE_SAME_SITE", "strict")

	if got := cookieSameSite(); got != http.SameSiteStrictMode {
		t.Fatalf("cookieSameSite() = %v, want %v", got, http.SameSiteStrictMode)
	}
}

func TestExtractAccessTokenIgnoresCookieWhenDisabled(t *testing.T) {
	t.Setenv("AUTH_COOKIES_ENABLED", "")

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: AccessTokenCookieName, Value: "cookie-token"})

	if got := ExtractAccessToken(c); got != "" {
		t.Fatalf("expected no access token from cookie when disabled, got %q", got)
	}
}

func TestExtractRefreshTokenReadsBodyWhenCookiesDisabled(t *testing.T) {
	t.Setenv("AUTH_COOKIES_ENABLED", "")

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"refresh_token":"body-token"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.AddCookie(&http.Cookie{Name: RefreshTokenCookieName, Value: "cookie-token"})

	if got := ExtractRefreshToken(c); got != "body-token" {
		t.Fatalf("expected refresh token from body when cookies are disabled, got %q", got)
	}
}
