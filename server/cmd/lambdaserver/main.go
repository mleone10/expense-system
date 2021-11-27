package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	api "github.com/mleone10/expense-system"
)

func main() {
	// authr, err := firebase.NewAuthenticator()
	// if err != nil {
	// 	log.Panicf("Failed to initialize authenticator: %v", err)
	// }
	server, err := api.NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	adapter := httpadapter.New(server)

	lambda.Start(serverHandler(adapter))
}

func serverHandler(adapter *httpadapter.HandlerAdapter) func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return adapter.Proxy(req)
	}
}
