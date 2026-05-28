package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	globalConfig *Config
	configOnce   sync.Once
)

// Config 全局配置
type Config struct {
	// JWT配置
	JWT JWTConfig `json:"jwt"`

	// 数据库配置
	Database DatabaseConfig `json:"database"`

	// 分页配置
	Pagination PaginationConfig `json:"pagination"`

	// AI模型配置
	AIModels AIModelsConfig `json:"ai_models"`

	// 防封引擎配置
	AntiBan AntiBanConfig `json:"anti_ban"`

	// 错峰调度配置
	Stagger StaggerConfig `json:"stagger"`

	// 重试配置
	Retry RetryConfig `json:"retry"`

	// Webhook配置
	Webhook WebhookConfig `json:"webhook"`

	// 健康检测配置
	HealthCheck HealthCheckConfig `json:"health_check"`

	// LLM客户端配置
	LLM LLMConfig `json:"llm"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	TokenExpire        time.Duration `json:"token_expire"`
	RefreshTokenExpire time.Duration `json:"refresh_token_expire"`
	Issuer             string        `json:"issuer"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MaxIdleConns    int           `json:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

// PaginationConfig 分页配置
type PaginationConfig struct {
	DefaultPageSize int `json:"default_page_size"`
	MaxPageSize     int `json:"max_page_size"`
}

// AIModelsConfig AI模型配置
type AIModelsConfig struct {
	DeepSeek AIModelConfig `json:"deepseek"`
	Kimi     AIModelConfig `json:"kimi"`
	Doubao   AIModelConfig `json:"doubao"`
	ChatGPT  AIModelConfig `json:"chatgpt"`
	Default  string        `json:"default"`
}

// AIModelConfig 单个AI模型配置
type AIModelConfig struct {
	BaseURL     string  `json:"base_url"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
}

// AntiBanConfig 防封引擎配置
type AntiBanConfig struct {
	// 平台延迟配置 [min, max] 秒
	PlatformDelays map[string][2]int `json:"platform_delays"`
	// 默认延迟 [min, max] 秒
	DefaultDelay [2]int `json:"default_delay"`
	// 最小延迟秒数
	MinDelay int `json:"min_delay"`
	// 正态分布偏移比例
	OffsetRatio float64 `json:"offset_ratio"`
	// 滚动间隔 [min, max] 毫秒
	ScrollInterval [2]int `json:"scroll_interval"`
	// 打字延迟 [min, max] 毫秒
	TypingDelay [2]int `json:"typing_delay"`
	// 打字停顿概率
	TypingPauseProbability float64 `json:"typing_pause_probability"`
	// 打字停顿延迟 [min, max] 毫秒
	TypingPauseDelay [2]int `json:"typing_pause_delay"`
	// 频率限制
	RateLimits map[string]RateLimitConfig `json:"rate_limits"`
	// 默认频率限制
	DefaultRateLimit RateLimitConfig `json:"default_rate_limit"`
	// 代理失败阈值
	ProxyFailureThreshold int `json:"proxy_failure_threshold"`
}

// RateLimitConfig 频率限制配置
type RateLimitConfig struct {
	MaxPerHour  int `json:"max_per_hour"`
	MaxPerDay   int `json:"max_per_day"`
	MinInterval int `json:"min_interval"` // 秒
}

// StaggerConfig 错峰调度配置
type StaggerConfig struct {
	MinInterval    time.Duration `json:"min_interval"`
	MaxInterval    time.Duration `json:"max_interval"`
	VarianceRatio  float32       `json:"variance_ratio"`
	MaxConcurrency int           `json:"max_concurrency"`
	BurstLimit     int           `json:"burst_limit"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries      int32   `json:"max_retries"`
	InitialDelay    int32   `json:"initial_delay"`    // 秒
	MaxDelay        int32   `json:"max_delay"`        // 秒
	BackoffFactor   float32 `json:"backoff_factor"`
	EnableFallback  bool    `json:"enable_fallback"`
	EnableAutoRetry bool    `json:"enable_auto_retry"`
}

// WebhookConfig Webhook配置
type WebhookConfig struct {
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
}

// HealthCheckConfig 健康检测配置
type HealthCheckConfig struct {
	HTTPTimeout     time.Duration `json:"http_timeout"`
	CheckInterval   int           `json:"check_interval"`   // 分钟
	BatchLimit      int           `json:"batch_limit"`
	AlertScoreThreshold float32 `json:"alert_score_threshold"`
	PauseScoreThreshold float32 `json:"pause_score_threshold"`
}

// LLMConfig LLM客户端配置
type LLMConfig struct {
	Timeout     time.Duration `json:"timeout"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float32       `json:"temperature"`
}

// GetConfig 获取全局配置（单例）
func GetConfig() *Config {
	configOnce.Do(func() {
		globalConfig = loadConfig()
	})
	return globalConfig
}

// loadConfig 加载配置（优先级：环境变量 > 配置文件 > 默认值）
func loadConfig() *Config {
	cfg := getDefaultConfig()

	// 尝试从配置文件加载
	if configPath := os.Getenv("CONFIG_FILE"); configPath != "" {
		if err := loadFromFile(cfg, configPath); err != nil {
			fmt.Printf("Warning: failed to load config file: %v\n", err)
		}
	}

	// 从环境变量覆盖
	loadFromEnv(cfg)

	return cfg
}

// getDefaultConfig 获取默认配置
func getDefaultConfig() *Config {
	return &Config{
		JWT: JWTConfig{
			TokenExpire:        24 * time.Hour,
			RefreshTokenExpire: 7 * 24 * time.Hour,
			Issuer:             "opengeo",
		},
		Database: DatabaseConfig{
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: time.Hour,
		},
		Pagination: PaginationConfig{
			DefaultPageSize: 20,
			MaxPageSize:     100,
		},
		AIModels: AIModelsConfig{
			DeepSeek: AIModelConfig{
				BaseURL:     "https://api.deepseek.com/v1",
				Model:       "deepseek-chat",
				MaxTokens:   4096,
				Temperature: 0.7,
			},
			Kimi: AIModelConfig{
				BaseURL:     "https://api.moonshot.cn/v1",
				Model:       "moonshot-v1-8k",
				MaxTokens:   4096,
				Temperature: 0.7,
			},
			Doubao: AIModelConfig{
				BaseURL:     "https://ark.cn-beijing.volces.com/api/v3",
				Model:       "doubao-pro-4k",
				MaxTokens:   4096,
				Temperature: 0.7,
			},
			ChatGPT: AIModelConfig{
				BaseURL:     "https://api.openai.com/v1",
				Model:       "gpt-4o-mini",
				MaxTokens:   4096,
				Temperature: 0.7,
			},
			Default: "deepseek",
		},
		AntiBan: AntiBanConfig{
			PlatformDelays: map[string][2]int{
				"wechat":      {30, 120},
				"weibo":       {10, 60},
				"douyin":      {20, 90},
				"xiaohongshu": {15, 75},
				"zhihu":       {20, 90},
				"toutiao":     {10, 60},
			},
			DefaultDelay:           [2]int{15, 90},
			MinDelay:               5,
			OffsetRatio:            0.2,
			ScrollInterval:         [2]int{500, 3000},
			TypingDelay:            [2]int{50, 200},
			TypingPauseProbability: 0.05,
			TypingPauseDelay:       [2]int{500, 2000},
			RateLimits: map[string]RateLimitConfig{
				"wechat":      {MaxPerHour: 5, MaxPerDay: 20, MinInterval: 300},
				"weibo":       {MaxPerHour: 10, MaxPerDay: 50, MinInterval: 60},
				"douyin":      {MaxPerHour: 3, MaxPerDay: 10, MinInterval: 600},
				"xiaohongshu": {MaxPerHour: 5, MaxPerDay: 20, MinInterval: 300},
				"zhihu":       {MaxPerHour: 5, MaxPerDay: 30, MinInterval: 300},
				"toutiao":     {MaxPerHour: 10, MaxPerDay: 50, MinInterval: 60},
			},
			DefaultRateLimit: RateLimitConfig{
				MaxPerHour:  10,
				MaxPerDay:   50,
				MinInterval: 60,
			},
			ProxyFailureThreshold: 3,
		},
		Stagger: StaggerConfig{
			MinInterval:    5 * time.Minute,
			MaxInterval:    30 * time.Minute,
			VarianceRatio:  0.3,
			MaxConcurrency: 10,
			BurstLimit:     5,
		},
		Retry: RetryConfig{
			MaxRetries:      3,
			InitialDelay:    30,
			MaxDelay:        3600,
			BackoffFactor:   2.0,
			EnableFallback:  true,
			EnableAutoRetry: true,
		},
		Webhook: WebhookConfig{
			Timeout:    30 * time.Second,
			MaxRetries: 3,
		},
		HealthCheck: HealthCheckConfig{
			HTTPTimeout:         10 * time.Second,
			CheckInterval:       60,
			BatchLimit:          100,
			AlertScoreThreshold: 60,
			PauseScoreThreshold: 30,
		},
		LLM: LLMConfig{
			Timeout:     30 * time.Second,
			MaxTokens:   4096,
			Temperature: 0.7,
		},
	}
}

// loadFromFile 从文件加载配置
func loadFromFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}
	return json.Unmarshal(data, cfg)
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(cfg *Config) {
	// JWT配置
	if v := os.Getenv("JWT_TOKEN_EXPIRE_HOURS"); v != "" {
		if hours, err := strconv.Atoi(v); err == nil {
			cfg.JWT.TokenExpire = time.Duration(hours) * time.Hour
		}
	}
	if v := os.Getenv("JWT_REFRESH_EXPIRE_DAYS"); v != "" {
		if days, err := strconv.Atoi(v); err == nil {
			cfg.JWT.RefreshTokenExpire = time.Duration(days) * 24 * time.Hour
		}
	}

	// 数据库配置
	if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Database.MaxIdleConns = n
		}
	}
	if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Database.MaxOpenConns = n
		}
	}

	// 分页配置
	if v := os.Getenv("DEFAULT_PAGE_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Pagination.DefaultPageSize = n
		}
	}

	// AI模型配置
	if v := os.Getenv("DEEPSEEK_BASE_URL"); v != "" {
		cfg.AIModels.DeepSeek.BaseURL = v
	}
	if v := os.Getenv("DEEPSEEK_MODEL"); v != "" {
		cfg.AIModels.DeepSeek.Model = v
	}
	if v := os.Getenv("KIMI_BASE_URL"); v != "" {
		cfg.AIModels.Kimi.BaseURL = v
	}
	if v := os.Getenv("KIMI_MODEL"); v != "" {
		cfg.AIModels.Kimi.Model = v
	}
	if v := os.Getenv("DOUBAO_BASE_URL"); v != "" {
		cfg.AIModels.Doubao.BaseURL = v
	}
	if v := os.Getenv("DOUBAO_MODEL"); v != "" {
		cfg.AIModels.Doubao.Model = v
	}
	if v := os.Getenv("OPENAI_BASE_URL"); v != "" {
		cfg.AIModels.ChatGPT.BaseURL = v
	}
	if v := os.Getenv("OPENAI_MODEL"); v != "" {
		cfg.AIModels.ChatGPT.Model = v
	}
	if v := os.Getenv("DEFAULT_AI_MODEL"); v != "" {
		cfg.AIModels.Default = v
	}

	// Webhook配置
	if v := os.Getenv("WEBHOOK_TIMEOUT_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Webhook.Timeout = time.Duration(n) * time.Second
		}
	}
	if v := os.Getenv("WEBHOOK_MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Webhook.MaxRetries = n
		}
	}

	// LLM配置
	if v := os.Getenv("LLM_TIMEOUT_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.LLM.Timeout = time.Duration(n) * time.Second
		}
	}
	if v := os.Getenv("LLM_MAX_TOKENS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.LLM.MaxTokens = n
		}
	}
}
