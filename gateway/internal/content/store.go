package content

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
	if err := s.db.AutoMigrate(&Content{}); err != nil {
		if !strings.Contains(err.Error(), "Duplicate key name") {
			return err
		}
	}
	return nil
}

func (s *Store) Create(ctx context.Context, c *Content) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return s.db.WithContext(ctx).Create(c).Error
}

func (s *Store) GetByID(ctx context.Context, id int64) (*Content, error) {
	var c Content
	err := s.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) Update(ctx context.Context, c *Content) error {
	c.UpdatedAt = time.Now()
	return s.db.WithContext(ctx).Model(c).Select("title", "body", "schema_markup", "updated_at").Updates(c).Error
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Delete(&Content{}, id).Error
}

func (s *Store) List(ctx context.Context, userID int64, status int, contentType string, page, pageSize int) ([]Content, int64, error) {
	query := s.db.WithContext(ctx).Model(&Content{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if contentType != "" {
		query = query.Where("content_type = ?", contentType)
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

	var items []Content
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
