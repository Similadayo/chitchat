package config

import (
	"fmt"
	"github/similadayo/chitchat/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDB establishes a connection to the database
func ConnectDB() (*gorm.DB, error) {
	LoadEnv()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// MigrateDB performs database migration
func MigrateDB(db *gorm.DB) {

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
}

// CloseDB closes the database connection
