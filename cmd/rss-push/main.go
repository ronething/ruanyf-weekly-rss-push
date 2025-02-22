package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ashing/ruanyf-weekly-rss-push/pkg/cache"
	"github.com/ashing/ruanyf-weekly-rss-push/pkg/feed"
	"github.com/ashing/ruanyf-weekly-rss-push/pkg/notifier"
)

const (
	RSSFeedURL = "https://github.com/ruanyf/weekly/releases.atom"
	FeedLimit  = 2
)

func main() {
	slackURL := os.Getenv("SLACK_WEBHOOK_URL")
	wecomToken := os.Getenv("WECOM_TOKEN")

	// Initialize cache
	cacheDir := os.Getenv("CACHE_DIR")
	if cacheDir == "" {
		cacheDir = "cache"
	}
	articleCache, err := cache.NewArticleCache(filepath.Join(cacheDir, "articles.json"))
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

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

	if wecomToken != "" {
		notifiers = append(notifiers, notifier.NewWeComNotifier(wecomToken))
	}

	if len(notifiers) == 0 {
		log.Println("Warning: No notifiers configured")
	}

	// Process each entry
	for _, entry := range entries {
		mdURL := feed.ConvertToMarkdownURL(entry.Link.Href)

		// Skip if already pushed
		if articleCache.IsArticlePushed(mdURL) {
			log.Printf("Skipping already pushed article: %s\n", entry.Title)
			continue
		}

		msg := notifier.Message{
			Title: entry.Title,
			URL:   mdURL,
		}
		log.Printf("Pushing new article: %+v\n", msg)

		// Send notifications
		for _, n := range notifiers {
			if err := n.Notify(msg); err != nil {
				log.Printf("Failed to send notification: %v\n", err)
				continue
			}
		}

		// Mark as pushed after successful notification
		if err := articleCache.MarkArticlePushed(mdURL); err != nil {
			log.Printf("Failed to mark article as pushed: %v\n", err)
		}
	}
}
