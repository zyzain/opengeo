package knowledge

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AutoMigrate() error {
	if err := s.db.AutoMigrate(&Entity{}); err != nil {
		if !strings.Contains(err.Error(), "Duplicate key name") {
			return err
		}
	}
	return nil
}

func (s *Store) Seed(ctx context.Context, adminUserID int64) error {
	var count int64
	s.db.WithContext(ctx).Model(&Entity{}).Count(&count)
	if count > 0 {
		return nil
	}

	now := time.Now()
	entities := []Entity{
		{UserID: adminUserID, EntityName: "OpenGEO", EntityType: "brand", EntityData: `{"description":"智能发布平台","website":"https://opengeo.com"}`, AuthorityLinks: `["https://opengeo.com/about"]`, ContentCount: 15, CreatedAt: now, UpdatedAt: now},
		{UserID: adminUserID, EntityName: "DeepSeek", EntityType: "product", EntityData: `{"description":"AI大模型","company":"深度求索"}`, AuthorityLinks: `["https://deepseek.com"]`, ContentCount: 12, CreatedAt: now, UpdatedAt: now},
		{UserID: adminUserID, EntityName: "GEO优化", EntityType: "concept", EntityData: `{"description":"生成式引擎优化"}`, AuthorityLinks: `["https://arxiv.org/abs/geo-optimization"]`, ContentCount: 8, CreatedAt: now, UpdatedAt: now},
		{UserID: adminUserID, EntityName: "张三", EntityType: "person", EntityData: `{"description":"技术专家","title":"首席架构师"}`, AuthorityLinks: `[]`, ContentCount: 3, CreatedAt: now, UpdatedAt: now},
		{UserID: adminUserID, EntityName: "北京市", EntityType: "place", EntityData: `{"description":"中国首都"}`, AuthorityLinks: `["https://en.wikipedia.org/wiki/Beijing"]`, ContentCount: 5, CreatedAt: now, UpdatedAt: now},
	}
	return s.db.WithContext(ctx).Create(&entities).Error
}

func (s *Store) Create(ctx context.Context, e *Entity) error {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	return s.db.WithContext(ctx).Create(e).Error
}

func (s *Store) GetByID(ctx context.Context, id int64) (*Entity, error) {
	var e Entity
	err := s.db.WithContext(ctx).First(&e, id).Error
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *Store) Update(ctx context.Context, e *Entity) error {
	e.UpdatedAt = time.Now()
	return s.db.WithContext(ctx).Model(e).Select("entity_name", "entity_type", "entity_data", "authority_links", "updated_at").Updates(e).Error
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Delete(&Entity{}, id).Error
}

func (s *Store) List(ctx context.Context, userID int64, entityType string, page, pageSize int) ([]Entity, int64, error) {
	query := s.db.WithContext(ctx).Model(&Entity{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var items []Entity
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *Store) SearchByName(ctx context.Context, userID int64, keyword string, limit int) ([]Entity, error) {
	query := s.db.WithContext(ctx).Model(&Entity{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if keyword != "" {
		query = query.Where("entity_name LIKE ?", "%"+keyword+"%")
	}
	if limit <= 0 {
		limit = 20
	}
	var items []Entity
	if err := query.Limit(limit).Order("content_count DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
