package utils

import (
	"crypto/tls"

	"github.com/onosannnnt/bonbaan-BE/src/Config"
	"gopkg.in/gomail.v2"
)

func SendingMail(m *gomail.Message) {
	d := gomail.NewDialer(Config.MailHost, Config.MailPort, Config.MailUser, Config.MailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}
