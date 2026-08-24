package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func EmailSender(receipient, subject, body string) {
	var (
		from     = os.Getenv("SOURCE_EMAIL")
		password = os.Getenv("SOURCE_EMAIL_PASSWORD")

		smtpHost = os.Getenv("SMTP_HOST")
		smtpPort = os.Getenv("SMTP_PORT")
	)

	to := []string{receipient}
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpPort+":"+smtpHost, auth, from, to, message)
	if err != nil {
		fmt.Println("Error handling email: " + err.Error())
		return
	}
	log.Printf("Email sent successfully")
}
