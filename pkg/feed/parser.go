package feed

import (
	"encoding/xml"
	"net/http"
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
	issueNum := releaseURL[len(releaseURL)-3:]
	return "https://github.com/ruanyf/weekly/blob/master/docs/issue-" + issueNum + ".md"
}
