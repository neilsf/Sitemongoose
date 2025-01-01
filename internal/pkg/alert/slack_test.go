package alert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setEnvVarsSlack() {
	os.Setenv("SLACK_WEBHOOK_URL", "123456789")
}

func unSetEnvVarsSlack() {
	os.Unsetenv("SLACK_WEBHOOK_URL")
}

type MockSlackSender struct{}

var mockSlackResult string

func (m *MockSlackSender) SendMsg(message string) error {
	mockSlackResult = "Sent a slack message: " + message
	return nil
}

func TestValidate_Slack(t *testing.T) {
	setEnvVarsSlack()
	alert := SlackAlerter{Alert{"slack", "", "", "", "", nil, nil}, nil}
	valid, err := alert.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
	alert.AlertMessage = "alert"
	alert.ResolutionMessage = "resolution"
	valid, err = alert.Validate()
	assert.True(t, valid)
	assert.Nil(t, err)
	unSetEnvVarsSlack()
	valid, err = alert.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestSendAlert_Slack(t *testing.T) {
	setEnvVarsSlack()
	alert := SlackAlerter{Alert{"slack", "alert", "resolution", "", "", nil, nil}, &MockSlackSender{}}
	alert.SendAlert()
	expected := "Sent a slack message: alert"
	assert.Equal(t, expected, mockSlackResult)
}

func TestSendResolution_Slack(t *testing.T) {
	setEnvVarsSlack()
	alert := SlackAlerter{Alert{"slack", "alert", "resolution", "", "", nil, nil}, &MockSlackSender{}}
	alert.SendResolution()
	expected := "Sent a slack message: resolution"
	assert.Equal(t, expected, mockSlackResult)
}
