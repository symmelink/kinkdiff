package main

import (
	"embed"
	"html/template"
)

//go:embed static/templates/*
var templateFS embed.FS
var templates = template.Must(template.ParseFS(templateFS, "static/templates/*"))
