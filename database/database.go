// Package database provides functionality for connecting to and interacting with
// the PostgreSQL database used by the Filia application.
package database

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm/logger"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DbConnection is a global variable that holds the database connection.
// It is initialized by the Connect function and used throughout the application
// to interact with the database.
var DbConnection *gorm.DB

func buildDSN() string {
	if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
		return databaseURL
	}

	dbName := strings.TrimSpace(os.Getenv("DB_NAME"))
	if dbName == "" {
		dbName = "filia"
	}

	dbPort := strings.TrimSpace(os.Getenv("DB_PORT"))
	if dbPort == "" {
		dbPort = "5432"
	}

	sslMode := strings.TrimSpace(os.Getenv("DB_SSLMODE"))
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		dbName,
		dbPort,
		sslMode,
	)
}

// Connect establishes a connection to the PostgreSQL database.
// It prefers DATABASE_URL and falls back to DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT, and DB_SSLMODE.
// If the connection fails, it logs a fatal error and terminates the application.
// Upon successful connection, it assigns the database connection to the global DbConnection variable.
func Connect() {
	log.Debug("Connecting to the filia database.....")
	dsn := buildDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
		return
	}

	log.Debug("Connected to the filia database.")
	DbConnection = db
	generateSeeds()
}
