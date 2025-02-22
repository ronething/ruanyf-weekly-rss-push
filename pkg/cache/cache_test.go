package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewArticleCache(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("create new cache without existing file", func(t *testing.T) {
		cachePath := filepath.Join(tmpDir, "new-cache.json")
		cache, err := NewArticleCache(cachePath)
		if err != nil {
			t.Fatalf("Failed to create new cache: %v", err)
		}
		if len(cache.Articles) != 0 {
			t.Errorf("Expected empty cache, got %d entries", len(cache.Articles))
		}
	})

	t.Run("load existing cache file", func(t *testing.T) {
		cachePath := filepath.Join(tmpDir, "existing-cache.json")
		testData := map[string]time.Time{
			"https://example.com/1": time.Now(),
			"https://example.com/2": time.Now().Add(-24 * time.Hour),
		}

		// Create test cache file
		data, _ := json.Marshal(testData)
		if err := os.WriteFile(cachePath, data, 0644); err != nil {
			t.Fatalf("Failed to create test cache file: %v", err)
		}

		cache, err := NewArticleCache(cachePath)
		if err != nil {
			t.Fatalf("Failed to load existing cache: %v", err)
		}
		if len(cache.Articles) != len(testData) {
			t.Errorf("Expected %d entries, got %d", len(testData), len(cache.Articles))
		}
	})

	t.Run("handle invalid cache file", func(t *testing.T) {
		cachePath := filepath.Join(tmpDir, "invalid-cache.json")
		if err := os.WriteFile(cachePath, []byte("invalid json"), 0644); err != nil {
			t.Fatalf("Failed to create invalid cache file: %v", err)
		}

		_, err := NewArticleCache(cachePath)
		if err == nil {
			t.Error("Expected error for invalid cache file, got nil")
		}
	})
}

func TestArticleCache_Operations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-ops-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cachePath := filepath.Join(tmpDir, "test-cache.json")
	cache, err := NewArticleCache(cachePath)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	testURL := "https://example.com/article"

	t.Run("check non-existent article", func(t *testing.T) {
		if cache.IsArticlePushed(testURL) {
			t.Error("Expected false for non-existent article")
		}
	})

	t.Run("mark article as pushed", func(t *testing.T) {
		if err := cache.MarkArticlePushed(testURL); err != nil {
			t.Fatalf("Failed to mark article as pushed: %v", err)
		}

		if !cache.IsArticlePushed(testURL) {
			t.Error("Expected true for pushed article")
		}
	})

	t.Run("verify persistence", func(t *testing.T) {
		// Create new cache instance with same file
		newCache, err := NewArticleCache(cachePath)
		if err != nil {
			t.Fatalf("Failed to create new cache instance: %v", err)
		}

		if !newCache.IsArticlePushed(testURL) {
			t.Error("Expected true for pushed article after reload")
		}
	})
}

func TestArticleCache_Save(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-save-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("save to read-only directory", func(t *testing.T) {
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		if err := os.MkdirAll(readOnlyDir, 0444); err != nil {
			t.Fatalf("Failed to create read-only directory: %v", err)
		}

		cachePath := filepath.Join(readOnlyDir, "cache.json")
		cache, _ := NewArticleCache(cachePath)

		err := cache.MarkArticlePushed("https://example.com")
		if err == nil {
			t.Error("Expected error when saving to read-only directory")
		}
	})

	t.Run("save and load multiple entries", func(t *testing.T) {
		cachePath := filepath.Join(tmpDir, "multi-cache.json")
		cache, _ := NewArticleCache(cachePath)

		urls := []string{
			"https://example.com/1",
			"https://example.com/2",
			"https://example.com/3",
		}

		for _, url := range urls {
			if err := cache.MarkArticlePushed(url); err != nil {
				t.Fatalf("Failed to mark article as pushed: %v", err)
			}
		}

		// Load cache again and verify
		newCache, _ := NewArticleCache(cachePath)
		for _, url := range urls {
			if !newCache.IsArticlePushed(url) {
				t.Errorf("Expected article %s to be pushed", url)
			}
		}
	})
}
