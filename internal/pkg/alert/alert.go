package alert

import (
	"errors"
)

const (
	ALERT_TYPE_EMAIL      = "email"
	ALERT_TYPE_SLACK      = "slack"
	ALETR_TYPE_PUSHOVER   = "pushover"
	ALERT_TYPE_CUSTOM_CMD = "command"
)

var validAlertTypes = map[string]bool{
	ALERT_TYPE_EMAIL:      true,
	ALERT_TYPE_SLACK:      true,
	ALETR_TYPE_PUSHOVER:   true,
	ALERT_TYPE_CUSTOM_CMD: true,
}

type Alert struct {
	Type              string
	AlertMessage      string   `yaml:"alert_message"`
	ResolutionMessage string   `yaml:"resolution_message"`
	From              string   `yaml:"from"`
	To                string   `yaml:"to"`
	AlertCommand      []string `yaml:"alert_command"`
	ResolutionCommand []string `yaml:"resolution_command"`
}

type IAlerter interface {
	SendAlert()
	SendResolution()
	Validate() (bool, error)
}

func GetAlerter(alert Alert) IAlerter {
	switch alert.Type {
	case ALERT_TYPE_EMAIL:
		return &EmailAlerter{alert, &GomailSender{}}
	case ALERT_TYPE_CUSTOM_CMD:
		return &CmdAlerter{alert}
	}
	return nil
}

func (a *Alert) Validate() (bool, error) {
	if _, ok := validAlertTypes[a.Type]; !ok {
		return false, errors.New("invalid alert type, must be one of: email, slack, pushover, command")
	}
	return GetAlerter(*a).Validate()
}
