package template

import "time"

type Template struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	UserID       int64     `gorm:"index;not null;default:0"`
	Name         string    `gorm:"size:256;not null"`
	Description  string    `gorm:"size:512;default:''"`
	TemplateType string    `gorm:"size:64;not null;default:'prompt'"`
	Content      string    `gorm:"type:text;not null"`
	Tags         string    `gorm:"size:512;default:''"`
	IsPublic     bool      `gorm:"default:true"`
	UsageCount   int64     `gorm:"default:0"`
	Rating       float64   `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"type:datetime"`
	UpdatedAt    time.Time `gorm:"type:datetime"`
}
