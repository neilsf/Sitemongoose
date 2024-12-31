package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetService(t *testing.T) {
	service := GetService()
	assert.NotNil(t, service)
	service2 := GetService()
	assert.Equal(t, service, service2)
}

func TestAddMonitor(t *testing.T) {
	service := GetService()
	monitor := Monitor{}
	service.AddMonitor(monitor)
	assert.Equal(t, 1, len(service.GetMonitors()))
	service.AddMonitor(monitor)
	assert.Equal(t, 2, len(service.GetMonitors()))
}

func TestGetMonitors(t *testing.T) {
	service := GetService()
	monitor := Monitor{}
	service.AddMonitor(monitor)
	assert.Equal(t, 3, len(service.GetMonitors()))
	service.AddMonitor(monitor)
	assert.Equal(t, 4, len(service.GetMonitors()))
}
