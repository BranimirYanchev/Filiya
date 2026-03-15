package tests

import (
	"github.com/Marionvd/filia-project-backend/config"
	"github.com/Marionvd/filia-project-backend/database"
	"github.com/Marionvd/filia-project-backend/internal/helper"
	"github.com/Marionvd/filia-project-backend/internal/model"
	"github.com/Marionvd/filia-project-backend/router"
	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func mockApplicationRoutes() *gin.Engine {
	log.SetLevel(log.DebugLevel)
	gin.SetMode(gin.DebugMode)
	log.Debug("Loading env...")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	config.GoogleConfig = config.SetUpGoogleConfig()

	database.Connect()

	engine := gin.Default()

	log.Debug("Setting up routers...")
	engine.Use(cors.Default())

	router.SetUpRoutes(engine)

	return engine
}

func mockStandardTestToken() string {
	var user model.User
	if err := database.DbConnection.Where("role_id=?", 3).First(&user).Error; err != nil {
		log.Error(err.Error())
		return ""
	}

	return helper.GenerateJWT(user.GenerateJwtUser())
}

func strptr(str string) *string {
	return &str
}
