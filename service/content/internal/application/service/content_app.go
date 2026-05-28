package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"opengeo/pkg/ai"
	"opengeo/service/content/internal/domain/model"
	"opengeo/service/content/internal/domain/service"
	"opengeo/service/content/internal/port/outbound"
)

// ContentAppService 内容应用服务（应用层）
// 编排领域逻辑，协调领域服务和仓储
type ContentAppService struct {
	contentRepo      outbound.ContentRepository
	templateRepo     outbound.ContentTemplateRepository
	entityRepo       outbound.KnowledgeEntityRepository
	aiService        ai.AIService
	cacheService     outbound.CacheService
	geoOptimService  *service.GEOOptimizationService
}

// NewContentAppService 创建内容应用服务
func NewContentAppService(
	contentRepo outbound.ContentRepository,
	templateRepo outbound.ContentTemplateRepository,
	entityRepo outbound.KnowledgeEntityRepository,
	aiService ai.AIService,
	cacheService outbound.CacheService,
) *ContentAppService {
	return &ContentAppService{
		contentRepo:     contentRepo,
		templateRepo:    templateRepo,
		entityRepo:      entityRepo,
		aiService:       aiService,
		cacheService:    cacheService,
		geoOptimService: service.NewGEOOptimizationService(),
	}
}

// CreateContent 创建内容
func (s *ContentAppService) CreateContent(ctx context.Context, req *CreateContentRequest) (*CreateContentResponse, error) {
	// 创建领域模型
	content := &model.Content{
		UserID:      req.UserID,
		Title:       req.Title,
		Body:        req.Body,
		ContentType: req.ContentType,
		Status:      model.ContentStatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 验证内容
	if !content.IsValid() {
		return nil, fmt.Errorf("invalid content")
	}

	// 持久化
	if err := s.contentRepo.Create(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}

	return &CreateContentResponse{Content: content}, nil
}

// GetContent 获取内容
func (s *ContentAppService) GetContent(ctx context.Context, req *GetContentRequest) (*GetContentResponse, error) {
	// 先尝试从缓存获取
	cacheKey := fmt.Sprintf("content:%d", req.ContentID)
	if cached, err := s.cacheService.Get(ctx, cacheKey); err == nil && cached != "" {
		// TODO: 反序列化缓存内容
	}

	// 从仓储获取
	content, err := s.contentRepo.GetByID(ctx, req.ContentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// 写入缓存
	// TODO: 序列化内容并缓存

	return &GetContentResponse{Content: content}, nil
}

// UpdateContent 更新内容
func (s *ContentAppService) UpdateContent(ctx context.Context, req *UpdateContentRequest) (*UpdateContentResponse, error) {
	// 获取现有内容
	content, err := s.contentRepo.GetByID(ctx, req.ContentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// 更新字段
	if req.Title != "" {
		content.Title = req.Title
	}
	if req.Body != "" {
		content.Body = req.Body
	}
	if req.SchemaMarkup != "" {
		content.SchemaMarkup = req.SchemaMarkup
	}
	content.UpdatedAt = time.Now()

	// 持久化
	if err := s.contentRepo.Update(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("content:%d", req.ContentID)
	s.cacheService.Delete(ctx, cacheKey)

	return &UpdateContentResponse{Content: content}, nil
}

// DeleteContent 删除内容
func (s *ContentAppService) DeleteContent(ctx context.Context, req *DeleteContentRequest) (*DeleteContentResponse, error) {
	if err := s.contentRepo.Delete(ctx, req.ContentID); err != nil {
		return nil, fmt.Errorf("failed to delete content: %w", err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("content:%d", req.ContentID)
	s.cacheService.Delete(ctx, cacheKey)

	return &DeleteContentResponse{Success: true}, nil
}

// ListContents 列出内容
func (s *ContentAppService) ListContents(ctx context.Context, req *ListContentsRequest) (*ListContentsResponse, error) {
	filter := &outbound.ContentFilter{
		UserID:      req.UserID,
		Status:      req.Status,
		ContentType: req.ContentType,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}

	contents, total, err := s.contentRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list contents: %w", err)
	}

	return &ListContentsResponse{
		Contents: contents,
		Total:    total,
	}, nil
}

// OptimizeContentForAI 为AI优化内容
func (s *ContentAppService) OptimizeContentForAI(ctx context.Context, req *OptimizeContentForAIRequest) (*OptimizeContentForAIResponse, error) {
	// 获取内容
	content, err := s.contentRepo.GetByID(ctx, req.ContentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// 调用AI服务优化
	aiResult, err := s.aiService.OptimizeContent(ctx, &ai.OptimizeRequest{
		ContentID:        content.ID,
		Title:            content.Title,
		Body:             content.Body,
		ContentType:      content.ContentType,
		OptimizationType: req.OptimizationType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to optimize content: %w", err)
	}

	// 更新内容AI分数
	content.UpdateAIScore(aiResult.Score)
	if err := s.contentRepo.Update(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	return &OptimizeContentForAIResponse{
		Success:          aiResult.Success,
		OptimizedContent: aiResult.OptimizedBody,
		Score:            aiResult.Score,
		Details:          aiResult.SchemaMarkup,
	}, nil
}

// AdaptContentForModel 为特定模型适配内容
func (s *ContentAppService) AdaptContentForModel(ctx context.Context, req *AdaptContentForModelRequest) (*AdaptContentForModelResponse, error) {
	// 获取内容
	content, err := s.contentRepo.GetByID(ctx, req.ContentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// 调用AI服务适配
	aiResult, err := s.aiService.AdaptForModel(ctx, &ai.AdaptRequest{
		ContentID:   content.ID,
		Title:       content.Title,
		Body:        content.Body,
		TargetModel: req.TargetModel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to adapt content: %w", err)
	}

	return &AdaptContentForModelResponse{
		Success:           aiResult.Success,
		AdaptedContent:    aiResult.AdaptedBody,
		ModelSpecificData: fmt.Sprintf("%v", aiResult.ModelSpecificData),
	}, nil
}

// CheckContentCompliance 检查内容合规性
func (s *ContentAppService) CheckContentCompliance(ctx context.Context, req *CheckContentComplianceRequest) (*CheckContentComplianceResponse, error) {
	// 获取内容
	content, err := s.contentRepo.GetByID(ctx, req.ContentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// 调用AI服务检查合规
	aiResult, err := s.aiService.CheckCompliance(ctx, &ai.ComplianceRequest{
		ContentID: content.ID,
		Title:     content.Title,
		Body:      content.Body,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check compliance: %w", err)
	}

	return &CheckContentComplianceResponse{
		Compliant: aiResult.Compliant,
		Issues:    aiResult.Issues,
		Report:    aiResult.Report,
	}, nil
}

// CreateContentVersion 创建内容版本
func (s *ContentAppService) CreateContentVersion(ctx context.Context, req *CreateContentVersionRequest) (*CreateContentVersionResponse, error) {
	// 获取内容
	content, err := s.contentRepo.GetByID(ctx, req.ContentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// 获取当前版本数量
	_, total, err := s.contentRepo.ListVersions(ctx, req.ContentID, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	// 创建新版本
	version := &model.ContentVersion{
		ContentID:         req.ContentID,
		Version:           total + 1,
		Title:             req.Title,
		Body:              req.Body,
		SchemaMarkup:      req.SchemaMarkup,
		AIModelAdaptation: req.AIModelAdaptation,
		CreatedAt:         time.Now(),
	}

	if err := s.contentRepo.CreateVersion(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// 更新内容
	content.Title = req.Title
	content.Body = req.Body
	content.SchemaMarkup = req.SchemaMarkup
	content.UpdatedAt = time.Now()

	if err := s.contentRepo.Update(ctx, content); err != nil {
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	return &CreateContentVersionResponse{Version: version}, nil
}

// GetContentVersion 获取内容版本
func (s *ContentAppService) GetContentVersion(ctx context.Context, req *GetContentVersionRequest) (*GetContentVersionResponse, error) {
	version, err := s.contentRepo.GetVersionByID(ctx, req.VersionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	return &GetContentVersionResponse{Version: version}, nil
}

// ListContentVersions 列出内容版本
func (s *ContentAppService) ListContentVersions(ctx context.Context, req *ListContentVersionsRequest) (*ListContentVersionsResponse, error) {
	versions, total, err := s.contentRepo.ListVersions(ctx, req.ContentID, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	return &ListContentVersionsResponse{
		Versions: versions,
		Total:    total,
	}, nil
}

// CreateContentTemplate 创建内容模板
func (s *ContentAppService) CreateContentTemplate(ctx context.Context, req *CreateContentTemplateRequest) (*CreateContentTemplateResponse, error) {
	template := &model.ContentTemplate{
		UserID:       req.UserID,
		Name:         req.Name,
		Description:  req.Description,
		TemplateType: req.TemplateType,
		TemplateData: req.TemplateData,
		IsPublic:     req.IsPublic,
		UsageCount:   0,
		Rating:       0,
		CreatedAt:    time.Now(),
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return &CreateContentTemplateResponse{Template: template}, nil
}

// GetContentTemplate 获取内容模板
func (s *ContentAppService) GetContentTemplate(ctx context.Context, req *GetContentTemplateRequest) (*GetContentTemplateResponse, error) {
	template, err := s.templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return &GetContentTemplateResponse{Template: template}, nil
}

// ListContentTemplates 列出内容模板
func (s *ContentAppService) ListContentTemplates(ctx context.Context, req *ListContentTemplatesRequest) (*ListContentTemplatesResponse, error) {
	filter := &outbound.TemplateFilter{
		UserID:       req.UserID,
		TemplateType: req.TemplateType,
		IsPublic:     req.IsPublic,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}

	templates, total, err := s.templateRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	return &ListContentTemplatesResponse{
		Templates: templates,
		Total:     total,
	}, nil
}

// CreateKnowledgeEntity 创建知识实体
func (s *ContentAppService) CreateKnowledgeEntity(ctx context.Context, req *CreateKnowledgeEntityRequest) (*CreateKnowledgeEntityResponse, error) {
	entity := &model.KnowledgeEntity{
		UserID:         req.UserID,
		EntityName:     req.EntityName,
		EntityType:     req.EntityType,
		EntityData:     req.EntityData,
		AuthorityLinks: req.AuthorityLinks,
		CreatedAt:      time.Now(),
	}

	if err := s.entityRepo.Create(ctx, entity); err != nil {
		return nil, fmt.Errorf("failed to create entity: %w", err)
	}

	return &CreateKnowledgeEntityResponse{Entity: entity}, nil
}

// GetKnowledgeEntity 获取知识实体
func (s *ContentAppService) GetKnowledgeEntity(ctx context.Context, req *GetKnowledgeEntityRequest) (*GetKnowledgeEntityResponse, error) {
	entity, err := s.entityRepo.GetByID(ctx, req.EntityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity: %w", err)
	}

	return &GetKnowledgeEntityResponse{Entity: entity}, nil
}

// ListKnowledgeEntities 列出知识实体
func (s *ContentAppService) ListKnowledgeEntities(ctx context.Context, req *ListKnowledgeEntitiesRequest) (*ListKnowledgeEntitiesResponse, error) {
	filter := &outbound.EntityFilter{
		UserID:     req.UserID,
		EntityType: req.EntityType,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	entities, total, err := s.entityRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list entities: %w", err)
	}

	return &ListKnowledgeEntitiesResponse{
		Entities: entities,
		Total:    total,
	}, nil
}

// LinkEntityToContent 关联实体到内容
func (s *ContentAppService) LinkEntityToContent(ctx context.Context, req *LinkEntityToContentRequest) (*LinkEntityToContentResponse, error) {
	if err := s.entityRepo.LinkToContent(ctx, req.ContentID, req.EntityID); err != nil {
		return nil, fmt.Errorf("failed to link entity: %w", err)
	}

	return &LinkEntityToContentResponse{Success: true}, nil
}

// GetContentEntities 获取内容关联的实体
func (s *ContentAppService) GetContentEntities(ctx context.Context, contentID int64) ([]*model.KnowledgeEntity, error) {
	entities, err := s.entityRepo.GetContentEntities(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content entities: %w", err)
	}

	return entities, nil
}

// SeedBuiltinTemplates 内置10个官方GEO优化模板
func (s *ContentAppService) SeedBuiltinTemplates(ctx context.Context, userID int64) error {
	templates := []struct {
		name         string
		description  string
		templateType string
		templateData string
		tags         string
	}{
		{
			name:         "GEO文章优化模板",
			description:  "适用于长文章的GEO优化，包含结构化标题、权威引用、数据支撑等要素",
			templateType: "article",
			templateData: `## {{title}}

### 概述
{{introduction}}

### 核心要点
{{key_points}}

### 详细分析
{{analysis}}

### 数据支撑
{{data_evidence}}

### 权威引用
{{authority_references}}

### 总结
{{conclusion}}

> 本文由AI辅助生成（AIGC），仅供参考。`,
			tags: "GEO,文章,优化,长文",
		},
		{
			name:         "FAQ问答模板",
			description:  "适用于FAQ类内容，问答结构便于AI搜索引擎提取",
			templateType: "faq",
			templateData: `## {{topic}} 常见问题

### Q1: {{question_1}}
**A:** {{answer_1}}

### Q2: {{question_2}}
**A:** {{answer_2}}

### Q3: {{question_3}}
**A:** {{answer_3}}

---

**参考来源：**
- {{source_1}}
- {{source_2}}`,
			tags: "FAQ,问答,结构化",
		},
		{
			name:         "产品评测模板",
			description:  "适用于产品评测类内容，包含评分、对比、优缺点分析",
			templateType: "review",
			templateData: `## {{product_name}} 评测报告

### 产品概述
- **品牌**: {{brand}}
- **型号**: {{model}}
- **价格**: {{price}}
- **评分**: {{rating}}/5

### 核心优势
{{advantages}}

### 不足之处
{{disadvantages}}

### 使用体验
{{user_experience}}

### 竞品对比
{{comparison}}

### 购买建议
{{recommendation}}`,
			tags: "评测,产品,对比,评分",
		},
		{
			name:         "行业分析模板",
			description:  "适用于行业趋势分析类内容，数据驱动，适合AI引用",
			templateType: "analysis",
			templateData: `## {{industry}} 行业分析报告（{{year}}）

### 行业概况
{{overview}}

### 市场规模
- 总规模：{{market_size}}
- 增长率：{{growth_rate}}
- 预测趋势：{{trend}}

### 主要玩家
{{major_players}}

### 技术趋势
{{tech_trends}}

### 机遇与挑战
**机遇：**
{{opportunities}}

**挑战：**
{{challenges}}

### 未来展望
{{outlook}}

> 数据来源：{{data_sources}}`,
			tags: "行业,分析,趋势,数据",
		},
		{
			name:         "教程指南模板",
			description:  "适用于教程类内容，步骤清晰，便于AI提取操作步骤",
			templateType: "tutorial",
			templateData: `## {{title}} - 完整教程

### 前置条件
{{prerequisites}}

### 步骤详解

#### 步骤1: {{step_1_title}}
{{step_1_detail}}

#### 步骤2: {{step_2_title}}
{{step_2_detail}}

#### 步骤3: {{step_3_title}}
{{step_3_detail}}

### 常见问题
{{faq}}

### 总结
{{summary}}`,
			tags: "教程,指南,步骤,操作",
		},
		{
			name:         "新闻资讯模板",
			description:  "适用于新闻类内容，5W1H结构，便于AI提取关键信息",
			templateType: "news",
			templateData: `## {{headline}}

**发布时间**: {{publish_time}}
**来源**: {{source}}

### 事件概要
{{summary}}

### 详细报道
{{details}}

### 关键信息
- **时间**: {{when}}
- **地点**: {{where}}
- **涉及方**: {{who}}
- **事件**: {{what}}
- **原因**: {{why}}
- **影响**: {{impact}}

### 后续发展
{{follow_up}}`,
			tags: "新闻,资讯,5W1H",
		},
		{
			name:         "对比分析模板",
			description:  "适用于A vs B对比类内容，结构化对比便于AI理解",
			templateType: "comparison",
			templateData: `## {{item_a}} vs {{item_b}} 全面对比

### 基本信息对比

| 对比项 | {{item_a}} | {{item_b}} |
|--------|-----------|-----------|
| {{feature_1}} | {{a_value_1}} | {{b_value_1}} |
| {{feature_2}} | {{a_value_2}} | {{b_value_2}} |
| {{feature_3}} | {{a_value_3}} | {{b_value_3}} |

### 优势分析
**{{item_a}}优势**: {{a_advantages}}
**{{item_b}}优势**: {{b_advantages}}

### 适用场景
- 选择{{item_a}}：{{a_scenarios}}
- 选择{{item_b}}：{{b_scenarios}}

### 最终建议
{{recommendation}}`,
			tags: "对比,分析,选择,评测",
		},
		{
			name:         "数据报告模板",
			description:  "适用于数据驱动的报告类内容，数字密度高，AI友好",
			templateType: "report",
			templateData: `## {{report_title}}（{{period}}）

### 关键指标
- **{{kpi_1_name}}**: {{kpi_1_value}}
- **{{kpi_2_name}}**: {{kpi_2_value}}
- **{{kpi_3_name}}**: {{kpi_3_value}}

### 数据分析
{{analysis}}

### 趋势变化
{{trend_analysis}}

### 核心发现
1. {{finding_1}}
2. {{finding_2}}
3. {{finding_3}}

### 行动建议
{{recommendations}}

> 数据截止时间：{{data_cutoff}}`,
			tags: "报告,数据,指标,分析",
		},
		{
			name:         "GEO优化清单模板",
			description:  "适用于GEO优化检查清单类内容，便于AI提取要点",
			templateType: "checklist",
			templateData: `## {{title}} GEO优化清单

### 内容结构
- [ ] 使用清晰的标题层级（H2/H3）
- [ ] 段落长度控制在50-200字
- [ ] 使用列表/要点提炼关键信息
- [ ] 添加总结/结论段落

### 数据支撑
- [ ] 包含具体数字和百分比
- [ ] 引用权威数据来源
- [ ] 添加日期和时间标记

### AI友好性
- [ ] 包含问答式段落
- [ ] 使用Schema Markup结构化数据
- [ ] 添加权威信源引用链接
- [ ] 包含AIGC标识

### 优化检查
- [ ] 标题关键词在正文中多次出现
- [ ] 内容长度≥1000字
- [ ] 包含列表/表格等结构化元素`,
			tags: "清单,GEO,优化,检查",
		},
		{
			name:         "竞品分析模板",
			description:  "适用于竞品分析类内容，便于AI提取对比信息",
			templateType: "competitor",
			templateData: `## {{company}} 竞品分析报告

### 竞品概览
{{competitor_overview}}

### 产品对比

| 维度 | {{company}} | 竞品A | 竞品B |
|------|-----------|-------|-------|
| {{dimension_1}} | {{self_1}} | {{a_1}} | {{b_1}} |
| {{dimension_2}} | {{self_2}} | {{a_2}} | {{b_2}} |

### SWOT分析
**优势**: {{strengths}}
**劣势**: {{weaknesses}}
**机会**: {{opportunities}}
**威胁**: {{threats}}

### 策略建议
{{strategy_recommendations}}

### 市场定位
{{market_positioning}}`,
			tags: "竞品,分析,市场,SWOT",
		},
	}

	for _, tmpl := range templates {
		existing, _ := s.ListContentTemplates(ctx, &ListContentTemplatesRequest{
			TemplateType: tmpl.templateType,
			IsPublic:     true,
			Page:         1,
			PageSize:     1,
		})
		skip := false
		if existing != nil {
			for _, e := range existing.Templates {
				if e.Name == tmpl.name && e.IsOfficial {
					skip = true
					break
				}
			}
		}
		if skip {
			continue
		}

		template := &model.ContentTemplate{
			UserID:       0,
			Name:         tmpl.name,
			Description:  tmpl.description,
			TemplateType: tmpl.templateType,
			TemplateData: tmpl.templateData,
			Tags:         tmpl.tags,
			Author:       "OpenGEO Official",
			IsPublic:     true,
			IsOfficial:   true,
			UsageCount:   0,
			Rating:       4.5,
			RatingCount:  10,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.templateRepo.Create(ctx, template); err != nil {
			return fmt.Errorf("failed to seed template %s: %w", tmpl.name, err)
		}
	}

	return nil
}

// RateTemplate 评分模板
func (s *ContentAppService) RateTemplate(ctx context.Context, templateID int64, rating float32) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	totalScore := template.Rating * float32(template.RatingCount)
	template.RatingCount++
	template.Rating = (totalScore + rating) / float32(template.RatingCount)
	template.UpdatedAt = time.Now()

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return fmt.Errorf("failed to update template rating: %w", err)
	}

	return nil
}

// IncrementUsageCount 增加模板使用次数
func (s *ContentAppService) IncrementUsageCount(ctx context.Context, templateID int64) error {
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	template.UsageCount++
	template.UpdatedAt = time.Now()

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return fmt.Errorf("failed to update usage count: %w", err)
	}

	return nil
}

// ExportTemplate 导出模板为JSON
func (s *ContentAppService) ExportTemplate(ctx context.Context, templateID int64) (string, error) {
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return "", fmt.Errorf("failed to get template: %w", err)
	}

	exportData := model.TemplateExportData{
		Name:         template.Name,
		Description:  template.Description,
		TemplateType: template.TemplateType,
		TemplateData: template.TemplateData,
		Tags:         template.Tags,
		Author:       template.Author,
		IsOfficial:   template.IsOfficial,
	}

	data, err := json.Marshal(exportData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %w", err)
	}

	return string(data), nil
}

// ImportTemplate 从JSON导入模板
func (s *ContentAppService) ImportTemplate(ctx context.Context, userID int64, jsonData string) (*model.ContentTemplate, error) {
	var exportData model.TemplateExportData
	if err := json.Unmarshal([]byte(jsonData), &exportData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template: %w", err)
	}

	if exportData.Name == "" || exportData.TemplateData == "" {
		return nil, fmt.Errorf("template name and data are required")
	}

	template := &model.ContentTemplate{
		UserID:       userID,
		Name:         exportData.Name,
		Description:  exportData.Description,
		TemplateType: exportData.TemplateType,
		TemplateData: exportData.TemplateData,
		Tags:         exportData.Tags,
		Author:       exportData.Author,
		IsPublic:     false,
		UsageCount:   0,
		Rating:       0,
		RatingCount:  0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return template, nil
}

// 辅助函数

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
	Success           bool   `json:"success"`
	AdaptedContent    string `json:"adapted_content"`
	ModelSpecificData string `json:"model_specific_data"`
}

type CheckContentComplianceRequest struct {
	ContentID int64 `json:"content_id"`
}

type CheckContentComplianceResponse struct {
	Compliant bool                 `json:"compliant"`
	Issues    []ai.ComplianceIssue `json:"issues"`
	Report    string               `json:"report"`
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