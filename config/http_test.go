package config

import "testing"

func TestAllowedOriginsPrefersConfiguredOrigins(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://frontend.example.com/, https://www.frontend.example.com")
	t.Setenv("FRONTEND_URL", "https://fallback.example.com")

	origins := AllowedOrigins()
	if len(origins) != 2 {
		t.Fatalf("expected 2 origins, got %d (%v)", len(origins), origins)
	}

	if origins[0] != "https://frontend.example.com" {
		t.Fatalf("unexpected first origin: %s", origins[0])
	}

	if origins[1] != "https://www.frontend.example.com" {
		t.Fatalf("unexpected second origin: %s", origins[1])
	}
}

func TestAllowedOriginsFallsBackToFrontendURL(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	t.Setenv("FRONTEND_URL", "https://frontend.example.com/app")

	origins := AllowedOrigins()
	if len(origins) != 1 {
		t.Fatalf("expected 1 origin, got %d (%v)", len(origins), origins)
	}

	if origins[0] != "https://frontend.example.com" {
		t.Fatalf("unexpected origin: %s", origins[0])
	}
}

func TestFrontendAuthURLs(t *testing.T) {
	t.Setenv("FRONTEND_URL", "https://frontend.example.com")
	t.Setenv("FRONTEND_AUTH_SUCCESS_URL", "https://frontend.example.com/auth/callback/")
	t.Setenv("FRONTEND_AUTH_ERROR_URL", "")

	if got := FrontendAuthSuccessURL(); got != "https://frontend.example.com/auth/callback" {
		t.Fatalf("unexpected success redirect url: %s", got)
	}

	if got := FrontendAuthErrorURL(); got != "https://frontend.example.com/auth/callback" {
		t.Fatalf("unexpected error redirect url: %s", got)
	}
}
