package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_Valid(t *testing.T) {
	monitor := Monitor{"Example", "http://example.com", 10, 1000, nil}
	valid, err := monitor.Validate()
	assert.True(t, valid)
	assert.Nil(t, err)
}

func TestValidate_MissingName(t *testing.T) {
	monitor := Monitor{"", "http://example.com", 10, 1000, nil}
	valid, err := monitor.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestValidate_MissingURL(t *testing.T) {
	monitor := Monitor{"test", "", 10, 1000, nil}
	valid, err := monitor.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestValidate_InvalidIntervalSec(t *testing.T) {
	monitor := Monitor{"test", "http://example.com", 0, 1000, nil}
	valid, err := monitor.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestSetDefaults(t *testing.T) {
	monitor := Monitor{"test", "http://example.com", 10, 0, nil}
	monitor.SetDefaults()
	assert.Less(t, 0, monitor.TimeoutMs)
}

func TestSetDefaults_NoChange(t *testing.T) {
	monitor := Monitor{"test", "http://example.com", 10, 1000, nil}
	monitor.SetDefaults()
	assert.Equal(t, 1000, monitor.TimeoutMs)
}
