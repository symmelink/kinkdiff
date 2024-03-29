package main

import (
	"log/slog"
	"net/http"
)

func logger(req *http.Request) *slog.Logger {
	return slog.Default().With(
		"client_ip", req.RemoteAddr,
		"path", req.URL.RawPath,
	)
}
