package main

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
	"net/http"
)

func renderQuiz(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("content-type", ContentTypeHtml)
	err := xray.Capture(req.Context(), "render template", func(ctx context.Context) error {
		return templates.ExecuteTemplate(rw, "view_quiz.html.tmpl", map[string]any{"Quiz": Quiz})
	})
	if err != nil {
		logger(req).Error("could not render template: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
