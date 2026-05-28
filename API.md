# OpenGEO API 文档

> 版本：v1.0.0  
> Base URL：`http://localhost:8080/api/v1`

---

## 认证

所有需要认证的 API 都需要在请求头中携带 `Authorization: Bearer <token>`

---

## 1. 认证 API

### 1.1 用户登录

**POST** `/auth/login`

**请求体：**
```json
{
  "username": "string (必填)",
  "password": "string (必填)"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": 1,
    "username": "admin",
    "email": "admin@example.com"
  }
}
```

**错误码：**
| 错误码 | 说明 |
|--------|------|
| 20001 | 用户不存在 |
| 20003 | 密码错误 |
| 20004 | 用户已禁用 |

---

### 1.2 用户注册

**POST** `/auth/register`

**请求体：**
```json
{
  "username": "string (必填，3-20字符)",
  "password": "string (必填，至少8字符)",
  "email": "string (必填，有效邮箱)"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": 1,
    "username": "admin",
    "email": "admin@example.com"
  }
}
```

---

### 1.3 刷新 Token

**POST** `/auth/refresh`

**请求体：**
```json
{
  "refresh_token": "string (必填)"
}
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "token": "新 token",
    "refresh_token": "新 refresh_token"
  }
}
```

---

## 2. 用户 API

### 2.1 获取用户信息

**GET** `/users/{id}`

**路径参数：**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | int64 | 用户 ID |

**响应：**
```json
{
  "code": 0,
  "data": {
    "user_id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "status": 1,
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

---

### 2.2 更新用户

**PUT** `/users/{id}`

**请求体：**
```json
{
  "email": "string (可选)",
  "status": 1
}
```

---

### 2.3 删除用户

**DELETE** `/users/{id}`

---

### 2.4 用户列表

**GET** `/users`

**查询参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 10 |
| keyword | string | 否 | 搜索关键词 |

---

## 3. 内容 API

### 3.1 创建内容

**POST** `/contents`

**请求体：**
```json
{
  "title": "string (必填)",
  "body": "string (必填)",
  "content_type": "article|video|image (必填)",
  "schema_markup": "string (可选，JSON格式)"
}
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "示例内容",
    "body": "内容正文...",
    "content_type": "article",
    "status": 0,
    "ai_optimization_score": 0,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

---

### 3.2 获取内容

**GET** `/contents/{id}`

---

### 3.3 更新内容

**PUT** `/contents/{id}`

**请求体：**
```json
{
  "title": "string (可选)",
  "body": "string (可选)",
  "schema_markup": "string (可选)"
}
```

---

### 3.4 删除内容

**DELETE** `/contents/{id}`

---

### 3.5 内容列表

**GET** `/contents`

**查询参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| user_id | int64 | 否 | 用户 ID |
| content_type | string | 否 | 内容类型 |
| status | int | 否 | 状态 (0:草稿, 1:已发布, 2:已归档) |

---

### 3.6 AI 优化内容

**POST** `/contents/{id}/optimize`

**请求体：**
```json
{
  "ai_model": "deepseek|kimi|doubao|chatgpt",
  "optimization_type": "geo_semantic|structure|readability"
}
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "success": true,
    "optimized_title": "优化后的标题",
    "optimized_body": "优化后的正文",
    "schema_markup": "结构化数据",
    "score": 85.5,
    "suggestions": ["建议1", "建议2"],
    "structural_changes": ["变化1", "变化2"]
  }
}
```

---

### 3.7 发布内容

**POST** `/contents/{id}/publish`

---

## 4. 账号 API

### 4.1 创建账号

**POST** `/accounts`

**请求体：**
```json
{
  "platform": "wechat|weibo|douyin|xiaohongshu|zhihu|toutiao",
  "account_name": "string (必填)",
  "account_id": "string (必填)"
}
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "user_id": 1,
    "platform": "wechat",
    "account_name": "我的公众号",
    "account_id": "wx_123456",
    "status": 1,
    "health_score": 100,
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

---

### 4.2 获取账号

**GET** `/accounts/{id}`

---

### 4.3 更新账号

**PUT** `/accounts/{id}`

**请求体：**
```json
{
  "account_name": "string (可选)",
  "status": 1
}
```

---

### 4.4 删除账号

**DELETE** `/accounts/{id}`

---

### 4.5 账号列表

**GET** `/accounts`

**查询参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| user_id | int64 | 否 | 用户 ID |
| platform | string | 否 | 平台类型 |

---

### 4.6 账号健康检查

**GET** `/accounts/{id}/health`

**响应：**
```json
{
  "code": 0,
  "data": {
    "account_id": 1,
    "health_score": 95.5,
    "status": "normal",
    "check_details": "{}",
    "checked_at": "2026-01-01T00:00:00Z"
  }
}
```

---

## 5. 账号分组 API

### 5.1 创建分组

**POST** `/account-groups`

**请求体：**
```json
{
  "name": "string (必填)",
  "group_type": "authority|professional|ecology (必填)",
  "description": "string (可选)",
  "parent_id": 1
}
```

---

### 5.2 获取分组

**GET** `/account-groups/{id}`

---

### 5.3 分组列表

**GET** `/account-groups`

---

### 5.4 添加账号到分组

**POST** `/account-groups/{id}/accounts`

**请求体：**
```json
{
  "account_id": 1
}
```

---

### 5.5 从分组移除账号

**DELETE** `/account-groups/{groupId}/accounts/{accountId}`

---

## 6. 发布任务 API

### 6.1 创建发布任务

**POST** `/publish/tasks`

**请求体：**
```json
{
  "content_id": 1,
  "channel_id": 1,
  "scheduled_time": "2026-01-01T00:00:00Z (可选，留空表示立即发布)"
}
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "user_id": 1,
    "content_id": 1,
    "channel_id": 1,
    "status": 0,
    "scheduled_time": null,
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

---

### 6.2 获取发布任务

**GET** `/publish/tasks/{id}`

---

### 6.3 发布任务列表

**GET** `/publish/tasks`

**查询参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| user_id | int64 | 否 | 用户 ID |
| status | int | 否 | 状态 (0:待发布, 1:发布中, 2:已发布, 3:失败, 4:已取消) |

---

### 6.4 取消发布任务

**POST** `/publish/tasks/{id}/cancel`

---

### 6.5 重试发布任务

**POST** `/publish/tasks/{id}/retry`

---

### 6.6 发布预览

**POST** `/publish/preview`

**请求体：**
```json
{
  "channel_id": 1,
  "title": "string",
  "body": "string",
  "media_urls": ["url1", "url2"]
}
```

**响应：**
```json
{
  "code": 0,
  "data": {
    "success": true,
    "html": "<div>预览HTML</div>",
    "markdown": "预览Markdown",
    "warnings": []
  }
}
```

---

### 6.7 发布校验

**POST** `/publish/validate`

---

## 7. 渠道 API

### 7.1 创建渠道

**POST** `/channels`

**请求体：**
```json
{
  "channel_type": "wechat|weibo|douyin|xiaohongshu",
  "channel_name": "string (必填)",
  "channel_config": "string (JSON格式配置)"
}
```

---

### 7.2 获取渠道

**GET** `/channels/{id}`

---

### 7.3 渠道列表

**GET** `/channels`

---

### 7.4 支持的平台

**GET** `/channels/platforms`

**响应：**
```json
{
  "code": 0,
  "data": ["wechat", "weibo", "douyin", "xiaohongshu", "zhihu", "toutiao"]
}
```

---

## 8. 调度 API

### 8.1 创建调度

**POST** `/schedules`

**请求体：**
```json
{
  "schedule_name": "string (必填)",
  "schedule_type": "fixed|interval|event|heat (必填)",
  "cron_expression": "string (可选)",
  "config": "string (可选，JSON格式)"
}
```

---

### 8.2 获取调度

**GET** `/schedules/{id}`

---

### 8.3 更新调度

**PUT** `/schedules/{id}`

---

### 8.4 删除调度

**DELETE** `/schedules/{id}`

---

### 8.5 调度列表

**GET** `/schedules`

---

### 8.6 启用调度

**POST** `/schedules/{id}/enable`

---

### 8.7 禁用调度

**POST** `/schedules/{id}/disable`

---

## 9. 监测 API

### 9.1 AI 引用列表

**GET** `/monitor/citations`

**查询参数：**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| content_id | int64 | 否 | 内容 ID |
| ai_model | string | 否 | AI 模型 |

**响应：**
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "id": 1,
        "content_id": 1,
        "ai_model": "deepseek",
        "query_text": "什么是GEO优化？",
        "is_cited": true,
        "citation_position": 2,
        "citation_text": "GEO优化是...",
        "sentiment": "positive",
        "tracked_at": "2026-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

### 9.2 信源评分

**GET** `/monitor/scores`

---

### 9.3 竞品监测

**GET** `/monitor/competitors`

---

### 9.4 ROI 指标

**GET** `/monitor/roi`

---

### 9.5 生成优化建议

**POST** `/monitor/suggestions/generate`

**请求体：**
```json
{
  "content_id": 1
}
```

---

## 10. 系统 API

### 10.1 系统配置列表

**GET** `/system/configs`

---

### 10.2 更新系统配置

**PUT** `/system/configs/{key}`

**请求体：**
```json
{
  "config_value": "string"
}
```

---

### 10.3 插件列表

**GET** `/system/plugins`

---

### 10.4 Webhook 列表

**GET** `/system/webhooks`

---

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10001 | 内部错误 |
| 10002 | 参数错误 |
| 10003 | 未授权 |
| 10004 | 禁止访问 |
| 10005 | 资源不存在 |
| 10006 | 资源已存在 |
| 20001 | 用户不存在 |
| 20002 | 用户已存在 |
| 20003 | 密码错误 |
| 30001 | 内容不存在 |
| 40001 | 账号不存在 |
| 50001 | 发布任务不存在 |
| 60001 | 调度不存在 |
| 70001 | AI引用不存在 |
| 80001 | 配置不存在 |

---

## 分页响应格式

```json
{
  "code": 0,
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 认证流程

1. 调用 `/auth/login` 获取 `token` 和 `refresh_token`
2. 在后续请求头中携带 `Authorization: Bearer <token>`
3. Token 过期后使用 `/auth/refresh` 刷新
4. 如果 refresh_token 也过期，需要重新登录

---

> 文档版本：v1.0.0  
> 最后更新：2026-05-26