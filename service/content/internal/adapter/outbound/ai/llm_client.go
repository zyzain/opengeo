package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"opengeo/pkg/ai"
	"opengeo/pkg/config"
)

// LLMConfig LLM配置
type LLMConfig struct {
	BaseURL    string        `json:"base_url"`
	APIKey     string        `json:"api_key"`
	Model      string        `json:"model"`
	Timeout    time.Duration `json:"timeout"`
	MaxTokens  int           `json:"max_tokens"`
	Temperature float32      `json:"temperature"`
}

// DefaultConfig 返回默认配置
func (c *LLMConfig) DefaultConfig() {
	cfg := config.GetConfig()
	if c.Timeout == 0 {
		c.Timeout = cfg.LLM.Timeout
	}
	if c.MaxTokens == 0 {
		c.MaxTokens = cfg.LLM.MaxTokens
	}
	if c.Temperature == 0 {
		c.Temperature = cfg.LLM.Temperature
	}
}

// chatRequest OpenAI兼容的chat请求
type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float32       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatResponse OpenAI兼容的chat响应
type chatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// LLMClient 通用LLM HTTP客户端
type LLMClient struct {
	config     LLMConfig
	httpClient *http.Client
}

// NewLLMClient 创建LLM客户端
func NewLLMClient(config LLMConfig) *LLMClient {
	config.DefaultConfig()
	return &LLMClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Chat 发送chat请求
func (c *LLMClient) Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	var lastErr error
	maxRetries := 3
	baseDelay := time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := baseDelay * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(delay):
			}
		}

		result, err, retry := c.doChat(ctx, systemPrompt, userPrompt)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !retry {
			return "", err
		}
	}

	return "", fmt.Errorf("chat failed after %d retries: %w", maxRetries, lastErr)
}

// doChat 执行单次chat请求，返回结果、错误、是否应重试
func (c *LLMClient) doChat(ctx context.Context, systemPrompt, userPrompt string) (string, error, bool) {
	req := chatRequest{
		Model: c.config.Model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
		Stream:      false,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err), false
	}

	url := strings.TrimRight(c.config.BaseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err), false
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err), true
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err), true
	}

	if resp.StatusCode != http.StatusOK {
		retry := resp.StatusCode >= 500
		return "", fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(respBody)), retry
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err), false
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("llm error: %s", chatResp.Error.Message), false
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from llm"), false
	}

	return chatResp.Choices[0].Message.Content, nil, false
}

// buildOptimizePrompt 构建GEO优化prompt
func buildOptimizePrompt(req *ai.OptimizeRequest) string {
	return fmt.Sprintf(`你是一个GEO（Generative Engine Optimization）优化专家。请对以下内容进行AI搜索优化分析。

标题:
---BEGIN CONTENT---
%s
---END CONTENT---
内容类型:
---BEGIN CONTENT---
%s
---END CONTENT---
正文:
---BEGIN CONTENT---
%s
---END CONTENT---

请以JSON格式返回分析结果，格式如下：
{
  "optimized_title": "优化后的标题",
  "optimized_body": "优化后的正文（保持原意，优化结构）",
  "schema_markup": "建议的JSON-LD结构化数据",
  "score": 85,
  "suggestions": ["建议1", "建议2"],
  "structural_changes": ["变更1", "变更2"]
}

优化要求：
1. 保持原文核心信息不变
2. 添加清晰的标题层级结构（## / ###）
3. 确保段落长度适中（50-200字）
4. 添加列表/要点提炼
5. 建议权威引用位置
6. 优化关键词密度
7. 建议Schema Markup结构化数据`, req.Title, req.ContentType, req.Body)
}

// buildAdaptPrompt 构建模型适配prompt
func buildAdaptPrompt(req *ai.AdaptRequest) string {
	modelHints := map[string]string{
		"deepseek": "DeepSeek偏好：结构化分点论述、有明确结论、有数据支撑",
		"kimi":     "Kimi偏好：长文本详细分析、数据密度高、有权威引用",
		"doubao":   "豆包偏好：简洁明了、要点提炼、数据支撑、直击要点",
		"chatgpt":  "ChatGPT偏好：对话式表达、FAQ格式、上下文连贯、问答结构",
	}

	hint, ok := modelHints[req.TargetModel]
	if !ok {
		hint = "通用优化"
	}

	return fmt.Sprintf(`你是一个内容适配专家。请将以下内容适配为适合%s模型的格式。

适配提示: %s

标题:
---BEGIN CONTENT---
%s
---END CONTENT---
正文:
---BEGIN CONTENT---
%s
---END CONTENT---

请以JSON格式返回：
{
  "adapted_title": "适配后的标题",
  "adapted_body": "适配后的正文",
  "format_changes": ["变更说明1", "变更说明2"]
}`, req.TargetModel, hint, req.Title, req.Body)
}

// buildCompliancePrompt 构建合规检测prompt
func buildCompliancePrompt(req *ai.ComplianceRequest) string {
	return fmt.Sprintf(`你是一个内容合规检测专家。请对以下内容进行全面的合规检测。

标题:
---BEGIN CONTENT---
%s
---END CONTENT---
正文:
---BEGIN CONTENT---
%s
---END CONTENT---

检测项目：
1. 敏感词检测（政治、色情、暴力、赌博、毒品等违法违规内容）
2. 广告法合规（绝对化用语、虚假宣传、夸大效果等）
3. AIGC标识（是否需要添加AI生成内容标识）

请以JSON格式返回：
{
  "compliant": true,
  "issues": [
    {
      "issue_type": "sensitive|ad_law|aigc_label",
      "description": "问题描述",
      "severity": "high|medium|low",
      "suggestion": "修改建议",
      "location": "问题位置"
    }
  ],
  "sensitive_words": ["敏感词1"],
  "aigc_label_required": true,
  "report": "合规检测报告摘要",
  "score": 90
}`, req.Title, req.Body)
}

// parseJSONFromResponse 从LLM响应中提取JSON
func parseJSONFromResponse(response string) string {
	// 尝试直接解析
	response = strings.TrimSpace(response)

	// 处理markdown代码块包裹的JSON
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}

	// 查找第一个 { 和最后一个 }
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start >= 0 && end > start {
		return response[start : end+1]
	}

	return response
}
