package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/lucasvictor3/bookingsbackend/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listerForMail() {

	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()

}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "172.21.176.1"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 15 * time.Second
	server.SendTimeout = 15 * time.Second

	log.Println("Connecting to email server...")

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, "Hello, <strong>world</strong>!")
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}

		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}
}
