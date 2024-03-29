package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"net/http"
	"net/url"
)

const (
	ContentTypeHtml = "text/html; charset=utf-8"
	ContentTypeText = "text/plain; charset=utf-8"
)

type Request struct {
	RawPath               string            `json:"rawPath"`
	RawQueryString        string            `json:"rawQueryString"`
	Cookies               []string          `json:"cookies"`
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	Body                  []byte            `json:"body"`
	RequestContext        struct {
		HTTP struct {
			Method   string `json:"method"`
			SourceIp string `json:"sourceIp"`
		} `json:"http"`
	} `json:"requestContext"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Headers    http.Header `json:"-"`
	Body       string      `json:"body"`
	Cookies    []string    `json:"cookies"`
}

func (r *Response) MarshalJSON() ([]byte, error) {
	headers := map[string]string{}
	for header, values := range r.Headers {
		if len(values) == 0 {
			continue
		}
		headers[header] = values[len(values)-1]
	}
	return json.Marshal(map[string]any{
		"statusCode": r.StatusCode,
		"body":       r.Body,
		"cookies":    r.Cookies,
		"headers":    headers,
	})
}

func (r *Response) Header() http.Header {
	return r.Headers
}

func (r *Response) Write(bytes []byte) (int, error) {
	r.Body += string(bytes)
	return len(bytes), nil
}

func (r *Response) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
}

var _ http.ResponseWriter = (*Response)(nil)
var _ json.Marshaler = (*Response)(nil)

func HandleRequest(request *Request) (res *Response, err error) {
	req := &http.Request{
		Method: request.RequestContext.HTTP.Method,
		URL: &url.URL{
			Scheme:   "https",
			Path:     request.RawPath,
			RawPath:  request.RawPath,
			RawQuery: request.RawQueryString,
		},
		Header:        map[string][]string{},
		Body:          io.NopCloser(bytes.NewBuffer(request.Body)),
		ContentLength: int64(len(request.Body)),
		RemoteAddr:    request.RequestContext.HTTP.SourceIp,
	}

	router := http.NewServeMux()
	router.HandleFunc("/", renderQuiz)

	res = new(Response)
	res.Headers = http.Header{}
	router.ServeHTTP(res, req)

	if res.StatusCode == 0 {
		res.StatusCode = 200
	}

	if _, ok := res.Headers["content-type"]; !ok {
		req.Header.Set("content-type", http.DetectContentType([]byte(res.Body)))
	}

	return res, nil
}

func main() {
	lambda.Start(HandleRequest)
}
