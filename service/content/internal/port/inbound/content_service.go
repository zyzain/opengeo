package inbound

import (
	"context"

	"opengeo/service/content/internal/domain/model"
)

// ContentService 内容服务接口（入站端口）
// 这个接口将由Kitex生成的代码实现
type ContentService interface {
	// 内容管理
	CreateContent(ctx context.Context, req *CreateContentRequest) (*CreateContentResponse, error)
	GetContent(ctx context.Context, req *GetContentRequest) (*GetContentResponse, error)
	UpdateContent(ctx context.Context, req *UpdateContentRequest) (*UpdateContentResponse, error)
	DeleteContent(ctx context.Context, req *DeleteContentRequest) (*DeleteContentResponse, error)
	ListContents(ctx context.Context, req *ListContentsRequest) (*ListContentsResponse, error)

	// 内容版本管理
	CreateContentVersion(ctx context.Context, req *CreateContentVersionRequest) (*CreateContentVersionResponse, error)
	GetContentVersion(ctx context.Context, req *GetContentVersionRequest) (*GetContentVersionResponse, error)
	ListContentVersions(ctx context.Context, req *ListContentVersionsRequest) (*ListContentVersionsResponse, error)

	// AI内容优化
	OptimizeContentForAI(ctx context.Context, req *OptimizeContentForAIRequest) (*OptimizeContentForAIResponse, error)
	AdaptContentForModel(ctx context.Context, req *AdaptContentForModelRequest) (*AdaptContentForModelResponse, error)
	CheckContentCompliance(ctx context.Context, req *CheckContentComplianceRequest) (*CheckContentComplianceResponse, error)

	// 内容模板管理
	CreateContentTemplate(ctx context.Context, req *CreateContentTemplateRequest) (*CreateContentTemplateResponse, error)
	GetContentTemplate(ctx context.Context, req *GetContentTemplateRequest) (*GetContentTemplateResponse, error)
	ListContentTemplates(ctx context.Context, req *ListContentTemplatesRequest) (*ListContentTemplatesResponse, error)

	// 知识图谱实体管理
	CreateKnowledgeEntity(ctx context.Context, req *CreateKnowledgeEntityRequest) (*CreateKnowledgeEntityResponse, error)
	GetKnowledgeEntity(ctx context.Context, req *GetKnowledgeEntityRequest) (*GetKnowledgeEntityResponse, error)
	ListKnowledgeEntities(ctx context.Context, req *ListKnowledgeEntitiesRequest) (*ListKnowledgeEntitiesResponse, error)
	LinkEntityToContent(ctx context.Context, req *LinkEntityToContentRequest) (*LinkEntityToContentResponse, error)
}

// 请求/响应模型
type CreateContentRequest struct {
	UserID       int64  `json:"user_id"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	ContentType  string `json:"content_type"`
	SchemaMarkup string `json:"schema_markup"`
}

type CreateContentResponse struct {
	Content *model.Content `json:"content"`
}

type GetContentRequest struct {
	ContentID int64 `json:"content_id"`
}

type GetContentResponse struct {
	Content *model.Content `json:"content"`
}

type UpdateContentRequest struct {
	ContentID    int64  `json:"content_id"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	SchemaMarkup string `json:"schema_markup"`
}

type UpdateContentResponse struct {
	Content *model.Content `json:"content"`
}

type DeleteContentRequest struct {
	ContentID int64 `json:"content_id"`
}

type DeleteContentResponse struct {
	Success bool `json:"success"`
}

type ListContentsRequest struct {
	UserID      int64  `json:"user_id"`
	Page        int32  `json:"page"`
	PageSize    int32  `json:"page_size"`
	Status      int32  `json:"status"`
	ContentType string `json:"content_type"`
}

type ListContentsResponse struct {
	Contents []*model.Content `json:"contents"`
	Total    int32            `json:"total"`
}

type CreateContentVersionRequest struct {
	ContentID        int64  `json:"content_id"`
	Title            string `json:"title"`
	Body             string `json:"body"`
	SchemaMarkup     string `json:"schema_markup"`
	AIModelAdaptation string `json:"ai_model_adaptation"`
}

type CreateContentVersionResponse struct {
	Version *model.ContentVersion `json:"version"`
}

type GetContentVersionRequest struct {
	VersionID int64 `json:"version_id"`
}

type GetContentVersionResponse struct {
	Version *model.ContentVersion `json:"version"`
}

type ListContentVersionsRequest struct {
	ContentID int64 `json:"content_id"`
	Page      int32 `json:"page"`
	PageSize  int32 `json:"page_size"`
}

type ListContentVersionsResponse struct {
	Versions []*model.ContentVersion `json:"versions"`
	Total    int32                   `json:"total"`
}

type OptimizeContentForAIRequest struct {
	ContentID        int64  `json:"content_id"`
	OptimizationType string `json:"optimization_type"`
}

type OptimizeContentForAIResponse struct {
	Success          bool    `json:"success"`
	OptimizedContent string  `json:"optimized_content"`
	Score            float32 `json:"score"`
	Details          string  `json:"details"`
}

type AdaptContentForModelRequest struct {
	ContentID   int64  `json:"content_id"`
	TargetModel string `json:"target_model"`
}

type AdaptContentForModelResponse struct {
	Success          bool   `json:"success"`
	AdaptedContent   string `json:"adapted_content"`
	ModelSpecificData string `json:"model_specific_data"`
}

type CheckContentComplianceRequest struct {
	ContentID int64 `json:"content_id"`
}

type CheckContentComplianceResponse struct {
	Compliant bool               `json:"compliant"`
	Issues    []ComplianceIssue  `json:"issues"`
	Report    string             `json:"report"`
}

type ComplianceIssue struct {
	IssueType   string `json:"issue_type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Suggestion  string `json:"suggestion"`
}

type CreateContentTemplateRequest struct {
	UserID       int64  `json:"user_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	TemplateType string `json:"template_type"`
	TemplateData string `json:"template_data"`
	IsPublic     bool   `json:"is_public"`
}

type CreateContentTemplateResponse struct {
	Template *model.ContentTemplate `json:"template"`
}

type GetContentTemplateRequest struct {
	TemplateID int64 `json:"template_id"`
}

type GetContentTemplateResponse struct {
	Template *model.ContentTemplate `json:"template"`
}

type ListContentTemplatesRequest struct {
	UserID       int64  `json:"user_id"`
	TemplateType string `json:"template_type"`
	IsPublic     bool   `json:"is_public"`
	Page         int32  `json:"page"`
	PageSize     int32  `json:"page_size"`
}

type ListContentTemplatesResponse struct {
	Templates []*model.ContentTemplate `json:"templates"`
	Total     int32                    `json:"total"`
}

type CreateKnowledgeEntityRequest struct {
	UserID         int64  `json:"user_id"`
	EntityName     string `json:"entity_name"`
	EntityType     string `json:"entity_type"`
	EntityData     string `json:"entity_data"`
	AuthorityLinks string `json:"authority_links"`
}

type CreateKnowledgeEntityResponse struct {
	Entity *model.KnowledgeEntity `json:"entity"`
}

type GetKnowledgeEntityRequest struct {
	EntityID int64 `json:"entity_id"`
}

type GetKnowledgeEntityResponse struct {
	Entity *model.KnowledgeEntity `json:"entity"`
}

type ListKnowledgeEntitiesRequest struct {
	UserID     int64  `json:"user_id"`
	EntityType string `json:"entity_type"`
	Page       int32  `json:"page"`
	PageSize   int32  `json:"page_size"`
}

type ListKnowledgeEntitiesResponse struct {
	Entities []*model.KnowledgeEntity `json:"entities"`
	Total    int32                    `json:"total"`
}

type LinkEntityToContentRequest struct {
	ContentID int64 `json:"content_id"`
	EntityID  int64 `json:"entity_id"`
}

type LinkEntityToContentResponse struct {
	Success bool `json:"success"`
}