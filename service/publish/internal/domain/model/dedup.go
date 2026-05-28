package model

import (
	"time"
)

// ContentFingerprint 内容指纹（用于快速查重）
type ContentFingerprint struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int64     `json:"user_id" gorm:"index;not null"`
	ContentID     int64     `json:"content_id" gorm:"index"`
	TitleHash     string    `json:"title_hash" gorm:"size:64;index"`      // 标题SimHash
	BodyHash      string    `json:"body_hash" gorm:"size:64;index"`       // 正文SimHash
	TitleFingerprint string `json:"title_fingerprint" gorm:size:512"`     // 标题分词指纹
	BodyFingerprint  string `json:"body_fingerprint" gorm:type:text"`     // 正文分词指纹
	Keywords      string    `json:"keywords" gorm:type:text"`             // 关键词JSON
	WordCount     int       `json:"word_count"`                           // 字数
	ContentType   string    `json:"content_type" gorm:size:32;index"`     // 内容类型
	CreatedAt     time.Time `json:"created_at"`
}

func (ContentFingerprint) TableName() string {
	return "content_fingerprints"
}

// SynonymDict 同义词词典
type SynonymDict struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Word      string    `json:"word" gorm:"size:64;uniqueIndex;not null"`  // 原词
	Synonyms  string    `json:"synonyms" gorm:type:text;not null"`         // 同义词JSON数组
	Category  string    `json:"category" gorm:size:32;index"`              // 分类: verb, adj, noun, etc
	Weight    int       `json:"weight" gorm:"default:1"`                   // 权重
	IsActive  bool      `json:"is_active" gorm:"default:true;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SynonymDict) TableName() string {
	return "synonym_dict"
}

// DedupHistory 去重历史记录
type DedupHistory struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          int64     `json:"user_id" gorm:"index;not null"`
	ContentID       int64     `json:"content_id" gorm:"index"`
	OriginalHash    string    `json:"original_hash" gorm:"size:64;index"`   // 原始内容hash
	DedupedHash     string    `json:"deduped_hash" gorm:"size:64"`         // 去重后hash
	Similarity      float32   `json:"similarity"`                           // 与原文相似度
	DuplicateCount  int       `json:"duplicate_count"`                      // 发现的重复内容数量
	DuplicateIDs    string    `json:"duplicate_ids" gorm:type:text"`        // 重复内容ID列表JSON
	Strategy        string    `json:"strategy" gorm:"size:32"`              // 去重策略
	AITransformed   bool      `json:"ai_transformed"`                       // 是否使用AI改写
	CreatedAt       time.Time `json:"created_at"`
}

func (DedupHistory) TableName() string {
	return "dedup_history"
}
