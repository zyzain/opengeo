# 前端品牌功能实现文档

## 概述

本文档描述了 OpenGEO BrandOS 前端品牌功能的完整实现，包括品牌管理、元数据编辑、术语表管理、知识图谱可视化和品牌快照等功能。

## 目录结构

```
web/src/
├── components/brand/              # 品牌组件
│   ├── BrandForm.tsx              # 品牌创建/编辑表单
│   ├── BrandSelector.tsx          # 品牌选择器
│   ├── MetadataEditor.tsx         # 元数据编辑器
│   ├── GlossaryManager.tsx        # 术语表管理
│   ├── KnowledgeGraph.tsx         # 知识图谱可视化
│   ├── BrandSnapshot.tsx          # 品牌快照
│   ├── BulkImportModal.tsx        # 批量导入弹窗
│   └── index.ts                   # 组件导出
│
├── hooks/                         # React Hooks
│   ├── useBrand.ts                # 品牌数据 Hooks
│   ├── useKnowledge.ts            # 知识图谱 Hooks
│   ├── useSnapshot.ts             # 快照 Hooks
│   ├── useStream.ts               # 流式数据 Hook
│   └── index.ts                   # Hooks 导出
│
├── pages/brand/                   # 品牌页面
│   ├── index.tsx                  # 品牌列表页
│   ├── detail.tsx                 # 品牌详情页
│   ├── [id].tsx                   # 品牌详情页（动态路由）
│   ├── brand.module.css           # 品牌页面样式
│   └── index.module.css           # 列表页样式
│
├── types/brand.ts                 # 类型定义
│
└── lib/api/brand.ts               # API 服务
```

## 功能模块

### 1. 品牌列表页

**文件**: `pages/brand/index.tsx`

**功能**:
- 品牌列表展示（表格形式）
- 品牌搜索（按名称、标识）
- 行业筛选
- 状态筛选
- 品牌统计（总数、活跃、归档、禁用）
- 创建品牌
- 编辑品牌
- 删除品牌（带确认）

**组件**:
- `Table` - 品牌列表表格
- `Card` - 统计卡片
- `Modal` - 创建/编辑弹窗
- `BrandForm` - 品牌表单

### 2. 品牌详情页

**文件**: `pages/brand/detail.tsx`

**功能**:
- 品牌信息展示
- 品牌元数据编辑
- 术语表管理
- 知识图谱可视化
- 品牌快照管理

**组件**:
- `Tabs` - 功能标签页
- `Descriptions` - 品牌信息展示
- `MetadataEditor` - 元数据编辑器
- `GlossaryManager` - 术语表管理
- `KnowledgeGraph` - 知识图谱
- `BrandSnapshot` - 品牌快照

### 3. 品牌表单组件

**文件**: `components/brand/BrandForm.tsx`

**功能**:
- 品牌名称输入
- 品牌标识输入（URL友好，创建后不可修改）
- 品牌描述
- 行业选择
- 品牌官网
- Logo URL
- 成立年份
- 总部所在地

**验证规则**:
- 品牌名称：必填
- 品牌标识：必填，只能包含小写字母、数字和连字符

### 4. 元数据编辑器

**文件**: `components/brand/MetadataEditor.tsx`

**功能**:
- VI 规范编辑
  - 主色/副色
  - 字体
  - 品牌口号
  - 品牌关键词（标签形式）
- 语调规范编辑
  - 正式度（正式/随意/技术）
  - 个性（友好/专业/活泼/权威）
  - 风格指南
  - 偏好短语（标签形式）
  - 禁用词（标签形式）
- 受众画像编辑
  - 受众名称
  - 年龄范围
  - 兴趣标签
  - 痛点标签
- 品牌价值观（标签形式）
- 独特卖点（标签形式）

### 5. 术语表管理

**文件**: `components/brand/GlossaryManager.tsx`

**功能**:
- 术语列表展示
- 术语搜索
- 分类筛选
- 禁用词筛选
- 术语统计（总数、禁用词、首选术语、分类数）
- 添加术语
- 编辑术语
- 删除术语（带确认）
- 批量导入（CSV格式）
- 导出术语
- 下载导入模板

**术语字段**:
- 术语名称
- 定义
- 分类（产品/技术/概念/人物/地点）
- 别名（多个）
- 使用上下文
- 是否禁用词
- 是否首选术语

### 6. 知识图谱可视化

**文件**: `components/brand/KnowledgeGraph.tsx`

**功能**:
- 实体列表展示
- 实体类型标识（颜色区分）
- 添加知识实体
- 添加知识关系
- 查看实体详情
- 删除实体

**实体类型**:
- 品牌（蓝色）
- 产品（绿色）
- 人物（紫色）
- 组织（青色）
- 事件（橙色）
- 概念（粉色）
- 地点（黄色）
- 技术（靛蓝色）

**关系类型**:
- 是一种（is_a）
- 属于（part_of）
- 相关（related_to）
- 竞争（competes_with）
- 提及（mentions）
- 依赖（depends_on）
- 拥有（owns）
- 创建者（created_by）

### 7. 品牌快照

**文件**: `components/brand/BrandSnapshot.tsx`

**功能**:
- 快照列表展示
- 创建快照
- 对比快照
- 快照版本管理

**快照字段**:
- 版本号
- 变更说明
- 创建时间

### 8. 品牌选择器

**文件**: `components/brand/BrandSelector.tsx`

**功能**:
- 品牌下拉选择
- 品牌搜索
- 品牌 Logo 展示
- 行业标签展示

### 9. 批量导入弹窗

**文件**: `components/brand/BulkImportModal.tsx`

**功能**:
- 文件上传（CSV/TSV）
- 数据预览
- 下载模板
- 导入结果展示

## Hooks

### useBrands

```typescript
const { brands, loading, error, refetch } = useBrands();
```

获取品牌列表数据。

### useBrand

```typescript
const { brand, loading, error } = useBrand(brandId);
```

获取单个品牌详情。

### useBrandMetadata

```typescript
const { metadata, loading, refetch } = useBrandMetadata(brandId);
```

获取品牌元数据。

### useGlossary

```typescript
const { entries, loading, refetch } = useGlossary(brandId);
```

获取品牌术语表。

### useKnowledgeEntities

```typescript
const { entities, loading, refetch, addEntity, updateEntity, deleteEntity } = useKnowledgeEntities(brandId);
```

获取和管理知识实体。

### useKnowledgeRelations

```typescript
const { relations, loading, refetch, addRelation, deleteRelation } = useKnowledgeRelations(brandId);
```

获取和管理知识关系。

### useKnowledgeGraph

```typescript
const { entities, relations, loading, refetch } = useKnowledgeGraph(brandId);
```

获取完整的知识图谱数据。

### useSnapshots

```typescript
const { snapshots, loading, refetch, createSnapshot, compareSnapshots } = useSnapshots(brandId);
```

获取和管理品牌快照。

### useStream

```typescript
const { content, issues, isStreaming, start, stop } = useStream(url, options);
```

处理 SSE 流式数据。

## API 接口

### 品牌管理

- `GET /api/v1/brands` - 列出品牌
- `GET /api/v1/brand/:id` - 获取品牌
- `POST /api/v1/brands` - 创建品牌
- `PUT /api/v1/brand/:id` - 更新品牌
- `DELETE /api/v1/brand/:id` - 删除品牌

### 品牌元数据

- `GET /api/v1/brand/:id/metadata` - 获取元数据
- `PUT /api/v1/brand/:id/metadata` - 更新元数据

### 术语表

- `GET /api/v1/brand/:id/glossary` - 列出术语
- `POST /api/v1/brand/:id/glossary` - 创建术语
- `PUT /api/v1/brand/:id/glossary/:entryId` - 更新术语
- `DELETE /api/v1/brand/:id/glossary/:entryId` - 删除术语
- `POST /api/v1/brand/:id/glossary/bulk-import` - 批量导入
- `GET /api/v1/brand/:id/glossary/export` - 导出术语

### 知识图谱

- `GET /api/v1/brand/:id/knowledge/entities` - 列出实体
- `POST /api/v1/brand/:id/knowledge/entities` - 创建实体
- `PUT /api/v1/brand/:id/knowledge/entities/:entityId` - 更新实体
- `DELETE /api/v1/brand/:id/knowledge/entities/:entityId` - 删除实体
- `GET /api/v1/brand/:id/knowledge/entities/search` - 搜索实体
- `GET /api/v1/brand/:id/knowledge/relations` - 列出关系
- `POST /api/v1/brand/:id/knowledge/relations` - 创建关系
- `DELETE /api/v1/brand/:id/knowledge/relations/:relationId` - 删除关系
- `POST /api/v1/brand/:id/knowledge/graph` - 查询图谱

### 品牌快照

- `GET /api/v1/brand/:id/snapshots` - 列出快照
- `POST /api/v1/brand/:id/snapshots` - 创建快照
- `POST /api/v1/brand/:id/snapshots/compare` - 对比快照

## 类型定义

详见 `types/brand.ts` 文件，包含以下主要类型：

- `Brand` - 品牌
- `BrandMetadata` - 品牌元数据
- `VIProfile` - VI 规范
- `ToneProfile` - 语调规范
- `AudienceProfile` - 受众画像
- `GlossaryEntry` - 术语条目
- `KnowledgeEntity` - 知识实体
- `KnowledgeRelation` - 知识关系
- `BrandSnapshot` - 品牌快照

## 样式

品牌页面使用 CSS Modules 进行样式管理：

- `brand.module.css` - 品牌详情页样式
- `index.module.css` - 品牌列表页样式

## 路由配置

在 `App.tsx` 中配置品牌相关路由：

```typescript
{ path: "brand", element: <BrandListPage /> },
{ path: "brand/:id", element: <BrandDetailPage /> },
```

## Mock 数据

前端使用 MSW (Mock Service Worker) 进行 API Mock：

- `lib/mock/handlers/brand.handlers.ts` - 品牌 API Mock
- `lib/mock/fixtures/brand.ts` - 品牌 Mock 数据

## 开发指南

### 添加新的品牌功能

1. 在 `types/brand.ts` 中添加类型定义
2. 在 `lib/api/brand.ts` 中添加 API 接口
3. 在 `hooks/` 中添加自定义 Hook
4. 在 `components/brand/` 中创建组件
5. 在 `pages/brand/` 中创建页面
6. 在 `App.tsx` 中添加路由

### 组件开发规范

1. 使用 TypeScript 进行类型定义
2. 使用 Ant Design 组件库
3. 使用 React Hooks 管理状态
4. 使用 CSS Modules 进行样式管理
5. 遵循 Proto First 原则，确保类型一致

### 测试

1. 使用 MSW Mock 进行前端开发
2. 使用 Vitest 进行单元测试
3. 使用 Playwright 进行 E2E 测试
