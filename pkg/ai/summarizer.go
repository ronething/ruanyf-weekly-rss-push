package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ashing/ruanyf-weekly-rss-push/pkg/feed"
	ark "github.com/sashabaranov/go-openai"
)

// Recommendation represents a single recommended item extracted from the newsletter
type Recommendation struct {
	Title       string
	Description string
	Link        string
}

// ContentSummarizer handles the AI-powered summarization of newsletter content
type ContentSummarizer struct {
	client *ark.Client
}

// NewContentSummarizer creates a new content summarizer using Ark API
func NewContentSummarizer() (*ContentSummarizer, error) {
	apiKey := os.Getenv("ARK_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ARK_API_KEY environment variable is not set")
	}

	config := ark.DefaultConfig(apiKey)
	baseURL := os.Getenv("ARK_BASE_URL")
	if baseURL != "" {
		config.BaseURL = baseURL
	} else {
		config.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	}

	return &ContentSummarizer{
		client: ark.NewClientWithConfig(config),
	}, nil
}

// ExtractRecommendations processes RSS entry content and extracts recommended items
func (s *ContentSummarizer) ExtractRecommendations(entry feed.Entry) ([]Recommendation, error) {
	prompt := buildPrompt(entry.Content.Src)

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		ark.ChatCompletionRequest{
			Model: os.Getenv("ARK_MODEL_ID"),
			Messages: []ark.ChatCompletionMessage{
				{
					Role:    ark.ChatMessageRoleSystem,
					Content: "你是一位内容分析专家，擅长从技术文章中提取最有价值的信息。",
				},
				{
					Role:    ark.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("AI completion error: %v", err)
	}

	log.Printf("AI response: %+v\n", resp.Choices[0].Message.Content)

	return parseRecommendations(resp.Choices[0].Message.Content)
}

// FormatRecommendationsText formats the recommendations as text for the notification
func FormatRecommendationsText(recs []Recommendation) string {
	if len(recs) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("本期推荐：\n\n")

	for i, rec := range recs {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec.Title))
		if rec.Description != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", rec.Description))
		}
		if rec.Link != "" {
			sb.WriteString(fmt.Sprintf("   链接：%s\n", rec.Link))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("AI Recommendation Powered By DeepSeek")

	return sb.String()
}

// buildPrompt creates the prompt for the AI to extract recommendations
func buildPrompt(content string) string {
	return `我需要你分析阮一峰技术周刊的内容，并提取出最有价值的5个推荐项目。

技术周刊内容如下:
"""
` + content + `
"""

请遵循以下规则:
1. 只选择5个最值得推荐的项目
2. 优先选择有实用价值、新颖性或重要技术趋势的项目
3. 每个推荐项应包含标题、简短说明和链接(如有)
4. 输出格式必须严格按照如下JSON格式:
[
  {
    "title": "项目标题",
    "description": "简短说明(不超过100字)",
    "link": "项目链接(如有)"
  },
  ...
]

不要添加任何额外的评论、解释或引言，只输出符合上述格式的JSON。`
}

// parseRecommendations parses the AI response into structured recommendations
func parseRecommendations(aiResponse string) ([]Recommendation, error) {
	// 简单解析AI返回的JSON内容
	// 实际使用时可能需要更复杂的解析逻辑
	var recommendations []Recommendation

	// 查找JSON内容的开始和结束
	startIdx := strings.Index(aiResponse, "[")
	endIdx := strings.LastIndex(aiResponse, "]")

	if startIdx == -1 || endIdx == -1 || startIdx > endIdx {
		return nil, fmt.Errorf("invalid AI response format: %s", aiResponse)
	}

	jsonContent := aiResponse[startIdx : endIdx+1]

	// 使用标准json库解析
	err := json.Unmarshal([]byte(jsonContent), &recommendations)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recommendations: %v", err)
	}

	return recommendations, nil
}
