package template

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
	if err := s.db.AutoMigrate(&Template{}); err != nil {
		if !strings.Contains(err.Error(), "Duplicate key name") {
			return err
		}
	}
	return nil
}

func (s *Store) Create(ctx context.Context, t *Template) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return s.db.WithContext(ctx).Create(t).Error
}

func (s *Store) GetByID(ctx context.Context, id int64) (*Template, error) {
	var t Template
	err := s.db.WithContext(ctx).First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Store) Update(ctx context.Context, t *Template) error {
	t.UpdatedAt = time.Now()
	return s.db.WithContext(ctx).Model(t).Select("name", "description", "template_type", "content", "tags", "is_public", "updated_at").Updates(t).Error
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Delete(&Template{}, id).Error
}

func (s *Store) List(ctx context.Context, templateType string, page, pageSize int) ([]Template, int64, error) {
	query := s.db.WithContext(ctx).Model(&Template{})
	if templateType != "" {
		query = query.Where("template_type = ?", templateType)
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

	var items []Template
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *Store) IncrementUsage(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Model(&Template{}).Where("id = ?", id).UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error
}
