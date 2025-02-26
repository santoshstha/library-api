package main

import (
	"log"
	"net/http"
	"library-api/database"
	"library-api/models"
	"library-api/routes"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(&models.User{}, &models.Book{}) // Creates tables
	router := routes.SetupRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}