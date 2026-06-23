// Package main is the entry point for the Filia application.
// It initializes the application, sets up the database connection,
// configures the API routes, and starts the HTTP server.
//
// @title           Filia Project Backend API
// @version         1.0
// @description     A comprehensive social application backend API built with Go and Gin framework. Provides endpoints for user management, authentication, posts, comments, categories, roles, and friend management.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/Marionvd/filia-project-backend
// @contact.email  support@filia.example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"fmt"
	"time"

	"github.com/Marionvd/filia-project-backend/config"
	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/middleware"
	_ "github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/router"
	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.SetLevel(log.DebugLevel)
	gin.SetMode(gin.DebugMode)
	log.Debug("Loading env...")
	err := godotenv.Load()
	if err != nil {
		log.Warnf("No .env file loaded: %v", err)
	}

	config.GoogleConfig = config.SetUpGoogleConfig()

	database.Connect()

	engine := gin.Default()

	log.Debug("Setting up routers...")
	engine.Use(cors.New(cors.Config{
		AllowOrigins: config.AllowedOrigins(),
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept", "Authorization"},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-Id",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add middleware
	engine.Use(middleware.RequestIDMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.ErrorHandlerMiddleware())
	engine.Static("/uploads", "./uploads")

	router.SetUpRoutes(engine)

	port := config.ServerPort()
	log.Debugf("API running on port %s with allowed origins %v", port, config.AllowedOrigins())
	engine.Run(fmt.Sprintf(":%s", port))
}
