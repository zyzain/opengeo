package auth

import "time"

// Tenant 租户（平台账号）
type Tenant struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	Domain    string    `json:"domain" gorm:"size:256"`
	Status    int32     `json:"status" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime"`
}

// User 用户（归属租户）
type User struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	TenantID    int64     `gorm:"index;not null;default:0"`
	Username    string    `gorm:"uniqueIndex:idx_tenant_username;size:64;not null"`
	Password    string    `gorm:"size:256;not null"`
	Email       string    `gorm:"size:128"`
	Status      int32     `gorm:"default:1"`
	LastLoginAt time.Time `gorm:"type:datetime"`
	CreatedAt   time.Time `gorm:"type:datetime"`
	UpdatedAt   time.Time `gorm:"type:datetime"`
}

// Role 角色（租户隔离）
type Role struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID    int64     `json:"tenant_id" gorm:"index;not null;default:0"`
	Name        string    `json:"name" gorm:"size:64;not null"`
	Description string    `json:"description" gorm:"size:256"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime"`
}

// Permission 权限（全局定义）
type Permission struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:64;not null"`
	Description string    `json:"description" gorm:"size:256"`
	Resource    string    `json:"resource" gorm:"size:128"`
	Action      string    `json:"action" gorm:"size:64"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime"`
}

// UserRole 用户角色关联
type UserRole struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey"`
}

// RolePermission 角色权限关联
type RolePermission struct {
	RoleID       int64 `gorm:"primaryKey"`
	PermissionID int64 `gorm:"primaryKey"`
}
