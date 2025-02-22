package feed

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestConvertToMarkdownURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "正常的 release URL",
			input:    "https://github.com/ruanyf/weekly/releases/tag/issue-123",
			expected: "https://github.com/ruanyf/weekly/blob/master/docs/issue-123.md",
		},
		{
			name:     "短 release URL",
			input:    "issue-456",
			expected: "https://github.com/ruanyf/weekly/blob/master/docs/issue-456.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToMarkdownURL(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertToMarkdownURL(%s) = %s, 期望 %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	// 模拟 RSS 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rssXML := `<?xml version="1.0" encoding="UTF-8"?>
<feed>
  <entry>
    <title>科技爱好者周刊（第 123 期）</title>
    <link href="https://github.com/ruanyf/weekly/releases/tag/issue-123"/>
    <updated>2023-01-01T00:00:00Z</updated>
  </entry>
  <entry>
    <title>科技爱好者周刊（第 124 期）</title>
    <link href="https://github.com/ruanyf/weekly/releases/tag/issue-124"/>
    <updated>2023-01-08T00:00:00Z</updated>
  </entry>
  <entry>
    <title>科技爱好者周刊（第 125 期）</title>
    <link href="https://github.com/ruanyf/weekly/releases/tag/issue-125"/>
    <updated>2023-01-15T00:00:00Z</updated>
  </entry>
</feed>`
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(rssXML))
	}))
	defer server.Close()

	tests := []struct {
		name        string
		limit       int
		wantEntries int
	}{
		{
			name:        "限制为2条",
			limit:       2,
			wantEntries: 2,
		},
		{
			name:        "限制为1条",
			limit:       1,
			wantEntries: 1,
		},
		{
			name:        "限制为5条（超过实际条数）",
			limit:       5,
			wantEntries: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(server.URL, tt.limit)
			entries, err := parser.Parse()

			if err != nil {
				t.Fatalf("Parse() 错误 = %v", err)
			}

			if len(entries) != tt.wantEntries {
				t.Errorf("Parse() 返回 %d 条记录, 期望 %d 条", len(entries), tt.wantEntries)
			}

			// 检查第一条记录的内容
			if len(entries) > 0 {
				firstEntry := entries[0]
				if firstEntry.Title != "科技爱好者周刊（第 123 期）" {
					t.Errorf("第一条记录标题错误，得到: %s", firstEntry.Title)
				}
				if firstEntry.Link.Href != "https://github.com/ruanyf/weekly/releases/tag/issue-123" {
					t.Errorf("第一条记录链接错误，得到: %s", firstEntry.Link.Href)
				}
				expectedTime, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
				if !firstEntry.Updated.Equal(expectedTime) {
					t.Errorf("第一条记录时间错误，得到: %v, 期望: %v", firstEntry.Updated, expectedTime)
				}
			}
		})
	}
}

func TestParser_Parse_Error(t *testing.T) {
	// 测试无效的 URL
	t.Run("无效的URL", func(t *testing.T) {
		parser := NewParser("http://invalid-url", 2)
		_, err := parser.Parse()
		if err == nil {
			t.Error("期望得到错误，但是没有")
		}
	})

	// 测试无效的 XML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte("invalid xml"))
	}))
	defer server.Close()

	t.Run("无效的XML", func(t *testing.T) {
		parser := NewParser(server.URL, 2)
		_, err := parser.Parse()
		if err == nil {
			t.Error("期望得到 XML 解析错误，但是没有")
		}
	})
}
