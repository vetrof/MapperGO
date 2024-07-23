package main

import (
	"log"
	"net/http"
	"os"

	"gomap/internal/bot"
	"gomap/internal/db"
	"gomap/internal/router"
)

func main() {

	// db init
	db.InitPostgresDB()

	// Bot init
	bot.Bot()

	// Server init
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
