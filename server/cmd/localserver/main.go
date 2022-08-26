package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mleone10/expense-system/adapters/auth"
	"github.com/mleone10/expense-system/adapters/rest"
	"github.com/mleone10/expense-system/adapters/stdlogger"
)

func main() {
	authClient := auth.NewAuthClient(
		auth.WithClientHostname("localhost:3000"),
		auth.WithClientScheme("http"),
		auth.WithCognitoClientId("6ka3m790cv5hrhjdqt2ju89v45"),
		auth.WithCognitoClientSecret(os.Getenv("COGNITO_CLIENT_SECRET")),
	)

	logger := stdlogger.Logger{}

	server, _ := rest.NewServer(
		rest.WithAuthClient(authClient),
		rest.WithLogger(logger),
	)

	fmt.Println(http.ListenAndServe(":8080", server))
}
