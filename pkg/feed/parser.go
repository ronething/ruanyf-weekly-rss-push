package feed

import (
	"encoding/xml"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title   string    `xml:"title"`
	Link    Link      `xml:"link"`
	Updated time.Time `xml:"updated"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

type Parser struct {
	URL   string
	Limit int
}

func NewParser(url string, limit int) *Parser {
	return &Parser{
		URL:   url,
		Limit: limit,
	}
}

func (p *Parser) Parse() ([]Entry, error) {
	resp, err := http.Get(p.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var feed Feed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	if len(feed.Entries) > p.Limit {
		feed.Entries = feed.Entries[:p.Limit]
	}

	return feed.Entries, nil
}

func ConvertToMarkdownURL(releaseURL string) string {
	// Example: https://github.com/ruanyf/weekly/releases/tag/issue-xxx
	// to: https://github.com/ruanyf/weekly/blob/master/docs/issue-xxx.md

	// Split by "issue-" and take the last part
	parts := strings.Split(releaseURL, "issue-")
	if len(parts) != 2 {
		return releaseURL // Return original URL if format doesn't match
	}

	// Extract only numbers from the issue part
	re := regexp.MustCompile(`^\d+`)
	issueNum := re.FindString(parts[1])
	if issueNum == "" {
		return releaseURL // Return original URL if no number found
	}

	// Construct the new URL
	return "https://github.com/ruanyf/weekly/blob/master/docs/issue-" + issueNum + ".md"
}
