package alert

import (
	"errors"
	"maps"
	"slices"
	"strings"
)

const (
	ALERT_TYPE_EMAIL      = "email"
	ALERT_TYPE_SLACK      = "slack"
	ALERT_TYPE_PUSHOVER   = "pushover"
	ALERT_TYPE_CUSTOM_CMD = "command"
)

var validAlertTypes = map[string]bool{
	ALERT_TYPE_EMAIL:      true,
	ALERT_TYPE_SLACK:      true,
	ALERT_TYPE_PUSHOVER:   true,
	ALERT_TYPE_CUSTOM_CMD: true,
}

type Alert struct {
	Type              string   `yaml:"channel"`
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
	case ALERT_TYPE_PUSHOVER:
		return &PushoverAlerter{alert, &LivePushoverSender{}}
	case ALERT_TYPE_SLACK:
		return &SlackAlerter{alert, &LiveSlackSender{}}
	}
	return nil
}

func (a *Alert) Validate() (bool, error) {
	if _, ok := validAlertTypes[a.Type]; !ok {
		return false, errors.New("invalid alert channel, must be one of: " + strings.Join(slices.Collect(maps.Keys(validAlertTypes)), ", "))
	}
	return GetAlerter(*a).Validate()
}
