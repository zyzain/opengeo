package service

import (
	"context"
	"fmt"
	"time"

	"opengeo/pkg/crypto"
	jwtUtil "opengeo/pkg/jwt"
	"opengeo/service/account/internal/dal"
)

// UserService 用户服务
type UserService struct {
	userRepo *dal.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo *dal.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, username, password, email string, tenantID int64) (*dal.User, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// 验证密码强度
	if err := crypto.ValidatePasswordStrength(password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	// 加密密码
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &dal.User{
		TenantID:  tenantID,
		Username:  username,
		Password:  hashedPassword,
		Email:     email,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx context.Context, userID int64) (*dal.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, userID int64, email string, status int32) (*dal.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if email != "" {
		user.Email = email
	}
	if status > 0 {
		user.Status = status
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers 列出用户
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int32, keyword string, tenantID int64) ([]*dal.User, int32, error) {
	users, total, err := s.userRepo.List(ctx, page, pageSize, keyword, tenantID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, username, password string) (string, string, *dal.User, error) {
	// 根据用户名获取用户
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", nil, fmt.Errorf("invalid username or password")
	}

	// 检查用户状态
	if user.Status != 1 {
		return "", "", nil, fmt.Errorf("user account is disabled")
	}

	// 验证密码
	if !crypto.CheckPassword(password, user.Password) {
		return "", "", nil, fmt.Errorf("invalid username or password")
	}

	// 生成JWT Token
	token, err := jwtUtil.GenerateToken(user.ID, user.Username, user.Email, "user", user.TenantID)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 生成Refresh Token
	refreshToken, err := jwtUtil.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now()
	s.userRepo.Update(ctx, user)

	return token, refreshToken, user, nil
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, username, password, email string) (string, string, *dal.User, error) {
	// 创建用户（使用租户ID 0，需调用方指定具体租户）
	user, err := s.CreateUser(ctx, username, password, email, 0)
	if err != nil {
		return "", "", nil, err
	}

	// 生成JWT Token
	token, err := jwtUtil.GenerateToken(user.ID, user.Username, user.Email, "user", user.TenantID)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 生成Refresh Token
	refreshToken, err := jwtUtil.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return token, refreshToken, user, nil
}

// RefreshToken 刷新Token
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// 验证并刷新Token
	newToken, newRefreshToken, err := jwtUtil.RefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, newRefreshToken, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 验证旧密码
	if !crypto.CheckPassword(oldPassword, user.Password) {
		return fmt.Errorf("invalid old password")
	}

	// 验证新密码强度
	if err := crypto.ValidatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("invalid new password: %w", err)
	}

	// 加密新密码
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}