# OpenGEO 智能发布平台

> AI 时代的 GEO（Generative Engine Optimization）内容优化与多平台智能发布系统

**[English](README.en.md)** | 中文

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)](https://react.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.7-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white)](https://www.docker.com)

---

## 功能特性

### GEO 内容优化
- AI 语义增强（结构 / 可读性 / 意图 / Schema Markup 四维评分）
- 多模型适配（DeepSeek / Kimi / 豆包 / ChatGPT）
- 中文分词关键词提取（字典最长匹配 + bigram 回退）
- 合规检测（敏感词 / 广告法 / AIGC 标识）
- 知识图谱实体管理，提升 AI 引用权威性

### 多平台发布
- 微信 / 微博 / 抖音 / 小红书 / 知乎 / 头条适配
- 去重引擎（同义词替换 / 段落重排 / 句式变换）
- 错峰调度（Worker Pool + 交错延迟 + 热力图）
- 重试 / 降级 / 人工审核三级容错

### 防封号引擎
- 代理池管理（加载 / 健康检测 / 失败自动剔除）
- 浏览器指纹绑定（每账号固定指纹，最少使用分配）
- 行为模拟（正态分布延迟 / 打字模拟 / 滚动模拟）
- 平台级限流（每小时 / 每日 / 最小间隔）

### 账号与权限
- 多租户隔离
- RBAC 权限控制（已接入全部路由，admin 自动放行）
- 资源归属校验（IDOR 防护）

### 监测分析
- AI 引用追踪（引用位置 / 情感 / 模型来源）
- 信源评分与竞品对比
- ROI 归因分析

---

## 技术栈

| 层 | 技术 |
|---|------|
| **前端** | React 19 + Vite + TypeScript + Ant Design 5 + Zustand + TanStack Query |
| **网关** | Hertz (CloudWeGo) + RBAC 中间件 + RateLimiter |
| **服务** | Kitex RPC + 六边形架构 + CQRS + EDA |
| **存储** | MySQL 8.0 + Redis 7 |
| **监控** | Prometheus + Grafana + Jaeger (OpenTelemetry) |
| **部署** | Docker Compose + Nginx 负载均衡 |
| **Proto** | Buf + Protobuf + gRPC |

---

## 项目结构

```
opengeo/
├── gateway/                          # HTTP 网关（入口）
│   ├── cmd/main.go                   # 启动入口
│   └── internal/
│       ├── auth/                     # 认证服务（JWT + RBAC）
│       ├── handler/                  # HTTP 处理器（按领域拆分）
│       │   ├── handler.go            # 公共工具 + Handler 结构体
│       │   ├── auth_handler.go       # 登录 / 注册 / 刷新
│       │   ├── user_handler.go       # 用户 / 角色 / 租户
│       │   ├── content_handler.go    # 内容 CRUD + 优化
│       │   ├── account_handler.go    # 账号管理
│       │   ├── knowledge_handler.go  # 知识图谱
│       │   ├── publish_handler.go    # 发布 + 渠道 + 平台
│       │   ├── schedule_handler.go   # 调度
│       │   ├── monitor_handler.go    # 监测
│       │   └── system_handler.go     # 系统配置 + 插件 + Webhook
│       ├── client/                   # 下游客户端（按领域拆分）
│       ├── middleware/               # 中间件（CORS / JWT / RBAC / RateLimiter）
│       ├── router/                   # 路由注册
│       └── dal/                      # 数据访问层
│
├── service/                          # 微服务层
│   ├── account/                      # 账号服务（Kitex 分层）
│   ├── content/                      # 内容服务（六边形架构）
│   │   └── internal/
│   │       ├── domain/               # 领域模型 + GEO 优化逻辑
│   │       ├── application/          # 应用服务（用例编排）
│   │       ├── port/                 # 端口接口（inbound / outbound）
│   │       └── adapter/              # 适配器（数据库 / AI / 事件）
│   ├── publish/                      # 发布服务（六边形 + EDA）
│   │   └── internal/
│   │       ├── service/              # 核心服务（防封号 / 去重 / 重试 / Worker Pool）
│   │       ├── domain/               # 领域模型 + 事件
│   │       ├── adapter/              # 平台适配器 + Kafka
│   │       └── port/                 # 端口接口
│   ├── scheduler/                    # 调度服务（热力图 / 错峰 / 优先队列）
│   ├── monitor/                      # 监测服务（CQRS）
│   └── system/                       # 系统服务（插件 SDK / Webhook / 审计）
│
├── pkg/                              # 公共组件
│   ├── ai/                           # AI 服务接口（统一类型定义）
│   ├── similarity/                   # 文本相似度（SimHash / 余弦 / Jaccard）
│   ├── crypto/                       # 加密工具（bcrypt）
│   ├── jwt/                          # JWT 工具
│   ├── config/                       # 配置加载
│   ├── database/                     # 数据库连接
│   ├── errcode/                      # 错误码
│   └── eventbus/                     # 事件总线
│
├── proto/                            # Protobuf 接口定义
│   ├── opengeo/
│   │   ├── common/v1/                # 公共类型（分页 / 租户 / 链路追踪）
│   │   ├── internal/v1/              # 内部 RPC（account / content / publish / ...）
│   │   ├── cloud/v1/                 # 对外 API（套餐 / 订阅 / API Key）
│   │   ├── brand/v1/                 # 品牌 API
│   │   ├── publish/v1/               # 发布 API
│   │   ├── monitor/v1/               # 监测 API
│   │   └── tenant/v1/                # 租户 API
│   ├── buf.yaml                      # Buf 配置
│   └── buf.gen.yaml                  # 代码生成配置
│
├── web/                              # 前端（React 19 + Vite）
│   ├── src/
│   │   ├── pages/                    # 页面（30+ 页面覆盖所有功能）
│   │   ├── components/               # 组件
│   │   ├── hooks/                    # React Query Hooks
│   │   ├── stores/                   # Zustand 状态管理
│   │   ├── lib/                      # API 客户端（80+ 接口）
│   │   ├── types/                    # TypeScript 类型
│   │   └── i18n/                     # 国际化（zh-CN / en-US）
│   └── package.json
│
├── configs/                          # 配置文件
│   ├── config.example.json           # 应用配置示例
│   ├── prometheus.yml                # Prometheus 配置
│   └── nginx.conf                    # Nginx 负载均衡配置
│
├── scripts/                          # 脚本
│   └── init.sql                      # 数据库初始化
│
├── deployments/                      # 部署配置
│   └── Dockerfile.gateway            # Gateway Dockerfile
│
├── docker-compose.yml                # 开发环境编排
├── docker-compose.prod.yml           # 生产环境编排（多副本 + 滚动更新）
├── Makefile                          # 构建命令
└── go.mod                            # Go 模块定义
```

---

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- MySQL 8.0
- Redis 7

### 1. 启动基础设施

```bash
docker compose up -d mysql redis consul jaeger prometheus grafana
```

### 2. 配置环境变量

```bash
export MYSQL_DSN="root:root@tcp(127.0.0.1:3306)/opengeo?charset=utf8mb4&parseTime=True&loc=Local"
export ADMIN_PASSWORD="YourSecurePassword@123"  # 生产环境必须设置
export JWT_SECRET_KEY="your-secret-key"          # 生产环境必须设置
```

### 3. 启动后端

```bash
make dev-gateway
# 管理员: admin / <ADMIN_PASSWORD>
# API: http://localhost:8080
```

### 4. 启动前端

```bash
cd web && npm install && npm run dev
# 访问: http://localhost:3000
```

---

## Makefile 命令

```bash
# 基础设施
make dev-up              # 启动 MySQL / Redis / Consul / Jaeger / Prometheus / Grafana
make dev-down            # 停止所有容器

# 开发运行
make dev-gateway         # 启动 Gateway
make dev-content         # 启动 Content 服务
make dev-publish         # 启动 Publish 服务
make dev-scheduler       # 启动 Scheduler 服务
make dev-monitor         # 启动 Monitor 服务
make dev-system          # 启动 System 服务

# 构建
make build               # 构建所有服务
make build-gateway       # 构建 Gateway

# 测试
make test                # 运行所有测试
make test-unit           # 运行单元测试（22 个测试文件，14 个包通过）

# Proto
make proto-gen           # 生成 Go / TS 代码
make proto-lint          # Proto 规范检查

# 代码质量
make lint                # golangci-lint
make fmt                 # gofmt
make vet                 # go vet

# Docker
make docker              # Docker 全栈启动
make docker-down         # 停止 Docker
```

---

## 安全特性

| 特性 | 说明 |
|------|------|
| RBAC 权限 | 所有路由按 resource:action 校验，admin 自动放行 |
| IDOR 防护 | Content / Account / Entity / Task / Schedule 均校验 user_id 归属 |
| XSS 防护 | PreviewPublish 使用 html.EscapeString 转义 |
| SSRF 防护 | Webhook URL 禁止内网地址（localhost / 169.254 / 私有 IP） |
| RateLimiter | 全局限流 + 登录接口每 IP 5次/分钟 |
| 输入校验 | page_size 上限 100，错误信息统一 safeError 脱敏 |
| 密码安全 | bcrypt 加密，代理密码响应脱敏，Admin 密码从环境变量读取 |
| Prompt 注入 | AI 调用使用分隔符包裹用户内容 |

---

## 架构模式

| 服务 | 模式 | 说明 |
|------|------|------|
| Gateway | 标准分层 | router / handler / middleware / client |
| Content | 六边形 | domain / application / port / adapter |
| Publish | 六边形 + EDA | 事件驱动 + Worker Pool + 防封号引擎 |
| Monitor | CQRS | 读写分离（command / query） |
| Scheduler | 六边形 | 热力图 + 错峰 + 优先队列 |
| System | 注册表 | 插件 SDK + Webhook + 审计日志 |

---

## API 概览

### 认证
```
POST   /api/v1/auth/login          # 登录
POST   /api/v1/auth/register       # 注册
POST   /api/v1/auth/refresh        # 刷新 Token
```

### 内容管理
```
GET    /api/v1/contents             # 列表（分页）
POST   /api/v1/contents             # 创建
GET    /api/v1/contents/:id         # 详情
PUT    /api/v1/contents/:id         # 更新
DELETE /api/v1/contents/:id         # 删除
POST   /api/v1/contents/:id/optimize  # AI 优化
POST   /api/v1/contents/:id/publish   # 发布
```

### 发布管理
```
GET    /api/v1/publish/tasks        # 任务列表
POST   /api/v1/publish/tasks        # 创建任务
POST   /api/v1/publish/tasks/:id/cancel  # 取消
POST   /api/v1/publish/preview      # 预览
POST   /api/v1/publish/dedup/check  # 去重检测
```

### 知识图谱
```
GET    /api/v1/knowledge/entities    # 实体列表
POST   /api/v1/knowledge/entities    # 创建实体
GET    /api/v1/knowledge/entities/search  # 搜索
```

### 系统管理
```
GET    /api/v1/system/configs       # 系统配置
GET    /api/v1/system/plugins       # 插件列表
GET    /api/v1/system/webhooks      # Webhook 列表
```

完整 API 文档见 [API.md](API.md)

---

## 测试

```bash
# 运行所有测试
go test ./... -short -v

# 运行指定包测试
go test ./service/publish/internal/service/... -v
go test ./service/content/internal/domain/service/... -v
go test ./pkg/similarity/... -v
```

测试覆盖：
- 22 个测试文件
- 14 个包通过
- 包含单元测试 + Benchmark 性能基准

---

## 部署

### 开发环境

```bash
docker compose up -d          # 启动基础设施
make dev-gateway              # 启动 Gateway
cd web && npm run dev         # 启动前端
```

### 生产环境

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

生产配置特性：
- Gateway 3 副本，Content / Publish 2 副本
- 滚动更新（start-first，零停机）
- Nginx 负载均衡 + 限流
- 资源限制（CPU / 内存）

---

## 文档

| 文档 | 说明 |
|------|------|
| [README.md](README.md) | 本文档 |
| [API.md](API.md) | RESTful API 接口 |
| [proto/README.md](proto/README.md) | Protobuf 接口定义与代码生成 |

---

## License

MIT
