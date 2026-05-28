package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"opengeo/service/publish/internal/domain/model"
)

// ContentFingerprintRepository 内容指纹仓储
type ContentFingerprintRepository struct {
	db *gorm.DB
}

// NewContentFingerprintRepository 创建内容指纹仓储
func NewContentFingerprintRepository(db *gorm.DB) *ContentFingerprintRepository {
	return &ContentFingerprintRepository{db: db}
}

// Create 创建内容指纹
func (r *ContentFingerprintRepository) Create(ctx context.Context, fp *model.ContentFingerprint) error {
	if err := r.db.WithContext(ctx).Create(fp).Error; err != nil {
		return fmt.Errorf("failed to create content fingerprint: %w", err)
	}
	return nil
}

// FindSimilarByHash 通过hash查找相似内容
func (r *ContentFingerprintRepository) FindSimilarByHash(ctx context.Context, userID int64, titleHash, bodyHash string, limit int) ([]*model.ContentFingerprint, error) {
	var fps []*model.ContentFingerprint

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// 查找标题或正文hash相同的内容
	if titleHash != "" && bodyHash != "" {
		query = query.Where("(title_hash = ? OR body_hash = ?)", titleHash, bodyHash)
	} else if titleHash != "" {
		query = query.Where("title_hash = ?", titleHash)
	} else if bodyHash != "" {
		query = query.Where("body_hash = ?", bodyHash)
	} else {
		return fps, nil
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&fps).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find similar content: %w", err)
	}

	return fps, nil
}

// FindByContentID 通过内容ID查找指纹
func (r *ContentFingerprintRepository) FindByContentID(ctx context.Context, contentID int64) (*model.ContentFingerprint, error) {
	var fp model.ContentFingerprint
	if err := r.db.WithContext(ctx).Where("content_id = ?", contentID).First(&fp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find fingerprint: %w", err)
	}
	return &fp, nil
}

// ListByUser 列出用户的内容指纹
func (r *ContentFingerprintRepository) ListByUser(ctx context.Context, userID int64, contentType string, limit int) ([]*model.ContentFingerprint, error) {
	var fps []*model.ContentFingerprint

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if contentType != "" {
		query = query.Where("content_type = ?", contentType)
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&fps).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list fingerprints: %w", err)
	}

	return fps, nil
}

// DeleteOlderThan 删除指定时间之前的指纹
func (r *ContentFingerprintRepository) DeleteOlderThan(ctx context.Context, userID int64, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at < ?", userID, before).
		Delete(&model.ContentFingerprint{})

	return result.RowsAffected, result.Error
}

// SynonymDictRepository 同义词词典仓储
type SynonymDictRepository struct {
	db *gorm.DB
}

// NewSynonymDictRepository 创建同义词词典仓储
func NewSynonymDictRepository(db *gorm.DB) *SynonymDictRepository {
	return &SynonymDictRepository{db: db}
}

// GetAll 获取所有启用的同义词
func (r *SynonymDictRepository) GetAll(ctx context.Context) ([]*model.SynonymDict, error) {
	var dicts []*model.SynonymDict
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&dicts).Error; err != nil {
		return nil, fmt.Errorf("failed to get synonym dict: %w", err)
	}
	return dicts, nil
}

// GetByCategory 按分类获取同义词
func (r *SynonymDictRepository) GetByCategory(ctx context.Context, category string) ([]*model.SynonymDict, error) {
	var dicts []*model.SynonymDict
	if err := r.db.WithContext(ctx).Where("category = ? AND is_active = ?", category, true).Find(&dicts).Error; err != nil {
		return nil, fmt.Errorf("failed to get synonym dict by category: %w", err)
	}
	return dicts, nil
}

// BuildSynonymMap 构建同义词映射
func (r *SynonymDictRepository) BuildSynonymMap(ctx context.Context) (map[string][]string, error) {
	dicts, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for _, d := range dicts {
		var synonyms []string
		if err := json.Unmarshal([]byte(d.Synonyms), &synonyms); err != nil {
			continue
		}
		if len(synonyms) > 0 {
			result[d.Word] = synonyms
		}
	}

	return result, nil
}

// BatchCreate 批量创建同义词
func (r *SynonymDictRepository) BatchCreate(ctx context.Context, dicts []*model.SynonymDict) error {
	if len(dicts) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(dicts, 100).Error
}

// DedupHistoryRepository 去重历史仓储
type DedupHistoryRepository struct {
	db *gorm.DB
}

// NewDedupHistoryRepository 创建去重历史仓储
func NewDedupHistoryRepository(db *gorm.DB) *DedupHistoryRepository {
	return &DedupHistoryRepository{db: db}
}

// Create 创建去重历史
func (r *DedupHistoryRepository) Create(ctx context.Context, history *model.DedupHistory) error {
	if err := r.db.WithContext(ctx).Create(history).Error; err != nil {
		return fmt.Errorf("failed to create dedup history: %w", err)
	}
	return nil
}

// ListByUser 列出用户的去重历史
func (r *DedupHistoryRepository) ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]*model.DedupHistory, int32, error) {
	var histories []*model.DedupHistory
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count dedup history: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&histories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list dedup history: %w", err)
	}

	return histories, int32(total), nil
}

// GetStats 获取去重统计
func (r *DedupHistoryRepository) GetStats(ctx context.Context, userID int64) (*DedupStats, error) {
	var stats DedupStats

	// 总去重次数
	r.db.WithContext(ctx).Model(&model.DedupHistory{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalCount)

	// 使用AI改写的次数
	r.db.WithContext(ctx).Model(&model.DedupHistory{}).
		Where("user_id = ? AND ai_transformed = ?", userID, true).
		Count(&stats.AICount)

	// 平均相似度
	r.db.WithContext(ctx).Model(&model.DedupHistory{}).
		Where("user_id = ?", userID).
		Select("AVG(similarity)").
		Scan(&stats.AvgSimilarity)

	// 发现的重复内容总数
	r.db.WithContext(ctx).Model(&model.DedupHistory{}).
		Where("user_id = ?", userID).
		Select("SUM(duplicate_count)").
		Scan(&stats.TotalDuplicates)

	return &stats, nil
}

// DedupStats 去重统计
type DedupStats struct {
	TotalCount      int64   `json:"total_count"`
	AICount         int64   `json:"ai_count"`
	AvgSimilarity   float32 `json:"avg_similarity"`
	TotalDuplicates int64   `json:"total_duplicates"`
}
