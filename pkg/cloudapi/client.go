package cloudapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 云端 API 客户端
type Client struct {
	baseURL    string
	apiKey     string
	apiSecret  string
	httpClient *http.Client
	cache      *Cache
}

// NewClient 创建云端 API 客户端
func NewClient(baseURL, apiKey, apiSecret string) *Client {
	return &Client{
		baseURL:   baseURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: NewCache(1000, 1*time.Hour),
	}
}

// AttributionRequest 归因分析请求
type AttributionRequest struct {
	TenantID  int64  `json:"tenant_id"`
	ContentID int64  `json:"content_id"`
	BrandID   int64  `json:"brand_id"`
	QueryText string `json:"query_text"`
}

// AttributionResponse 归因分析响应
type AttributionResponse struct {
	RequestID            string   `json:"request_id"`
	AttributionScore     float64  `json:"attribution_score"`
	CitedFragments       []string `json:"cited_fragments"`
	SourceAnalysis       string   `json:"source_analysis"`
	CompetitorComparison string   `json:"competitor_comparison"`
	Recommendations      []string `json:"recommendations"`
	TokensUsed           int64    `json:"tokens_used"`
	CostCents            int64    `json:"cost_cents"`
}

// TrustScoreRequest 可信度评分请求
type TrustScoreRequest struct {
	TenantID int64 `json:"tenant_id"`
	BrandID  int64 `json:"brand_id"`
}

// TrustScoreResponse 可信度评分响应
type TrustScoreResponse struct {
	RequestID     string  `json:"request_id"`
	OverallScore  float64 `json:"overall_score"`
	SearchScore   float64 `json:"search_score"`
	SocialScore   float64 `json:"social_score"`
	ComplianceScore float64 `json:"compliance_score"`
	CitationScore float64 `json:"citation_score"`
	Factors       string  `json:"factors"`
	TokensUsed    int64   `json:"tokens_used"`
	CostCents     int64   `json:"cost_cents"`
}

// ComplianceCheckRequest 合规校验请求
type ComplianceCheckRequest struct {
	TenantID  int64    `json:"tenant_id"`
	ContentID int64    `json:"content_id"`
	BrandID   int64    `json:"brand_id"`
	CheckTypes []string `json:"check_types"`
}

// ComplianceCheckResponse 合规校验响应
type ComplianceCheckResponse struct {
	RequestID   string             `json:"request_id"`
	IsCompliant bool               `json:"is_compliant"`
	RiskLevel   string             `json:"risk_level"`
	Issues      []ComplianceIssue  `json:"issues"`
	Suggestions []string           `json:"suggestions"`
	TokensUsed  int64              `json:"tokens_used"`
	CostCents   int64              `json:"cost_cents"`
}

// ComplianceIssue 合规问题
type ComplianceIssue struct {
	IssueType   string `json:"issue_type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Location    string `json:"location"`
	Suggestion  string `json:"suggestion"`
}

// AnalyzeAttribution 分析归因（带缓存和离线降级）
func (c *Client) AnalyzeAttribution(ctx context.Context, req *AttributionRequest) (*AttributionResponse, error) {
	cacheKey := fmt.Sprintf("attribution:%d:%d:%s", req.BrandID, req.ContentID, req.QueryText)

	// 检查缓存
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.(*AttributionResponse), nil
	}

	// 调用云端 API
	resp := &AttributionResponse{}
	if err := c.doRequest(ctx, "POST", "/api/v1/attribution", req, resp); err != nil {
		// 离线降级：返回基础结果
		return c.fallbackAttribution(req), nil
	}

	// 缓存结果
	c.cache.Set(cacheKey, resp)

	return resp, nil
}

// GetBrandTrustScore 获取品牌可信度评分（带缓存）
func (c *Client) GetBrandTrustScore(ctx context.Context, req *TrustScoreRequest) (*TrustScoreResponse, error) {
	cacheKey := fmt.Sprintf("trust_score:%d:%d", req.TenantID, req.BrandID)

	// 检查缓存
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.(*TrustScoreResponse), nil
	}

	// 调用云端 API
	resp := &TrustScoreResponse{}
	if err := c.doRequest(ctx, "POST", "/api/v1/trust-score", req, resp); err != nil {
		// 离线降级：返回基础评分
		return c.fallbackTrustScore(req), nil
	}

	// 缓存结果
	c.cache.Set(cacheKey, resp)

	return resp, nil
}

// CheckCompliance 合规校验（带缓存）
func (c *Client) CheckCompliance(ctx context.Context, req *ComplianceCheckRequest) (*ComplianceCheckResponse, error) {
	cacheKey := fmt.Sprintf("compliance:%d:%d", req.BrandID, req.ContentID)

	// 检查缓存
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.(*ComplianceCheckResponse), nil
	}

	// 调用云端 API
	resp := &ComplianceCheckResponse{}
	if err := c.doRequest(ctx, "POST", "/api/v1/compliance", req, resp); err != nil {
		// 离线降级：返回基础校验结果
		return c.fallbackCompliance(req), nil
	}

	// 缓存结果
	c.cache.Set(cacheKey, resp)

	return resp, nil
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(ctx context.Context, method, path string, req interface{}, resp interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)
	httpReq.Header.Set("X-API-Secret", c.apiSecret)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %d", httpResp.StatusCode)
	}

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(respBody, resp)
}

// fallbackAttribution 离线降级：归因分析
func (c *Client) fallbackAttribution(req *AttributionRequest) *AttributionResponse {
	return &AttributionResponse{
		RequestID:        "offline",
		AttributionScore: 50.0,
		Recommendations:  []string{"无法连接云端服务，请检查网络连接"},
	}
}

// fallbackTrustScore 离线降级：可信度评分
func (c *Client) fallbackTrustScore(req *TrustScoreRequest) *TrustScoreResponse {
	return &TrustScoreResponse{
		RequestID:    "offline",
		OverallScore: 50.0,
		Factors:      "离线模式，使用本地缓存数据",
	}
}

// fallbackCompliance 离线降级：合规校验
func (c *Client) fallbackCompliance(req *ComplianceCheckRequest) *ComplianceCheckResponse {
	return &ComplianceCheckResponse{
		RequestID:   "offline",
		IsCompliant: true,
		RiskLevel:   "unknown",
		Suggestions: []string{"无法连接云端服务，使用本地规则校验"},
	}
}
