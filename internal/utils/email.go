package utils

import (
	"crypto/tls"

	"github.com/jordanknott/monitor/internal/config"
	gomail "gopkg.in/mail.v2"
)

type Email struct {
	HTML  string
	Plain string
	To    string
}

func SendMail(cfg config.EmailConfig, email Email, subject string) error {

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", cfg.From)

	// Set E-Mail receivers
	m.SetHeader("To", email.To)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", email.HTML)

	// Settings for SMTP server
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{ServerName: "email-smtp.us-west-2.amazonaws.com"}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
