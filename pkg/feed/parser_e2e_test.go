package feed

import (
	"strings"
	"testing"
	"time"
)

// 如果不想运行 e2e 测试，可以使用 -short 标志运行测试
// go test -short ./...
func TestRealRSSFeed(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 e2e 测试")
	}

	const (
		feedURL = "https://github.com/ruanyf/weekly/releases.atom"
		limit   = 2
	)

	parser := NewParser(feedURL, limit)
	entries, err := parser.Parse()
	if err != nil {
		t.Fatalf("解析 RSS feed 失败: %v", err)
	}

	// 验证返回条目数量
	if len(entries) != limit {
		t.Errorf("期望获取 %d 条记录，实际获取 %d 条", limit, len(entries))
	}

	// 验证每个条目的格式
	for i, entry := range entries {
		t.Run("验证条目格式", func(t *testing.T) {
			// 标题格式检查
			if !strings.Contains(entry.Title, "issue-") {
				t.Errorf("条目 %d 标题格式错误: %s", i, entry.Title)
			}

			// 链接格式检查
			expectedPrefix := "https://github.com/ruanyf/weekly/releases/tag/issue-"
			if !strings.HasPrefix(entry.Link.Href, expectedPrefix) {
				t.Errorf("条目 %d 链接格式错误: %s", i, entry.Link.Href)
			}

			// 更新时间检查
			if entry.Updated.IsZero() {
				t.Errorf("条目 %d 更新时间为空", i)
			}

			// 确保更新时间不是未来时间
			if entry.Updated.After(time.Now().AddDate(0, 0, 1)) { // 允许1天的时间差异
				t.Errorf("条目 %d 更新时间是未来时间: %v", i, entry.Updated)
			}

			// 测试 URL 转换
			mdURL := ConvertToMarkdownURL(entry.Link.Href)
			expectedMDPrefix := "https://github.com/ruanyf/weekly/blob/master/docs/issue-"
			if !strings.HasPrefix(mdURL, expectedMDPrefix) {
				t.Errorf("条目 %d Markdown URL 转换错误: %s", i, mdURL)
			}

			t.Logf("成功验证第 %d 条记录:", i+1)
			t.Logf("  标题: %s", entry.Title)
			t.Logf("  链接: %s", entry.Link.Href)
			t.Logf("  更新时间: %v", entry.Updated)
			t.Logf("  Markdown URL: %s", mdURL)
		})
	}
}
