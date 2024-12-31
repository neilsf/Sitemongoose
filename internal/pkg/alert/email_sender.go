package alert

import (
	"log"
	"os"
	"strconv"
	"time"

	gomail "gopkg.in/mail.v2"
)

type GomailSender struct{}

func (s *GomailSender) SendEmail(to, subject, body string) error {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	dialer := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))
	dialer.Timeout = 30 * time.Second

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("SMTP_USER"))
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	var err error
	// Send the email
	if err = dialer.DialAndSend(message); err != nil {
		log.Println("Error sending alert", err)
	} else {
		log.Printf("Sent alert message to %s", to)
	}

	return err
}
