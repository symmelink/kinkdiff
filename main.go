package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	CommitHash     string
	BuildTimestamp string
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
	StatusCode int
	Headers    http.Header
	Body       *bytes.Buffer
	Cookies    []string
	ctx        context.Context
}

func (r *Response) MarshalJSON() ([]byte, error) {
	if r.Body == nil {
		r.Body = &bytes.Buffer{}
	}

	headers := map[string]string{}
	for header, values := range r.Headers {
		if len(values) == 0 {
			continue
		}
		headers[header] = values[len(values)-1]
	}
	_, seg := xray.BeginSubsegment(r.ctx, "*Response.MarshalJSON")
	jsonBytes, err := json.Marshal(map[string]any{
		"statusCode": r.StatusCode,
		"body":       r.Body.String(),
		"cookies":    r.Cookies,
		"headers":    headers,
	})
	seg.CloseAndStream(err)
	return jsonBytes, err
}

func (r *Response) Header() http.Header {
	return r.Headers
}

func (r *Response) Write(data []byte) (int, error) {
	if r.Body == nil {
		r.Body = &bytes.Buffer{}
	}
	return r.Body.Write(data)
}

func (r *Response) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
}

var _ http.ResponseWriter = (*Response)(nil)
var _ json.Marshaler = (*Response)(nil)

func withXray(name string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx, seg := xray.BeginSubsegment(req.Context(), name)
		handlerFunc(rw, req.WithContext(ctx))
		seg.CloseAndStream(nil)
	}
}

func HandleRequest(ctx context.Context, request *Request) (res *Response, err error) {
	ctx, seg := xray.BeginSubsegment(ctx, "HandleRequest")
	defer seg.CloseAndStream(nil)

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
	req = req.WithContext(ctx)

	router := http.NewServeMux()
	router.HandleFunc("/", withXray("renderQuiz", renderQuiz))
	router.HandleFunc("/static/css/", withXray("render css", http.FileServerFS(cssFs).ServeHTTP))

	res = new(Response)
	res.ctx = ctx
	res.Headers = http.Header{}
	router.ServeHTTP(res, req)

	if res.StatusCode == 0 {
		res.StatusCode = 200
	}

	if _, ok := res.Headers["content-type"]; !ok {
		req.Header.Set("content-type", http.DetectContentType(res.Body.Bytes()))
	}

	return res, nil
}

func init() {
	err := xray.Configure(xray.Config{
		ServiceVersion: CommitHash + "-" + BuildTimestamp,
	})
	if err != nil {
		panic(err)
	}
	xray.SetLogger(xraylog.NewDefaultLogger(os.Stderr, xraylog.LogLevelError))
}

func main() {
	lambda.Start(HandleRequest)
}
