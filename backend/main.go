package main

import (
	"JovaCentral/database"
	"JovaCentral/goControllers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	database.InitializeDatabase()
	router := mux.NewRouter()
	goControllers.SetupUserRoutes(router)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Println(":: Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", cors(router)))
}
