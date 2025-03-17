package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"library-api/config"
	"library-api/logger" // Correct import
	"time"
)

var DB *gorm.DB

func InitDB(cfg *config.Config, logger *logger.AsyncLogger) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	var err error
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := DB.DB()
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)
			logger.Log("Connected to remote database!") // Line 25
			return
		}
		logger.Log(fmt.Sprintf("Retrying DB connection (%d/10): %v", i+1, err)) // Line 28
		time.Sleep(2 * time.Second)
	}
	panic("Failed to connect to database!")
}