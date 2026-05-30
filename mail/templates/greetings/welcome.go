package greetings

import (
	"sync"
	"text/template"
)

import _ "embed"

//go:embed welcome.html
var welcomeHTML string

var (
	welcomeOnce sync.Once
	welcomeTmpl *template.Template
)

type WelcomeData struct {
 	Name string
	Link string
}

func (w WelcomeData) GetSubject() string {
	return "Welcome to our email campaign"
}

func (w WelcomeData) GetTemplate() *template.Template {
	welcomeOnce.Do(func() {
		tmpl := template.New("welcome")
		welcomeTmpl, _ = tmpl.Parse(welcomeHTML)
	})
	return welcomeTmpl
}