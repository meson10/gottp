package gottp

import (
	"log"
	"net/smtp"
	"strings"
)

type MailConn struct {
	Hostname   string
	Username   string
	Password   string
	SenderName string
	Port       string
	Host       string
}

type Message struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func (conn *MailConn) getAuth() smtp.Auth {
	return smtp.PlainAuth("", conn.Username, conn.Password, conn.Hostname)
}

func (conn *MailConn) MessageBytes(message Message) []byte {
	subject := "Subject: "
	subject += message.Subject

	subject = strings.TrimSpace(subject)
	from := strings.TrimSpace("From: " + conn.SenderName + " <" + message.From + ">")
	to := strings.TrimSpace("To: " + strings.Join(message.To, ", "))
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	return []byte(subject + "\n" + from + "\n" + to + "\n" + mime + "\n\n" + strings.TrimSpace(message.Body))

}

func (conn *MailConn) SendEmail(message Message) {
	if settings.Gottp.EmailDummy == true {
		log.Println("Not sending email as dummy set to true")
		return
	}

	err := smtp.SendMail(conn.Host,
		conn.getAuth(),
		message.From,
		message.To,
		conn.MessageBytes(message))

	if err != nil {
		log.Panic(err)
	}
}

func MakeConn() *MailConn {
	mailconn := &MailConn{
		settings.Gottp.EmailHost,
		settings.Gottp.EmailUsername,
		settings.Gottp.EmailPassword,
		settings.Gottp.EmailSender,
		settings.Gottp.EmailPort,
		settings.Gottp.EmailHost + ":" + settings.Gottp.EmailPort,
	}
	return mailconn
}
