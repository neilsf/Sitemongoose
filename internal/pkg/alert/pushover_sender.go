package alert

import (
	"net/http"
	"net/url"
	"os"
)

const PUSHOVER_API_ENDPOINT = "https://api.pushover.net/1/messages.json"

type LivePushoverSender struct{}

func (s *LivePushoverSender) SendPush(message string) error {

	data := url.Values{}
	data.Set("token", os.Getenv("PUSHOVER_APP_TOKEN"))
	data.Set("user", os.Getenv("PUSHOVER_USER_KEY"))
	data.Set("message", message)

	resp, err := http.PostForm(PUSHOVER_API_ENDPOINT, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
