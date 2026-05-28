package dal

import (
	"context"
	"fmt"
	"time"

	"opengeo/pkg/crypto"
)

// SeedAdmin 创建超级管理员账号
// 默认账号: admin
// 默认密码: Admin@123456
func SeedAdmin(ctx context.Context) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// 检查是否已存在
	var count int64
	DB.WithContext(ctx).Model(&User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		fmt.Println("Admin user already exists, skipping...")
		return nil
	}

	// 生成 bcrypt 密码哈希
	hashedPassword, err := crypto.HashPassword("Admin@123456")
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	// 创建超级管理员用户（tenant_id=0 表示超级管理员，不受租户限制）
	admin := &User{
		TenantID:    0,
		Username:    "admin",
		Password:    hashedPassword,
		Email:       "admin@opengeo.com",
		Status:      1,
		LastLoginAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := DB.WithContext(ctx).Create(admin).Error; err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}

	// 分配管理员角色
	var adminRole Role
	if err := DB.WithContext(ctx).Where("name = ?", "admin").First(&adminRole).Error; err == nil {
		userRole := &UserRole{
			UserID: admin.ID,
			RoleID: adminRole.ID,
		}
		DB.WithContext(ctx).Create(userRole)
	}

	fmt.Println("===========================================")
	fmt.Println("超级管理员账号创建成功")
	fmt.Println("===========================================")
	fmt.Println("用户名: admin")
	fmt.Println("密  码: Admin@123456")
	fmt.Println("邮  箱: admin@opengeo.com")
	fmt.Println("===========================================")
	fmt.Println("请登录后立即修改默认密码！")
	fmt.Println("===========================================")

	return nil
}

// SeedTenant 创建默认租户
// 委托给 rbac_dal.go 中的实现
// (已在 rbac_dal.go 中定义)
