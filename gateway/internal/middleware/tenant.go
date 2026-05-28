package middleware

import (
	"context"
	"net/http"
	"strings"

	appctx "opengeo/pkg/context"
	"opengeo/pkg/logger"
)

// TenantMiddleware 租户上下文中间件
// 从 JWT/Session 中提取租户信息，注入到请求上下文
func TenantMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 从 Header 中提取 tenant_id（实际应该从 JWT 中解析）
			tenantIDStr := r.Header.Get("X-Tenant-ID")
			userIDStr := r.Header.Get("X-User-ID")
			username := r.Header.Get("X-Username")

			// 如果没有 tenant_id，使用默认值（开源版单租户模式）
			var tenantID, userID int64
			if tenantIDStr != "" {
				// 解析 tenant_id
				tenantID = parseInt64(tenantIDStr)
			} else {
				// 开源版默认租户
				tenantID = 1
			}

			if userIDStr != "" {
				userID = parseInt64(userIDStr)
			}

			// 创建租户上下文
			tenantCtx := &appctx.TenantContext{
				TenantID: tenantID,
				UserID:   userID,
				Username: username,
			}

			// 注入到 Context
			ctx = appctx.WithTenant(ctx, tenantCtx)

			// 记录日志
			log := logger.Default().WithContext(ctx)
			log.Debug("Tenant context injected",
				logger.Default().WithTenant(tenantID).Logger,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TraceMiddleware 追踪上下文中间件
func TraceMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 从 Header 中提取 TraceID
			traceID := r.Header.Get("X-Trace-ID")
			if traceID == "" {
				// 生成新的 TraceID
				traceID = generateTraceID()
			}

			// 创建追踪上下文
			traceCtx := &appctx.TraceContext{
				TraceID: traceID,
			}

			// 注入到 Context
			ctx = appctx.WithTrace(ctx, traceCtx)

			// 设置响应 Header
			w.Header().Set("X-Trace-ID", traceID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// parseInt64 解析 int64
func parseInt64(s string) int64 {
	var result int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int64(c-'0')
		}
	}
	return result
}

// generateTraceID 生成 TraceID
func generateTraceID() string {
	// 简化实现，实际应该使用 UUID 或 OTel 生成
	return "trace-" + strings.Replace("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "x", "0", -1)
}
