package templates

import (
	"bytes"
	"text/template"
)

type Template interface {
	GetTemplate() *template.Template
	GetSubject() string
}

func Render(data Template) (string,string,error) {
	var buf bytes.Buffer
	err := data.GetTemplate().Execute(&buf, data)
	return buf.String(), data.GetSubject(), err
}