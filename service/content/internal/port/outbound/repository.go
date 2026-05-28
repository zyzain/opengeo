package outbound

import (
	"context"

	"opengeo/service/content/internal/domain/model"
)

// ContentRepository 内容仓储接口（出站端口）
type ContentRepository interface {
	// 内容CRUD
	Create(ctx context.Context, content *model.Content) error
	GetByID(ctx context.Context, id int64) (*model.Content, error)
	Update(ctx context.Context, content *model.Content) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter *ContentFilter) ([]*model.Content, int32, error)

	// 内容版本
	CreateVersion(ctx context.Context, version *model.ContentVersion) error
	GetVersionByID(ctx context.Context, id int64) (*model.ContentVersion, error)
	ListVersions(ctx context.Context, contentID int64, page, pageSize int32) ([]*model.ContentVersion, int32, error)
}

// ContentFilter 内容过滤器
type ContentFilter struct {
	UserID      int64
	Status      int32
	ContentType string
	Page        int32
	PageSize    int32
}

// ContentTemplateRepository 内容模板仓储接口（出站端口）
type ContentTemplateRepository interface {
	Create(ctx context.Context, template *model.ContentTemplate) error
	GetByID(ctx context.Context, id int64) (*model.ContentTemplate, error)
	List(ctx context.Context, filter *TemplateFilter) ([]*model.ContentTemplate, int32, error)
	Update(ctx context.Context, template *model.ContentTemplate) error
	Delete(ctx context.Context, id int64) error
}

// TemplateFilter 模板过滤器
type TemplateFilter struct {
	UserID       int64
	TemplateType string
	IsPublic     bool
	Page         int32
	PageSize     int32
}

// KnowledgeEntityRepository 知识实体仓储接口（出站端口）
type KnowledgeEntityRepository interface {
	Create(ctx context.Context, entity *model.KnowledgeEntity) error
	GetByID(ctx context.Context, id int64) (*model.KnowledgeEntity, error)
	List(ctx context.Context, filter *EntityFilter) ([]*model.KnowledgeEntity, int32, error)
	Update(ctx context.Context, entity *model.KnowledgeEntity) error
	Delete(ctx context.Context, id int64) error
	LinkToContent(ctx context.Context, contentID, entityID int64) error
	GetContentEntities(ctx context.Context, contentID int64) ([]*model.KnowledgeEntity, error)
}

// EntityFilter 实体过滤器
type EntityFilter struct {
	UserID     int64
	EntityType string
	Page       int32
	PageSize   int32
}

// CacheService 缓存服务接口（出站端口）
type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration int64) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// EventBus 事件总线接口（出站端口）
type EventBus interface {
	Publish(ctx context.Context, topic string, event interface{}) error
	Subscribe(ctx context.Context, topic string, handler EventHandler) error
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event interface{}) error