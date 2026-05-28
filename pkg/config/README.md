# 配置服务使用指南

## 概述

所有硬编码的配置项已外部化，支持以下配置方式（按优先级从高到低）：

1. **环境变量** - 最高优先级
2. **配置文件** - 中等优先级
3. **默认值** - 最低优先级

## 配置文件

配置文件为JSON格式，示例见 `configs/config.example.json`。

### 使用方式

设置环境变量 `CONFIG_FILE` 指向配置文件路径：

```bash
export CONFIG_FILE=/path/to/config.json
```

## 环境变量

### JWT配置
| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `JWT_TOKEN_EXPIRE_HOURS` | Token过期时间（小时） | 24 |
| `JWT_REFRESH_EXPIRE_DAYS` | Refresh Token过期时间（天） | 7 |

### 数据库配置
| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `DB_MAX_IDLE_CONNS` | 最大空闲连接数 | 10 |
| `DB_MAX_OPEN_CONNS` | 最大打开连接数 | 100 |

### 分页配置
| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `DEFAULT_PAGE_SIZE` | 默认分页大小 | 20 |

### AI模型配置
| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `DEEPSEEK_BASE_URL` | DeepSeek API地址 | https://api.deepseek.com/v1 |
| `DEEPSEEK_MODEL` | DeepSeek模型名称 | deepseek-chat |
| `KIMI_BASE_URL` | Kimi API地址 | https://api.moonshot.cn/v1 |
| `KIMI_MODEL` | Kimi模型名称 | moonshot-v1-8k |
| `DOUBAO_BASE_URL` | 豆包 API地址 | https://ark.cn-beijing.volces.com/api/v3 |
| `DOUBAO_MODEL` | 豆包模型名称 | doubao-pro-4k |
| `OPENAI_BASE_URL` | OpenAI API地址 | https://api.openai.com/v1 |
| `OPENAI_MODEL` | OpenAI模型名称 | gpt-4o-mini |
| `DEFAULT_AI_MODEL` | 默认AI模型 | deepseek |

### Webhook配置
| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `WEBHOOK_TIMEOUT_SECONDS` | Webhook超时时间（秒） | 30 |
| `WEBHOOK_MAX_RETRIES` | Webhook最大重试次数 | 3 |

### LLM配置
| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `LLM_TIMEOUT_SECONDS` | LLM超时时间（秒） | 30 |
| `LLM_MAX_TOKENS` | LLM最大Token数 | 4096 |

## 代码中使用

```go
import "opengeo/pkg/config"

// 获取配置
cfg := config.GetConfig()

// 使用配置
tokenExpire := cfg.JWT.TokenExpire
maxRetries := cfg.Retry.MaxRetries
deepseekBaseURL := cfg.AIModels.DeepSeek.BaseURL
```

## 可配置项列表

### JWT配置
- Token过期时间
- Refresh Token过期时间
- 签发者

### 数据库配置
- 最大空闲连接数
- 最大打开连接数
- 连接最大生命周期

### 分页配置
- 默认分页大小
- 最大分页大小

### AI模型配置
- DeepSeek: BaseURL、Model、MaxTokens、Temperature
- Kimi: BaseURL、Model、MaxTokens、Temperature
- Doubao: BaseURL、Model、MaxTokens、Temperature
- ChatGPT: BaseURL、Model、MaxTokens、Temperature
- 默认AI模型

### 防封引擎配置
- 各平台延迟范围
- 默认延迟范围
- 最小延迟
- 正态分布偏移比例
- 滚动间隔范围
- 打字延迟范围
- 打字停顿概率
- 各平台频率限制
- 代理失败阈值

### 错峰调度配置
- 最小间隔
- 最大间隔
- 浮动比例
- 最大并发数
- 突发限制

### 重试配置
- 最大重试次数
- 初始延迟
- 最大延迟
- 退避因子
- 是否启用降级
- 是否启用自动重试

### Webhook配置
- 超时时间
- 最大重试次数

### 健康检测配置
- HTTP超时时间
- 检测间隔
- 批量限制
- 告警分数阈值
- 暂停分数阈值

### LLM配置
- 超时时间
- 最大Token数
- 温度
