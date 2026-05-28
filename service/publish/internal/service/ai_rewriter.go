package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"opengeo/pkg/ai"
	"opengeo/pkg/similarity"
)

// AIRewriter AI改写服务
type AIRewriter struct {
	aiService ai.AIService
}

// NewAIRewriter 创建AI改写服务
func NewAIRewriter(aiService ai.AIService) *AIRewriter {
	return &AIRewriter{
		aiService: aiService,
	}
}

// AIRewriteRequest AI改写请求
type AIRewriteRequest struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	Style   string `json:"style"`   // light, medium, heavy
	Purpose string `json:"purpose"` // dedup, optimize, adapt
}

// AIRewriteResponse AI改写响应
type AIRewriteResponse struct {
	Success bool   `json:"success"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Details string `json:"details"`
}

// Rewrite 使用AI改写内容
func (r *AIRewriter) Rewrite(ctx context.Context, req *AIRewriteRequest) (*AIRewriteResponse, error) {
	if r.aiService == nil {
		return &AIRewriteResponse{
			Success: false,
			Title:   req.Title,
			Body:    req.Body,
			Details: "AI service not available",
		}, nil
	}

	// 调用AI服务
	resp, err := r.aiService.OptimizeContent(ctx, &ai.OptimizeRequest{
		Title:            req.Title,
		Body:             req.Body,
		ContentType:      "article",
		OptimizationType: "dedup",
	})

	if err != nil {
		return &AIRewriteResponse{
			Success: false,
			Title:   req.Title,
			Body:    req.Body,
			Details: fmt.Sprintf("AI rewrite failed: %v", err),
		}, nil
	}

	return &AIRewriteResponse{
		Success: resp.Success,
		Title:   resp.OptimizedTitle,
		Body:    resp.OptimizedBody,
		Details: "AI rewrite completed",
	}, nil
}

// buildRewritePrompt 构建改写提示词
func (r *AIRewriter) buildRewritePrompt(req *AIRewriteRequest) string {
	var styleDesc string
	switch req.Style {
	case "light":
		styleDesc = "轻微改写，保持原文结构，只替换少量词汇"
	case "medium":
		styleDesc = "中等改写，调整段落顺序，替换部分表达"
	case "heavy":
		styleDesc = "深度改写，重新组织内容结构，大幅调整表达方式"
	default:
		styleDesc = "中等改写"
	}

	return fmt.Sprintf(`请对以下内容进行去重改写，要求：
1. %s
2. 保持原文核心信息不变
3. 确保改写后的内容与原文有明显差异
4. 保持内容的可读性和专业性

标题: %s
正文:
%s

请以JSON格式返回：
{
  "title": "改写后的标题",
  "body": "改写后的正文",
  "changes": ["变更说明1", "变更说明2"]
}`, styleDesc, req.Title, req.Body)
}

// AIRewriteWithComparison AI改写并对比
func (r *AIRewriter) RewriteWithComparison(ctx context.Context, req *AIRewriteRequest, originalSimilarity *similarity.CombinedSimilarity) (*AIRewriteResponse, error) {
	// 先进行AI改写
	resp, err := r.Rewrite(ctx, req)
	if err != nil {
		return resp, err
	}

	if !resp.Success {
		return resp, nil
	}

	// 检查改写后的内容是否足够不同
	if originalSimilarity != nil {
		simResult := originalSimilarity.Compute(req.Body, resp.Body, 0.7)
		if simResult.CombinedSimilarity > 0.8 {
			// 相似度太高，需要再次改写
			resp.Details = fmt.Sprintf("AI rewrite completed but similarity too high (%.2f), may need more changes", simResult.CombinedSimilarity)
		}
	}

	return resp, nil
}

// BatchRewrite 批量AI改写
func (r *AIRewriter) BatchRewrite(ctx context.Context, requests []*AIRewriteRequest) ([]*AIRewriteResponse, error) {
	responses := make([]*AIRewriteResponse, len(requests))

	for i, req := range requests {
		resp, err := r.Rewrite(ctx, req)
		if err != nil {
			responses[i] = &AIRewriteResponse{
				Success: false,
				Title:   req.Title,
				Body:    req.Body,
				Details: fmt.Sprintf("Batch rewrite failed: %v", err),
			}
		} else {
			responses[i] = resp
		}
	}

	return responses, nil
}

// RewriteForPlatform 为特定平台改写内容
func (r *AIRewriter) RewriteForPlatform(ctx context.Context, req *AIRewriteRequest, platform string) (*AIRewriteResponse, error) {
	// 根据平台调整改写策略
	platformStyles := map[string]string{
		"wechat":      "适合微信公众号的风格，段落清晰，适当使用emoji",
		"weibo":       "适合微博的风格，简洁明了，突出重点",
		"zhihu":       "适合知乎的风格，专业深入，有理有据",
		"toutiao":     "适合今日头条的风格，吸引眼球，通俗易懂",
		"douyin":      "适合抖音的风格，简短有力，节奏感强",
		"xiaohongshu": "适合小红书的风格，亲切自然，分享感强",
	}

	// 构建平台特定的提示词
	platformStyle := platformStyles[platform]
	if platformStyle != "" {
		req.Style = req.Style + "，" + platformStyle
	}

	return r.Rewrite(ctx, req)
}

// ParseAIRewriteResponse 解析AI改写响应
func ParseAIRewriteResponse(raw string) (*AIRewriteResponse, error) {
	// 尝试解析JSON
	var result struct {
		Title   string   `json:"title"`
		Body    string   `json:"body"`
		Changes []string `json:"changes"`
	}

	// 提取JSON部分
	jsonStr := extractJSON(raw)
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("parse AI response: %w", err)
	}

	return &AIRewriteResponse{
		Success: true,
		Title:   result.Title,
		Body:    result.Body,
		Details: strings.Join(result.Changes, "; "),
	}, nil
}

func extractJSON(text string) string {
	// 查找第一个 { 和最后一个 }
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start >= 0 && end > start {
		return text[start : end+1]
	}
	return text
}
