# Proto 接口定义

OpenGEO 平台的 Protocol Buffers 接口定义，用于服务间 RPC 通信和对外 API 暴露。

## 目录结构

```
proto/
├── buf.gen.yaml                    # Buf 代码生成配置
├── buf.yaml                        # Buf 项目配置
└── opengeo/
    ├── common/v1/                  # 公共类型定义
    │   ├── pagination.proto        # 分页请求/响应
    │   ├── tenant_context.proto    # 租户上下文
    │   └── trace.proto             # 链路追踪
    │
    ├── internal/v1/                # 内部 RPC 服务定义（Kitex）
    │   ├── account.proto           # 账号服务
    │   ├── content.proto           # 内容服务
    │   ├── publish.proto           # 发布服务
    │   ├── scheduler.proto         # 调度服务
    │   ├── monitor.proto           # 监测服务
    │   ├── system.proto            # 系统服务
    │   └── events/                 # Avro 领域事件 Schema
    │       ├── content_optimized.avsc
    │       ├── publish_requested.avsc
    │       └── publish_success.avsc
    │
    ├── cloud/v1/                   # 云端 API（套餐/订阅/API Key）
    ├── brand/v1/                   # 品牌管理 API
    ├── publish/v1/                 # 发布 API
    ├── monitor/v1/                 # 监测 API
    └── tenant/v1/                  # 租户管理 API
```

## 命名约定

| 目录 | 用途 | 框架 | 说明 |
|------|------|------|------|
| `internal/v1/` | 服务间 RPC | Kitex | 微服务内部通信，高性能 |
| 其他目录 | 对外 API | gRPC + REST Gateway | 面向客户端/第三方 |

## 环境准备

### 安装 Buf

```bash
# macOS
brew install bufbuild/buf/buf

# Linux
curl -sSL https://github.com/bufbuild/buf/releases/latest/download/buf-Linux-x86_64 \
  -o /usr/local/bin/buf && chmod +x /usr/local/bin/buf

# 验证
buf --version
```

### 安装 Kitex 代码生成工具（用于 internal RPC）

```bash
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```

## 代码生成

### 使用 Buf 生成（推荐）

```bash
cd proto/

# 生成所有代码
buf generate

# 生成指定模块
buf generate --path opengeo/cloud/v1
buf generate --path opengeo/internal/v1
```

生成产物位置：
- Go 代码 → `gen/go/opengeo/...`
- TypeScript 类型 → `web/src/types/gen/opengeo/...`

### 使用 Kitex 生成内部 RPC（可选）

```bash
# 生成 Content 服务的 Kitex 代码
kitex -module opengeo -service content \
  -gen-path ../gen/kitex \
  proto/opengeo/internal/v1/content.proto

# 生成 Publish 服务的 Kitex 代码
kitex -module opengeo -service publish \
  -gen-path ../gen/kitex \
  proto/opengeo/internal/v1/publish.proto
```

## 服务清单

### 内部 RPC 服务 (`internal/v1/`)

| 服务 | Proto 文件 | 功能 |
|------|-----------|------|
| AccountService | `account.proto` | 用户认证、RBAC 权限、租户管理 |
| ContentService | `content.proto` | 内容 CRUD、版本管理、AI 优化、模板、知识图谱 |
| PublishService | `publish.proto` | 发布任务、渠道管理、平台适配 |
| SchedulerService | `scheduler.proto` | 定时调度、错峰策略、热度图 |
| MonitorService | `monitor.proto` | AI 引用监测、来源评分、竞品分析、ROI |
| SystemService | `system.proto` | 系统配置、插件管理、Webhook、审计日志 |

### 对外 API 服务

| 服务 | Proto 文件 | 功能 |
|------|-----------|------|
| CloudAPIService | `cloud/v1/cloud_api.proto` | 套餐订阅、API Key、用量计费 |
| BrandService | `brand/v1/brand.proto` | 品牌管理 |
| KnowledgeService | `brand/v1/knowledge.proto` | 知识图谱实体 |
| PublishAPIService | `publish/v1/publish.proto` | 发布 API |
| MonitorAPIService | `monitor/v1/monitor.proto` | 监测 API |
| TenantService | `tenant/v1/tenant.proto` | 租户 API |

### 领域事件 (`internal/v1/events/`)

| 事件 | Schema | 触发时机 |
|------|--------|---------|
| ContentOptimized | `content_optimized.avsc` | 内容 AI 优化完成 |
| PublishRequested | `publish_requested.avsc` | 发布任务创建 |
| PublishSuccess | `publish_success.avsc` | 发布成功 |

## 开发流程

### 1. 修改 Proto 定义

```bash
# 编辑对应 proto 文件
vim proto/opengeo/internal/v1/content.proto
```

### 2. Lint 检查

```bash
cd proto/
buf lint
```

### 3. Breaking 变更检测

```bash
buf breaking --against '.git#branch=main'
```

### 4. 生成代码

```bash
buf generate
```

### 5. 提交

```bash
git add proto/ gen/
git commit -m "feat(content): add batch optimize RPC"
```

## Go 项目中使用

### 作为 gRPC 客户端

```go
import (
    cloudv1 "opengeo/gen/go/opengeo/cloud/v1"
    "google.golang.org/grpc"
)

conn, _ := grpc.Dial("localhost:9090", grpc.WithInsecure())
client := cloudv1.NewCloudAPIServiceClient(conn)

resp, _ := client.ListPlans(ctx, &cloudv1.ListPlansRequest{})
```

### 作为 Kitex 服务端

```go
import (
    content "opengeo/gen/kitex/opengeo/internal/v1/content"
)

type ContentServiceImpl struct{}

func (s *ContentServiceImpl) CreateContent(ctx context.Context, req *content.CreateContentRequest) (*content.Content, error) {
    // 实现逻辑
}

func main() {
    svr := content.NewServer(&ContentServiceImpl{}, server.WithServiceAddr(":8081"))
    svr.Run()
}
```

## 前端 TypeScript 类型

Buf 会自动生成 TypeScript 类型到 `web/src/types/gen/`：

```typescript
import { CloudAPIService } from '@/types/gen/opengeo/cloud/v1/cloud_api'

// 类型安全的 API 调用
const plans: CloudAPIService.ListPlansRequest = {}
```

## 注意事项

1. **版本管理**：所有 proto 文件使用 `v1` 版本号，破坏性变更需创建 `v2`
2. **向后兼容**：新增字段使用新的字段号，不要修改或删除已有字段
3. **命名规范**：服务名用 PascalCase，字段名用 snake_case
4. **注释**：所有 RPC 方法和消息字段必须添加中文注释
5. **公共类型**：分页、租户上下文等复用 `common/v1/` 下的定义
