package otp

import (
	"sync"
	"text/template"
)

import _ "embed"

//go:embed otp.html
var otpHTML string

var (
	otpOnce sync.Once
	otpTmpl *template.Template
)

type Data struct {
	Name   string
	OTP    string
	Expire int
}

func (Data) GetTemplate() *template.Template {
	otpOnce.Do(func() {
		tmpl := template.New("otp")
		otpTmpl, _ = tmpl.Parse(otpHTML)
	})
	return otpTmpl
}