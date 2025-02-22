package main

import (
	"log"
	"os"

	"github.com/ashing/ruanyf-weekly-rss-push/pkg/feed"
	"github.com/ashing/ruanyf-weekly-rss-push/pkg/notifier"
)

const (
	RSSFeedURL = "https://github.com/ruanyf/weekly/releases.atom"
	FeedLimit  = 2
)

func main() {
	slackURL := os.Getenv("SLACK_WEBHOOK_URL")
	wecomURL := os.Getenv("WECOM_WEBHOOK_URL")

	// Initialize RSS parser
	parser := feed.NewParser(RSSFeedURL, FeedLimit)
	entries, err := parser.Parse()
	if err != nil {
		log.Fatalf("Failed to parse RSS feed: %v", err)
	}

	// Initialize notifiers
	notifiers := []notifier.Notifier{}

	if slackURL != "" {
		notifiers = append(notifiers, notifier.NewSlackNotifier(slackURL))
	}

	if wecomURL != "" {
		notifiers = append(notifiers, notifier.NewWeComNotifier(wecomURL))
	}

	if len(notifiers) == 0 {
		log.Println("Warning: No notifiers configured")
	}

	// Process each entry
	for _, entry := range entries {
		mdURL := feed.ConvertToMarkdownURL(entry.Link.Href)
		msg := notifier.Message{
			Title: entry.Title,
			URL:   mdURL,
		}
		log.Printf("msg: %+v\n", msg)

		// Send notifications
		for _, n := range notifiers {
			if err := n.Notify(msg); err != nil {
				log.Printf("Failed to send notification: %v", err)
			}
		}
	}
}
