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
	OTP           string
	ExpireMinutes int
}

func SendOTPEmail(to, otp string) error {
	form := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := 587

	// Cek environment variable
	if form == "" || password == "" || host == "" {
		fmt.Println("⚠️  SMTP not configured properly")
		return fmt.Errorf("smtp not configured")
	}

	wd, _ := os.Getwd()
	templatePath := filepath.Join(wd, "internal", "helper", "templates", "reset_password.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("⚠️ Failed parse template:", err)
		return fmt.Errorf("failed to load email template")
	}

	data := ResetPasswordTemplateData{
		OTP:           otp,
		ExpireMinutes: 5,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		fmt.Println("⚠️ Failed execute template:", err)
		return fmt.Errorf("failed to render email template")
	}

	// Build message
	m := gomail.NewMessage()
	m.SetHeader("From", form)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset Password OTP")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(host, port, form, password)

	// Kirim Email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("⚠️ Failed sending email:", err)
		return fmt.Errorf("failed to send OTP email")
	}

	return nil
}
