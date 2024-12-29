package alert

const (
	ALERT_TYPE_EMAIL      = "email"
	ALERT_TYPE_SLACK      = "slack"
	ALETR_TYPE_PUSHOVER   = "pushover"
	ALERT_TYPE_CUSTOM_CMD = "command"
)

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
}

func GetAlerter(alert Alert) IAlerter {
	switch alert.Type {
	case ALERT_TYPE_EMAIL:
		return &EmailAlerter{alert}
	case ALERT_TYPE_CUSTOM_CMD:
		return &CmdAlerter{alert}
	}
	return nil
}