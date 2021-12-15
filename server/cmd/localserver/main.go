package main

import (
	"log"
	"net/http"
	"os"

	api "github.com/mleone10/expense-system"
)

func main() {
	server, err := api.NewServer(api.Config{
		CognitoClientId:     "6ka3m790cv5hrhjdqt2ju89v45",
		CognitoClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
		ClientHostname:      "localhost:3000",
		ClientScheme:        "http",
	})
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	http.ListenAndServe(":8080", server)
}
