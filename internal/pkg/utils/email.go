package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, body string) error {
	// Get SMTP config
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	senderName := os.Getenv("SMTP_SENDER_NAME")
	
	// Format message
	msg := []byte(fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
		"\r\n"+
		"%s\r\n", senderName, "no-reply@eshop.com", to, subject, body))

	// Send via SMTP
	addr := fmt.Sprintf("%s:%s", host, port)
	
	// Auth
	var auth smtp.Auth
	if os.Getenv("SMTP_AUTH_EMAIL") != "" {
		auth = smtp.PlainAuth("", os.Getenv("SMTP_AUTH_EMAIL"), os.Getenv("SMTP_AUTH_PASSWORD"), host)
	}

	err := smtp.SendMail(addr, auth, "no-reply@eshop.com", []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}