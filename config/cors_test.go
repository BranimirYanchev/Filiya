package config

import (
	"reflect"
	"testing"
)

func TestAllowedOriginsFallsBackToLocalDefaults(t *testing.T) {
	t.Setenv("FRONTEND_URL", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	got := allowedOrigins()

	if !reflect.DeepEqual(got, defaultAllowedOrigins) {
		t.Fatalf("allowedOrigins() = %v, want %v", got, defaultAllowedOrigins)
	}
}

func TestAllowedOriginsUsesConfiguredOrigins(t *testing.T) {
	t.Setenv("FRONTEND_URL", "https://filiya-frontend.onrender.com")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://www.filiya.com, https://app.filiya.com, https://www.filiya.com")

	got := allowedOrigins()
	want := []string{
		"https://filiya-frontend.onrender.com",
		"https://www.filiya.com",
		"https://app.filiya.com",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("allowedOrigins() = %v, want %v", got, want)
	}
}
