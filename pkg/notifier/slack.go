package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SlackNotifier struct {
	WebhookURL string
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
	}
}

func (s *SlackNotifier) Notify(msg Message) error {
	payload := map[string]string{
		"text": msg.Title + " " + msg.URL,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
