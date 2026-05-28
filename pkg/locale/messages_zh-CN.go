package locale

var zhCNMessages = map[string]string{
	// General
	"success":                 "成功",
	"failed":                  "失败",
	"operation_success":       "操作成功",
	"operation_failed":        "操作失败",
	"internal_error":          "内部错误",
	"invalid_parameter":       "参数错误",
	"unauthorized":            "未授权",
	"forbidden":               "禁止访问",
	"not_found":               "资源不存在",
	"already_exists":          "资源已存在",
	"rate_limit_exceeded":     "请求频率超限，请稍后再试",
	"server_busy":             "服务器繁忙，请稍后再试",

	// Auth
	"login_success":           "登录成功",
	"login_failed":            "登录失败",
	"register_success":        "注册成功",
	"register_failed":         "注册失败",
	"token_refreshed":         "Token已刷新",
	"invalid_credentials":     "用户名或密码错误",
	"user_disabled":           "用户账号已被禁用",
	"user_not_found":          "用户不存在",
	"username_exists":         "用户名已存在",
	"missing_auth_header":     "缺少认证信息",
	"invalid_auth_format":     "认证格式无效",
	"invalid_token":           "Token无效或已过期",
	"password_mismatch":       "两次输入的密码不一致",

	// Content
	"content_created":         "内容创建成功",
	"content_updated":         "内容更新成功",
	"content_deleted":         "内容删除成功",
	"content_not_found":       "内容不存在",
	"content_optimized":       "AI优化完成",
	"content_optimize_failed": "AI优化失败",
	"compliance_passed":       "合规检测通过",
	"compliance_failed":       "合规检测失败",
	"compliance_issues":       "合规检测发现问题，请查看详情",

	// Publish
	"publish_task_created":    "发布任务已创建",
	"publish_cancelled":       "任务已取消",
	"publish_cancel_failed":   "取消失败",
	"publish_retry_submitted": "重试任务已提交",
	"publish_retry_failed":    "重试失败",
	"publish_preview_failed":  "获取预览数据失败",

	// Account
	"account_created":         "账号创建成功",
	"account_updated":         "账号更新成功",
	"account_deleted":         "账号删除成功",
	"account_not_found":       "账号不存在",

	// Brand
	"brand_created":           "品牌创建成功",
	"brand_updated":           "品牌更新成功",
	"brand_deleted":           "品牌删除成功",
	"brand_not_found":         "品牌不存在",
	"metadata_saved":          "元数据保存成功",
	"metadata_save_failed":    "元数据保存失败",
	"entity_added":            "实体添加成功",
	"entity_add_failed":       "添加实体失败",
	"relation_added":          "关系添加成功",
	"relation_add_failed":     "添加关系失败",
	"entity_deleted":          "实体删除成功",
	"entity_delete_failed":    "删除实体失败",
	"snapshot_created":        "快照创建成功",
	"snapshot_create_failed":  "创建快照失败",

	// Monitor
	"sync_created":            "同步任务已创建",
	"sync_failed":             "同步失败",

	// Tenant
	"tenant_created":          "租户创建成功",
	"tenant_updated":          "租户更新成功",
	"tenant_deleted":          "租户已删除",
	"tenant_not_found":        "租户不存在",

	// System
	"plugin_installed":        "插件安装成功",
	"config_updated":          "配置更新成功",
	"role_created":            "角色创建成功",
	"role_updated":            "角色更新成功",
	"role_deleted":            "角色删除成功",
	"permission_updated":      "权限更新成功",
	"permission_update_failed":"权限更新失败",
	"user_created":            "用户创建成功",
	"user_updated":            "用户更新成功",
	"user_deleted":            "用户删除成功",

	// Knowledge
	"entity_name_required":    "实体名称不能为空",
	"entity_type_required":    "实体类型不能为空",

	// Publish validation
	"title_required":          "标题不能为空",
	"body_required":           "正文不能为空",
	"title_too_long":          "标题超过64字符可能被截断",
	"content_library":         "内容库",
	"published_content":       "已发布内容",
	"similar_content_found":   "发现相似内容，建议添加独特的观点和分析",
	"use_synonyms":            "可使用同义词替换部分高频词汇",
	"restructure_content":     "建议调整段落顺序或重新组织内容结构",
	"content_original":        "内容原创度较高，建议保持",
	"content_too_short":       "内容较短，建议补充更多细节",

	// Roles
	"role_admin":              "系统管理员",
	"role_operator":           "运营人员",
	"role_viewer":             "只读用户",

	// Permissions
	"perm_user_create":        "创建用户",
	"perm_user_read":          "查看用户",
	"perm_user_update":        "更新用户",
	"perm_user_delete":        "删除用户",
	"perm_content_create":     "创建内容",
	"perm_content_read":       "查看内容",
	"perm_content_update":     "编辑内容",
	"perm_content_delete":     "删除内容",
	"perm_publish_create":     "创建发布",
	"perm_publish_read":       "查看发布",
	"perm_publish_execute":    "执行发布",
	"perm_role_create":        "创建角色",
	"perm_role_read":          "查看角色",
	"perm_role_update":        "更新角色",
	"perm_role_delete":        "删除角色",
	"perm_tenant_read":        "查看租户",
	"perm_tenant_update":      "更新租户",

	// Permission groups
	"perm_group_content":      "内容管理",
	"perm_group_account":      "账号管理",
	"perm_group_publish":      "发布管理",
	"perm_group_monitor":      "监测管理",
	"perm_group_system":       "系统管理",

	// Plans
	"plan_starter":            "入门版",
	"plan_professional":       "专业版",
	"plan_enterprise":         "企业版",

	// Platforms
	"platform_wechat":         "微信公众号",
	"platform_weibo":          "微博",
	"platform_douyin":         "抖音",
	"platform_xiaohongshu":    "小红书",
	"platform_zhihu":          "知乎",
	"platform_toutiao":        "今日头条",

	// Validation
	"title_empty":             "标题不能为空",
	"body_empty":              "正文不能为空",

	// Startup messages
	"connecting_db":           "正在连接数据库...",
	"db_connect_failed":       "数据库连接失败",
	"db_connect_hint":         "请确保 MySQL 已启动，连接信息可通过 MYSQL_DSN 环境变量配置",
	"db_connected":            "数据库连接成功",
	"db_instance_failed":      "获取数据库实例失败",
	"db_migration_failed":     "数据库迁移失败",
	"tenant_create_failed":    "创建默认租户失败",
	"tenant_check_done":       "默认租户检查完成",
	"perm_create_failed":      "创建默认权限失败",
	"perm_check_done":         "默认权限检查完成",
	"role_create_failed":      "创建默认角色失败",
	"role_check_done":         "默认角色检查完成",
	"role_perm_failed":        "创建角色权限关联失败",
	"role_perm_check_done":    "角色权限关联检查完成",
	"admin_password_required": "错误: 生产环境必须设置 ADMIN_PASSWORD 环境变量",
	"admin_password_failed":   "生成管理员密码失败",
	"admin_create_failed":     "创建管理员账号失败",
	"admin_role_failed":       "分配管理员角色失败",
	"admin_check_done":        "管理员账号检查完成",
	"content_table_failed":    "内容表迁移失败",
	"knowledge_table_failed":  "知识图谱表迁移失败",
	"knowledge_seed_failed":   "知识图谱种子数据创建失败",
	"knowledge_seed_done":     "知识图谱种子数据检查完成",
	"gateway_shutting_down":   "正在关闭Gateway服务...",

	// Content type
	"optimized_content":       "优化后的内容",
	"content_structure_good":  "内容结构良好",
	"publish_task_created_msg":"发布任务已创建",
	"content_compliant":       "内容合规",
	"plugin_install_success":  "插件安装成功",
}
