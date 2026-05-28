package content

import "time"

type Content struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	UserID       int64     `gorm:"index;not null"`
	Title        string    `gorm:"size:256;not null"`
	Body         string    `gorm:"type:text;not null"`
	ContentType  string    `gorm:"size:64;default:''"`
	SchemaMarkup string    `gorm:"type:text"`
	Status       int32     `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"type:datetime"`
	UpdatedAt    time.Time `gorm:"type:datetime"`
}
