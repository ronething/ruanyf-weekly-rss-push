package notifier

import (
	"github.com/cloud-org/msgpush"
)

type WeComNotifier struct {
	Token string
}

func NewWeComNotifier(token string) *WeComNotifier {
	return &WeComNotifier{
		Token: token,
	}
}

func (w *WeComNotifier) Notify(msg Message) error {
	msgpush.NewWeCom(w.Token).SendText(msg.Title + "\n" + msg.URL)
	return nil
}

type SlackNotifier struct {
	WebhookURL string
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
	}
}

func (s *SlackNotifier) Notify(msg Message) error {
	msgpush.NewSlack(s.WebhookURL).Send(msg.Title + "\n" + msg.URL)
	return nil
}
