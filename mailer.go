package main

import (
	"crypto/tls"
	"fmt"
	"net"
	netmail "net/mail"
	"net/smtp"
)

type Mailer struct {
	mailsrv, mailfrom, mailpass string
}

func NewMailer(mailsrv, mailfrom, mailpass string) *Mailer {
	return &Mailer{
		mailsrv:  mailsrv,
		mailfrom: mailfrom,
		mailpass: mailpass,
	}
}

func (m *Mailer) Send(email, subj, body string) error {
	from := netmail.Address{"", m.mailfrom}
	to := netmail.Address{"", email}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	host, _, _ := net.SplitHostPort(m.mailsrv)

	auth := smtp.PlainAuth("", m.mailfrom, m.mailpass, host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", m.mailsrv, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
