package main

import (
	"net/http"
)

func renderQuiz(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("content-type", ContentTypeHtml)
	if err := templates.ExecuteTemplate(rw, "view_quiz.html.tmpl", map[string]any{"Quiz": Quiz}); err != nil {
		logger(req).Error("could not render template: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
