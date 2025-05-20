package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashing/ruanyf-weekly-rss-push/pkg/ai"
	"github.com/ashing/ruanyf-weekly-rss-push/pkg/cache"
	"github.com/ashing/ruanyf-weekly-rss-push/pkg/feed"
	"github.com/ashing/ruanyf-weekly-rss-push/pkg/notifier"
)

const (
	RSSFeedURL = "https://github.com/ruanyf/weekly/releases.atom"
	FeedLimit  = 2
)

func main() {
	slackURLs := getMultipleValues("SLACK_WEBHOOK_URL")
	wecomTokens := getMultipleValues("WECOM_TOKEN")

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

	// Initialize AI content summarizer
	summarizer, err := ai.NewContentSummarizer()
	if err != nil {
		log.Printf("Failed to initialize AI summarizer: %v, will continue without recommendations", err)
	}

	// Initialize notifiers
	notifiers := []notifier.Notifier{}

	// Initialize Slack notifiers
	if len(slackURLs) > 0 {
		for _, url := range slackURLs {
			notifiers = append(notifiers, notifier.NewSlackNotifier(url))
		}
		log.Printf("Configured %d Slack notifier(s)", len(slackURLs))
	} else {
		log.Println("Warning: No Slack notifier configured")
	}

	// Initialize WeCom notifiers
	if len(wecomTokens) > 0 {
		for _, token := range wecomTokens {
			notifiers = append(notifiers, notifier.NewWeComNotifier(token))
		}
		log.Printf("Configured %d WeCom notifier(s)", len(wecomTokens))
	} else {
		log.Println("Warning: No WeCom notifier configured")
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

		// prepare the message content
		msg := notifier.Message{
			Title: entry.Title,
			URL:   mdURL,
		}

		// if there is content and the AI summarizer is available, try to extract recommendations
		if summarizer != nil && entry.Content.Src != "" {
			log.Println("Extracting recommendations using AI...")
			recommendations, err := summarizer.ExtractRecommendations(entry)
			if err != nil {
				log.Printf("Failed to extract recommendations: %v\n", err)
			} else {
				recommendationsText := ai.FormatRecommendationsText(recommendations)
				msg.Recommendations = recommendationsText
				log.Printf("Added %d recommendations to the message\n", len(recommendations))
			}
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

	log.Println("All articles pushed")
}

// getMultipleValues gets multiple values from a single environment variable,
// values are expected to be comma-separated
func getMultipleValues(envKey string) []string {
	value := os.Getenv(envKey)
	if value == "" {
		return []string{}
	}

	values := strings.Split(value, ",")
	result := make([]string, 0, len(values))

	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
