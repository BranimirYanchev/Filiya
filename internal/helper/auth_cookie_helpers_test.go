package helper

import (
	"net/http"
	"testing"
)

func TestCookieSameSiteDefaultsToLax(t *testing.T) {
	t.Setenv("COOKIE_SAME_SITE", "")

	if got := cookieSameSite(); got != http.SameSiteLaxMode {
		t.Fatalf("expected lax same-site mode, got %v", got)
	}
}

func TestCookieSameSiteSupportsNone(t *testing.T) {
	t.Setenv("COOKIE_SAME_SITE", "none")

	if got := cookieSameSite(); got != http.SameSiteNoneMode {
		t.Fatalf("expected none same-site mode, got %v", got)
	}
}
