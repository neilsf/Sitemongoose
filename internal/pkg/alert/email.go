package alert

import (
	"errors"
	"os"
	"regexp"
)

const (
	EMAIL_ALERT_SUBJECT      = "Alert from Sitemongoose"
	EMAIL_RESOLUTION_SUBJECT = "Resolved"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type EmailAlerter struct {
	Alert
	Sender EmailSender
}

func (e *EmailAlerter) Validate() (bool, error) {
	if e.From == "" {
		return false, errors.New("'From' field is required")
	}
	if !isValidEmail(e.From) {
		return false, errors.New("'From' field is not a valid email address")
	}
	if e.To == "" {
		return false, errors.New("'To' field is required")
	}
	if !isValidEmail(e.To) {
		return false, errors.New("'To' field is not a valid email address")
	}
	if e.AlertMessage == "" {
		return false, errors.New("alert message is required")
	}
	if e.ResolutionMessage == "" {
		return false, errors.New("resolution message is required")
	}
	if os.Getenv("SMTP_HOST") == "" {
		return false, errors.New("the environment variable SMTP_HOST must be set in order to send email alerts")
	}
	if os.Getenv("SMTP_PORT") == "" {
		return false, errors.New("the environment variable SMTP_PORT must be set in order to send email alerts")
	}
	if os.Getenv("SMTP_USER") == "" {
		return false, errors.New("the environment variable SMTP_USER must be set in order to send email alerts")
	}
	return true, nil
}

func (e *EmailAlerter) SendAlert() {
	e.Sender.SendEmail(e.To, EMAIL_ALERT_SUBJECT, e.AlertMessage)
}

func (e *EmailAlerter) SendResolution() {
	e.Sender.SendEmail(e.To, EMAIL_RESOLUTION_SUBJECT, e.ResolutionMessage)
}

// IsValidEmail checks if the given string is a valid email address.
func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
