package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Cookies    []string          `json:"cookies"`
}

func HandleRequest(ctx context.Context) (Response, error) {
	return Response{
		StatusCode: 200,
		Headers: map[string]string{
			"content-type": "text/html; charset=utf-8",
		},
		Body:    "hello, world",
		Cookies: nil,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
