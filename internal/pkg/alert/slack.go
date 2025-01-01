package alert

import (
	"errors"
	"os"
)

type SlackSender interface {
	SendMsg(body string) error
}

type SlackAlerter struct {
	Alert
	Sender SlackSender
}

func (s *SlackAlerter) Validate() (bool, error) {
	if s.AlertMessage == "" {
		return false, errors.New("alert message is required")
	}
	if s.ResolutionMessage == "" {
		return false, errors.New("resolution message is required")
	}
	if os.Getenv("SLACK_WEBHOOK_URL") == "" {
		return false, errors.New("the environment variable SLACK_WEBHOOK_URL must be set in order to send slack alerts")
	}
	return true, nil
}

func (s *SlackAlerter) SendAlert() {
	s.Sender.SendMsg(s.AlertMessage)
}

func (s *SlackAlerter) SendResolution() {
	s.Sender.SendMsg(s.ResolutionMessage)
}
