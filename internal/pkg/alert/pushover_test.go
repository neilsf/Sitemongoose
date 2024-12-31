package alert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setEnvVarsPushover() {
	os.Setenv("PUSHOVER_APP_TOKEN", "123456789")
	os.Setenv("PUSHOVER_USER_KEY", "abcdefgh")
}

func unSetEnvVarsPushover() {
	os.Unsetenv("PUSHOVER_APP_TOKEN")
	os.Unsetenv("PUSHOVER_USER_KEY")
}

type MockPushoverSender struct{}

var mockPushoverResult string

func (m *MockPushoverSender) SendPush(message string) error {
	mockPushoverResult = "Sent a pushover message: " + message
	return nil
}

func TestValidate_Pushover(t *testing.T) {
	setEnvVarsPushover()
	alert := PushoverAlerter{Alert{"pushover", "", "", "", "", nil, nil}, nil}
	valid, err := alert.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
	alert.AlertMessage = "alert"
	alert.ResolutionMessage = "resolution"
	valid, err = alert.Validate()
	assert.True(t, valid)
	assert.Nil(t, err)
	unSetEnvVarsPushover()
	valid, err = alert.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestSendAlert_Pushover(t *testing.T) {
	setEnvVarsPushover()
	alert := PushoverAlerter{Alert{"pushover", "alert", "resolution", "", "", nil, nil}, &MockPushoverSender{}}
	alert.SendAlert()
	expected := "Sent a pushover message: alert"
	assert.Equal(t, expected, mockPushoverResult)
}

func TestSendResolution_Pushover(t *testing.T) {
	setEnvVarsPushover()
	alert := PushoverAlerter{Alert{"pushover", "alert", "resolution", "", "", nil, nil}, &MockPushoverSender{}}
	alert.SendResolution()
	expected := "Sent a pushover message: resolution"
	assert.Equal(t, expected, mockPushoverResult)
}
