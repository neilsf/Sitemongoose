package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_Invalid(t *testing.T) {
	alert := Alert{"invalid", "", "", "", "", nil, nil}
	valid, err := alert.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}
