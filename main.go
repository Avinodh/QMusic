package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	router := NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}
