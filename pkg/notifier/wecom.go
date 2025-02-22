package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type WeComNotifier struct {
	WebhookURL string
}

func NewWeComNotifier(webhookURL string) *WeComNotifier {
	return &WeComNotifier{
		WebhookURL: webhookURL,
	}
}

func (w *WeComNotifier) Notify(msg Message) error {
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg.Title + " " + msg.URL,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(w.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
