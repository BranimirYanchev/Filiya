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

func TestSetAuthCookiesUsesSameSiteNoneForSecureCookies(t *testing.T) {
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

func TestCookieSameSiteHonorsExplicitOverride(t *testing.T) {
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("COOKIE_SAME_SITE", "strict")

	if got := cookieSameSite(); got != http.SameSiteStrictMode {
		t.Fatalf("cookieSameSite() = %v, want %v", got, http.SameSiteStrictMode)
	}
}
