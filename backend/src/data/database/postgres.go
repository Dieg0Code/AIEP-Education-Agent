package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDB() *gorm.DB {
	// Get enviroment variables
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	//GORM configuration
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create vector extension if not exists
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Println("Vector extension already exists")
		} else {
			log.Fatalf("Failed to create vector extension: %v", err)
		}
	}

	if err := models.AutoMigrateAll(db); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	return db
}
