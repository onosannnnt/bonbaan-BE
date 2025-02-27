package utils

import (
	"crypto/tls"

	"github.com/onosannnnt/bonbaan-BE/src/config"
	"gopkg.in/gomail.v2"
)

func SendingMail(m *gomail.Message) {
	d := gomail.NewDialer(config.MailHost, config.MailPort, config.MailUser, config.MailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}
