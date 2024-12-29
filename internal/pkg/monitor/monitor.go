package monitor

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/neilsf/sitemongoose/internal/pkg/event"
)

type Monitor struct {
	Name        string
	URL         string
	IntervalSec int `yaml:"interval_sec"`
	TimeoutMs   int `yaml:"timeout_ms"`
	Events      []event.Event
}

func (m *Monitor) Validate() (bool, error) {
	if m.Name == "" {
		return false, errors.New("name is required")
	}
	if m.URL == "" {
		return false, errors.New("url is required")
	}
	if m.IntervalSec <= 0 {
		return false, errors.New("interval_sec must be greater than 0")
	}
	return true, nil
}

func (m *Monitor) SetDefaults() {
	if m.TimeoutMs == 0 {
		m.TimeoutMs = 30000
	}
}

func (m *Monitor) doCheck() {
	start := time.Now()
	client := http.Client{
		Timeout: time.Duration(m.TimeoutMs) * time.Millisecond,
	}
	resp, err := client.Get(m.URL)
	if err != nil {
		if os.IsTimeout(err) {
			for i := range m.Events {
				e := &m.Events[i]
				if e.TriggerType == event.TRIGGER_TYPE_RESPONSE_TIME {
					e.MonitorName = m.Name
					go e.CheckTrigger(0, m.TimeoutMs, nil)
				}
			}
		}
		log.Printf("Error checking %s: %v\n", m.URL, err)
		return
	}
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
	}
	defer resp.Body.Close()
	duration := time.Since(start)

	log.Printf("%s, Status code: %d, Response time: %v", m.URL, resp.StatusCode, duration)

	for i := range m.Events {
		e := &m.Events[i]
		e.MonitorName = m.Name
		go e.CheckTrigger(resp.StatusCode, int(duration.Milliseconds()), body)
	}
}
