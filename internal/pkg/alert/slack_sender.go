package alert

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type LiveSlackSender struct{}

func (s *LiveSlackSender) SendMsg(message string) error {

	payload := map[string]string{
		"text": message,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", os.Getenv("SLACK_WEBHOOK_URL"), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
