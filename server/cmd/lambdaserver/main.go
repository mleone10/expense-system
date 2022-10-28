package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	// server, err := api.NewServer(api.Config{
	// 	CognitoClientId:     "6ka3m790cv5hrhjdqt2ju89v45",
	// 	CognitoClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
	// 	ClientHostname:      "expense.mleone.dev",
	// 	ClientScheme:        "https",
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to initialize server: %v", err)
	// }

	// adapter := httpadapter.New(server)

	// lambda.Start(serverHandler(adapter))
}

func serverHandler(adapter *httpadapter.HandlerAdapter) func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return adapter.Proxy(req)
	}
}
