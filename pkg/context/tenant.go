package context

import (
	"context"
)

type contextKey string

const (
	tenantContextKey contextKey = "tenant_context"
	traceContextKey  contextKey = "trace_context"
)

// TenantContext 租户上下文
type TenantContext struct {
	TenantID int64
	UserID   int64
	Username string
	Plan     string
}

// TraceContext 追踪上下文
type TraceContext struct {
	TraceID      string
	SpanID       string
	ParentSpanID string
}

// WithTenant 将租户上下文注入到 Context
func WithTenant(ctx context.Context, tc *TenantContext) context.Context {
	return context.WithValue(ctx, tenantContextKey, tc)
}

// GetTenant 从 Context 中获取租户上下文
func GetTenant(ctx context.Context) *TenantContext {
	tc, _ := ctx.Value(tenantContextKey).(*TenantContext)
	return tc
}

// MustGetTenant 从 Context 中获取租户上下文（如果不存在则 panic）
func MustGetTenant(ctx context.Context) *TenantContext {
	tc := GetTenant(ctx)
	if tc == nil {
		panic("tenant context not found")
	}
	return tc
}

// WithTrace 将追踪上下文注入到 Context
func WithTrace(ctx context.Context, tc *TraceContext) context.Context {
	return context.WithValue(ctx, traceContextKey, tc)
}

// GetTrace 从 Context 中获取追踪上下文
func GetTrace(ctx context.Context) *TraceContext {
	tc, _ := ctx.Value(traceContextKey).(*TraceContext)
	return tc
}

// GetTenantID 从 Context 中获取租户 ID
func GetTenantID(ctx context.Context) int64 {
	tc := GetTenant(ctx)
	if tc != nil {
		return tc.TenantID
	}
	return 0
}

// GetUserID 从 Context 中获取用户 ID
func GetUserID(ctx context.Context) int64 {
	tc := GetTenant(ctx)
	if tc != nil {
		return tc.UserID
	}
	return 0
}
