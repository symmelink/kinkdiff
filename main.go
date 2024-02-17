package main

import (
	"bytes"
	_ "embed"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ghodss/yaml"
	"io"
	"log"
	"net/http"
	"net/url"
)

//go:embed quiz.yaml
var quizYaml []byte

var Quiz map[string]QuizCategory

type QuizCategory struct {
	Title string
	Verbs []string
	Items []QuizItem
}

type QuizItem struct {
	Name        string
	Description string
	Version     int
}

func init() {
	if err := yaml.Unmarshal(quizYaml, &Quiz); err != nil {
		log.Fatal(err)
	}
}

type Request struct {
	RawPath               string            `json:"rawPath"`
	RawQueryString        string            `json:"rawQueryString"`
	Cookies               []string          `json:"cookies"`
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	Body                  []byte            `json:"body"`
	RequestContext        struct {
		HTTP struct {
			Method string `json:"method"`
		} `json:"http"`
	} `json:"requestContext"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Cookies    []string          `json:"cookies"`
}

func HandleRequest(request Request) (Response, error) {
	req := http.Request{
		Method: request.RequestContext.HTTP.Method,
		URL: &url.URL{
			Path:     request.RawPath,
			RawQuery: request.RawQueryString,
		},
	}

	if len(request.Body) > 0 {
		req.Body = io.NopCloser(bytes.NewBuffer(request.Body))
		req.ContentLength = int64(len(request.Body))
	}

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
