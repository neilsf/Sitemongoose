package alert

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setEnvVarsEmail() {
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USER", "user")
	os.Setenv("SMTP_PASS", "pass")
}

func unSetEnvVarsEmail() {
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USER")
	os.Unsetenv("SMTP_PASS")
}

func TestValidate_Email(t *testing.T) {
	alert := EmailAlerter{Alert{"email", "alert", "resolution", "from", "to", nil, nil}, nil}
	valid, err := alert.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
	alert.From = "alerts@example.com"
	alert.To = "info@example.com"
	valid, err = alert.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
	setEnvVarsEmail()
	valid, err = alert.Validate()
	assert.True(t, valid)
	assert.Nil(t, err)
	unSetEnvVarsEmail()
}

type MockEmailSender struct{}

var mockEmailResult string

func (m *MockEmailSender) SendEmail(to, subject, body string) error {
	mockEmailResult = fmt.Sprintf("Sent an email to %s with subject %s and body %s", to, subject, body)
	return nil
}

func TestSendAlert(t *testing.T) {
	setEnvVarsEmail()
	alert := EmailAlerter{Alert{"email", "alert", "resolution", "from", "to", nil, nil}, &MockEmailSender{}}
	alert.From = "alerts@example.com"
	alert.To = "info@example.com"
	alert.SendAlert()
	expected := fmt.Sprintf("Sent an email to %s with subject %s and body %s", alert.To, EMAIL_ALERT_SUBJECT, "alert")
	assert.Equal(t, expected, mockEmailResult)
}

func TestSendResolution(t *testing.T) {
	setEnvVarsEmail()
	alert := EmailAlerter{Alert{"email", "alert", "resolution", "from", "to", nil, nil}, &MockEmailSender{}}
	alert.From = "alerts@example.com"
	alert.To = "info@example.com"
	alert.SendResolution()
	expected := fmt.Sprintf("Sent an email to %s with subject %s and body %s", alert.To, EMAIL_RESOLUTION_SUBJECT, "resolution")
	assert.Equal(t, expected, mockEmailResult)
}
