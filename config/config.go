package config

import (
	"os"
	"strconv"
)

type Config struct {
	JWTSecret  string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	RedisAddr  string
	SMTPHost   string
	SMTPPort   string
	SMTPUser   string
	SMTPPass   string
}

func LoadConfig() *Config {
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	return &Config{
		JWTSecret:  os.Getenv("JWT_SECRET"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		RedisAddr:  os.Getenv("REDIS_ADDR"),
		SMTPHost:   os.Getenv("SMTP_HOST"),
		SMTPPort:   os.Getenv("SMTP_PORT"),
		SMTPUser:   os.Getenv("SMTP_USER"),
		SMTPPass:   os.Getenv("SMTP_PASS"),
	}
}