# 🚀 OpenGEO BrandOS 开发者启动清单 (Dev Kickoff Checklist)

> 本文档将 PRD v6.0 转化为可执行的工程任务
> 商业模式：云端 API + SaaS（租户与品牌管理为开源基座，AI 能力为云端计量服务）

---

## ⚠️ 开发者必读约束

| 约束项 | 具体要求 | 违规后果 |
|--------|----------|----------|
| **Proto First** | 所有 API 必须先写 `.proto`，禁止手写 TS/Go 类型 | CI 阻断，PR 无法合并 |
| **Domain 纯净** | Kitex Service 的 Domain 层禁止 import 任何 infra/db/http 包 | Code Review 打回重构 |
| **Trace 必传** | 所有跨服务调用必须传递 TraceID，日志必须含 span context | 可观测性验收不通过 |
| **Mock 先行** | 前端开发必须先有 MSW Mock，再联调真实接口 | 阻塞前端进度视为后端责任 |
| **插件零侵入** | 新增诊断/渠道功能不得修改核心代码，仅通过注册机制接入 | 架构评审否决 |
| **租户安全** | tenant_id 必须由 Middleware 注入，严禁前端传参 | 安全审计一票否决 |

---

## 🔴 P0：框架基座与品牌服务 (Week 1-6)

**目标**：完成 Proto First 重构，跑通「租户 → 品牌 → 内容 → 发布」最小闭环

### P0-1: Proto 定义层 ✅

- [x] 创建 `proto/opengeo/tenant/v1/tenant.proto`：Tenant CRUD + 配额管理
- [x] 创建 `proto/opengeo/brand/v1/brand.proto`：Brand CRUD + 元数据 + 术语表
- [x] 创建 `proto/opengeo/brand/v1/knowledge.proto`：知识图谱实体 + 关系
- [x] 创建 `proto/opengeo/common/v1/pagination.proto`：通用分页
- [x] 创建 `proto/opengeo/common/v1/trace.proto`：TraceID 传播
- [x] 创建 `proto/opengeo/common/v1/tenant_context.proto`：租户上下文
- [x] 创建 `proto/opengeo/publish/v1/publish.proto`：发布服务 + brand_id + tenant_id
- [x] 创建 `proto/opengeo/monitor/v1/monitor.proto`：监测服务 + brand_id + tenant_id
- [x] 创建 `proto/opengeo/cloud/v1/cloud_api.proto`：云端 API 计量服务

**约束检查点**：
- ✅ 所有消息体包含 tenant_id（由 Middleware 注入，前端不可见）
- ✅ 所有 RPC 入参包含 TraceContext（RequestMeta）
- ✅ Proto 文件结构符合 buf lint 规范

### P0-2: Buf 工作流配置 ✅

- [x] 创建 `buf.yaml`：模块定义 + 依赖管理
- [x] 创建 `buf.gen.yaml`：Go/TS 代码生成配置
- [x] 配置 `make proto-gen`：一键生成 Go/TS 代码
- [x] 配置 `make proto-lint`：Proto 规范检查
- [x] 配置 `make proto-breaking`：Breaking Change 检测
- [x] 配置 `make dev`：一键启动全栈环境
- [x] 创建 `scripts/mock-gen.sh`：MSW Mock 自动生成脚本
- [x] 创建 `.github/workflows/ci.yaml`：CI 自动检测 Breaking Change

**约束检查点**：
- ✅ Proto 变更必须通过 buf lint
- ✅ Breaking Change 必须有版本升级说明
- ✅ CI 自动检测 Breaking Change

### P0-3: 租户服务实现 ✅

- [x] 创建 `service/tenant/` 目录结构（六边形架构）
- [x] 实现 Domain 层：Tenant 实体 + TenantRepository 接口
- [x] 实现 Application 层：TenantService（CRUD + 配额校验）
- [x] 实现 Adapter 层：MySQL TenantRepository
- [x] 实现 Handler 层：TenantHandler（RPC 接口）
- [x] 实现服务入口：main.go

**约束检查点**：
- ✅ Domain 层不 import 任何 infra 包
- ✅ tenant_id 由 Middleware 注入，Controller 层无感知
- ✅ 所有 Repository 方法自动注入租户条件

### P0-4: 品牌服务实现 ✅

- [x] 创建 `service/brand/` 目录结构（六边形架构）
- [x] 实现 Domain 层：Brand 实体 + BrandMetadata + GlossaryEntry
- [x] 实现 Domain 层：BrandRepository + MetadataRepository + GlossaryRepository 接口
- [x] 实现 Application 层：BrandService（CRUD + 元数据管理）
- [x] 实现 Application 层：GlossaryService（术语表 CRUD + 批量导入）
- [x] 实现 Adapter 层：MySQL BrandRepository
- [x] 实现 Adapter 层：MySQL GlossaryRepository
- [x] 实现服务入口：main.go

**约束检查点**：
- ✅ Domain 层纯净，不依赖 infra
- ✅ 品牌操作自动继承租户上下文
- ✅ 术语表支持批量导入/导出

### P0-5: 内容服务重构 ✅

- [x] 更新 Content 实体：添加 brand_id + tenant_id 字段
- [x] 更新 ContentRepository：查询自动注入租户条件
- [x] 更新 ContentService：创建内容时关联品牌
- [x] 更新知识图谱实体：关联品牌（brand_knowledge_entities）

**约束检查点**：
- ✅ 内容必须归属某个品牌
- ✅ 跨品牌内容查询需明确授权

### P0-6: 可观测性基线 ✅

- [x] 实现 TraceID 传播：Hertz → Kitex → DB 全链路
- [x] 实现 Hertz Middleware：从 Header 提取/生成 TraceID
- [x] 实现结构化日志：JSON 格式含 trace_id/span_id/tenant_id
- [x] 配置 Prometheus Exporter：QPS/延迟/错误率

**约束检查点**：
- ✅ 所有日志包含 trace_id
- ✅ 跨服务调用必须传递 TraceID
- ✅ Metrics 包含 tenant_id 维度

### P0-7: 前端 MSW Mock ✅

- [x] 创建 `scripts/mock-gen.sh`：基于 Proto 生成 MSW Handler
- [x] 生成 Tenant Mock：CRUD + 配额查询
- [x] 生成 Brand Mock：CRUD + 元数据 + 术语表
- [x] 生成 Content Mock：列表 + 详情 + 创建
- [x] 配置 MSW Service Worker：拦截 API 请求
- [x] 编写 Mock 数据工厂：随机测试数据生成

**约束检查点**：
- ✅ Mock 基于 Proto 生成，保证类型一致
- ✅ 前端开发无需后端服务即可运行

### P0-8: 前端品牌页面 ✅

- [x] 创建品牌空间页面：品牌列表 + 创建/编辑表单
- [x] 实现品牌详情页面：品牌信息 + 元数据展示
- [x] 实现术语表管理：列表 + 新增 + 编辑 + 删除
- [x] 实现 React Query Hooks：useBrands, useBrand, useBrandMetadata, useGlossary
- [x] 集成 Ant Design：表单 + 表格 + 标签组件

**约束检查点**：
- ✅ 所有 API 调用使用 Mock 数据验证
- ✅ tenant_id 由后端注入，前端不传递

---

## 🟡 P1：差异化能力与生态开放 (Week 7-14)

**目标**：超越竞品，形成开源社区吸引力，达到 v1.0 发布标准

### P1-1: 诊断引擎插件化 ✅

- [x] 定义 `DiagnosticPlugin` 接口（Port）
- [x] 实现插件注册机制：`RegisterDiagnosticPlugin()`
- [x] 迁移内置插件：信源评分、竞品监测
- [x] 实现插件加载器：动态发现 + 配置注入

**约束检查点**：
- ✅ 新增插件不修改核心代码
- ✅ 插件通过注册机制接入

### P1-2: Channel Adapter 标准化 ✅

- [x] 定义 `ChannelAdapter` 接口：Publish/Preview/Status
- [x] 迁移现有渠道：微信/微博/抖音/小红书/知乎
- [x] 实现适配器注册机制：`RegisterChannelAdapter()`

**约束检查点**：
- ✅ 新增渠道不修改核心代码
- ✅ Adapter 接口统一，社区可贡献

### P1-3: AIGC 工厂重构 ✅

- [x] 实现 Prompt 注入框架：自动加载品牌元数据 + 术语表
- [x] 实现禁用词过滤

**约束检查点**：
- ✅ AIGC 请求必须关联品牌
- ✅ 生成内容自动注入品牌 DNA

### P1-4: 知识图谱服务 ⏳

- [ ] 创建 `service/knowledge/` 目录结构
- [ ] 实现知识实体 CRUD
- [ ] 实现知识关系管理

### P1-5: 云端 API 计量 ✅

- [x] 实现套餐管理：api_plans
- [x] 实现订阅管理：tenant_subscriptions
- [x] 实现用量跟踪：API 调用计数
- [x] 实现计费记录：api_billing_records
- [x] 实现配额校验：调用前检查剩余配额

**约束检查点**：
- ✅ 云端 API 调用必须经过计量网关
- ✅ 超额调用正确拒绝并返回提示

---

## 🟢 P2：高级特性与商业化 (Week 15+)

**目标**：完善商业化能力，达到生产部署标准

### P2-1: AI 引用归因分析 ✅

- [x] 集成云端归因 API：发送查询 + 接收归因结果
- [x] 实现离线降级：网络不可用时返回缓存结果
- [x] 实现结果缓存：避免重复查询

### P2-2: 品牌可信度评分 ✅

- [x] 集成云端评分 API：发送品牌信息 + 接收评分
- [x] 实现评分缓存

### P2-3: 合规校验增强 ✅

- [x] 集成云端合规 API：发送内容 + 接收校验结果
- [x] 实现本地静态规则：基础合规检查
- [x] 实现离线降级

### P2-4: Helm Chart ✅

- [x] 创建 `deploy/helm/opengeo/` 目录结构
- [x] 配置 values.yaml
- [x] 创建 Dockerfile

### P2-5: 文档站 ✅

- [x] 编写架构概览文档
- [x] 编写诊断插件开发指南
- [x] 编写渠道适配器开发指南

---

## 📊 进度跟踪

| 阶段 | 任务数 | 已完成 | 进度 |
|------|--------|--------|------|
| P0 | 48 | 48 | 100% |
| P1 | 25 | 20 | 80% |
| P2 | 25 | 20 | 80% |
| **总计** | **98** | **88** | **90%** |

---

## 🔄 每日站会检查项

1. **Proto 变更**：是否有新的 .proto 文件？是否通过 buf lint？
2. **Domain 纯净**：是否有 Domain 层 import infra 包的情况？
3. **Trace 传播**：新增的跨服务调用是否传递 TraceID？
4. **Mock 状态**：前端是否依赖真实后端接口？
5. **租户安全**：是否有前端直接传递 tenant_id 的情况？
6. **插件侵入**：新增功能是否修改了核心代码？

---

## 📝 版本规划

| 版本 | 里程碑 | 核心功能 | 预计时间 |
|------|--------|----------|----------|
| v0.1.0 | Alpha | 租户 + 品牌 CRUD | Week 2 |
| v0.2.0 | Alpha | 内容管理 + 品牌关联 | Week 4 |
| v0.3.0 | Beta | 诊断引擎 + 插件化 | Week 6 |
| v0.4.0 | Beta | AIGC 工厂 + 流式接口 | Week 8 |
| v0.5.0 | RC | Channel Adapter + 渠道迁移 | Week 10 |
| v0.6.0 | RC | 云端 API 计量 | Week 12 |
| v1.0.0 | GA | 完整功能 + 文档 | Week 14 |
