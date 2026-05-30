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

type Data struct {
 	Name string
	Link string
}

func (w Data) GetSubject() string {
	return "Welcome to our email campaign"
}

func (w Data) GetTemplate() *template.Template {
	welcomeOnce.Do(func() {
		tmpl := template.New("welcome")
		welcomeTmpl, _ = tmpl.Parse(welcomeHTML)
	})
	return welcomeTmpl
}