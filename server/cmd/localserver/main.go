package main

import (
	"log"
	"net/http"

	api "github.com/mleone10/expense-system"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	http.ListenAndServe(":8080", server)
}
