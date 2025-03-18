package routes

import (
	"library-api/controllers"
	"library-api/middleware"
	"github.com/gorilla/mux"
)

func SetupRouter(userCtrl *controllers.UserController, bookCtrl *controllers.BookController) *mux.Router {
	router := mux.NewRouter()
	// Apply rate limiting to all routes
	router.Use(middleware.RateLimit(10, 20)) // 10 req/s, burst 20

	router.HandleFunc("/users", userCtrl.CreateUser).Methods("POST")
	router.HandleFunc("/login", userCtrl.Login).Methods("POST")
	router.HandleFunc("/books", bookCtrl.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", bookCtrl.GetBook).Methods("GET")
	router.HandleFunc("/books", middleware.Authenticate(bookCtrl.CreateBook)).Methods("POST")
	router.HandleFunc("/books/{id}", middleware.Authenticate(bookCtrl.UpdateBook)).Methods("PUT")
	router.HandleFunc("/books/{id}", middleware.Authenticate(bookCtrl.DeleteBook)).Methods("DELETE")
	router.HandleFunc("/bulk-emails", middleware.Authenticate(userCtrl.SendBulkEmails)).Methods("POST")
	router.HandleFunc("/bulk-emails/status", middleware.Authenticate(userCtrl.GetBulkEmailStatus)).Methods("GET")
	router.HandleFunc("/bulk-emails/stream", middleware.Authenticate(userCtrl.StreamBulkEmailStatus)).Methods("GET")
	return router
}