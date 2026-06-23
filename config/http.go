package config

import (
	"net/url"
	"os"
	"strings"
)

var defaultAllowedOrigins = []string{
	"http://localhost:3000",
	"http://127.0.0.1:3000",
}

func AllowedOrigins() []string {
	if origins := parseOrigins(os.Getenv("CORS_ALLOWED_ORIGINS")); len(origins) > 0 {
		return origins
	}

	if origins := parseOrigins(os.Getenv("FRONTEND_URL")); len(origins) > 0 {
		return origins
	}

	return append([]string(nil), defaultAllowedOrigins...)
}

func GoogleRedirectURL() string {
	if redirectURL := normalizeURL(os.Getenv("GOOGLE_REDIRECT_URL")); redirectURL != "" {
		return redirectURL
	}

	return "http://localhost:8080/api/auth/google/callback"
}

func FrontendURL() string {
	return normalizeURL(os.Getenv("FRONTEND_URL"))
}

func FrontendAuthSuccessURL() string {
	if redirectURL := normalizeURL(os.Getenv("FRONTEND_AUTH_SUCCESS_URL")); redirectURL != "" {
		return redirectURL
	}

	return FrontendURL()
}

func FrontendAuthErrorURL() string {
	if redirectURL := normalizeURL(os.Getenv("FRONTEND_AUTH_ERROR_URL")); redirectURL != "" {
		return redirectURL
	}

	return FrontendAuthSuccessURL()
}

func ServerPort() string {
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		return port
	}

	return "8080"
}

func parseOrigins(value string) []string {
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		origin := normalizeOrigin(part)
		if origin == "" {
			continue
		}

		if _, exists := seen[origin]; exists {
			continue
		}

		seen[origin] = struct{}{}
		origins = append(origins, origin)
	}

	return origins
}

func normalizeOrigin(value string) string {
	value = strings.TrimSpace(strings.TrimRight(value, "/"))
	if value == "" {
		return ""
	}

	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return value
	}

	return parsed.Scheme + "://" + parsed.Host
}

func normalizeURL(value string) string {
	return strings.TrimSpace(strings.TrimRight(value, "/"))
}
