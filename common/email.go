package common

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func NewSmtp(host string, port int, protocol string, account string, password string, email string) *Smtp {
	d := gomail.NewDialer(host, port, account, password)
	if protocol == "ssl" || protocol == "tls" {
		d.SSL = true
		d.TLSConfig = &tls.Config{ServerName: host}
	}
	return &Smtp{dialer: d, email: email}
}

type Smtp struct {
	email  string
	dialer *gomail.Dialer
}

func (s *Smtp) SendTo(sends []string, title, msg string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", s.email)
	m.SetHeader("To", sends...)
	m.SetHeader("Subject", title)
	m.SetBody("text/html;charset=utf-8", msg)

	return s.dialer.DialAndSend(m)
}
