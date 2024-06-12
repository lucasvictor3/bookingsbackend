package main

import (
	"log"
	"time"

	"github.com/lucasvictor3/bookingsbackend/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listerForMail() {

	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(&msg)
		}
	}()

}

func sendMsg(m *models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "172.18.0.1"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 20 * time.Second
	server.SendTimeout = 20 * time.Second

	log.Println("Connecting to email server...")

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom("me@here.com").AddTo("john@do.ca").SetSubject("Some subject")
	email.SetBody(mail.TextHTML, "Hello, <strong>world</strong>!")

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}
}
