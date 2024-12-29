package alert

import (
	"log"
	"os"
	"strconv"
	"time"

	gomail "gopkg.in/mail.v2"
)

type EmailAlerter struct {
	Alert
}

func getDialer() *gomail.Dialer {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	dialer := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))
	dialer.Timeout = 30 * time.Second
	return dialer
}

func (e *EmailAlerter) send(subject string, msg string) {
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", e.From)
	message.SetHeader("To", e.To)
	message.SetHeader("Subject", subject)

	// Set email body
	message.SetBody("text/plain", msg)

	// Set up the SMTP dialer
	dialer := getDialer()

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		log.Println("Error sending alert", err)
	} else {
		log.Printf("Sent alert message to %s", e.To)
	}
}

func (e *EmailAlerter) SendAlert() {
	e.send(e.AlertMessage, "Alert from Sitemongoose")
}

func (e *EmailAlerter) SendResolution() {
	e.send(e.ResolutionMessage, "Resolved")
}
