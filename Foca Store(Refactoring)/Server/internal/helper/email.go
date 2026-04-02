package helper

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type ResetPasswordTemplateData struct {
	OTP string
	ExpireMinutes int
}

func SendOTPEmail(to, otp string) error {
	form := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := 587

	templatePath := filepath.Join("helper", "templates", "reset_password.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed parse template: %w", err)
	}

	data := ResetPasswordTemplateData{
		OTP: otp,
		ExpireMinutes: 5,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed execute template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", form)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset Password OTP")

	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(host, port, form, password)

	return d.DialAndSend(m)
}
