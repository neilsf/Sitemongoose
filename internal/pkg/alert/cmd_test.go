package alert

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCmd_Noargs(t *testing.T) {
	alert := CmdAlerter{Alert{"command", "", "", "", "", nil, nil}}
	valid, err := alert.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestValidateCmd_Valid(t *testing.T) {
	alert := CmdAlerter{Alert{"command", "", "", "", "", []string{";"}, []string{";"}}}
	valid, err := alert.Validate()
	assert.True(t, valid)
	assert.Nil(t, err)
}

func TestSendAlertCmd(t *testing.T) {
	alert := CmdAlerter{Alert{"command", "", "", "", "", []string{"echo", "alert"}, []string{"echo", "resolution"}}}
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	alert.SendAlert()
	assert.Contains(t, buf.String(), "alert")
}

func TestSendResolutionCmd(t *testing.T) {
	alert := CmdAlerter{Alert{"command", "", "", "", "", []string{"echo", "alert"}, []string{"echo", "resolution"}}}
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	alert.SendResolution()
	assert.Contains(t, buf.String(), "resolution")
}
