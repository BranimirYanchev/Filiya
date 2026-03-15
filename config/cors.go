package config

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
)

var defaultAllowedOrigins = []string{
	"http://localhost:3000",
	"http://127.0.0.1:3000",
	"https://filiya-frontend.onrender.com",
}

// CORSConfig returns the shared CORS configuration for the API.
func CORSConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     allowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-Id", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

func allowedOrigins() []string {
	var origins []string

	if frontendURL := strings.TrimSpace(os.Getenv("FRONTEND_URL")); frontendURL != "" {
		origins = append(origins, frontendURL)
	}

	if configuredOrigins := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); configuredOrigins != "" {
		for _, origin := range strings.Split(configuredOrigins, ",") {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				origins = append(origins, origin)
			}
		}
	}

	if len(origins) == 0 {
		return append([]string(nil), defaultAllowedOrigins...)
	}

	seen := make(map[string]struct{}, len(origins))
	deduped := make([]string, 0, len(origins))
	for _, origin := range origins {
		if _, exists := seen[origin]; exists {
			continue
		}
		seen[origin] = struct{}{}
		deduped = append(deduped, origin)
	}

	return deduped
}
