package service

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendMail(to []string, message []byte) error {
	auth := smtp.PlainAuth("", os.Getenv("MAIL_DOMAIN"), os.Getenv("MAIL_PASSWORD"), os.Getenv("MAIL_HOST"))

	err := smtp.SendMail(os.Getenv("MAIL_HOST")+":"+os.Getenv("MAIL_PORT"), auth, os.Getenv("MAIL_DOMAIN"), to, message)
	if err != nil {
		return err
	}
	fmt.Print("Mail Sent!")
	return nil
}
