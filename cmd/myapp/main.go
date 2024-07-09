package main

import (
	"gomap/internal/db"
	"gomap/internal/router"
	"log"
	"net/http"
	"os"
)

func main() {

	db.InitPostgresDB()

	r := router.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5555"
	}

	log.Println("Starting server at port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
