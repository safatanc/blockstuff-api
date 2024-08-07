package mail

import (
	"crypto/tls"
	"os"

	"gopkg.in/gomail.v2"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(to []string, subject string, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	d := gomail.NewDialer("smtp.safatanc.com", 587, os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"))

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	err := d.DialAndSend(m)
	return err
}
