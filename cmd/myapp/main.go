package main

import (
	"framrless/internal/router"
	"log"
	"net/http"
	"os"
)

func main() {
	r := router.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server at port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
