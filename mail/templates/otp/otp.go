package greetings

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

type OTPData struct {
	Name   string
	OTP    string
	Expire int
}

func (OTPData) GetTemplate() *template.Template {
	otpOnce.Do(func() {
		tmpl := template.New("otp")
		otpTmpl, _ = tmpl.Parse(otpHTML)
	})
	return otpTmpl
}