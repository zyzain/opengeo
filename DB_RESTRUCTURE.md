# OpenGEO BrandOS 数据库设计文档

## 文档信息

- 版本：v6.0.0
- 日期：2026-05-28
- 状态：已执行
- 架构：云端 API + SaaS 商业模式

---

## 一、设计原则

### 1.1 核心理念

根据 PRD v6.0 的定义，OpenGEO BrandOS 采用 **"云端 API + SaaS"** 商业模式：

| 组件 | 定位 | 说明 |
|------|------|------|
| **租户管理** | 开源基座 | 所有版本原生支持多租户，降低获客门槛 |
| **品牌治理** | 开源基座 | 完整的品牌元数据、知识图谱、术语表 |
| **AI 能力** | 云端计量 | 归因分析、可信度评分、合规校验等高价值能力 |

### 1.2 商业化锚点

| 维度 | 开源基座（免费） | 云端增强（付费） |
|------|------------------|------------------|
| 租户管理 | 完整支持 | 跨租户聚合分析 |
| 品牌治理 | 完整 CRUD + 知识图谱 | AI 驱动的洞察/评分 |
| 内容生成 | 集成主流 LLM | 智能优化建议 |
| 诊断引擎 | 规则型指标 | AI 引用归因分析 |
| 合规校验 | 本地静态规则 | 实时规则库更新 |

### 1.3 数据隔离策略

```
┌─────────────────────────────────────────────────────────────┐
│                    框架层原生支持                              │
├─────────────────────────────────────────────────────────────┤
│  1. Context 传递 tenant_id                                   │
│  2. Repository 自动注入 WHERE tenant_id = ?                  │
│  3. PostgreSQL RLS 模板（DB层安全兜底）                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 二、表结构概览

### 2.1 总览（58 张表）

| 模块 | 表数量 | 核心表 |
|------|--------|--------|
| **租户服务** | 3 | tenants, tenant_api_usage |
| **账号服务** | 8 | users, roles, permissions, user_roles, role_permissions, accounts, account_groups, account_group_relations |
| **品牌服务** | 7 | brands, brand_metadata, glossary_entries, brand_snapshots, brand_knowledge_entities, knowledge_relations, brand_trust_scores |
| **内容服务** | 6 | contents, content_versions, content_templates, content_entities, compliance_checks |
| **发布服务** | 5 | platforms, publish_channels, publish_tasks, fallback_queues |
| **调度服务** | 6 | schedules, schedule_tasks, ai_activity_heatmaps, publish_calendars, schedule_logs |
| **监测服务** | 9 | ai_citations, citation_attributions, source_scores, competitor_monitors, competitor_analyses, roi_metrics, optimization_suggestions, citation_trends |
| **系统服务** | 8 | system_configs, plugins, webhooks, webhook_events, translations, audit_logs, notifications |
| **云端API计量** | 6 | api_plans, tenant_subscriptions, api_billing_records, api_request_logs, api_offline_cache |

### 2.2 ER 关系图（核心）

```
┌─────────────┐
│   tenants   │─────────────────────────────────────────────┐
└─────────────┘                                             │
      │                                                     │
      │ 1:N                                                 │ 1:N
      ▼                                                     ▼
┌─────────────┐       ┌─────────────────┐       ┌─────────────────┐
│    users    │───M:N─│     brands      │───1:N─│  brand_metadata │
└─────────────┘       └─────────────────┘       └─────────────────┘
      │                       │
      │                 1:N   │   1:N
      │             ┌─────────┴─────────┐
      │             ▼                   ▼
      │   ┌─────────────────┐   ┌─────────────────┐
      │   │ glossary_entries│   │brand_knowledge_ │
      │   │                 │   │   entities      │
      │   └─────────────────┘   └─────────────────┘
      │                                   │
      │                             N:N   │
      │                                   ▼
      │                           ┌─────────────────┐
      │                           │knowledge_       │
      │                           │  relations      │
      │                           └─────────────────┘
      │
      ▼
┌─────────────┐       ┌─────────────────┐
│  accounts   │       │    contents     │
└─────────────┘       └─────────────────┘
                              │
                    ┌─────────┼─────────┐
                    ▼         ▼         ▼
            ┌───────────┐ ┌─────────┐ ┌─────────────┐
            │ai_citations│ │roi_     │ │optimization_│
            │           │ │metrics  │ │suggestions  │
            └───────────┘ └─────────┘ └─────────────┘
                    │
                    ▼
            ┌─────────────────┐
            │citation_        │
            │attributions     │
            │(云端API)        │
            └─────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    云端API计量服务                             │
├─────────────────────────────────────────────────────────────┤
│  api_plans → tenant_subscriptions → api_billing_records     │
│                                    → api_request_logs       │
│                                    → api_offline_cache      │
└─────────────────────────────────────────────────────────────┘
```

---

## 三、核心表结构详解

### 3.1 租户表（tenants）

```sql
CREATE TABLE IF NOT EXISTS tenants (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(128) NOT NULL,
    slug            VARCHAR(64)  NOT NULL,
    domain          VARCHAR(256) DEFAULT NULL,
    plan            VARCHAR(32)  DEFAULT 'free',
    status          TINYINT      DEFAULT 1,
    brand_limit     INT          DEFAULT 5,
    user_limit      INT          DEFAULT 10,
    storage_limit   BIGINT       DEFAULT 1073741824,
    api_quota       INT          DEFAULT 1000,
    api_used        INT          DEFAULT 0,
    settings        JSON,
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_slug (slug)
);
```

**设计说明**：
- `slug` 用于 URL 友好标识
- `brand_limit` / `user_limit` 控制资源配额
- `api_quota` / `api_used` 跟踪云端 API 用量
- `settings` 存储租户级配置

### 3.2 品牌表（brands）

```sql
CREATE TABLE IF NOT EXISTS brands (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT,
    tenant_id       BIGINT       NOT NULL,
    name            VARCHAR(128) NOT NULL,
    slug            VARCHAR(64)  NOT NULL,
    description     TEXT,
    logo_url        VARCHAR(512),
    website         VARCHAR(256),
    industry        VARCHAR(64),
    status          TINYINT      DEFAULT 1,
    settings        JSON,
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_tenant_slug (tenant_id, slug),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);
```

**设计说明**：
- `tenant_id` 为必填字段，确保数据隔离
- `slug` 在租户内唯一，用于 URL
- `settings` 存储品牌级配置

### 3.3 云端API计量表

#### 套餐表（api_plans）

```sql
CREATE TABLE IF NOT EXISTS api_plans (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(64)  NOT NULL,
    code            VARCHAR(32)  NOT NULL,
    monthly_quota   INT          NOT NULL,
    price_cents     INT          NOT NULL,
    overage_price   INT          DEFAULT 0,
    features        JSON,
    is_active       BOOLEAN      DEFAULT TRUE,
    UNIQUE INDEX idx_code (code)
);
```

#### 订阅表（tenant_subscriptions）

```sql
CREATE TABLE IF NOT EXISTS tenant_subscriptions (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT,
    tenant_id       BIGINT       NOT NULL,
    plan_id         BIGINT       NOT NULL,
    status          VARCHAR(32)  DEFAULT 'active',
    starts_at       DATETIME     NOT NULL,
    ends_at         DATETIME,
    api_key         VARCHAR(64),
    api_secret      VARCHAR(128),
    INDEX idx_tenant (tenant_id),
    INDEX idx_api_key (api_key)
);
```

#### 计费记录表（api_billing_records）

```sql
CREATE TABLE IF NOT EXISTS api_billing_records (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT,
    tenant_id       BIGINT       NOT NULL,
    subscription_id BIGINT       NOT NULL,
    api_type        VARCHAR(64)  NOT NULL,
    endpoint        VARCHAR(256) NOT NULL,
    request_id      VARCHAR(64)  NOT NULL,
    tokens_input    INT          DEFAULT 0,
    tokens_output   INT          DEFAULT 0,
    cost_cents      INT          DEFAULT 0,
    is_overage      BOOLEAN      DEFAULT FALSE,
    INDEX idx_tenant (tenant_id),
    INDEX idx_created (created_at)
);
```

---

## 四、云端API类型定义

### 4.1 API 类型清单

| API 类型 | 说明 | 计费单位 | 开源版 | 云端版 |
|----------|------|----------|--------|--------|
| `attribution` | AI引用归因分析 | 次 | ✗ | ✓ |
| `attribution_basic` | 基础归因（本地缓存） | 次 | ✓ | ✓ |
| `trust_score` | 品牌可信度评分 | 次 | ✗ | ✓ |
| `compliance` | 实时合规校验 | 次 | ✗ | ✓ |
| `compliance_basic` | 基础合规（本地规则） | 次 | ✓ | ✓ |
| `optimization` | 智能内容优化建议 | Token | ✗ | ✓ |
| `bulk_analysis` | 批量分析 | 次 | ✗ | ✓ |

### 4.2 离线降级策略

```go
// 离线降级逻辑
func (c *CloudAPIClient) CallWithFallback(ctx context.Context, req *APIRequest) (*APIResponse, error) {
    // 1. 尝试调用云端API
    resp, err := c.callCloud(ctx, req)
    if err == nil {
        return resp, nil
    }
    
    // 2. 检查离线缓存
    cached, ok := c.cache.Get(req.CacheKey())
    if ok {
        return cached.(*APIResponse), nil
    }
    
    // 3. 降级到本地轻量版
    return c.fallbackLocal(ctx, req)
}
```

---

## 五、数据隔离实现

### 5.1 Context 传递

```go
// pkg/context/tenant.go
package context

type TenantContext struct {
    TenantID int64
    UserID   int64
    BrandID  int64 // 可选
}

func WithTenant(ctx context.Context, tc *TenantContext) context.Context {
    return context.WithValue(ctx, tenantKey, tc)
}

func GetTenant(ctx context.Context) *TenantContext {
    tc, _ := ctx.Value(tenantKey).(*TenantContext)
    return tc
}
```

### 5.2 Repository 自动注入

```go
// pkg/repository/base.go
type BaseRepository struct {
    db *gorm.DB
}

func (r *BaseRepository) WithTenant(ctx context.Context) *gorm.DB {
    tc := context.GetTenant(ctx)
    if tc == nil {
        panic("tenant context not found")
    }
    return r.db.Where("tenant_id = ?", tc.TenantID)
}

func (r *BaseRepository) FindByID(ctx context.Context, id int64, dest interface{}) error {
    return r.WithTenant(ctx).First(dest, id).Error
}
```

### 5.3 PostgreSQL RLS 模板

```sql
-- 启用 RLS
ALTER TABLE brands ENABLE ROW LEVEL SECURITY;

-- 创建策略
CREATE POLICY tenant_isolation ON brands
    USING (tenant_id = current_setting('app.tenant_id')::bigint);

-- 设置会话变量
SET app.tenant_id = '123';
```

---

## 六、初始数据

### 6.1 默认套餐

| 套餐 | 代码 | 月配额 | 价格 | 包含API |
|------|------|--------|------|---------|
| 免费版 | free | 100 | 0 | attribution_basic, compliance_basic |
| 基础版 | starter | 1,000 | ¥99 | attribution, trust_score, compliance |
| 专业版 | pro | 10,000 | ¥499 | + optimization |
| 企业版 | enterprise | 100,000 | ¥1,999 | + bulk_analysis, priority_support |

### 6.2 默认角色

| 角色 | 标识 | 说明 |
|------|------|------|
| 租户管理员 | tenant_admin | 拥有租户内全部权限 |
| 品牌编辑者 | brand_editor | 可管理品牌和内容 |
| 只读用户 | viewer | 仅查看权限 |

### 6.3 默认插件

| 插件 | 类型 | 说明 |
|------|------|------|
| source_score | diagnostic | 信源评分 |
| competitor_monitor | diagnostic | 竞品监测 |
| roi_analyzer | diagnostic | ROI分析 |
| semantic_relevance | diagnostic | 语义相关性 |
| citation_authority | diagnostic | 引用权威性 |
| wechat_adapter | channel | 微信适配器 |
| weibo_adapter | channel | 微博适配器 |
| douyin_adapter | channel | 抖音适配器 |
| xiaohongshu_adapter | channel | 小红书适配器 |
| zhihu_adapter | channel | 知乎适配器 |

---

## 七、迁移策略

### 7.1 从 v3.0 迁移到 v6.0

```sql
-- 1. 备份数据库
mysqldump --single-transaction opengeo > backup_v3_$(date +%Y%m%d%H%M%S).sql

-- 2. 创建租户表
CREATE TABLE IF NOT EXISTS tenants (...);

-- 3. 创建默认租户
INSERT INTO tenants (id, name, slug, plan) VALUES (1, '默认租户', 'default', 'free');

-- 4. 为相关表添加 tenant_id 列
ALTER TABLE users ADD COLUMN tenant_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE brands ADD COLUMN tenant_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE contents ADD COLUMN tenant_id BIGINT NOT NULL DEFAULT 1;
-- ... 其他表

-- 5. 创建云端API计量表
CREATE TABLE IF NOT EXISTS api_plans (...);
CREATE TABLE IF NOT EXISTS tenant_subscriptions (...);
CREATE TABLE IF NOT EXISTS api_billing_records (...);
-- ... 其他表

-- 6. 更新索引
ALTER TABLE users DROP INDEX idx_username;
ALTER TABLE users ADD UNIQUE INDEX idx_tenant_username (tenant_id, username);
-- ... 其他索引
```

---

## 八、验收标准

| 维度 | 指标 | 验证方式 |
|------|------|----------|
| 数据隔离 | 跨租户查询 0 容忍 | 安全扫描 + 渗透测试 |
| 配额控制 | 超额调用正确拒绝 | 单元测试 + 集成测试 |
| 离线降级 | 网络不可用时业务不中断 | 模拟网络故障测试 |
| 计费准确 | Token 计数误差 < 1% | 对比测试 |
| API 稳定性 | 升级后 API 不变 | 契约测试 |

---

## 九、表数量汇总

| 版本 | 表数量 | 说明 |
|------|--------|------|
| v3.0 | 46 | 移除租户，插件化架构 |
| v6.0 | 58 | 租户回归基座，新增云端API计量 |

### v6.0 新增表（12 张）

| 表名 | 模块 | 说明 |
|------|------|------|
| `tenants` | 租户 | 租户主表 |
| `tenant_api_usage` | 租户 | API用量明细 |
| `brand_trust_scores` | 品牌 | 可信度评分 |
| `compliance_checks` | 内容 | 合规校验记录 |
| `citation_attributions` | 监测 | 引用归因分析 |
| `api_plans` | 计量 | API套餐 |
| `tenant_subscriptions` | 计量 | 租户订阅 |
| `api_billing_records` | 计量 | 计费记录 |
| `api_request_logs` | 计量 | 请求日志 |
| `api_offline_cache` | 计量 | 离线缓存 |
| `notifications` | 系统 | 通知表 |
