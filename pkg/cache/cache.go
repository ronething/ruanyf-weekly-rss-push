package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ArticleCache struct {
	Articles map[string]time.Time `json:"articles"` // URL -> PushTime
	path     string
}

func NewArticleCache(cachePath string) (*ArticleCache, error) {
	cache := &ArticleCache{
		Articles: make(map[string]time.Time),
		path:     cachePath,
	}

	// Create cache directory if not exists
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return nil, err
	}

	// Try to load existing cache
	if _, err := os.Stat(cachePath); err == nil {
		data, err := os.ReadFile(cachePath)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &cache.Articles); err != nil {
			return nil, err
		}
	}

	return cache, nil
}

func (c *ArticleCache) IsArticlePushed(url string) bool {
	_, exists := c.Articles[url]
	return exists
}

func (c *ArticleCache) MarkArticlePushed(url string) error {
	c.Articles[url] = time.Now()
	return c.Save()
}

func (c *ArticleCache) Save() error {
	data, err := json.Marshal(c.Articles)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}
