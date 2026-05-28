package model

import "time"

// BrowserFingerprint 浏览器指纹配置
type BrowserFingerprint struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID        int64     `json:"user_id" gorm:"index;not null;comment:用户ID"`
	Name          string    `json:"name" gorm:"size:128;not null;comment:指纹名称"`
	UserAgent     string    `json:"user_agent" gorm:"size:512;not null;comment:浏览器UserAgent"`
	Platform      string    `json:"platform" gorm:"size:32;not null;comment:操作系统平台"`
	Screen        string    `json:"screen" gorm:"size:32;comment:屏幕分辨率"`
	Language      string    `json:"language" gorm:"size:32;comment:浏览器语言"`
	Timezone      string    `json:"timezone" gorm:"size:64;comment:时区"`
	WebGLVendor   string    `json:"webgl_vendor" gorm:"size:128;comment:WebGL供应商"`
	WebGLRenderer string   `json:"webgl_renderer" gorm:"size:256;comment:WebGL渲染器"`
	CanvasHash    string    `json:"canvas_hash" gorm:"size:128;comment:Canvas指纹哈希"`
	AudioHash     string    `json:"audio_hash" gorm:"size:128;comment:Audio指纹哈希"`
	Status        string    `json:"status" gorm:"size:32;default:active;index;comment:状态"`
	AccountCount  int       `json:"account_count" gorm:"default:0;comment:关联账号数"`
	IsEnabled     bool      `json:"is_enabled" gorm:"default:true;index;comment:是否启用"`
	CreatedAt     time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (BrowserFingerprint) TableName() string {
	return "browser_fingerprints"
}

// ProxyIP 代理IP配置
type ProxyIP struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID    int64      `json:"user_id" gorm:"index;not null;comment:用户ID"`
	IP        string     `json:"ip" gorm:"size:64;not null;comment:代理IP地址"`
	Port      int        `json:"port" gorm:"not null;comment:端口号"`
	Protocol  string     `json:"protocol" gorm:"size:16;not null;comment:协议类型"`
	Username  string     `json:"username" gorm:"size:128;comment:认证用户名"`
	Password  string     `json:"password" gorm:"size:128;comment:认证密码"`
	Location  string     `json:"location" gorm:"size:128;comment:地理位置"`
	Speed     int        `json:"speed" gorm:"comment:延迟(ms)"`
	Uptime    float64    `json:"uptime" gorm:"comment:可用率(%)"`
	Status    string     `json:"status" gorm:"size:32;default:active;index;comment:状态"`
	LastCheck *time.Time `json:"last_check" gorm:"comment:最后检查时间"`
	IsEnabled bool       `json:"is_enabled" gorm:"default:true;index;comment:是否启用"`
	CreatedAt time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"comment:更新时间"`
}

func (ProxyIP) TableName() string {
	return "proxy_ips"
}

// ContentTemplate 内容模板
type ContentTemplate struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID       int64     `json:"user_id" gorm:"index;not null;comment:用户ID"`
	Name         string    `json:"name" gorm:"size:128;not null;comment:模板名称"`
	TemplateType string    `json:"template_type" gorm:"size:32;index;comment:模板类型"`
	ChannelType  string    `json:"channel_type" gorm:"size:32;index;comment:渠道类型"`
	Content      string    `json:"content" gorm:"type:text;not null;comment:模板内容"`
	Variables    string    `json:"variables" gorm:"type:text;comment:变量列表JSON"`
	Description  string    `json:"description" gorm:"size:512;comment:模板描述"`
	UsageCount   int       `json:"usage_count" gorm:"default:0;comment:使用次数"`
	IsPublic     bool      `json:"is_public" gorm:"default:false;index;comment:是否公开"`
	IsEnabled    bool      `json:"is_enabled" gorm:"default:true;index;comment:是否启用"`
	CreatedAt    time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (ContentTemplate) TableName() string {
	return "content_templates"
}

// StaggerStrategy 错峰策略
type StaggerStrategy struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID      int64     `json:"user_id" gorm:"index;not null;comment:用户ID"`
	Name        string    `json:"name" gorm:"size:128;not null;comment:策略名称"`
	Accounts    int       `json:"accounts" gorm:"default:10;comment:账号数量"`
	Interval    int       `json:"interval" gorm:"default:5;comment:间隔(分钟)"`
	RandomRange int       `json:"random_range" gorm:"default:30;comment:随机范围(%)"`
	Status      string    `json:"status" gorm:"size:32;default:active;index;comment:状态"`
	IsEnabled   bool      `json:"is_enabled" gorm:"default:true;index;comment:是否启用"`
	CreatedAt   time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (StaggerStrategy) TableName() string {
	return "stagger_strategies"
}

// StaggerConfig 错峰配置
type StaggerConfig struct {
	ID               int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID           int64     `json:"user_id" gorm:"uniqueIndex;not null;comment:用户ID"`
	MinInterval      int       `json:"min_interval" gorm:"default:5;comment:最小间隔(分钟)"`
	MaxInterval      int       `json:"max_interval" gorm:"default:15;comment:最大间隔(分钟)"`
	RandomRange      int       `json:"random_range" gorm:"default:30;comment:随机范围(%)"`
	BatchSize        int       `json:"batch_size" gorm:"default:10;comment:批量大小"`
	CooldownAfter    int       `json:"cooldown_after" gorm:"default:50;comment:冷却触发数"`
	CooldownDuration int       `json:"cooldown_duration" gorm:"default:30;comment:冷却时长(分钟)"`
	CreatedAt        time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (StaggerConfig) TableName() string {
	return "stagger_configs"
}

// ContentFingerprint 内容指纹（用于去重）
type ContentFingerprint struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID          int64     `json:"user_id" gorm:"index;not null;comment:用户ID"`
	ContentID       int64     `json:"content_id" gorm:"index;comment:内容ID"`
	TitleHash       string    `json:"title_hash" gorm:"size:64;index;comment:标题SimHash"`
	BodyHash        string    `json:"body_hash" gorm:"size:64;index;comment:正文SimHash"`
	TitleFingerprint string   `json:"title_fingerprint" gorm:"size:512;comment:标题分词指纹"`
	BodyFingerprint string    `json:"body_fingerprint" gorm:"type:text;comment:正文分词指纹"`
	Keywords        string    `json:"keywords" gorm:"type:text;comment:关键词JSON"`
	WordCount       int       `json:"word_count" gorm:"comment:字数"`
	ContentType     string    `json:"content_type" gorm:"size:32;index;comment:内容类型"`
	CreatedAt       time.Time `json:"created_at" gorm:"comment:创建时间"`
}

func (ContentFingerprint) TableName() string {
	return "content_fingerprints"
}

// SynonymDict 同义词词典
type SynonymDict struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Word      string    `json:"word" gorm:"size:64;uniqueIndex;not null;comment:原词"`
	Synonyms  string    `json:"synonyms" gorm:"type:text;not null;comment:同义词JSON数组"`
	Category  string    `json:"category" gorm:"size:32;index;comment:分类"`
	Weight    int       `json:"weight" gorm:"default:1;comment:权重"`
	IsActive  bool      `json:"is_active" gorm:"default:true;index;comment:是否启用"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (SynonymDict) TableName() string {
	return "synonym_dict"
}

// DedupHistory 去重历史记录
type DedupHistory struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID         int64     `json:"user_id" gorm:"index;not null;comment:用户ID"`
	ContentID      int64     `json:"content_id" gorm:"index;comment:内容ID"`
	OriginalHash   string    `json:"original_hash" gorm:"size:64;index;comment:原始内容hash"`
	DedupedHash    string    `json:"deduped_hash" gorm:"size:64;comment:去重后hash"`
	Similarity     float32   `json:"similarity" gorm:"comment:与原文相似度"`
	DuplicateCount int       `json:"duplicate_count" gorm:"comment:发现的重复内容数量"`
	DuplicateIDs   string    `json:"duplicate_ids" gorm:"type:text;comment:重复内容ID列表JSON"`
	Strategy       string    `json:"strategy" gorm:"size:32;comment:去重策略"`
	AITransformed  bool      `json:"ai_transformed" gorm:"comment:是否使用AI改写"`
	CreatedAt      time.Time `json:"created_at" gorm:"comment:创建时间"`
}

func (DedupHistory) TableName() string {
	return "dedup_history"
}
