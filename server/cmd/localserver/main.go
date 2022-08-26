package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mleone10/expense-system/adapters/dynamodb"
	"github.com/mleone10/expense-system/adapters/googleauth"
	"github.com/mleone10/expense-system/adapters/rest"
	"github.com/mleone10/expense-system/adapters/stdlogger"
	"github.com/mleone10/expense-system/service"
)

func main() {
	authClient := googleauth.NewAuthClient(
		googleauth.WithClientHostname("localhost:3000"),
		googleauth.WithClientScheme("http"),
		googleauth.WithCognitoClientId("6ka3m790cv5hrhjdqt2ju89v45"),
		googleauth.WithCognitoClientSecret(os.Getenv("COGNITO_CLIENT_SECRET")),
	)

	orgRepo, _ := dynamodb.NewClient()

	orgService := service.NewOrgService(orgRepo)

	server, _ := rest.NewServer(
		rest.WithAuthClient(authClient),
		rest.WithOrgService(orgService),
		rest.WithLogger(stdlogger.Logger{}),
	)

	fmt.Println(http.ListenAndServe(":8080", server))
}
