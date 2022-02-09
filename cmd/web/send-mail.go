package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/kjfs/dieffe_sensor/internal/models"
	mail "github.com/xhit/go-simple-mail"
)

func listenForMail() {

	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 465
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.Username = "XXX"
	server.Password = "XXX"
	server.Encryption = mail.EncryptionSSL

	client, err := server.Connect()
	if err != nil {
		log.Println("Log: Cant connect to SMTP server: ", err)
		errorLog.Println("Errorlog: Cant connect to SMTP server: ", err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
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
		log.Println("Email send!")
	}

}
