package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
	Path      string `json:"path,omitempty"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error"`
	RequestID string      `json:"request_id,omitempty"`
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// ErrorHandlerMiddleware provides centralized error handling
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors if any were set
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			requestID, _ := c.Get("request_id")

			// Determine status code based on error type
			statusCode := http.StatusInternalServerError
			errorMsg := err.Error()

			if strings.Contains(errorMsg, "external") {
				statusCode = http.StatusBadRequest
				errorMsg = errorMsg[10:] // Remove "external: " prefix
			} else if strings.Contains(errorMsg, "internal") {
				statusCode = http.StatusInternalServerError
				errorMsg = errorMsg[10:] // Remove "internal: " prefix
				log.WithFields(log.Fields{
					"request_id": requestID,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}).Error("Internal error: ", err.Error())
			}

			c.JSON(statusCode, ErrorResponse{
				Error:     errorMsg,
				RequestID: requestID.(string),
				Path:      c.Request.URL.Path,
			})
			return
		}
	}
}

// RecoveryMiddleware recovers from panics and returns a standardized error response
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, exists := c.Get("request_id")
				requestIDStr := ""
				if exists {
					requestIDStr = requestID.(string)
				}

				log.WithFields(log.Fields{
					"request_id": requestIDStr,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
					"panic":      err,
				}).Error("Panic recovered")

				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:     "Internal server error",
					RequestID: requestIDStr,
					Path:      c.Request.URL.Path,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
