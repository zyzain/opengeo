package auth

import (
	"testing"
	"time"
)

func TestUser_Fields(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	u := User{
		ID:          1,
		TenantID:    10,
		Username:    "admin",
		Password:    "$2a$10$hashedpassword",
		Email:       "admin@opengeo.com",
		Status:      1,
		LastLoginAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if u.ID != 1 {
		t.Errorf("expected ID=1, got %d", u.ID)
	}
	if u.TenantID != 10 {
		t.Errorf("expected TenantID=10, got %d", u.TenantID)
	}
	if u.Username != "admin" {
		t.Errorf("expected Username='admin', got '%s'", u.Username)
	}
	if u.Status != 1 {
		t.Errorf("expected Status=1, got %d", u.Status)
	}
}

func TestTenant_Fields(t *testing.T) {
	now := time.Now()
	tenant := Tenant{
		ID:        1,
		Name:      "Default",
		Domain:    "default.opengeo.com",
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if tenant.Name != "Default" {
		t.Errorf("expected Name='Default', got '%s'", tenant.Name)
	}
	if tenant.Domain != "default.opengeo.com" {
		t.Errorf("expected Domain='default.opengeo.com', got '%s'", tenant.Domain)
	}
}

func TestRole_Fields(t *testing.T) {
	now := time.Now()
	r := Role{
		ID:          1,
		TenantID:    10,
		Name:        "admin",
		Description: "系统管理员",
		CreatedAt:   now,
	}

	if r.Name != "admin" {
		t.Errorf("expected Name='admin', got '%s'", r.Name)
	}
	if r.TenantID != 10 {
		t.Errorf("expected TenantID=10, got %d", r.TenantID)
	}
}

func TestPermission_Fields(t *testing.T) {
	now := time.Now()
	p := Permission{
		ID:          1,
		Name:        "user:create",
		Description: "创建用户",
		Resource:    "user",
		Action:      "create",
		CreatedAt:   now,
	}

	if p.Name != "user:create" {
		t.Errorf("expected Name='user:create', got '%s'", p.Name)
	}
	if p.Resource != "user" {
		t.Errorf("expected Resource='user', got '%s'", p.Resource)
	}
	if p.Action != "create" {
		t.Errorf("expected Action='create', got '%s'", p.Action)
	}
}

func TestUserRole_Fields(t *testing.T) {
	ur := UserRole{
		UserID: 1,
		RoleID: 2,
	}

	if ur.UserID != 1 {
		t.Errorf("expected UserID=1, got %d", ur.UserID)
	}
	if ur.RoleID != 2 {
		t.Errorf("expected RoleID=2, got %d", ur.RoleID)
	}
}

func TestRolePermission_Fields(t *testing.T) {
	rp := RolePermission{
		RoleID:       1,
		PermissionID: 5,
	}

	if rp.RoleID != 1 {
		t.Errorf("expected RoleID=1, got %d", rp.RoleID)
	}
	if rp.PermissionID != 5 {
		t.Errorf("expected PermissionID=5, got %d", rp.PermissionID)
	}
}

func BenchmarkUserCreation(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		_ = User{
			ID:        int64(i),
			TenantID:  10,
			Username:  "testuser",
			Password:  "hashed",
			Email:     "test@test.com",
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
}
