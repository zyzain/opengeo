package model

import "time"

// Platform 发布平台
type Platform struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code         string    `json:"code" gorm:"uniqueIndex;size:32;not null"`   // 平台代码：wechat/weibo/douyin
	Name         string    `json:"name" gorm:"size:64;not null"`               // 平台名称：微信公众号
	Icon         string    `json:"icon" gorm:"size:32"`                        // 图标标识
	Color        string    `json:"color" gorm:"size:16"`                       // 标签颜色
	Description  string    `json:"description" gorm:"size:256"`                // 平台描述
	ConfigSchema string    `json:"config_schema" gorm:"type:text"`             // 配置项JSON Schema
	ConfigTemplate []PlatformConfigTemplate `json:"config_template" gorm:"-"`  // 配置模板（不存储，从ConfigSchema解析）
	IsEnabled    bool      `json:"is_enabled" gorm:"default:true;index"`       // 是否启用
	SortOrder    int32     `json:"sort_order" gorm:"default:0"`                // 排序权重
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PlatformConfigTemplate 平台配置模板
type PlatformConfigTemplate struct {
	Field       string `json:"field"`       // 字段名
	Label       string `json:"label"`       // 显示名称
	Type        string `json:"type"`        // 类型：text/password/textarea
	Required    bool   `json:"required"`    // 是否必填
	Placeholder string `json:"placeholder"` // 占位提示
	Default     string `json:"default"`     // 默认值
}

// DefaultPlatforms 默认平台列表
var DefaultPlatforms = []Platform{
	{
		Code:        "wechat",
		Name:        "微信公众号",
		Icon:        "wechat",
		Color:       "green",
		Description: "微信公众号内容发布",
		ConfigTemplate: []PlatformConfigTemplate{
			{Field: "app_id", Label: "AppID", Type: "text", Required: true, Placeholder: "请输入AppID"},
			{Field: "app_secret", Label: "AppSecret", Type: "password", Required: true, Placeholder: "请输入AppSecret"},
		},
		IsEnabled: true,
		SortOrder: 100,
	},
	{
		Code:        "weibo",
		Name:        "微博",
		Icon:        "weibo",
		Color:       "red",
		Description: "新浪微博内容发布",
		ConfigTemplate: []PlatformConfigTemplate{
			{Field: "access_token", Label: "Access Token", Type: "password", Required: true, Placeholder: "请输入Access Token"},
		},
		IsEnabled: true,
		SortOrder: 90,
	},
	{
		Code:        "douyin",
		Name:        "抖音",
		Icon:        "douyin",
		Color:       "purple",
		Description: "抖音短视频/图文发布",
		ConfigTemplate: []PlatformConfigTemplate{
			{Field: "access_token", Label: "Access Token", Type: "password", Required: true, Placeholder: "请输入Access Token"},
		},
		IsEnabled: true,
		SortOrder: 80,
	},
	{
		Code:        "xiaohongshu",
		Name:        "小红书",
		Icon:        "xiaohongshu",
		Color:       "pink",
		Description: "小红书笔记发布",
		ConfigTemplate: []PlatformConfigTemplate{
			{Field: "access_token", Label: "Access Token", Type: "password", Required: true, Placeholder: "请输入Access Token"},
		},
		IsEnabled: true,
		SortOrder: 70,
	},
	{
		Code:        "zhihu",
		Name:        "知乎",
		Icon:        "zhihu",
		Color:       "blue",
		Description: "知乎文章/回答发布",
		ConfigTemplate: []PlatformConfigTemplate{
			{Field: "access_token", Label: "Access Token", Type: "password", Required: true, Placeholder: "请输入Access Token"},
		},
		IsEnabled: true,
		SortOrder: 60,
	},
	{
		Code:        "toutiao",
		Name:        "今日头条",
		Icon:        "toutiao",
		Color:       "orange",
		Description: "今日头条文章发布",
		ConfigTemplate: []PlatformConfigTemplate{
			{Field: "access_token", Label: "Access Token", Type: "password", Required: true, Placeholder: "请输入Access Token"},
		},
		IsEnabled: true,
		SortOrder: 50,
	},
}
