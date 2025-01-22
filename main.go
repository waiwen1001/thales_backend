package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"waiwen.com/thales-backend/config"
	"waiwen.com/thales-backend/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Init Database failed: %v", err)
	}

	defer db.Close()

	router := routes.InitRoutes(db)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Content-Type", "Authorization", "application/json"}),
		handlers.AllowCredentials(),
	)(router))
}
