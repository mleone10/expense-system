package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/mleone10/expense-system/adapters/cognito"
	"github.com/mleone10/expense-system/adapters/dynamodb"
	"github.com/mleone10/expense-system/adapters/rest"
	"github.com/mleone10/expense-system/adapters/stdlogger"
	"github.com/mleone10/expense-system/service"
)

func main() {
	authClient := cognito.NewAuthClient(
		cognito.WithClientHostname("expense.mleone.dev"),
		cognito.WithClientScheme("https"),
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
	)

	adapter := httpadapter.New(server)

	lambda.Start(serverHandler(adapter))
}

func serverHandler(adapter *httpadapter.HandlerAdapter) func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return adapter.Proxy(req)
	}
}
