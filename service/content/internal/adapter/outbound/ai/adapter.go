package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"opengeo/pkg/ai"
	"opengeo/pkg/config"
)

// ==================== LLM适配器基类 ====================

// llmAdapter 通用LLM适配器基类
type llmAdapter struct {
	client *LLMClient
}

// ==================== DeepSeek适配器 ====================

// DeepSeekAdapter DeepSeek AI适配器
type DeepSeekAdapter struct {
	llmAdapter
}

// NewDeepSeekAdapter 创建DeepSeek适配器
func NewDeepSeekAdapter(apiKey string) *DeepSeekAdapter {
	cfg := config.GetConfig()
	return &DeepSeekAdapter{
		llmAdapter: llmAdapter{
			client: NewLLMClient(LLMConfig{
				BaseURL: cfg.AIModels.DeepSeek.BaseURL,
				APIKey:  apiKey,
				Model:   cfg.AIModels.DeepSeek.Model,
			}),
		},
	}
}

func (a *DeepSeekAdapter) OptimizeContent(ctx context.Context, req *ai.OptimizeRequest) (*ai.OptimizeResponse, error) {
	resp, err := a.client.Chat(ctx,
		"你是一个GEO优化专家，擅长优化内容以提升在AI搜索引擎中的可见性。",
		buildOptimizePrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("deepseek optimize: %w", err)
	}
	return parseOptimizeResponse(resp)
}

func (a *DeepSeekAdapter) AdaptForModel(ctx context.Context, req *ai.AdaptRequest) (*ai.AdaptResponse, error) {
	resp, err := a.client.Chat(ctx,
		"你是一个内容适配专家，擅长将内容适配为不同AI模型偏好的格式。",
		buildAdaptPrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("deepseek adapt: %w", err)
	}
	return parseAdaptResponse(resp)
}

func (a *DeepSeekAdapter) CheckCompliance(ctx context.Context, req *ai.ComplianceRequest) (*ai.ComplianceResponse, error) {
	resp, err := a.client.Chat(ctx,
		"你是一个内容合规检测专家，擅长检测敏感词、广告法违规和AIGC标识问题。",
		buildCompliancePrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("deepseek compliance: %w", err)
	}
	return parseComplianceResponse(resp)
}

// ==================== Kimi适配器 ====================

// KimiAdapter Kimi AI适配器
type KimiAdapter struct {
	llmAdapter
}

// NewKimiAdapter 创建Kimi适配器
func NewKimiAdapter(apiKey string) *KimiAdapter {
	cfg := config.GetConfig()
	return &KimiAdapter{
		llmAdapter: llmAdapter{
			client: NewLLMClient(LLMConfig{
				BaseURL: cfg.AIModels.Kimi.BaseURL,
				APIKey:  apiKey,
				Model:   cfg.AIModels.Kimi.Model,
			}),
		},
	}
}

func (a *KimiAdapter) OptimizeContent(ctx context.Context, req *ai.OptimizeRequest) (*ai.OptimizeResponse, error) {
	resp, err := a.client.Chat(ctx,
		"你是一个GEO优化专家，擅长优化内容以提升在AI搜索引擎中的可见性。你擅长长文本分析和详细的数据支撑。",
		buildOptimizePrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("kimi optimize: %w", err)
	}
	return parseOptimizeResponse(resp)
}

func (a *KimiAdapter) AdaptForModel(ctx context.Context, req *ai.AdaptRequest) (*ai.AdaptResponse, error) {
	resp, err := a.client.Chat(ctx,
		"你是一个内容适配专家，擅长将内容适配为不同AI模型偏好的格式。",
		buildAdaptPrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("kimi adapt: %w", err)
	}
	return parseAdaptResponse(resp)
}

func (a *KimiAdapter) CheckCompliance(ctx context.Context, req *ai.ComplianceRequest) (*ai.ComplianceResponse, error) {
	resp, err := a.client.Chat(ctx,
		"你是一个内容合规检测专家，擅长检测敏感词、广告法违规和AIGC标识问题。",
		buildCompliancePrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("kimi compliance: %w", err)
	}
	return parseComplianceResponse(resp)
}

// ==================== 豆包适配器 ====================

// DoubaoAdapter 豆包 AI适配器
type DoubaoAdapter struct {
	llmAdapter
}

// NewDoubaoAdapter 创建豆包适配器
func NewDoubaoAdapter(apiKey string) *DoubaoAdapter {
	cfg := config.GetConfig()
	return &DoubaoAdapter{
		llmAdapter: llmAdapter{
			client: NewLLMClient(LLMConfig{
				BaseURL: cfg.AIModels.Doubao.BaseURL,
				APIKey:  apiKey,
				Model:   cfg.AIModels.Doubao.Model,
			}),
		},
	}
}

func (a *DoubaoAdapter) OptimizeContent(ctx context.Context, req *ai.OptimizeRequest) (*ai.OptimizeResponse, error) {
	resp, err := a.client.Chat(ctx,
		"你是一个GEO优化专家，擅长优化内容以提升在AI搜索引擎中的可见性。你偏好简洁明了、数据驱动的表达方式。",
		buildOptimizePrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("doubao optimize: %w", err)
	}
	return parseOptimizeResponse(resp)
}

func (a *DoubaoAdapter) AdaptForModel(ctx context.Context, req *ai.AdaptRequest) (*ai.AdaptResponse, error) {
	resp, err := a.client.Chat(ctx,
		"你是一个内容适配专家，擅长将内容适配为不同AI模型偏好的格式。",
		buildAdaptPrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("doubao adapt: %w", err)
	}
	return parseAdaptResponse(resp)
}

func (a *DoubaoAdapter) CheckCompliance(ctx context.Context, req *ai.ComplianceRequest) (*ai.ComplianceResponse, error) {
	resp, err := a.client.Chat(ctx,
		"你是一个内容合规检测专家，擅长检测敏感词、广告法违规和AIGC标识问题。",
		buildCompliancePrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("doubao compliance: %w", err)
	}
	return parseComplianceResponse(resp)
}

// ==================== ChatGPT适配器 ====================

// ChatGPTAdapter ChatGPT AI适配器
type ChatGPTAdapter struct {
	llmAdapter
}

// NewChatGPTAdapter 创建ChatGPT适配器
func NewChatGPTAdapter(apiKey string) *ChatGPTAdapter {
	cfg := config.GetConfig()
	return &ChatGPTAdapter{
		llmAdapter: llmAdapter{
			client: NewLLMClient(LLMConfig{
				BaseURL: cfg.AIModels.ChatGPT.BaseURL,
				APIKey:  apiKey,
				Model:   cfg.AIModels.ChatGPT.Model,
			}),
		},
	}
}

func (a *ChatGPTAdapter) OptimizeContent(ctx context.Context, req *ai.OptimizeRequest) (*ai.OptimizeResponse, error) {
	resp, err := a.client.Chat(ctx,
		"You are a GEO optimization expert. Optimize content for AI search engine visibility. Respond in the same language as the input content.",
		buildOptimizePrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("chatgpt optimize: %w", err)
	}
	return parseOptimizeResponse(resp)
}

func (a *ChatGPTAdapter) AdaptForModel(ctx context.Context, req *ai.AdaptRequest) (*ai.AdaptResponse, error) {
	resp, err := a.client.Chat(ctx,
		"You are a content adaptation expert. Adapt content for different AI model preferences. Respond in the same language as the input content.",
		buildAdaptPrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("chatgpt adapt: %w", err)
	}
	return parseAdaptResponse(resp)
}

func (a *ChatGPTAdapter) CheckCompliance(ctx context.Context, req *ai.ComplianceRequest) (*ai.ComplianceResponse, error) {
	resp, err := a.client.Chat(ctx,
		"You are a content compliance expert. Check for sensitive words, advertising law violations, and AIGC labeling requirements. Respond in the same language as the input content.",
		buildCompliancePrompt(req),
	)
	if err != nil {
		return nil, fmt.Errorf("chatgpt compliance: %w", err)
	}
	return parseComplianceResponse(resp)
}

// ==================== 响应解析 ====================

// optimizeResult LLM返回的优化结果JSON结构
type optimizeResult struct {
	OptimizedTitle    string   `json:"optimized_title"`
	OptimizedBody     string   `json:"optimized_body"`
	SchemaMarkup      string   `json:"schema_markup"`
	Score             float32  `json:"score"`
	Suggestions       []string `json:"suggestions"`
	StructuralChanges []string `json:"structural_changes"`
}

func parseOptimizeResponse(raw string) (*ai.OptimizeResponse, error) {
	jsonStr := parseJSONFromResponse(raw)
	var result optimizeResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("parse optimize response: %w, raw: %s", err, truncate(raw, 200))
	}
	if result.OptimizedTitle == "" {
		return nil, fmt.Errorf("optimize response missing optimized_title, raw: %s", truncate(raw, 200))
	}
	if result.OptimizedBody == "" {
		return nil, fmt.Errorf("optimize response missing optimized_body, raw: %s", truncate(raw, 200))
	}
	return &ai.OptimizeResponse{
		Success:           true,
		OptimizedTitle:    result.OptimizedTitle,
		OptimizedBody:     result.OptimizedBody,
		SchemaMarkup:      result.SchemaMarkup,
		Score:             result.Score,
		Suggestions:       result.Suggestions,
		StructuralChanges: result.StructuralChanges,
	}, nil
}

// adaptResult LLM返回的适配结果JSON结构
type adaptResult struct {
	AdaptedTitle  string   `json:"adapted_title"`
	AdaptedBody   string   `json:"adapted_body"`
	FormatChanges []string `json:"format_changes"`
}

func parseAdaptResponse(raw string) (*ai.AdaptResponse, error) {
	jsonStr := parseJSONFromResponse(raw)
	var result adaptResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("parse adapt response: %w, raw: %s", err, truncate(raw, 200))
	}
	if result.AdaptedTitle == "" {
		return nil, fmt.Errorf("adapt response missing adapted_title, raw: %s", truncate(raw, 200))
	}
	if result.AdaptedBody == "" {
		return nil, fmt.Errorf("adapt response missing adapted_body, raw: %s", truncate(raw, 200))
	}
	return &ai.AdaptResponse{
		Success:       true,
		AdaptedTitle:  result.AdaptedTitle,
		AdaptedBody:   result.AdaptedBody,
		FormatChanges: result.FormatChanges,
		ModelSpecificData: map[string]interface{}{
			"source": "llm",
		},
	}, nil
}

// complianceResult LLM返回的合规结果JSON结构
type complianceResult struct {
	Compliant         bool             `json:"compliant"`
	Issues            []ai.ComplianceIssue `json:"issues"`
	SensitiveWords    []string         `json:"sensitive_words"`
	AIGCLabelRequired bool             `json:"aigc_label_required"`
	Report            string           `json:"report"`
	Score             float32          `json:"score"`
}

func parseComplianceResponse(raw string) (*ai.ComplianceResponse, error) {
	jsonStr := parseJSONFromResponse(raw)
	var result complianceResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("parse compliance response: %w, raw: %s", err, truncate(raw, 200))
	}
	return &ai.ComplianceResponse{
		Compliant:         result.Compliant,
		Issues:            result.Issues,
		SensitiveWords:    result.SensitiveWords,
		AIGCLabelRequired: result.AIGCLabelRequired,
		Report:            result.Report,
		Score:             result.Score,
	}, nil
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// ==================== AI服务工厂 ====================

// AIServiceFactory AI服务工厂
type AIServiceFactory struct {
	adapters map[string]ai.AIService
}

// NewAIServiceFactory 创建AI服务工厂
func NewAIServiceFactory() *AIServiceFactory {
	return &AIServiceFactory{
		adapters: make(map[string]ai.AIService),
	}
}

// RegisterAdapter 注册适配器
func (f *AIServiceFactory) RegisterAdapter(model string, adapter ai.AIService) {
	f.adapters[model] = adapter
}

// GetAdapter 获取适配器
func (f *AIServiceFactory) GetAdapter(model string) (ai.AIService, error) {
	adapter, ok := f.adapters[model]
	if !ok {
		return nil, fmt.Errorf("unsupported AI model: %s", model)
	}
	return adapter, nil
}

// GetDefaultAdapter 获取默认适配器
func (f *AIServiceFactory) GetDefaultAdapter() (ai.AIService, error) {
	cfg := config.GetConfig()
	if adapter, ok := f.adapters[cfg.AIModels.Default]; ok {
		return adapter, nil
	}
	for _, adapter := range f.adapters {
		return adapter, nil
	}
	return nil, fmt.Errorf("no AI adapter registered")
}
