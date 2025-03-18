package main

import (
	"context"
	"log"
	"net/http"
	"library-api/cache"
	"library-api/config"
	"library-api/controllers"
	"library-api/database"
	"library-api/email"
	"library-api/logger"
	"library-api/models"
	"library-api/repositories"
	"library-api/routes"
	"library-api/services"
)

var Logger *logger.AsyncLogger

func main() {
	cfg := config.LoadConfig()
	Logger = logger.NewAsyncLogger()
	defer Logger.Close()

	database.InitDB(cfg, Logger)
	database.DB.AutoMigrate(&models.User{}, &models.Book{})

	redisCache := cache.NewCache(cfg.RedisAddr)
	emailService := email.NewEmailService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, 10, Logger) // Pass Logger
	defer emailService.Shutdown(context.Background())

	userRepo := repositories.NewUserRepository(database.DB)
	bookRepo := repositories.NewBookRepository(database.DB)
	userService := services.NewUserService(userRepo, emailService)
	bookService := services.NewBookService(bookRepo, redisCache)
	userCtrl := controllers.NewUserController(userService)
	bookCtrl := controllers.NewBookController(bookService)

	router := routes.SetupRouter(userCtrl, bookCtrl)
	Logger.Log("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}