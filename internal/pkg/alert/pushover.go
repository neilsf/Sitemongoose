package alert

import (
	"errors"
	"os"
)

type PushoverSender interface {
	SendPush(message string) error
}

type PushoverAlerter struct {
	Alert
	Sender PushoverSender
}

func (p *PushoverAlerter) Validate() (bool, error) {
	if p.AlertMessage == "" {
		return false, errors.New("alert message is required")
	}
	if p.ResolutionMessage == "" {
		return false, errors.New("resolution message is required")
	}
	if os.Getenv("PUSHOVER_APP_TOKEN") == "" {
		return false, errors.New("the environment variable PUSHOVER_APP_TOKEN must be set in order to send Pushover alerts")
	}
	if os.Getenv("PUSHOVER_USER_KEY") == "" {
		return false, errors.New("the environment variable PUSHOVER_USER_KEY must be set in order to send Pushover alerts")
	}
	return true, nil
}

func (p *PushoverAlerter) SendAlert() {
	p.Sender.SendPush(p.AlertMessage)
}

func (p *PushoverAlerter) SendResolution() {
	p.Sender.SendPush(p.ResolutionMessage)
}
