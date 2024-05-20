package main

import (
	"embed"
	_ "embed"
	"github.com/ghodss/yaml"
	"log"
)

//go:embed static/css
var cssFs embed.FS

//go:embed static/quiz.yaml
var quizYaml []byte

var Quiz []*QuizCategory

type QuizCategory struct {
	Title string
	Verbs []string
	Items []*QuizItem
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
	for _, cat := range Quiz {
		if len(cat.Verbs) == 0 {
			cat.Verbs = []string{"giving", "receiving"}
		}
	}
}
