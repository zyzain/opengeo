package domain

import (
	"errors"
	"time"
)

// BrandStatus 品牌状态
type BrandStatus int32

const (
	BrandStatusUnspecified BrandStatus = 0
	BrandStatusActive      BrandStatus = 1
	BrandStatusArchived    BrandStatus = 2
	BrandStatusDisabled    BrandStatus = 3
)

// Brand 品牌聚合根
type Brand struct {
	ID           int64
	TenantID     int64
	Name         string
	Slug         string
	Description  string
	LogoURL      string
	Website      string
	Industry     string
	FoundedYear  int32
	Headquarters string
	Status       BrandStatus
	Settings     map[string]string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// BrandMetadata 品牌元数据值对象
type BrandMetadata struct {
	BrandID             int64
	VIProfile           VIProfile
	ToneProfile         ToneProfile
	AudienceProfiles    []AudienceProfile
	CompetitorList      []CompetitorInfo
	BrandValues         []string
	UniqueSellingPoints []string
	SchemaVersion       string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// VIProfile VI 规范值对象
type VIProfile struct {
	PrimaryColor   string
	SecondaryColor string
	LogoURL        string
	FontFamily     string
	BrandKeywords  []string
	Slogan         string
}

// ToneProfile 语调规范值对象
type ToneProfile struct {
	Formality       string
	Personality     string
	AvoidWords      []string
	PreferredPhrases []string
	StyleGuide      string
}

// AudienceProfile 受众画像值对象
type AudienceProfile struct {
	Name             string
	AgeRange         string
	Interests        []string
	PainPoints       []string
	PreferredChannels []string
	Locations        []string
	Languages        []string
}

// CompetitorInfo 竞品信息值对象
type CompetitorInfo struct {
	Name        string
	Domain      string
	Description string
	Strengths   []string
	Weaknesses  []string
}

// GlossaryEntry 术语条目实体
type GlossaryEntry struct {
	ID          int64
	BrandID     int64
	Term        string
	Definition  string
	Category    string
	Aliases     []string
	Context     string
	IsForbidden bool
	IsPreferred bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// BrandSnapshot 品牌快照实体
type BrandSnapshot struct {
	ID           int64
	BrandID      int64
	Version      string
	SnapshotData string
	ChangeLog    string
	CreatedBy    int64
	CreatedAt    time.Time
}

// 领域错误
var (
	ErrBrandNotFound      = errors.New("brand not found")
	ErrBrandSlugExists    = errors.New("brand slug already exists")
	ErrBrandQuotaExceeded = errors.New("brand quota exceeded")
	ErrGlossaryNotFound   = errors.New("glossary entry not found")
	ErrSnapshotNotFound   = errors.New("brand snapshot not found")
)

// Validate 验证品牌实体
func (b *Brand) Validate() error {
	if b.TenantID <= 0 {
		return errors.New("tenant_id is required")
	}
	if b.Name == "" {
		return errors.New("brand name is required")
	}
	if b.Slug == "" {
		return errors.New("brand slug is required")
	}
	return nil
}

// IsActive 检查品牌是否活跃
func (b *Brand) IsActive() bool {
	return b.Status == BrandStatusActive
}

// Validate 验证术语条目
func (g *GlossaryEntry) Validate() error {
	if g.BrandID <= 0 {
		return errors.New("brand_id is required")
	}
	if g.Term == "" {
		return errors.New("term is required")
	}
	if g.Definition == "" {
		return errors.New("definition is required")
	}
	return nil
}
