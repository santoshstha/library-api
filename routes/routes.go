package routes

import (
	"library-api/controllers"
	"library-api/middleware"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/books", controllers.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", controllers.GetBook).Methods("GET")
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	// Protected routes
	router.HandleFunc("/books", middleware.Authenticate(controllers.CreateBook)).Methods("POST")
	router.HandleFunc("/books/{id}", middleware.Authenticate(controllers.UpdateBook)).Methods("PUT")
	router.HandleFunc("/books/{id}", middleware.Authenticate(controllers.DeleteBook)).Methods("DELETE")

	return router
}