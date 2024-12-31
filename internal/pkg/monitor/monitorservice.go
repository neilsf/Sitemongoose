package monitor

import (
	"log"
	"reflect"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type MonitorService struct {
	monitors []Monitor
}

var instance *MonitorService

func GetService() *MonitorService {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &MonitorService{}
		}
	}

	return instance
}

func (ms *MonitorService) AddMonitor(m Monitor) {
	m.SetDefaults()
	ms.monitors = append(ms.monitors, m)
}

func (m *MonitorService) GetMonitors() []Monitor {
	return m.monitors
}

func (ms *MonitorService) Start() {
	monitorChannels := ms.createMonitorChannels()
	ms.runMonitorLoop(monitorChannels)
}

func (ms *MonitorService) createMonitorChannels() map[string]<-chan time.Time {
	monitorChannels := make(map[string]<-chan time.Time)
	for _, monitor := range ms.monitors {
		monitorChannels[monitor.Name] = time.Tick(time.Duration(monitor.IntervalSec) * time.Second)
	}
	return monitorChannels
}

func (ms *MonitorService) runMonitorLoop(monitorChannels map[string]<-chan time.Time) {
	log.Printf("Starting monitoring service with %d monitors\n", len(ms.monitors))

	cases := make([]reflect.SelectCase, len(monitorChannels))
	i := 0
	for _, ch := range monitorChannels {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
		i++
	}

	for {
		chosen, _, _ := reflect.Select(cases)
		monitor := ms.monitors[chosen]
		monitor.doCheck()
	}
}
