package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"

	jwtUtil "opengeo/pkg/jwt"
)

// ==================== CORS ====================

func CORS() app.HandlerFunc {
	allowedOrigins := getAllowedOrigins()

	return func(ctx context.Context, c *app.RequestContext) {
		origin := string(c.GetHeader("Origin"))

		if isOriginAllowed(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if string(c.Request.Method()) == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next(ctx)
	}
}

func getAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		// 开发环境默认允许 localhost
		if os.Getenv("GO_ENV") != "production" {
			return []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080"}
		}
		return []string{}
	}
	return strings.Split(origins, ",")
}

func isOriginAllowed(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}
	for _, a := range allowed {
		if strings.TrimSpace(a) == origin {
			return true
		}
	}
	return false
}

// ==================== Logger ====================

func Logger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.Path())
		method := string(c.Request.Method())

		c.Next(ctx)

		latency := time.Since(start)
		status := c.Response.StatusCode()

		fmt.Printf("[%s] %s %s %d %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			status,
			latency,
		)
	}
}

// ==================== Metrics ====================

func Metrics() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		_ = time.Since(start)
		_ = c.Response.StatusCode()
	}
}

// ==================== Recovery ====================

func Recovery() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				fmt.Printf("[PANIC] %v\n%s\n", err, stack)
				c.JSON(http.StatusInternalServerError, utils.H{
					"code":    500,
					"message": "internal server error",
					"data":    nil,
				})
				c.Abort()
			}
		}()
		c.Next(ctx)
	}
}

// ==================== JWT ====================

func JWT() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.H{
				"code":    401,
				"message": "missing authorization header",
				"data":    nil,
			})
			c.Abort()
			return
		}

		tokenString, err := jwtUtil.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.H{
				"code":    401,
				"message": "invalid authorization format",
				"data":    nil,
			})
			c.Abort()
			return
		}

		claims, err := jwtUtil.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.H{
				"code":    401,
				"message": "invalid or expired token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("tenant_id", claims.TenantID)

		c.Next(ctx)
	}
}

// ==================== Rate Limiter (Token Bucket) ====================

type tokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

func newTokenBucket(maxTokens, refillRate float64) *tokenBucket {
	return &tokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}
	tb.lastRefill = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

var (
	globalLimiter  = newTokenBucket(100, 10) // 100 burst, 10/sec refill
	loginLimiters  = make(map[string]*tokenBucket)
	loginMu        sync.Mutex
)

func getLoginLimiter(ip string) *tokenBucket {
	loginMu.Lock()
	defer loginMu.Unlock()

	limiter, exists := loginLimiters[ip]
	if !exists {
		// 每IP: 5次/分钟
		limiter = newTokenBucket(5, 0.08)
		loginLimiters[ip] = limiter
	}
	return limiter
}

func RateLimiter() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.Path())

		// 登录接口更严格
		if strings.HasSuffix(path, "/auth/login") || strings.HasSuffix(path, "/auth/register") {
			ip := c.ClientIP()
			limiter := getLoginLimiter(ip)
			if !limiter.allow() {
				c.JSON(http.StatusTooManyRequests, utils.H{
					"code":    429,
					"message": "rate limit exceeded, please try again later",
					"data":    nil,
				})
				c.Abort()
				return
			}
		}

		// 全局限流
		if !globalLimiter.allow() {
			c.JSON(http.StatusTooManyRequests, utils.H{
				"code":    429,
				"message": "server is busy, please try again later",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// ==================== RBAC Permission Check ====================

// PermissionChecker 权限检查接口
type PermissionChecker interface {
	CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error)
}

// RequirePermission 返回一个中间件，检查当前用户是否拥有指定权限
// admin 角色自动放行
func RequirePermission(checker PermissionChecker, resource, action string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// admin 角色自动放行
		role, _ := c.Get("role")
		if roleStr, ok := role.(string); ok && roleStr == "admin" {
			c.Next(ctx)
			return
		}

		// 获取当前用户ID
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, utils.H{
				"code":    401,
				"message": "unauthorized",
				"data":    nil,
			})
			c.Abort()
			return
		}

		userID, ok := userIDVal.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, utils.H{
				"code":    500,
				"message": "invalid user context",
				"data":    nil,
			})
			c.Abort()
			return
		}

		allowed, err := checker.CheckPermission(ctx, userID, resource, action)
		if err != nil {
			fmt.Printf("[RBAC] permission check error: user=%d resource=%s action=%s err=%v\n", userID, resource, action, err)
			c.JSON(http.StatusInternalServerError, utils.H{
				"code":    500,
				"message": "permission check failed",
				"data":    nil,
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, utils.H{
				"code":    403,
				"message": fmt.Sprintf("permission denied: %s:%s", resource, action),
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// ==================== Request ID ====================

func RequestID() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		requestID := string(c.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next(ctx)
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + cryptoRandomString(8)
}

func cryptoRandomString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
