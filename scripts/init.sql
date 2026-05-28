-- ============================================================
-- OpenGEO BrandOS - 数据库初始化脚本
-- 版本: v6.0.0
-- 字符集: utf8mb4 (支持完整 Unicode，包括 Emoji)
-- 架构: 云端 API + SaaS 商业模式
-- 说明: 租户与品牌管理为开源基座，AI 能力为云端计量服务
-- ============================================================

SET NAMES utf8mb4;
SET CHARACTER_SET_CLIENT = utf8mb4;
SET CHARACTER_SET_CONNECTION = utf8mb4;
SET CHARACTER_SET_RESULTS = utf8mb4;

CREATE DATABASE IF NOT EXISTS opengeo
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE opengeo;

-- ============================================================
-- 1. 租户服务 (Tenant Service) - 开源基座
-- ============================================================

-- 租户表：支持多租户 SaaS 模式，开源基座标配
CREATE TABLE IF NOT EXISTS tenants (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '租户ID',
    name            VARCHAR(128) NOT NULL COMMENT '租户名称',
    slug            VARCHAR(64)  NOT NULL COMMENT '租户标识（URL友好）',
    domain          VARCHAR(256) DEFAULT NULL COMMENT '租户绑定域名',
    logo_url        VARCHAR(512) DEFAULT NULL COMMENT '租户Logo',
    plan            VARCHAR(32)  DEFAULT 'free' COMMENT '套餐：free/pro/enterprise',
    status          TINYINT      DEFAULT 1 COMMENT '状态：1=正常 0=禁用 2=过期',
    brand_limit     INT          DEFAULT 5 COMMENT '品牌数上限',
    user_limit      INT          DEFAULT 10 COMMENT '用户数上限',
    storage_limit   BIGINT       DEFAULT 1073741824 COMMENT '存储上限（字节），默认1GB',
    api_quota       INT          DEFAULT 1000 COMMENT '云端API月调用配额',
    api_used        INT          DEFAULT 0 COMMENT '本月已用API调用次数',
    quota_reset_at  DATETIME     DEFAULT NULL COMMENT '配额重置时间',
    settings        JSON         COMMENT '租户配置JSON',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_slug (slug) COMMENT '租户标识唯一索引',
    INDEX idx_domain (domain) COMMENT '域名索引',
    INDEX idx_plan (plan) COMMENT '套餐索引',
    INDEX idx_status (status) COMMENT '状态索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- 租户配额表：记录云端 API 用量明细
CREATE TABLE IF NOT EXISTS tenant_api_usage (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    api_type        VARCHAR(64)  NOT NULL COMMENT 'API类型：attribution/trust_score/compliance/optimization',
    endpoint        VARCHAR(256) NOT NULL COMMENT 'API端点',
    tokens_used     INT          DEFAULT 0 COMMENT '消耗Token数',
    cost_cents      INT          DEFAULT 0 COMMENT '费用（分）',
    request_id      VARCHAR(64)  DEFAULT NULL COMMENT '请求ID',
    response_time   INT          DEFAULT 0 COMMENT '响应时间(ms)',
    status_code     INT          DEFAULT 200 COMMENT 'HTTP状态码',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_api_type (tenant_id, api_type) COMMENT '租户+API类型联合索引',
    INDEX idx_created (created_at) COMMENT '创建时间索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户API用量表';

-- ============================================================
-- 2. 账号服务 (Account Service) - 开源基座
-- ============================================================

-- 用户表：系统用户，归属租户
CREATE TABLE IF NOT EXISTS users (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    username        VARCHAR(64)  NOT NULL COMMENT '用户名',
    password        VARCHAR(256) NOT NULL COMMENT '密码（bcrypt加密存储）',
    email           VARCHAR(128) DEFAULT NULL COMMENT '邮箱地址',
    display_name    VARCHAR(64)  DEFAULT NULL COMMENT '显示名称',
    avatar_url      VARCHAR(512) DEFAULT NULL COMMENT '头像URL',
    phone           VARCHAR(32)  DEFAULT NULL COMMENT '手机号',
    status          TINYINT      DEFAULT 1 COMMENT '状态：1=正常 0=禁用',
    email_verified  BOOLEAN      DEFAULT FALSE COMMENT '邮箱是否验证',
    last_login_at   DATETIME     DEFAULT NULL COMMENT '最后登录时间',
    last_login_ip   VARCHAR(64)  DEFAULT NULL COMMENT '最后登录IP',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_tenant_username (tenant_id, username) COMMENT '租户内用户名唯一索引',
    UNIQUE INDEX idx_tenant_email (tenant_id, email) COMMENT '租户内邮箱唯一索引',
    INDEX idx_status (status) COMMENT '状态索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 角色表：RBAC 角色定义
CREATE TABLE IF NOT EXISTS roles (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '角色ID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    name            VARCHAR(64)  NOT NULL COMMENT '角色标识',
    display_name    VARCHAR(64)  NOT NULL COMMENT '角色显示名称',
    description     VARCHAR(256) DEFAULT NULL COMMENT '角色描述',
    is_system       BOOLEAN      DEFAULT FALSE COMMENT '是否系统内置角色',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE INDEX idx_tenant_name (tenant_id, name) COMMENT '租户内角色名唯一索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 权限表：细粒度权限定义
CREATE TABLE IF NOT EXISTS permissions (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '权限ID',
    name            VARCHAR(128) NOT NULL COMMENT '权限标识，如 brand:create',
    module          VARCHAR(64)  NOT NULL COMMENT '所属模块：brand/content/publish/monitor/system',
    action          VARCHAR(64)  NOT NULL COMMENT '操作类型：create/read/update/delete/manage',
    description     VARCHAR(256) DEFAULT NULL COMMENT '权限描述',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE INDEX idx_name (name) COMMENT '权限标识唯一索引',
    INDEX idx_module (module) COMMENT '模块索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    user_id         BIGINT NOT NULL COMMENT '用户ID',
    role_id         BIGINT NOT NULL COMMENT '角色ID',
    brand_id        BIGINT DEFAULT NULL COMMENT '品牌ID（NULL=全局角色）',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (user_id, role_id, brand_id),
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id         BIGINT NOT NULL COMMENT '角色ID',
    permission_id   BIGINT NOT NULL COMMENT '权限ID',
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 第三方平台账号表：管理自媒体/新闻源等外部平台账号
CREATE TABLE IF NOT EXISTS accounts (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '账号ID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    user_id         BIGINT       NOT NULL COMMENT '所属用户ID',
    platform        VARCHAR(64)  NOT NULL COMMENT '平台类型：wechat/weibo/douyin/xiaohongshu/zhihu',
    account_name    VARCHAR(128) NOT NULL COMMENT '账号名称',
    account_id      VARCHAR(128) DEFAULT NULL COMMENT '第三方平台账号ID',
    access_token    TEXT         COMMENT '访问令牌（加密存储）',
    refresh_token   TEXT         COMMENT '刷新令牌（加密存储）',
    token_expires_at DATETIME    DEFAULT NULL COMMENT '令牌过期时间',
    status          TINYINT      DEFAULT 1 COMMENT '状态：1=正常 2=限流 3=封禁 0=禁用',
    health_score    DECIMAL(5,2) DEFAULT 100.00 COMMENT '健康度评分 0-100',
    last_check_time DATETIME     DEFAULT NULL COMMENT '最后健康检查时间',
    metadata        JSON         COMMENT '平台特有元数据',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_user (user_id) COMMENT '用户索引',
    INDEX idx_platform (tenant_id, platform) COMMENT '租户+平台联合索引',
    INDEX idx_status (status) COMMENT '状态索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='第三方平台账号表';

-- 账号分组表：支持三层分组架构（权威背书层/专业认证层/生态渗透层）
CREATE TABLE IF NOT EXISTS account_groups (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '分组ID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    brand_id        BIGINT       DEFAULT NULL COMMENT '关联品牌ID（NULL=全局分组）',
    name            VARCHAR(128) NOT NULL COMMENT '分组名称',
    parent_id       BIGINT       DEFAULT NULL COMMENT '父分组ID，NULL=顶级分组',
    group_type      VARCHAR(64)  DEFAULT NULL COMMENT '分组类型：authority/professional/ecology',
    description     VARCHAR(256) DEFAULT NULL COMMENT '分组描述',
    sort_order      INT          DEFAULT 0 COMMENT '排序权重',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_parent (parent_id) COMMENT '父分组索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账号分组表';

-- 账号分组关联表
CREATE TABLE IF NOT EXISTS account_group_relations (
    account_id      BIGINT NOT NULL COMMENT '账号ID',
    group_id        BIGINT NOT NULL COMMENT '分组ID',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (account_id, group_id),
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES account_groups(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账号分组关联表';

-- ============================================================
-- 3. 品牌服务 (Brand Service) - 开源基座
-- ============================================================

-- 品牌主表：支持多品牌治理，开源基座标配
CREATE TABLE IF NOT EXISTS brands (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '品牌ID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    name            VARCHAR(128) NOT NULL COMMENT '品牌名称',
    slug            VARCHAR(64)  NOT NULL COMMENT '品牌标识（URL友好）',
    description     TEXT         COMMENT '品牌描述',
    logo_url        VARCHAR(512) DEFAULT NULL COMMENT '品牌Logo URL',
    website         VARCHAR(256) DEFAULT NULL COMMENT '品牌官网',
    industry        VARCHAR(64)  DEFAULT NULL COMMENT '所属行业',
    founded_year    INT          DEFAULT NULL COMMENT '成立年份',
    headquarters    VARCHAR(128) DEFAULT NULL COMMENT '总部所在地',
    status          TINYINT      DEFAULT 1 COMMENT '状态：1=活跃 2=归档 0=禁用',
    settings        JSON         COMMENT '品牌配置JSON',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_tenant_slug (tenant_id, slug) COMMENT '租户内品牌标识唯一索引',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_industry (industry) COMMENT '行业索引',
    INDEX idx_status (status) COMMENT '状态索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌主表';

-- 品牌元数据表：存储品牌 VI 规范、语调、受众画像
CREATE TABLE IF NOT EXISTS brand_metadata (
    id                BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '元数据ID',
    brand_id          BIGINT       NOT NULL COMMENT '品牌ID',
    vi_profile        JSON         COMMENT 'VI规范：主色/副色/Logo/字体/关键词',
    tone_profile      JSON         COMMENT '语调规范：正式度/个性/禁用词/偏好短语',
    audience_profiles JSON         COMMENT '受众画像：年龄/兴趣/痛点/偏好渠道',
    competitor_list   JSON         COMMENT '竞品列表',
    brand_values      JSON         COMMENT '品牌价值观',
    unique_selling_points JSON     COMMENT '独特卖点',
    schema_version    VARCHAR(32)  DEFAULT '1.0' COMMENT 'Schema版本号',
    created_at        DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at        DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_brand (brand_id) COMMENT '品牌唯一索引',
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌元数据表';

-- 品牌术语表：管理品牌专属术语
CREATE TABLE IF NOT EXISTS glossary_entries (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '术语ID',
    brand_id        BIGINT       NOT NULL COMMENT '所属品牌ID',
    term            VARCHAR(128) NOT NULL COMMENT '术语名称',
    definition      TEXT         NOT NULL COMMENT '术语定义',
    category        VARCHAR(64)  DEFAULT NULL COMMENT '术语分类：product/technology/concept/person/place',
    aliases         JSON         COMMENT '别名列表',
    context         TEXT         COMMENT '使用上下文示例',
    is_forbidden    BOOLEAN      DEFAULT FALSE COMMENT '是否禁用词',
    is_preferred    BOOLEAN      DEFAULT FALSE COMMENT '是否首选术语',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_term (term) COMMENT '术语索引',
    INDEX idx_category (brand_id, category) COMMENT '品牌+分类联合索引',
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌术语表';

-- 品牌快照表：支持品牌配置版本管理和对比诊断
CREATE TABLE IF NOT EXISTS brand_snapshots (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '快照ID',
    brand_id        BIGINT       NOT NULL COMMENT '品牌ID',
    version         VARCHAR(32)  NOT NULL COMMENT '版本号',
    snapshot_data   JSON         NOT NULL COMMENT '快照数据（完整品牌配置）',
    change_log      TEXT         COMMENT '变更说明',
    created_by      BIGINT       DEFAULT NULL COMMENT '创建人ID',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_version (brand_id, version) COMMENT '品牌+版本联合索引',
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌快照表';

-- 品牌知识实体表：品牌专属知识图谱
CREATE TABLE IF NOT EXISTS brand_knowledge_entities (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '实体ID',
    brand_id        BIGINT       NOT NULL COMMENT '所属品牌ID',
    entity_name     VARCHAR(128) NOT NULL COMMENT '实体名称',
    entity_type     VARCHAR(32)  NOT NULL COMMENT '实体类型：brand/product/person/org/event/concept',
    entity_data     JSON         COMMENT '实体属性JSON',
    authority_links JSON         COMMENT '权威链接JSON',
    embedding_id    VARCHAR(128) DEFAULT NULL COMMENT 'Milvus向量ID',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_type (brand_id, entity_type) COMMENT '品牌+类型联合索引',
    INDEX idx_name (entity_name) COMMENT '实体名称索引',
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌知识实体表';

-- 知识图谱关系表
CREATE TABLE IF NOT EXISTS knowledge_relations (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '关系ID',
    from_entity_id  BIGINT       NOT NULL COMMENT '源实体ID',
    to_entity_id    BIGINT       NOT NULL COMMENT '目标实体ID',
    relation_type   VARCHAR(64)  NOT NULL COMMENT '关系类型：is_a/part_of/related_to/competes_with/mentions',
    weight          DECIMAL(3,2) DEFAULT 1.00 COMMENT '关系权重 0-1',
    metadata        JSON         COMMENT '关系元数据',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_from (from_entity_id) COMMENT '源实体索引',
    INDEX idx_to (to_entity_id) COMMENT '目标实体索引',
    INDEX idx_type (relation_type) COMMENT '关系类型索引',
    FOREIGN KEY (from_entity_id) REFERENCES brand_knowledge_entities(id) ON DELETE CASCADE,
    FOREIGN KEY (to_entity_id) REFERENCES brand_knowledge_entities(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识图谱关系表';

-- 品牌可信度评分表：存储云端 API 返回的评分结果
CREATE TABLE IF NOT EXISTS brand_trust_scores (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '评分ID',
    brand_id        BIGINT       NOT NULL COMMENT '品牌ID',
    score           DECIMAL(5,2) NOT NULL COMMENT '综合可信度评分 0-100',
    dimensions      JSON         COMMENT '各维度评分：search/social/compliance/citation',
    factors         JSON         COMMENT '评分因子详情',
    api_request_id  VARCHAR(64)  DEFAULT NULL COMMENT '云端API请求ID',
    expires_at      DATETIME     DEFAULT NULL COMMENT '评分过期时间',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_score (brand_id, score) COMMENT '品牌+评分联合索引',
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌可信度评分表';

-- ============================================================
-- 4. 内容服务 (Content Service) - 开源基座
-- ============================================================

-- 内容表：存储文章/视频/图片等内容
CREATE TABLE IF NOT EXISTS contents (
    id                    BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '内容ID',
    tenant_id             BIGINT       NOT NULL COMMENT '所属租户ID',
    user_id               BIGINT       NOT NULL COMMENT '所属用户ID',
    brand_id              BIGINT       DEFAULT NULL COMMENT '关联品牌ID',
    title                 VARCHAR(256) NOT NULL COMMENT '内容标题',
    body                  TEXT         NOT NULL COMMENT '内容正文',
    summary               VARCHAR(512) DEFAULT NULL COMMENT '内容摘要',
    content_type          VARCHAR(32)  DEFAULT 'article' COMMENT '内容类型：article/video/image/infographic',
    status                TINYINT      DEFAULT 0 COMMENT '状态：0=草稿 1=已发布 2=已归档 3=待审核',
    visibility            VARCHAR(32)  DEFAULT 'private' COMMENT '可见性：private/team/public',
    schema_markup         TEXT         COMMENT 'JSON-LD结构化数据',
    ai_optimization_score DECIMAL(5,2) DEFAULT NULL COMMENT 'AI优化评分 0-100',
    word_count            INT          DEFAULT 0 COMMENT '字数',
    reading_time          INT          DEFAULT 0 COMMENT '预计阅读时间（分钟）',
    tags                  JSON         COMMENT '标签列表',
    metadata              JSON         COMMENT '扩展元数据',
    published_at          DATETIME     DEFAULT NULL COMMENT '发布时间',
    created_at            DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at            DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_user (user_id) COMMENT '用户索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_status (tenant_id, status) COMMENT '租户+状态联合索引',
    INDEX idx_type (content_type) COMMENT '内容类型索引',
    INDEX idx_published (published_at) COMMENT '发布时间索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内容表';

-- 内容版本表：记录内容的历史版本
CREATE TABLE IF NOT EXISTS content_versions (
    id                  BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '版本ID',
    content_id          BIGINT       NOT NULL COMMENT '所属内容ID',
    version             INT          NOT NULL COMMENT '版本号',
    title               VARCHAR(256) NOT NULL COMMENT '版本标题',
    body                TEXT         NOT NULL COMMENT '版本正文',
    summary             VARCHAR(512) DEFAULT NULL COMMENT '版本摘要',
    schema_markup       TEXT         COMMENT '版本结构化数据',
    ai_model_adaptation VARCHAR(64)  DEFAULT NULL COMMENT '适配的AI模型',
    change_summary      VARCHAR(512) DEFAULT NULL COMMENT '变更摘要',
    created_by          BIGINT       DEFAULT NULL COMMENT '创建人ID',
    created_at          DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE INDEX idx_content_version (content_id, version) COMMENT '内容+版本号唯一索引',
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内容版本表';

-- 内容模板表：GEO优化Prompt模板市场
CREATE TABLE IF NOT EXISTS content_templates (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '模板ID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    user_id         BIGINT       DEFAULT NULL COMMENT '创建者用户ID',
    brand_id        BIGINT       DEFAULT NULL COMMENT '关联品牌ID（NULL=全局模板）',
    name            VARCHAR(128) NOT NULL COMMENT '模板名称',
    description     VARCHAR(256) DEFAULT NULL COMMENT '模板描述',
    template_type   VARCHAR(32)  DEFAULT NULL COMMENT '模板类型：article/faq/review/analysis/social',
    template_data   TEXT         NOT NULL COMMENT '模板内容（支持变量占位符）',
    variables       JSON         COMMENT '变量定义JSON',
    is_public       BOOLEAN      DEFAULT FALSE COMMENT '是否公开可见',
    usage_count     INT          DEFAULT 0 COMMENT '使用次数',
    rating          DECIMAL(3,2) DEFAULT 0.00 COMMENT '平均评分 0-5',
    rating_count    INT          DEFAULT 0 COMMENT '评分人数',
    tags            VARCHAR(512) DEFAULT NULL COMMENT '标签，逗号分隔',
    author          VARCHAR(128) DEFAULT NULL COMMENT '作者名称',
    is_official     BOOLEAN      DEFAULT FALSE COMMENT '是否官方模板',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_type_public (template_type, is_public) COMMENT '类型+公开联合索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内容模板表';

-- 内容实体关联表：内容与知识实体的多对多关系
CREATE TABLE IF NOT EXISTS content_entities (
    content_id      BIGINT NOT NULL COMMENT '内容ID',
    entity_id       BIGINT NOT NULL COMMENT '实体ID',
    relevance_score DECIMAL(3,2) DEFAULT 1.00 COMMENT '相关性评分 0-1',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (content_id, entity_id),
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE,
    FOREIGN KEY (entity_id) REFERENCES brand_knowledge_entities(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内容实体关联表';

-- 合规校验记录表：存储内容合规校验结果
CREATE TABLE IF NOT EXISTS compliance_checks (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '校验ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    content_id      BIGINT       NOT NULL COMMENT '内容ID',
    brand_id        BIGINT       DEFAULT NULL COMMENT '品牌ID',
    check_type      VARCHAR(32)  NOT NULL COMMENT '校验类型：local/cloud',
    status          VARCHAR(32)  DEFAULT 'pending' COMMENT '状态：pending/passed/failed/warning',
    risk_level      VARCHAR(16)  DEFAULT 'low' COMMENT '风险级别：low/medium/high/critical',
    issues          JSON         COMMENT '问题列表JSON',
    suggestions     JSON         COMMENT '修改建议JSON',
    api_request_id  VARCHAR(64)  DEFAULT NULL COMMENT '云端API请求ID',
    checked_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '校验时间',
    expires_at      DATETIME     DEFAULT NULL COMMENT '结果过期时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_content (content_id) COMMENT '内容索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合规校验记录表';

-- ============================================================
-- 5. 发布服务 (Publish Service) - 开源基座
-- ============================================================

-- 平台表：管理发布平台配置
CREATE TABLE IF NOT EXISTS platforms (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '平台ID',
    code            VARCHAR(32)  NOT NULL COMMENT '平台代码：wechat/weibo/douyin',
    name            VARCHAR(64)  NOT NULL COMMENT '平台名称：微信公众号',
    icon            VARCHAR(32)  DEFAULT NULL COMMENT '图标标识',
    color           VARCHAR(16)  DEFAULT NULL COMMENT '标签颜色',
    description     VARCHAR(256) DEFAULT NULL COMMENT '平台描述',
    config_schema   TEXT         COMMENT '配置项JSON Schema',
    features        JSON         COMMENT '平台特性：支持的content_type等',
    is_enabled      BOOLEAN      DEFAULT TRUE COMMENT '是否启用',
    sort_order      INT          DEFAULT 0 COMMENT '排序权重',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_code (code) COMMENT '平台代码唯一索引',
    INDEX idx_enabled (is_enabled) COMMENT '启用状态索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发布平台表';

-- 发布渠道表：配置各平台发布渠道
CREATE TABLE IF NOT EXISTS publish_channels (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '渠道ID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    brand_id        BIGINT       DEFAULT NULL COMMENT '关联品牌ID',
    channel_type    VARCHAR(64)  NOT NULL COMMENT '渠道类型：wechat/weibo/douyin/xiaohongshu/zhihu',
    channel_name    VARCHAR(128) NOT NULL COMMENT '渠道名称',
    channel_config  TEXT         COMMENT '渠道配置JSON（API密钥等）',
    title_template  TEXT         COMMENT '标题模板（支持变量替换）',
    body_template   TEXT         COMMENT '正文模板',
    tags_template   VARCHAR(512) DEFAULT NULL COMMENT '标签模板',
    cover_template  VARCHAR(512) DEFAULT NULL COMMENT '封面模板',
    geo_config      TEXT         COMMENT 'GEO专属配置JSON',
    is_enabled      BOOLEAN      DEFAULT TRUE COMMENT '是否启用',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_type (tenant_id, channel_type) COMMENT '租户+类型联合索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发布渠道表';

-- 发布任务表：管理内容发布任务
CREATE TABLE IF NOT EXISTS publish_tasks (
    id                  BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '任务ID',
    tenant_id           BIGINT       NOT NULL COMMENT '所属租户ID',
    user_id             BIGINT       NOT NULL COMMENT '所属用户ID',
    content_id          BIGINT       NOT NULL COMMENT '内容ID',
    channel_id          BIGINT       NOT NULL COMMENT '渠道ID',
    fallback_channel_id BIGINT       DEFAULT NULL COMMENT '备用渠道ID',
    status              TINYINT      DEFAULT 0 COMMENT '状态：0=待发布 1=发布中 2=成功 3=失败 4=取消 5=重试中 6=降级中 7=人工审核',
    scheduled_time      DATETIME     DEFAULT NULL COMMENT '计划发布时间',
    published_time      DATETIME     DEFAULT NULL COMMENT '实际发布时间',
    external_id         VARCHAR(128) DEFAULT NULL COMMENT '外部平台内容ID',
    external_url        VARCHAR(512) DEFAULT NULL COMMENT '外部平台内容URL',
    retry_count         INT          DEFAULT 0 COMMENT '已重试次数',
    max_retries         INT          DEFAULT 3 COMMENT '最大重试次数',
    retry_delay         INT          DEFAULT 30 COMMENT '重试延迟(秒)',
    error_message       VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
    error_history       TEXT         COMMENT '错误历史记录',
    priority            TINYINT      DEFAULT 0 COMMENT '优先级：0=普通 1=中 2=高 3=紧急',
    is_manually_review  BOOLEAN      DEFAULT FALSE COMMENT '是否需要人工审核',
    review_note         VARCHAR(512) DEFAULT NULL COMMENT '审核备注',
    metadata            JSON         COMMENT '扩展元数据',
    created_at          DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at          DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_status_time (status, scheduled_time) COMMENT '状态+计划时间联合索引',
    INDEX idx_user_status (user_id, status) COMMENT '用户+状态联合索引',
    INDEX idx_priority (priority) COMMENT '优先级索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发布任务表';

-- 降级队列表：发布失败后的备用渠道队列
CREATE TABLE IF NOT EXISTS fallback_queues (
    id                  BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    task_id             BIGINT       NOT NULL COMMENT '原始任务ID',
    tenant_id           BIGINT       NOT NULL COMMENT '租户ID',
    user_id             BIGINT       NOT NULL COMMENT '用户ID',
    content_id          BIGINT       DEFAULT NULL COMMENT '内容ID',
    original_channel_id BIGINT       NOT NULL COMMENT '原始渠道ID',
    fallback_channel_id BIGINT       NOT NULL COMMENT '备用渠道ID',
    reason              TEXT         COMMENT '降级原因',
    status              VARCHAR(32)  DEFAULT 'pending' COMMENT '状态：pending/approved/rejected/processed',
    review_note         VARCHAR(512) DEFAULT NULL COMMENT '审核备注',
    reviewed_by         VARCHAR(128) DEFAULT NULL COMMENT '审核人',
    reviewed_at         DATETIME     DEFAULT NULL COMMENT '审核时间',
    created_at          DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at          DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_task (task_id) COMMENT '任务ID索引',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_status (status) COMMENT '状态索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='降级队列表';

-- ============================================================
-- 6. 调度服务 (Scheduler Service) - 开源基座
-- ============================================================

-- 调度表：定时任务配置
CREATE TABLE IF NOT EXISTS schedules (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '调度ID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    brand_id        BIGINT       DEFAULT NULL COMMENT '关联品牌ID',
    user_id         BIGINT       NOT NULL COMMENT '创建者用户ID',
    schedule_name   VARCHAR(128) NOT NULL COMMENT '调度名称',
    schedule_type   VARCHAR(32)  NOT NULL COMMENT '调度类型：fixed/interval/event/heat',
    cron_expression VARCHAR(128) DEFAULT NULL COMMENT 'Cron表达式或间隔表达式',
    config          TEXT         COMMENT '调度配置JSON',
    is_enabled      BOOLEAN      DEFAULT TRUE COMMENT '是否启用',
    next_run_time   DATETIME     DEFAULT NULL COMMENT '下次运行时间',
    last_run_time   DATETIME     DEFAULT NULL COMMENT '上次运行时间',
    run_count       BIGINT       DEFAULT 0 COMMENT '累计运行次数',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_next_run (is_enabled, next_run_time) COMMENT '待执行调度查询索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='调度表';

-- 调度任务关联表
CREATE TABLE IF NOT EXISTS schedule_tasks (
    schedule_id     BIGINT   NOT NULL COMMENT '调度ID',
    task_id         BIGINT   NOT NULL COMMENT '任务ID',
    scheduled_at    DATETIME NOT NULL COMMENT '计划执行时间',
    status          VARCHAR(32) DEFAULT 'pending' COMMENT '状态：pending/executed/failed',
    executed_at     DATETIME DEFAULT NULL COMMENT '实际执行时间',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (schedule_id, task_id),
    INDEX idx_schedule_status (schedule_id, status, scheduled_at) COMMENT '调度+状态+时间联合索引',
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='调度任务关联表';

-- AI活跃度热力图：记录各平台/模型的活跃度数据
CREATE TABLE IF NOT EXISTS ai_activity_heatmaps (
    id              BIGINT      PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    platform        VARCHAR(64) NOT NULL COMMENT '平台：wechat/weibo/douyin 等',
    ai_model        VARCHAR(64) NOT NULL COMMENT 'AI模型：deepseek/kimi/chatgpt 等',
    time_slot       TINYINT     NOT NULL COMMENT '小时段 0-23',
    day_of_week     TINYINT     NOT NULL COMMENT '星期几 1=周一 7=周日',
    activity_score  DECIMAL(5,2) DEFAULT 0 COMMENT '活跃度评分 0-100',
    sample_count    BIGINT      DEFAULT 0 COMMENT '采样次数',
    updated_at      DATETIME    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_heatmap_lookup (platform, ai_model, time_slot, day_of_week) COMMENT '热力图查询联合索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI活跃度热力图';

-- 发布日历表：可视化排期管理
CREATE TABLE IF NOT EXISTS publish_calendars (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '日历事件ID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    brand_id        BIGINT       DEFAULT NULL COMMENT '关联品牌ID',
    user_id         BIGINT       NOT NULL COMMENT '所属用户ID',
    title           VARCHAR(256) NOT NULL COMMENT '事件标题',
    description     TEXT         COMMENT '事件描述',
    start_time      DATETIME     NOT NULL COMMENT '开始时间',
    end_time        DATETIME     DEFAULT NULL COMMENT '结束时间',
    all_day         BOOLEAN      DEFAULT FALSE COMMENT '是否全天事件',
    schedule_id     BIGINT       DEFAULT NULL COMMENT '关联调度ID',
    content_id      BIGINT       DEFAULT NULL COMMENT '关联内容ID',
    channel_id      BIGINT       DEFAULT NULL COMMENT '关联渠道ID',
    status          VARCHAR(32)  DEFAULT 'pending' COMMENT '状态：pending/executing/completed/cancelled',
    color           VARCHAR(16)  DEFAULT NULL COMMENT '事件颜色',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_user (user_id) COMMENT '用户索引',
    INDEX idx_start_time (start_time) COMMENT '开始时间索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发布日历表';

-- 调度日志表：记录调度执行日志
CREATE TABLE IF NOT EXISTS schedule_logs (
    id              BIGINT      PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    schedule_id     BIGINT      NOT NULL COMMENT '调度ID',
    task_id         BIGINT      DEFAULT NULL COMMENT '关联任务ID',
    action          VARCHAR(32) NOT NULL COMMENT '操作：trigger/skip/error/complete',
    message         TEXT        COMMENT '日志消息',
    duration_ms     INT         DEFAULT 0 COMMENT '执行时长(ms)',
    created_at      DATETIME    DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_schedule (schedule_id) COMMENT '调度ID索引',
    INDEX idx_created (created_at) COMMENT '创建时间索引',
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='调度日志表';

-- ============================================================
-- 7. 监测服务 (Monitor Service) - 开源基座 + 云端API
-- ============================================================

-- AI引用追踪表：记录内容被AI模型引用的情况
CREATE TABLE IF NOT EXISTS ai_citations (
    id                  BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    tenant_id           BIGINT       NOT NULL COMMENT '租户ID',
    content_id          BIGINT       NOT NULL COMMENT '内容ID',
    brand_id            BIGINT       DEFAULT NULL COMMENT '品牌ID',
    ai_model            VARCHAR(64)  NOT NULL COMMENT 'AI模型：deepseek/kimi/chatgpt 等',
    query_text          VARCHAR(512) NOT NULL COMMENT '查询文本',
    is_cited            BOOLEAN      DEFAULT FALSE COMMENT '是否被引用',
    citation_position   INT          DEFAULT NULL COMMENT '引用位置（1=最靠前）',
    citation_text       TEXT         COMMENT '引用的文本片段',
    citation_url        VARCHAR(512) DEFAULT NULL COMMENT '引用URL',
    sentiment           VARCHAR(32)  DEFAULT 'neutral' COMMENT '情感倾向：positive/neutral/negative',
    confidence_score    DECIMAL(3,2) DEFAULT NULL COMMENT '置信度评分 0-1',
    attribution_data    JSON         DEFAULT NULL COMMENT '归因链路数据',
    api_request_id      VARCHAR(64)  DEFAULT NULL COMMENT '云端API请求ID',
    tracked_at          DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '追踪时间',
    created_at          DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_content (content_id) COMMENT '内容索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_model (ai_model) COMMENT 'AI模型索引',
    INDEX idx_citation_stats (content_id, is_cited) COMMENT '引用统计索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI引用追踪表';

-- AI引用归因分析表：存储云端API返回的归因结果
CREATE TABLE IF NOT EXISTS citation_attributions (
    id                  BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '归因ID',
    tenant_id           BIGINT       NOT NULL COMMENT '租户ID',
    content_id          BIGINT       NOT NULL COMMENT '内容ID',
    brand_id            BIGINT       DEFAULT NULL COMMENT '品牌ID',
    query_text          VARCHAR(512) NOT NULL COMMENT '查询文本',
    attribution_score   DECIMAL(5,2) NOT NULL COMMENT '归因评分 0-100',
    cited_fragments     JSON         COMMENT '被引用片段列表',
    source_analysis     JSON         COMMENT '信源分析',
    competitor_comparison JSON       COMMENT '竞品对比',
    recommendations     JSON         COMMENT '优化建议',
    api_request_id      VARCHAR(64)  NOT NULL COMMENT '云端API请求ID',
    tokens_used         INT          DEFAULT 0 COMMENT '消耗Token数',
    cost_cents          INT          DEFAULT 0 COMMENT '费用（分）',
    expires_at          DATETIME     DEFAULT NULL COMMENT '结果过期时间',
    created_at          DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_content (content_id) COMMENT '内容索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_api_request (api_request_id) COMMENT 'API请求ID索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI引用归因分析表';

-- 信源评分表：渠道/账号的权威度评分
CREATE TABLE IF NOT EXISTS source_scores (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '评分ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    channel_id      BIGINT       NOT NULL COMMENT '渠道ID',
    account_id      BIGINT       NOT NULL COMMENT '账号ID',
    score           DECIMAL(5,2) NOT NULL COMMENT '综合评分 0-100',
    score_dimensions JSON        COMMENT '各维度评分JSON',
    plugin_name     VARCHAR(64)  DEFAULT 'source_score' COMMENT '评分插件名称',
    api_request_id  VARCHAR(64)  DEFAULT NULL COMMENT '云端API请求ID',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_source_lookup (channel_id, account_id) COMMENT '渠道+账号联合索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='信源评分表';

-- 竞品监测表：配置竞品监测任务
CREATE TABLE IF NOT EXISTS competitor_monitors (
    id                  BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '监测ID',
    tenant_id           BIGINT       NOT NULL COMMENT '所属租户ID',
    brand_id            BIGINT       NOT NULL COMMENT '关联品牌ID',
    user_id             BIGINT       NOT NULL COMMENT '创建者用户ID',
    competitor_name     VARCHAR(128) NOT NULL COMMENT '竞品名称',
    competitor_domain   VARCHAR(256) DEFAULT NULL COMMENT '竞品域名',
    monitor_config      TEXT         COMMENT '监测配置JSON',
    last_check_time     DATETIME     DEFAULT NULL COMMENT '最后检查时间',
    is_active           BOOLEAN      DEFAULT TRUE COMMENT '是否激活',
    created_at          DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at          DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞品监测表';

-- 竞品分析结果表：存储竞品分析结果
CREATE TABLE IF NOT EXISTS competitor_analyses (
    id                  BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '分析ID',
    monitor_id          BIGINT       NOT NULL COMMENT '监测ID',
    visibility_score    DECIMAL(5,2) DEFAULT 0 COMMENT '可见性评分 0-100',
    top_queries         TEXT         COMMENT '热门查询JSON',
    content_gaps        TEXT         COMMENT '内容差距JSON',
    recommendations     TEXT         COMMENT '建议JSON',
    api_request_id      VARCHAR(64)  DEFAULT NULL COMMENT '云端API请求ID',
    analyzed_at         DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '分析时间',
    created_at          DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_monitor (monitor_id) COMMENT '监测ID索引',
    FOREIGN KEY (monitor_id) REFERENCES competitor_monitors(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='竞品分析结果表';

-- ROI指标表：记录转化数据
CREATE TABLE IF NOT EXISTS roi_metrics (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '指标ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    content_id      BIGINT       NOT NULL COMMENT '内容ID',
    brand_id        BIGINT       DEFAULT NULL COMMENT '品牌ID',
    channel_id      BIGINT       NOT NULL COMMENT '渠道ID',
    metric_type     VARCHAR(32)  NOT NULL COMMENT '指标类型：impression/click/inquiry/visit/consult/conversion',
    metric_value    DECIMAL(12,2) DEFAULT 0 COMMENT '指标值',
    currency        VARCHAR(3)   DEFAULT 'CNY' COMMENT '货币类型',
    utm_source      VARCHAR(64)  DEFAULT NULL COMMENT 'UTM来源',
    utm_medium      VARCHAR(64)  DEFAULT NULL COMMENT 'UTM媒介',
    utm_campaign    VARCHAR(64)  DEFAULT NULL COMMENT 'UTM活动',
    recorded_at     DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '记录时间',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_content (content_id) COMMENT '内容索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_recorded (recorded_at) COMMENT '记录时间索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ROI指标表';

-- 优化建议表：AI生成的优化建议
CREATE TABLE IF NOT EXISTS optimization_suggestions (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '建议ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    content_id      BIGINT       NOT NULL COMMENT '内容ID',
    brand_id        BIGINT       DEFAULT NULL COMMENT '品牌ID',
    suggestion_type VARCHAR(32)  DEFAULT NULL COMMENT '建议类型：content/structure/authority/compliance/citation',
    source          VARCHAR(32)  DEFAULT 'local' COMMENT '来源：local/cloud',
    suggestion_data TEXT         COMMENT '建议内容JSON',
    priority        TINYINT      DEFAULT 0 COMMENT '优先级：0=低 1=中 2=高',
    impact_score    DECIMAL(3,2) DEFAULT NULL COMMENT '预期影响评分 0-1',
    status          TINYINT      DEFAULT 0 COMMENT '状态：0=待处理 1=已应用 2=已忽略',
    applied_at      DATETIME     DEFAULT NULL COMMENT '应用时间',
    api_request_id  VARCHAR(64)  DEFAULT NULL COMMENT '云端API请求ID',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_content (content_id) COMMENT '内容索引',
    INDEX idx_status (content_id, status) COMMENT '内容+状态联合索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优化建议表';

-- 引用趋势表：按天统计引用趋势
CREATE TABLE IF NOT EXISTS citation_trends (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '趋势ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    content_id      BIGINT       NOT NULL COMMENT '内容ID',
    brand_id        BIGINT       DEFAULT NULL COMMENT '品牌ID',
    ai_model        VARCHAR(64)  NOT NULL COMMENT 'AI模型',
    citation_count  INT          DEFAULT 0 COMMENT '引用次数',
    citation_rate   DECIMAL(5,2) DEFAULT 0 COMMENT '引用率',
    trend_date      DATE         NOT NULL COMMENT '统计日期',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_content (content_id) COMMENT '内容索引',
    INDEX idx_brand (brand_id) COMMENT '品牌索引',
    INDEX idx_date (trend_date) COMMENT '日期索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='引用趋势表';

-- ============================================================
-- 8. 系统服务 (System Service) - 开源基座
-- ============================================================

-- 系统配置表：全局配置项
CREATE TABLE IF NOT EXISTS system_configs (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    config_key      VARCHAR(128) NOT NULL COMMENT '配置键',
    config_value    TEXT         COMMENT '配置值',
    config_type     VARCHAR(32)  DEFAULT 'string' COMMENT '值类型：string/number/json/boolean',
    description     VARCHAR(256) DEFAULT NULL COMMENT '配置描述',
    is_public       BOOLEAN      DEFAULT FALSE COMMENT '是否前端可见',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_config_key (config_key) COMMENT '配置键唯一索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 插件表：管理渠道适配器/诊断插件/AI模型连接器
CREATE TABLE IF NOT EXISTS plugins (
    id                  BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '插件ID',
    plugin_name         VARCHAR(64)  NOT NULL COMMENT '插件名称',
    plugin_type         VARCHAR(32)  NOT NULL COMMENT '插件类型：channel/diagnostic/ai_model/analyzer',
    display_name        VARCHAR(64)  NOT NULL COMMENT '插件显示名称',
    description         VARCHAR(256) DEFAULT NULL COMMENT '插件描述',
    version             VARCHAR(32)  DEFAULT NULL COMMENT '版本号',
    author              VARCHAR(64)  DEFAULT NULL COMMENT '作者',
    homepage            VARCHAR(256) DEFAULT NULL COMMENT '插件主页',
    config_schema       TEXT         COMMENT '配置Schema JSON',
    supported_metrics   JSON         DEFAULT NULL COMMENT '支持的指标列表（诊断插件）',
    supported_channels  JSON         DEFAULT NULL COMMENT '支持的渠道列表（渠道插件）',
    is_system           BOOLEAN      DEFAULT FALSE COMMENT '是否系统内置插件',
    is_enabled          BOOLEAN      DEFAULT TRUE COMMENT '是否启用',
    installed_at        DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '安装时间',
    updated_at          DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_plugin_name (plugin_name) COMMENT '插件名唯一索引',
    INDEX idx_type (plugin_type) COMMENT '插件类型索引',
    INDEX idx_enabled (is_enabled) COMMENT '启用状态索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='插件表';

-- Webhook表：Webhook配置
CREATE TABLE IF NOT EXISTS webhooks (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT 'WebhookID',
    tenant_id       BIGINT       NOT NULL COMMENT '所属租户ID',
    user_id         BIGINT       NOT NULL COMMENT '创建者用户ID',
    webhook_name    VARCHAR(128) NOT NULL COMMENT 'Webhook名称',
    url             VARCHAR(512) NOT NULL COMMENT '回调URL',
    secret          VARCHAR(128) DEFAULT NULL COMMENT '签名密钥',
    events          TEXT         COMMENT '订阅事件列表JSON',
    headers         JSON         COMMENT '自定义请求头',
    is_active       BOOLEAN      DEFAULT TRUE COMMENT '是否激活',
    last_trigger    DATETIME     DEFAULT NULL COMMENT '最后触发时间',
    fail_count      INT          DEFAULT 0 COMMENT '连续失败次数',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Webhook表';

-- Webhook事件表：记录Webhook触发事件
CREATE TABLE IF NOT EXISTS webhook_events (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '事件ID',
    webhook_id      BIGINT       NOT NULL COMMENT 'WebhookID',
    event_type      VARCHAR(64)  NOT NULL COMMENT '事件类型',
    payload         TEXT         COMMENT '事件负载JSON',
    status_code     INT          DEFAULT 0 COMMENT 'HTTP响应状态码',
    response_body   TEXT         COMMENT 'HTTP响应体',
    success         BOOLEAN      DEFAULT FALSE COMMENT '是否投递成功',
    duration_ms     INT          DEFAULT 0 COMMENT '请求耗时(ms)',
    triggered_at    DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '触发时间',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_webhook (webhook_id) COMMENT 'WebhookID索引',
    INDEX idx_triggered (triggered_at) COMMENT '触发时间索引',
    FOREIGN KEY (webhook_id) REFERENCES webhooks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Webhook事件表';

-- 国际化翻译表：多语言文案管理
CREATE TABLE IF NOT EXISTS translations (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '翻译ID',
    locale          VARCHAR(32)  NOT NULL COMMENT '语言代码：zh-CN/en-US/ja-JP 等',
    `key`           VARCHAR(256) NOT NULL COMMENT '文案键',
    value           TEXT         NOT NULL COMMENT '文案值',
    context         VARCHAR(256) DEFAULT NULL COMMENT '上下文说明',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_locale_key (locale, `key`) COMMENT '语言+键联合索引',
    INDEX idx_key (`key`) COMMENT '键索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='国际化翻译表';

-- 审计日志表：记录所有操作日志
CREATE TABLE IF NOT EXISTS audit_logs (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    tenant_id       BIGINT       DEFAULT NULL COMMENT '租户ID',
    user_id         BIGINT       DEFAULT NULL COMMENT '操作用户ID',
    username        VARCHAR(64)  DEFAULT NULL COMMENT '操作用户名',
    action          VARCHAR(64)  NOT NULL COMMENT '操作类型',
    resource_type   VARCHAR(64)  DEFAULT NULL COMMENT '资源类型',
    resource_id     BIGINT       DEFAULT NULL COMMENT '资源ID',
    resource_name   VARCHAR(128) DEFAULT NULL COMMENT '资源名称',
    details         TEXT         COMMENT '操作详情JSON',
    ip_address      VARCHAR(64)  DEFAULT NULL COMMENT 'IP地址',
    user_agent      VARCHAR(256) DEFAULT NULL COMMENT 'User-Agent',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_user (user_id) COMMENT '用户索引',
    INDEX idx_action (action) COMMENT '操作类型索引',
    INDEX idx_resource (resource_type, resource_id) COMMENT '资源类型+ID联合索引',
    INDEX idx_created (created_at) COMMENT '创建时间索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';

-- 通知表：系统通知
CREATE TABLE IF NOT EXISTS notifications (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '通知ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    user_id         BIGINT       NOT NULL COMMENT '接收用户ID',
    type            VARCHAR(32)  NOT NULL COMMENT '通知类型：system/alert/task/compliance',
    title           VARCHAR(256) NOT NULL COMMENT '通知标题',
    content         TEXT         COMMENT '通知内容',
    link            VARCHAR(512) DEFAULT NULL COMMENT '相关链接',
    is_read         BOOLEAN      DEFAULT FALSE COMMENT '是否已读',
    read_at         DATETIME     DEFAULT NULL COMMENT '阅读时间',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_tenant_user (tenant_id, user_id) COMMENT '租户+用户联合索引',
    INDEX idx_unread (user_id, is_read) COMMENT '未读通知索引',
    INDEX idx_created (created_at) COMMENT '创建时间索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知表';

-- ============================================================
-- 9. 云端API计量服务 (Cloud API Metering) - 商业化核心
-- ============================================================

-- 云端API套餐表：定义可用的API套餐
CREATE TABLE IF NOT EXISTS api_plans (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '套餐ID',
    name            VARCHAR(64)  NOT NULL COMMENT '套餐名称',
    code            VARCHAR(32)  NOT NULL COMMENT '套餐代码',
    description     VARCHAR(256) DEFAULT NULL COMMENT '套餐描述',
    monthly_quota   INT          NOT NULL COMMENT '月调用配额',
    price_cents     INT          NOT NULL COMMENT '月价格（分）',
    overage_price   INT          DEFAULT 0 COMMENT '超额单价（分/次）',
    features        JSON         COMMENT '包含的API类型列表',
    is_active       BOOLEAN      DEFAULT TRUE COMMENT '是否可用',
    sort_order      INT          DEFAULT 0 COMMENT '排序权重',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_code (code) COMMENT '套餐代码唯一索引',
    INDEX idx_active (is_active) COMMENT '可用状态索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='云端API套餐表';

-- 租户订阅表：记录租户的API订阅
CREATE TABLE IF NOT EXISTS tenant_subscriptions (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '订阅ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    plan_id         BIGINT       NOT NULL COMMENT '套餐ID',
    status          VARCHAR(32)  DEFAULT 'active' COMMENT '状态：active/cancelled/expired/trial',
    starts_at       DATETIME     NOT NULL COMMENT '开始时间',
    ends_at         DATETIME     DEFAULT NULL COMMENT '结束时间',
    trial_ends_at   DATETIME     DEFAULT NULL COMMENT '试用结束时间',
    auto_renew      BOOLEAN      DEFAULT TRUE COMMENT '是否自动续费',
    api_key         VARCHAR(64)  DEFAULT NULL COMMENT '云端API密钥',
    api_secret      VARCHAR(128) DEFAULT NULL COMMENT '云端API密钥（加密存储）',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_plan (plan_id) COMMENT '套餐索引',
    INDEX idx_status (status) COMMENT '状态索引',
    INDEX idx_api_key (api_key) COMMENT 'API密钥索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (plan_id) REFERENCES api_plans(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户订阅表';

-- API调用计费记录表：按次计费明细
CREATE TABLE IF NOT EXISTS api_billing_records (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    subscription_id BIGINT       NOT NULL COMMENT '订阅ID',
    api_type        VARCHAR(64)  NOT NULL COMMENT 'API类型',
    endpoint        VARCHAR(256) NOT NULL COMMENT 'API端点',
    request_id      VARCHAR(64)  NOT NULL COMMENT '请求ID',
    tokens_input    INT          DEFAULT 0 COMMENT '输入Token数',
    tokens_output   INT          DEFAULT 0 COMMENT '输出Token数',
    cost_cents      INT          DEFAULT 0 COMMENT '费用（分）',
    is_overage      BOOLEAN      DEFAULT FALSE COMMENT '是否超额调用',
    status          VARCHAR(32)  DEFAULT 'success' COMMENT '状态：success/failed/refunded',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_subscription (subscription_id) COMMENT '订阅索引',
    INDEX idx_api_type (tenant_id, api_type) COMMENT '租户+API类型联合索引',
    INDEX idx_created (created_at) COMMENT '创建时间索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (subscription_id) REFERENCES tenant_subscriptions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API调用计费记录表';

-- 云端API请求日志表：记录所有云端API请求
CREATE TABLE IF NOT EXISTS api_request_logs (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    request_id      VARCHAR(64)  NOT NULL COMMENT '请求ID',
    api_type        VARCHAR(64)  NOT NULL COMMENT 'API类型',
    endpoint        VARCHAR(256) NOT NULL COMMENT 'API端点',
    method          VARCHAR(10)  NOT NULL COMMENT 'HTTP方法',
    request_body    TEXT         COMMENT '请求体（脱敏）',
    response_body   TEXT         COMMENT '响应体',
    status_code     INT          DEFAULT 0 COMMENT 'HTTP状态码',
    response_time   INT          DEFAULT 0 COMMENT '响应时间(ms)',
    tokens_used     INT          DEFAULT 0 COMMENT '消耗Token数',
    error_message   VARCHAR(512) DEFAULT NULL COMMENT '错误信息',
    client_ip       VARCHAR(64)  DEFAULT NULL COMMENT '客户端IP',
    user_agent      VARCHAR(256) DEFAULT NULL COMMENT 'User-Agent',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_request (request_id) COMMENT '请求ID索引',
    INDEX idx_api_type (api_type) COMMENT 'API类型索引',
    INDEX idx_created (created_at) COMMENT '创建时间索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='云端API请求日志表';

-- 离线缓存表：云端API离线降级缓存
CREATE TABLE IF NOT EXISTS api_offline_cache (
    id              BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '缓存ID',
    cache_key       VARCHAR(256) NOT NULL COMMENT '缓存键',
    tenant_id       BIGINT       NOT NULL COMMENT '租户ID',
    api_type        VARCHAR(64)  NOT NULL COMMENT 'API类型',
    request_hash    VARCHAR(64)  NOT NULL COMMENT '请求参数哈希',
    response_data   TEXT         COMMENT '响应数据',
    hit_count       INT          DEFAULT 0 COMMENT '命中次数',
    expires_at      DATETIME     NOT NULL COMMENT '过期时间',
    created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE INDEX idx_cache_key (cache_key) COMMENT '缓存键唯一索引',
    INDEX idx_tenant (tenant_id) COMMENT '租户索引',
    INDEX idx_expires (expires_at) COMMENT '过期时间索引',
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='云端API离线缓存表';

-- ============================================================
-- 初始数据
-- ============================================================

-- 默认租户
INSERT INTO tenants (id, name, slug, plan, status, brand_limit, user_limit, api_quota) VALUES
(1, '默认租户', 'default', 'free', 1, 10, 20, 1000)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 默认角色
INSERT INTO roles (tenant_id, name, display_name, description, is_system) VALUES
(1, 'tenant_admin', '租户管理员', '拥有租户内全部权限', TRUE),
(1, 'brand_editor', '品牌编辑者', '可管理品牌和内容', TRUE),
(1, 'viewer', '只读用户', '仅查看权限', TRUE)
ON DUPLICATE KEY UPDATE display_name=VALUES(display_name);

-- 默认权限
INSERT INTO permissions (name, module, action, description) VALUES
-- 租户管理
('tenant:read', 'tenant', 'read', '查看租户信息'),
('tenant:update', 'tenant', 'update', '更新租户配置'),
('tenant:manage', 'tenant', 'manage', '管理租户设置'),
-- 用户管理
('user:create', 'user', 'create', '创建用户'),
('user:read', 'user', 'read', '查看用户'),
('user:update', 'user', 'update', '更新用户'),
('user:delete', 'user', 'delete', '删除用户'),
-- 品牌管理
('brand:create', 'brand', 'create', '创建品牌'),
('brand:read', 'brand', 'read', '查看品牌'),
('brand:update', 'brand', 'update', '更新品牌'),
('brand:delete', 'brand', 'delete', '删除品牌'),
('brand:manage', 'brand', 'manage', '管理品牌配置'),
-- 内容管理
('content:create', 'content', 'create', '创建内容'),
('content:read', 'content', 'read', '查看内容'),
('content:update', 'content', 'update', '更新内容'),
('content:delete', 'content', 'delete', '删除内容'),
('content:publish', 'content', 'publish', '发布内容'),
-- 发布管理
('publish:create', 'publish', 'create', '创建发布任务'),
('publish:read', 'publish', 'read', '查看发布任务'),
('publish:execute', 'publish', 'execute', '执行发布'),
('publish:cancel', 'publish', 'cancel', '取消发布'),
-- 监测分析
('monitor:read', 'monitor', 'read', '查看监测数据'),
('monitor:configure', 'monitor', 'configure', '配置监测任务'),
('monitor:export', 'monitor', 'export', '导出监测报告'),
-- 系统管理
('system:read', 'system', 'read', '查看系统信息'),
('system:update', 'system', 'update', '更新系统配置')
ON DUPLICATE KEY UPDATE description=VALUES(description);

-- 给租户管理员分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'tenant_admin' AND r.tenant_id = 1
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- 给品牌编辑者分配品牌和内容相关权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'brand_editor' AND r.tenant_id = 1 
AND p.name IN ('brand:create', 'brand:read', 'brand:update', 'content:create', 'content:read', 'content:update', 'content:delete', 'content:publish', 'publish:create', 'publish:read', 'publish:execute', 'monitor:read')
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- 给只读用户分配查看权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'viewer' AND r.tenant_id = 1 
AND p.name IN ('brand:read', 'content:read', 'publish:read', 'monitor:read')
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- 默认管理员用户（密码: Admin@123456 bcrypt hash）
INSERT INTO users (tenant_id, username, password, email, display_name, status, email_verified) VALUES
(1, 'admin', '$2a$10$EKD2ZjDX11ocpxA1V/3ZgOwkMUScdJAom6MdDhwOMsnZN4bNA4Aka', 'admin@opengeo.com', '系统管理员', 1, TRUE)
ON DUPLICATE KEY UPDATE password=VALUES(password), email=VALUES(email);

-- 分配管理员角色
INSERT INTO user_roles (user_id, role_id, brand_id)
SELECT u.id, r.id, NULL FROM users u, roles r WHERE u.username = 'admin' AND r.name = 'tenant_admin' AND r.tenant_id = 1
ON DUPLICATE KEY UPDATE user_id=VALUES(user_id);

-- 默认品牌
INSERT INTO brands (id, tenant_id, name, slug, description, website, industry, status) VALUES
(1, 1, 'OpenGEO默认品牌', 'opengeo-default', '系统自动创建的默认品牌', 'https://opengeo.com', '科技', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 默认品牌元数据
INSERT INTO brand_metadata (brand_id, vi_profile, tone_profile, audience_profiles, schema_version) VALUES
(1,
 '{"primary_color": "#1890ff", "secondary_color": "#52c41a", "logo_url": "", "font_family": "PingFang SC", "brand_keywords": ["GEO", "AI优化", "品牌治理"]}',
 '{"formality": "professional", "personality": "friendly", "avoid_words": [], "preferred_phrases": ["AI驱动", "智能化", "高效"]}',
 '{"primary": {"name": "企业用户", "age_range": "25-45", "interests": ["技术", "营销", "AI"], "pain_points": ["内容效率", "品牌一致性"]}}',
 '1.0'
) ON DUPLICATE KEY UPDATE schema_version=VALUES(schema_version);

-- 默认系统配置
INSERT INTO system_configs (config_key, config_value, config_type, description, is_public) VALUES
('system.name', 'OpenGEO BrandOS', 'string', '系统名称', TRUE),
('system.version', '6.0.0', 'string', '系统版本', TRUE),
('system.edition', 'oss', 'string', '版本类型：oss/enterprise', TRUE),
('cloud.api.enabled', 'false', 'boolean', '是否启用云端API', TRUE),
('cloud.api.endpoint', 'https://api.opengeo.com', 'string', '云端API端点', TRUE),
('publish.max_retry', '3', 'number', '最大重试次数', TRUE),
('publish.default_interval', '5', 'number', '默认发布间隔（分钟）', TRUE),
('ai.default_model', 'deepseek', 'string', '默认AI模型', TRUE),
('monitor.check_interval', '24', 'number', '监测检查间隔（小时）', TRUE),
('tenant.default_brand_limit', '5', 'number', '默认品牌数上限', TRUE),
('tenant.default_user_limit', '10', 'number', '默认用户数上限', TRUE)
ON DUPLICATE KEY UPDATE config_value=VALUES(config_value);

-- 默认云端API套餐
INSERT INTO api_plans (name, code, description, monthly_quota, price_cents, overage_price, features, is_active, sort_order) VALUES
('免费版', 'free', '适合个人用户和小型团队', 100, 0, 50, '["attribution_basic", "compliance_basic"]', TRUE, 0),
('基础版', 'starter', '适合成长型品牌', 1000, 9900, 30, '["attribution", "trust_score", "compliance"]', TRUE, 10),
('专业版', 'pro', '适合专业营销团队', 10000, 49900, 20, '["attribution", "trust_score", "compliance", "optimization"]', TRUE, 20),
('企业版', 'enterprise', '适合大型企业和代理商', 100000, 199900, 10, '["attribution", "trust_score", "compliance", "optimization", "bulk_analysis", "priority_support"]', TRUE, 30)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 默认订阅（免费版）
INSERT INTO tenant_subscriptions (tenant_id, plan_id, status, starts_at, api_key)
SELECT 1, id, 'active', NOW(), CONCAT('og_', REPLACE(UUID(), '-', ''))
FROM api_plans WHERE code = 'free'
ON DUPLICATE KEY UPDATE status=VALUES(status);

-- 默认渠道适配器
INSERT INTO publish_channels (tenant_id, channel_type, channel_name, is_enabled) VALUES
(1, 'wechat', '微信公众号', TRUE),
(1, 'weibo', '微博', TRUE),
(1, 'douyin', '抖音', TRUE),
(1, 'xiaohongshu', '小红书', TRUE),
(1, 'zhihu', '知乎', TRUE)
ON DUPLICATE KEY UPDATE channel_name=VALUES(channel_name);

-- 默认插件注册
INSERT INTO plugins (plugin_name, plugin_type, display_name, description, version, author, supported_metrics, is_system, is_enabled) VALUES
('source_score', 'diagnostic', '信源评分', '评估渠道和账号的权威度', '1.0.0', 'OpenGEO', '["authority", "freshness", "relevance"]', TRUE, TRUE),
('competitor_monitor', 'diagnostic', '竞品监测', '监测竞品动态和差距', '1.0.0', 'OpenGEO', '["visibility", "content_gap", "market_share"]', TRUE, TRUE),
('roi_analyzer', 'diagnostic', 'ROI分析', '分析内容投资回报', '1.0.0', 'OpenGEO', '["conversion", "attribution", "cost_analysis"]', TRUE, TRUE),
('semantic_relevance', 'diagnostic', '语义相关性', '评估内容与品牌的语义匹配度', '1.0.0', 'OpenGEO', '["semantic_similarity", "topic_coverage", "intent_match"]', TRUE, TRUE),
('citation_authority', 'diagnostic', '引用权威性', '评估被引用内容的权威性', '1.0.0', 'OpenGEO', '["source_trust", "domain_authority", "citation_quality"]', TRUE, TRUE),
('wechat_adapter', 'channel', '微信适配器', '微信公众号内容发布', '1.0.0', 'OpenGEO', NULL, TRUE, TRUE),
('weibo_adapter', 'channel', '微博适配器', '微博内容发布', '1.0.0', 'OpenGEO', NULL, TRUE, TRUE),
('douyin_adapter', 'channel', '抖音适配器', '抖音短视频/图文发布', '1.0.0', 'OpenGEO', NULL, TRUE, TRUE),
('xiaohongshu_adapter', 'channel', '小红书适配器', '小红书笔记发布', '1.0.0', 'OpenGEO', NULL, TRUE, TRUE),
('zhihu_adapter', 'channel', '知乎适配器', '知乎文章/回答发布', '1.0.0', 'OpenGEO', NULL, TRUE, TRUE)
ON DUPLICATE KEY UPDATE version=VALUES(version);

-- 默认平台
INSERT INTO platforms (code, name, icon, color, description, is_enabled, sort_order) VALUES
('wechat', '微信公众号', 'wechat', 'green', '微信公众号内容发布', TRUE, 100),
('weibo', '微博', 'weibo', 'red', '新浪微博内容发布', TRUE, 90),
('douyin', '抖音', 'douyin', 'purple', '抖音短视频/图文发布', TRUE, 80),
('xiaohongshu', '小红书', 'xiaohongshu', 'pink', '小红书笔记发布', TRUE, 70),
('zhihu', '知乎', 'zhihu', 'blue', '知乎文章/回答发布', TRUE, 60),
('toutiao', '今日头条', 'toutiao', 'orange', '今日头条文章发布', TRUE, 50)
ON DUPLICATE KEY UPDATE name=VALUES(name);
