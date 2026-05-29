package notifier

import (
	"log"
	"strings"

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
	content := formatMessageContent(msg)
	if err := msgpush.NewWeCom(w.Token).SendText(content); err != nil {
		log.Printf("WeCom: Failed to send notification: %v\n", err)
		return err
	}
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
	content := formatMessageContent(msg)
	if err := msgpush.NewSlack(s.WebhookURL).Send(content); err != nil {
		log.Printf("Slack: Failed to send notification: %v\n", err)
		return err
	}
	return nil
}

type FeiShuNotifier struct {
	WebhookURL string
}

func NewFeiShuNotifier(webhookURL string) *FeiShuNotifier {
	return &FeiShuNotifier{
		WebhookURL: webhookURL,
	}
}

func (f *FeiShuNotifier) Notify(msg Message) error {
	content := formatMessageContent(msg)
	if err := msgpush.NewFeiShu(f.WebhookURL).SendText(content); err != nil {
		log.Printf("FeiShu: Failed to send notification: %v\n", err)
		return err
	}
	return nil
}

// formatMessageContent formats the message content, including the title, URL, and recommendations (if any)
func formatMessageContent(msg Message) string {
	var sb strings.Builder

	sb.WriteString(msg.Title)
	sb.WriteString("\n")
	sb.WriteString(msg.URL)

	if msg.Recommendations != "" {
		sb.WriteString("\n\n")
		sb.WriteString(msg.Recommendations)
	}

	return sb.String()
}
