package inbound

import (
	"context"

	"opengeo/pkg/ai"
	"opengeo/service/content/internal/application/service"
	"opengeo/service/content/internal/port/inbound"
)

// ContentKitexHandler Content服务的Kitex处理器（入站适配器）
type ContentKitexHandler struct {
	appService *service.ContentAppService
}

// NewContentKitexHandler 创建Content服务的Kitex处理器
func NewContentKitexHandler(appService *service.ContentAppService) *ContentKitexHandler {
	return &ContentKitexHandler{
		appService: appService,
	}
}

// CreateContent 创建内容
func (h *ContentKitexHandler) CreateContent(ctx context.Context, req *inbound.CreateContentRequest) (*inbound.CreateContentResponse, error) {
	appReq := &service.CreateContentRequest{
		UserID:       req.UserID,
		Title:        req.Title,
		Body:         req.Body,
		ContentType:  req.ContentType,
		SchemaMarkup: req.SchemaMarkup,
	}

	appResp, err := h.appService.CreateContent(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.CreateContentResponse{
		Content: appResp.Content,
	}, nil
}

// GetContent 获取内容
func (h *ContentKitexHandler) GetContent(ctx context.Context, req *inbound.GetContentRequest) (*inbound.GetContentResponse, error) {
	appReq := &service.GetContentRequest{
		ContentID: req.ContentID,
	}

	appResp, err := h.appService.GetContent(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.GetContentResponse{
		Content: appResp.Content,
	}, nil
}

// UpdateContent 更新内容
func (h *ContentKitexHandler) UpdateContent(ctx context.Context, req *inbound.UpdateContentRequest) (*inbound.UpdateContentResponse, error) {
	appReq := &service.UpdateContentRequest{
		ContentID:    req.ContentID,
		Title:        req.Title,
		Body:         req.Body,
		SchemaMarkup: req.SchemaMarkup,
	}

	appResp, err := h.appService.UpdateContent(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.UpdateContentResponse{
		Content: appResp.Content,
	}, nil
}

// DeleteContent 删除内容
func (h *ContentKitexHandler) DeleteContent(ctx context.Context, req *inbound.DeleteContentRequest) (*inbound.DeleteContentResponse, error) {
	appReq := &service.DeleteContentRequest{
		ContentID: req.ContentID,
	}

	appResp, err := h.appService.DeleteContent(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.DeleteContentResponse{
		Success: appResp.Success,
	}, nil
}

// ListContents 列出内容
func (h *ContentKitexHandler) ListContents(ctx context.Context, req *inbound.ListContentsRequest) (*inbound.ListContentsResponse, error) {
	appReq := &service.ListContentsRequest{
		UserID:      req.UserID,
		Page:        req.Page,
		PageSize:    req.PageSize,
		Status:      req.Status,
		ContentType: req.ContentType,
	}

	appResp, err := h.appService.ListContents(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.ListContentsResponse{
		Contents: appResp.Contents,
		Total:    appResp.Total,
	}, nil
}

// CreateContentVersion 创建内容版本
func (h *ContentKitexHandler) CreateContentVersion(ctx context.Context, req *inbound.CreateContentVersionRequest) (*inbound.CreateContentVersionResponse, error) {
	appReq := &service.CreateContentVersionRequest{
		ContentID:        req.ContentID,
		Title:            req.Title,
		Body:             req.Body,
		SchemaMarkup:     req.SchemaMarkup,
		AIModelAdaptation: req.AIModelAdaptation,
	}

	appResp, err := h.appService.CreateContentVersion(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.CreateContentVersionResponse{
		Version: appResp.Version,
	}, nil
}

// GetContentVersion 获取内容版本
func (h *ContentKitexHandler) GetContentVersion(ctx context.Context, req *inbound.GetContentVersionRequest) (*inbound.GetContentVersionResponse, error) {
	appReq := &service.GetContentVersionRequest{
		VersionID: req.VersionID,
	}

	appResp, err := h.appService.GetContentVersion(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.GetContentVersionResponse{
		Version: appResp.Version,
	}, nil
}

// ListContentVersions 列出内容版本
func (h *ContentKitexHandler) ListContentVersions(ctx context.Context, req *inbound.ListContentVersionsRequest) (*inbound.ListContentVersionsResponse, error) {
	appReq := &service.ListContentVersionsRequest{
		ContentID: req.ContentID,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}

	appResp, err := h.appService.ListContentVersions(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.ListContentVersionsResponse{
		Versions: appResp.Versions,
		Total:    appResp.Total,
	}, nil
}

// OptimizeContentForAI 为AI优化内容
func (h *ContentKitexHandler) OptimizeContentForAI(ctx context.Context, req *inbound.OptimizeContentForAIRequest) (*inbound.OptimizeContentForAIResponse, error) {
	appReq := &service.OptimizeContentForAIRequest{
		ContentID:        req.ContentID,
		OptimizationType: req.OptimizationType,
	}

	appResp, err := h.appService.OptimizeContentForAI(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.OptimizeContentForAIResponse{
		Success:          appResp.Success,
		OptimizedContent: appResp.OptimizedContent,
		Score:            appResp.Score,
		Details:          appResp.Details,
	}, nil
}

// AdaptContentForModel 为特定模型适配内容
func (h *ContentKitexHandler) AdaptContentForModel(ctx context.Context, req *inbound.AdaptContentForModelRequest) (*inbound.AdaptContentForModelResponse, error) {
	appReq := &service.AdaptContentForModelRequest{
		ContentID:   req.ContentID,
		TargetModel: req.TargetModel,
	}

	appResp, err := h.appService.AdaptContentForModel(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.AdaptContentForModelResponse{
		Success:           appResp.Success,
		AdaptedContent:    appResp.AdaptedContent,
		ModelSpecificData: appResp.ModelSpecificData,
	}, nil
}

// CheckContentCompliance 检查内容合规性
func (h *ContentKitexHandler) CheckContentCompliance(ctx context.Context, req *inbound.CheckContentComplianceRequest) (*inbound.CheckContentComplianceResponse, error) {
	appReq := &service.CheckContentComplianceRequest{
		ContentID: req.ContentID,
	}

	appResp, err := h.appService.CheckContentCompliance(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.CheckContentComplianceResponse{
		Compliant: appResp.Compliant,
		Issues:    convertIssues(appResp.Issues),
		Report:    appResp.Report,
	}, nil
}

// CreateContentTemplate 创建内容模板
func (h *ContentKitexHandler) CreateContentTemplate(ctx context.Context, req *inbound.CreateContentTemplateRequest) (*inbound.CreateContentTemplateResponse, error) {
	appReq := &service.CreateContentTemplateRequest{
		UserID:       req.UserID,
		Name:         req.Name,
		Description:  req.Description,
		TemplateType: req.TemplateType,
		TemplateData: req.TemplateData,
		IsPublic:     req.IsPublic,
	}

	appResp, err := h.appService.CreateContentTemplate(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.CreateContentTemplateResponse{
		Template: appResp.Template,
	}, nil
}

// GetContentTemplate 获取内容模板
func (h *ContentKitexHandler) GetContentTemplate(ctx context.Context, req *inbound.GetContentTemplateRequest) (*inbound.GetContentTemplateResponse, error) {
	appReq := &service.GetContentTemplateRequest{
		TemplateID: req.TemplateID,
	}

	appResp, err := h.appService.GetContentTemplate(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.GetContentTemplateResponse{
		Template: appResp.Template,
	}, nil
}

// ListContentTemplates 列出内容模板
func (h *ContentKitexHandler) ListContentTemplates(ctx context.Context, req *inbound.ListContentTemplatesRequest) (*inbound.ListContentTemplatesResponse, error) {
	appReq := &service.ListContentTemplatesRequest{
		UserID:       req.UserID,
		TemplateType: req.TemplateType,
		IsPublic:     req.IsPublic,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}

	appResp, err := h.appService.ListContentTemplates(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.ListContentTemplatesResponse{
		Templates: appResp.Templates,
		Total:     appResp.Total,
	}, nil
}

// CreateKnowledgeEntity 创建知识实体
func (h *ContentKitexHandler) CreateKnowledgeEntity(ctx context.Context, req *inbound.CreateKnowledgeEntityRequest) (*inbound.CreateKnowledgeEntityResponse, error) {
	appReq := &service.CreateKnowledgeEntityRequest{
		UserID:         req.UserID,
		EntityName:     req.EntityName,
		EntityType:     req.EntityType,
		EntityData:     req.EntityData,
		AuthorityLinks: req.AuthorityLinks,
	}

	appResp, err := h.appService.CreateKnowledgeEntity(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.CreateKnowledgeEntityResponse{
		Entity: appResp.Entity,
	}, nil
}

// GetKnowledgeEntity 获取知识实体
func (h *ContentKitexHandler) GetKnowledgeEntity(ctx context.Context, req *inbound.GetKnowledgeEntityRequest) (*inbound.GetKnowledgeEntityResponse, error) {
	appReq := &service.GetKnowledgeEntityRequest{
		EntityID: req.EntityID,
	}

	appResp, err := h.appService.GetKnowledgeEntity(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.GetKnowledgeEntityResponse{
		Entity: appResp.Entity,
	}, nil
}

// ListKnowledgeEntities 列出知识实体
func (h *ContentKitexHandler) ListKnowledgeEntities(ctx context.Context, req *inbound.ListKnowledgeEntitiesRequest) (*inbound.ListKnowledgeEntitiesResponse, error) {
	appReq := &service.ListKnowledgeEntitiesRequest{
		UserID:     req.UserID,
		EntityType: req.EntityType,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	appResp, err := h.appService.ListKnowledgeEntities(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.ListKnowledgeEntitiesResponse{
		Entities: appResp.Entities,
		Total:    appResp.Total,
	}, nil
}

// LinkEntityToContent 关联实体到内容
func (h *ContentKitexHandler) LinkEntityToContent(ctx context.Context, req *inbound.LinkEntityToContentRequest) (*inbound.LinkEntityToContentResponse, error) {
	appReq := &service.LinkEntityToContentRequest{
		ContentID: req.ContentID,
		EntityID:  req.EntityID,
	}

	appResp, err := h.appService.LinkEntityToContent(ctx, appReq)
	if err != nil {
		return nil, err
	}

	return &inbound.LinkEntityToContentResponse{
		Success: appResp.Success,
	}, nil
}

// 辅助函数
func convertIssues(issues []ai.ComplianceIssue) []inbound.ComplianceIssue {
	result := make([]inbound.ComplianceIssue, len(issues))
	for i, issue := range issues {
		result[i] = inbound.ComplianceIssue{
			IssueType:   issue.IssueType,
			Description: issue.Description,
			Severity:    issue.Severity,
			Suggestion:  issue.Suggestion,
		}
	}
	return result
}