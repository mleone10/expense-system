package main

import (
	"net/http"
	"os"

	"github.com/mleone10/expense-system/adapters/cognito"
	"github.com/mleone10/expense-system/adapters/dynamodb"
	"github.com/mleone10/expense-system/adapters/rest"
	"github.com/mleone10/expense-system/adapters/stdlogger"
	"github.com/mleone10/expense-system/service"
)

func main() {
	authClient := cognito.NewAuthClient(
		cognito.WithClientHostname("localhost:3000"),
		cognito.WithClientScheme("http"),
		cognito.WithCognitoClientId("6ka3m790cv5hrhjdqt2ju89v45"),
		cognito.WithCognitoClientSecret(os.Getenv("COGNITO_CLIENT_SECRET")),
	)

	orgRepo, _ := dynamodb.NewClient()

	orgService := service.NewOrgService(orgRepo)

	authenticatedUserService := service.NewAuthenticatedUserService(authClient)

	logger := stdlogger.NewLogger()

	server, _ := rest.NewServer(
		rest.WithAuthClient(authClient),
		rest.WithOrgService(orgService),
		rest.WithAuthenticatedUserService(authenticatedUserService),
		rest.WithLogger(logger),
		rest.WithSkipAuth(),
	)

	http.ListenAndServe(":8080", server)
}
