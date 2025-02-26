package main

import (
	"log"
	"net/http"
	"library-api/cache"
	"library-api/config"
	"library-api/controllers"
	"library-api/database"
	"library-api/logger"
	"library-api/models"
	"library-api/repositories"
	"library-api/routes"
	"library-api/services"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger() // Initialize the global logger
	defer logger.Logger.Close()

	database.InitDB(cfg)
	database.DB.AutoMigrate(&models.User{}, &models.Book{})

	redisCache := cache.NewCache(cfg.RedisAddr)
	userRepo := repositories.NewUserRepository(database.DB)
	bookRepo := repositories.NewBookRepository(database.DB)
	userService := services.NewUserService(userRepo)
	bookService := services.NewBookService(bookRepo, redisCache)
	userCtrl := controllers.NewUserController(userService)
	bookCtrl := controllers.NewBookController(bookService)

	router := routes.SetupRouter(userCtrl, bookCtrl)
	logger.Logger.Log("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}