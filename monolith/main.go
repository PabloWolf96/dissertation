package main

import (
	"log"
	"net/http"

	"monolith/database"
	"monolith/handlers"
	"monolith/middleware"

	"github.com/gorilla/mux"
)

func main() {
    // Initialize database
    err := database.Initialize()
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    

    r := mux.NewRouter()

    // Public routes
    r.HandleFunc("/signup", handlers.SignUp).Methods("POST")
    r.HandleFunc("/login", handlers.Login).Methods("POST")
    r.HandleFunc("/products/search", handlers.SearchProducts).Methods("GET")
    r.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")
    r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

    // Protected routes
    s := r.PathPrefix("").Subrouter()
    s.Use(middleware.JwtAuthentication)
    s.HandleFunc("/cart", handlers.AddToCart).Methods("POST")
    s.HandleFunc("/cart/checkout", handlers.Checkout).Methods("POST")
    log.Println("Server is running on port 8000")
    log.Fatal(http.ListenAndServe(":8000", r))
}