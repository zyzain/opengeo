# OpenGEO BrandOS 架构概览

## 核心定位

OpenGEO BrandOS 是一个 AI-Native 品牌治理开源基础设施，采用「云端 API + SaaS」商业模式。

## 架构原则

### 1. 六边形架构

所有服务采用六边形架构（端口与适配器模式），确保 Domain 层纯净：

```
┌─────────────────────────────────────┐
│           Application Layer         │
│  ┌───────────────────────────────┐  │
│  │        Domain Layer           │  │
│  │  (纯净，无 infra 依赖)        │  │
│  └───────────────────────────────┘  │
│           Port Layer                │
│  ┌───────────────────────────────┐  │
│  │      Adapter Layer            │  │
│  │  (MySQL/Redis/HTTP/gRPC)      │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

### 2. Proto First

所有 API 必须先写 `.proto` 文件，CI 自动生成 Go/TS 代码：

```
proto/
├── opengeo/
│   ├── common/v1/     # 通用消息体
│   ├── tenant/v1/     # 租户服务
│   ├── brand/v1/      # 品牌服务
│   ├── content/v1/    # 内容服务
│   ├── publish/v1/    # 发布服务
│   ├── monitor/v1/    # 监测服务
│   └── cloud/v1/      # 云端 API
```

### 3. 租户隔离

所有数据表包含 `tenant_id` 字段，通过 Middleware 自动注入：

```go
// 获取租户上下文
tenantID := context.GetTenantID(ctx)

// Repository 自动注入租户条件
func (r *BrandRepository) FindByID(ctx context.Context, id int64) (*Brand, error) {
    tenantID := context.GetTenantID(ctx)
    return r.db.Where("tenant_id = ? AND id = ?", tenantID, id).First(&brand).Error
}
```

### 4. 插件化架构

诊断引擎和渠道适配器通过注册机制接入，核心代码零侵入：

```go
// 注册诊断插件
plugin.RegisterDiagnosticPlugin(NewSourceScorePlugin())

// 注册渠道适配器
plugin.RegisterChannelAdapter(NewWechatAdapter())
```

## 服务架构

```
┌─────────────────────────────────────────────────────────────┐
│                        Gateway (Hertz)                       │
│              路由 / 中间件 / 鉴权 / 限流 / 日志              │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│ Tenant Svc   │    │  Brand Svc   │    │ Content Svc  │
│   (Kitex)    │    │   (Kitex)    │    │   (Kitex)    │
└──────────────┘    └──────────────┘    └──────────────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│ Publish Svc  │    │ Monitor Svc  │    │  System Svc  │
│   (Kitex)    │    │   (Kitex)    │    │   (Kitex)    │
└──────────────┘    └──────────────┘    └──────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │   Cloud API      │
                    │ (计量 + 计费)    │
                    └──────────────────┘
```

## 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 前端 | React + Vite + Ant Design | SPA + 组件库 |
| 网关 | Hertz (CloudWeGo) | 高性能 HTTP 框架 |
| 服务 | Kitex (CloudWeGo) | 高性能 RPC 框架 |
| 数据库 | MySQL + Redis | 关系型 + 缓存 |
| 向量库 | Milvus | 知识图谱向量索引 |
| 监控 | Prometheus + Grafana + Jaeger | 指标 + 面板 + 追踪 |
| 部署 | Docker + Kubernetes | 容器化部署 |

## 商业模式

| 组件 | 定位 | 说明 |
|------|------|------|
| 租户管理 | 开源基座 | 所有版本原生支持多租户 |
| 品牌治理 | 开源基座 | 完整的品牌元数据、知识图谱、术语表 |
| AI 能力 | 云端计量 | 归因分析、可信度评分、合规校验 |
