package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestTotal 请求总量
	RequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "opengeo_request_total",
			Help: "Total number of requests",
		},
		[]string{"service", "method", "status"},
	)

	// RequestDuration 请求延迟
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "opengeo_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method"},
	)

	// AIGCFirstTokenLatency AIGC 首 Token 延迟
	AIGCFirstTokenLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "opengeo_aigc_first_token_seconds",
			Help:    "AIGC first token latency",
			Buckets: []float64{0.1, 0.2, 0.3, 0.5, 1.0, 2.0},
		},
	)

	// DiagnosticPluginDuration 诊断插件执行时间
	DiagnosticPluginDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "opengeo_diagnostic_plugin_seconds",
			Help:    "Diagnostic plugin execution time",
			Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.5},
		},
		[]string{"plugin_name"},
	)

	// APIUsageTotal API 使用总量
	APIUsageTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "opengeo_api_usage_total",
			Help: "Total API usage by tenant",
		},
		[]string{"tenant_id", "api_type"},
	)

	// TenantQuotaExceeded 租户配额超限次数
	TenantQuotaExceeded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "opengeo_tenant_quota_exceeded_total",
			Help: "Number of times tenant quota was exceeded",
		},
		[]string{"tenant_id"},
	)

	// BrandCount 品牌数量
	BrandCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "opengeo_brand_count",
			Help: "Number of brands per tenant",
		},
		[]string{"tenant_id"},
	)

	// ContentCount 内容数量
	ContentCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "opengeo_content_count",
			Help: "Number of contents per tenant",
		},
		[]string{"tenant_id", "brand_id"},
	)

	// CitationCount 引用数量
	CitationCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "opengeo_citation_total",
			Help: "Total citations tracked",
		},
		[]string{"tenant_id", "brand_id", "ai_model"},
	)
)

// IncRequest 增加请求计数
func IncRequest(service, method, status string) {
	RequestTotal.WithLabelValues(service, method, status).Inc()
}

// ObserveRequestDuration 观察请求延迟
func ObserveRequestDuration(service, method string, duration float64) {
	RequestDuration.WithLabelValues(service, method).Observe(duration)
}

// IncAPIUsage 增加 API 使用计数
func IncAPIUsage(tenantID, apiType string) {
	APIUsageTotal.WithLabelValues(tenantID, apiType).Inc()
}

// IncQuotaExceeded 增加配额超限计数
func IncQuotaExceeded(tenantID string) {
	TenantQuotaExceeded.WithLabelValues(tenantID).Inc()
}
