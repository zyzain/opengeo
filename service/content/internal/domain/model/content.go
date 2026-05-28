package model

import "time"

// Content 内容实体（领域模型）
type Content struct {
	ID                  int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID            int64     `json:"tenant_id" gorm:"not null;index:idx_tenant"`
	UserID              int64     `json:"user_id" gorm:"not null;index:idx_content_list"`
	BrandID             int64     `json:"brand_id" gorm:"index:idx_brand"`
	Title               string    `json:"title" gorm:"size:256;not null"`
	Body                string    `json:"body" gorm:"type:text;not null"`
	Summary             string    `json:"summary" gorm:"size:512"`
	ContentType         string    `json:"content_type" gorm:"size:32;default:article;index:idx_content_list"`
	Status              int32     `json:"status" gorm:"default:0;index:idx_content_list"`
	Visibility          string    `json:"visibility" gorm:"size:32;default:private"`
	SchemaMarkup        string    `json:"schema_markup" gorm:"type:text"`
	AIOptimizationScore float32   `json:"ai_optimization_score"`
	WordCount           int32     `json:"word_count"`
	ReadingTime         int32     `json:"reading_time"`
	Tags                string    `json:"tags" gorm:"type:json"`
	PublishedAt         *time.Time `json:"published_at"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ContentStatus 内容状态常量
const (
	ContentStatusDraft     = 0 // 草稿
	ContentStatusPublished = 1 // 已发布
	ContentStatusArchived  = 2 // 已归档
)

// ContentType 内容类型常量
const (
	ContentTypeArticle = "article"
	ContentTypeVideo   = "video"
	ContentTypeImage   = "image"
)

// IsValid 检查内容是否有效
func (c *Content) IsValid() bool {
	return c.Title != "" && c.Body != "" && c.UserID > 0
}

// IsDraft 检查是否为草稿
func (c *Content) IsDraft() bool {
	return c.Status == ContentStatusDraft
}

// IsPublished 检查是否已发布
func (c *Content) IsPublished() bool {
	return c.Status == ContentStatusPublished
}

// Publish 发布内容
func (c *Content) Publish() {
	c.Status = ContentStatusPublished
	c.UpdatedAt = time.Now()
}

// Archive 归档内容
func (c *Content) Archive() {
	c.Status = ContentStatusArchived
	c.UpdatedAt = time.Now()
}

// UpdateAIScore 更新AI优化分数
func (c *Content) UpdateAIScore(score float32) {
	c.AIOptimizationScore = score
	c.UpdatedAt = time.Now()
}

// ContentVersion 内容版本
type ContentVersion struct {
	ID                int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ContentID         int64     `json:"content_id" gorm:"index;not null"`
	Version           int32     `json:"version" gorm:"not null"`
	Title             string    `json:"title" gorm:"size:256;not null"`
	Body              string    `json:"body" gorm:"type:text;not null"`
	SchemaMarkup      string    `json:"schema_markup" gorm:"type:text"`
	AIModelAdaptation string    `json:"ai_model_adaptation" gorm:"size:64"`
	CreatedAt         time.Time `json:"created_at"`
}

// ContentTemplate 内容模板
type ContentTemplate struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID     int64     `json:"tenant_id" gorm:"not null;index:idx_tenant"`
	UserID       int64     `json:"user_id" gorm:"index"`
	BrandID      int64     `json:"brand_id" gorm:"index:idx_brand"`
	Name         string    `json:"name" gorm:"size:128;not null"`
	Description  string    `json:"description" gorm:"size:256"`
	TemplateType string    `json:"template_type" gorm:"size:32"`
	TemplateData string    `json:"template_data" gorm:"type:text;not null"`
	Variables    string    `json:"variables" gorm:"type:json"`
	IsPublic     bool      `json:"is_public" gorm:"default:false"`
	UsageCount   int32     `json:"usage_count" gorm:"default:0"`
	Rating       float32   `json:"rating" gorm:"default:0"`
	RatingCount  int32     `json:"rating_count" gorm:"default:0"`
	Tags         string    `json:"tags" gorm:"size:512"`
	Author       string    `json:"author" gorm:"size:128"`
	IsOfficial   bool      `json:"is_official" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TemplateExportData 模板导出数据结构
type TemplateExportData struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	TemplateType string `json:"template_type"`
	TemplateData string `json:"template_data"`
	Tags         string `json:"tags"`
	Author       string `json:"author"`
	IsOfficial   bool   `json:"is_official"`
}

// KnowledgeEntity 知识图谱实体
type KnowledgeEntity struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID       int64     `json:"tenant_id" gorm:"not null;index:idx_tenant"`
	BrandID        int64     `json:"brand_id" gorm:"index:idx_brand"`
	UserID         int64     `json:"user_id" gorm:"index;not null"`
	EntityName     string    `json:"entity_name" gorm:"size:128;not null"`
	EntityType     string    `json:"entity_type" gorm:"size:32"`
	EntityData     string    `json:"entity_data" gorm:"type:text"`
	AuthorityLinks string    `json:"authority_links" gorm:"type:text"`
	EmbeddingID    string    `json:"embedding_id" gorm:"size:128"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ContentEntity 内容实体关联
type ContentEntity struct {
	ContentID int64 `json:"content_id" gorm:"primaryKey"`
	EntityID  int64 `json:"entity_id" gorm:"primaryKey"`
}