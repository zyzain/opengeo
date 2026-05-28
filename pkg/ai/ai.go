package ai

import "context"

// AIService AI服务接口
type AIService interface {
	// OptimizeContent GEO语义增强
	OptimizeContent(ctx context.Context, req *OptimizeRequest) (*OptimizeResponse, error)
	// AdaptForModel 多模型适配改写
	AdaptForModel(ctx context.Context, req *AdaptRequest) (*AdaptResponse, error)
	// CheckCompliance 合规与安全检测
	CheckCompliance(ctx context.Context, req *ComplianceRequest) (*ComplianceResponse, error)
}

// OptimizeRequest 优化请求
type OptimizeRequest struct {
	ContentID        int64  `json:"content_id"`
	Title            string `json:"title"`
	Body             string `json:"body"`
	ContentType      string `json:"content_type"`
	OptimizationType string `json:"optimization_type"` // geo_semantic, structure, readability, dedup
}

// OptimizeResponse 优化响应
type OptimizeResponse struct {
	Success           bool     `json:"success"`
	OptimizedTitle    string   `json:"optimized_title"`
	OptimizedBody     string   `json:"optimized_body"`
	SchemaMarkup      string   `json:"schema_markup"`
	Score             float32  `json:"score"`
	Suggestions       []string `json:"suggestions"`
	StructuralChanges []string `json:"structural_changes"`
}

// AdaptRequest 适配请求
type AdaptRequest struct {
	ContentID   int64  `json:"content_id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	TargetModel string `json:"target_model"` // deepseek, kimi, doubao, chatgpt
}

// AdaptResponse 适配响应
type AdaptResponse struct {
	Success           bool                   `json:"success"`
	AdaptedTitle      string                 `json:"adapted_title"`
	AdaptedBody       string                 `json:"adapted_body"`
	FormatChanges     []string               `json:"format_changes"`
	ModelSpecificData map[string]interface{} `json:"model_specific_data"`
}

// ComplianceRequest 合规检测请求
type ComplianceRequest struct {
	ContentID int64  `json:"content_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
}

// ComplianceResponse 合规检测响应
type ComplianceResponse struct {
	Compliant         bool              `json:"compliant"`
	Issues            []ComplianceIssue `json:"issues"`
	SensitiveWords    []string          `json:"sensitive_words"`
	AIGCLabelRequired bool              `json:"aigc_label_required"`
	Report            string            `json:"report"`
	Score             float32           `json:"score"`
}

// ComplianceIssue 合规问题
type ComplianceIssue struct {
	IssueType   string `json:"issue_type"` // sensitive, ad_law, aigc_label
	Description string `json:"description"`
	Severity    string `json:"severity"` // high, medium, low
	Suggestion  string `json:"suggestion"`
	Location    string `json:"location"` // 问题位置
}
