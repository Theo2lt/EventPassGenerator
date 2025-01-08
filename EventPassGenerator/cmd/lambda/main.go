package main

import (
	"EventPassGenerator/internal/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.LambdaHandler)
}
