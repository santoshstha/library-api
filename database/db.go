package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func Connect() {
	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Build DSN for remote MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	
	var err error
	for i := 0; i < 10; i++ { // Retry up to 10 times
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to remote database!")
			return
		}
		log.Printf("Failed to connect to database, retrying (%d/10): %v", i+1, err)
		time.Sleep(2 * time.Second) // Wait 2 seconds before retrying
	}
	panic("Failed to connect to remote database after retries!")
}