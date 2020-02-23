package mailutils

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type SMTPSettings struct {
	Port     int
	Host     string
	User     string
	Password string
}

type SendInput struct {
	// Email address of the sender
	From string
	// Email address of the receiver
	To string
	// Subject of the message
	Subject string
	// Content of the message (mail body)
	Content string
	// Whether or not to render as html
	IsHTML bool
}

// MailSender is an interface for sending emails
type MailSender interface {
	Send(input *SendInput) error
}

// PlainMailSender is a plain smtp-based email sender
type PlainMailSender struct {
	smtp SMTPSettings
}

// Send sends an email to the target address
func (p PlainMailSender) Send(input *SendInput) error {
	m := gomail.NewMessage()
	m.SetHeader("From", input.From)
	m.SetHeader("To", input.To)
	m.SetHeader("Subject", input.Subject)

	if input.IsHTML {
		m.SetBody("text/html", input.Content)
	} else {
		m.SetBody("text/plain", input.Content)
	}

	d := gomail.NewDialer(p.smtp.Host, p.smtp.Port, p.smtp.User, p.smtp.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// NewPlainMailSender returns a new plain mail sender instance
func NewPlainMailSender(smtp SMTPSettings) PlainMailSender {
	return PlainMailSender{smtp: smtp}
}

// MockMailSender is a mock mail sender used while testing
type MockMailSender struct{}

// Send just logs the input message
func (p MockMailSender) Send(input *SendInput) error {
	logrus.Infof("MockMailSender [send]: %v", input)
	return nil
}

// NewMockMailSender returns a new mock mail sender instance
func NewMockMailSender() MockMailSender {
	return MockMailSender{}
}
