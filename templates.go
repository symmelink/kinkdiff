package main

import (
	"embed"
	"html/template"
)

//go:embed static/templates/*
var templateFS embed.FS
var funcs = map[string]any{
	"N": func(n int) []int {
		ret := make([]int, n)
		for i := 0; i < n; i++ {
			ret[i] = i
		}
		return ret
	},
}
var templates = template.Must(
	template.New("").
		Funcs(funcs).
		ParseFS(templateFS, "static/templates/*"),
)
