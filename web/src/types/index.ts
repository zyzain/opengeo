// API响应类型
export interface ApiResponse<T> {
	code: number;
	message: string;
	data: T;
}

// 分页响应
export interface PaginatedResponse<T> {
	items: T[];
	total: number;
	page: number;
	page_size: number;
}

// 用户类型
export interface User {
	id: number;
	username: string;
	email: string;
	status: number;
	created_at: string;
	updated_at: string;
}

// 登录请求
export interface LoginRequest {
	username: string;
	password: string;
}

// 登录响应
export interface LoginResponse {
	token: string;
	refresh_token: string;
	user_id: number;
	username: string;
	email: string;
}

// 注册请求
export interface RegisterRequest {
	username: string;
	password: string;
	email: string;
}

// 内容类型
export interface Content {
	id: number;
	user_id: number;
	title: string;
	body: string;
	content_type: string;
	status: number;
	schema_markup: string;
	ai_optimization_score: number;
	created_at: string;
	updated_at: string;
}

// 账号类型
export interface Account {
	id: number;
	user_id: number;
	platform: string;
	account_name: string;
	account_id: string;
	status: number;
	health_score: number;
	created_at: string;
	updated_at: string;
}

// 账号分组类型
export interface AccountGroup {
	id: number;
	user_id: number;
	name: string;
	parent_id: number | null;
	group_type: string;
	description: string;
	created_at: string;
}

// 发布任务类型
export interface PublishTask {
	id: number;
	user_id: number;
	content_id: number;
	channel_id: number;
	status: number;
	scheduled_time: string | null;
	published_time: string | null;
	retry_count: number;
	error_message: string;
	created_at: string;
}

// 渠道类型
export interface Channel {
	id: number;
	user_id: number;
	channel_type: string;
	channel_name: string;
	channel_config: string;
	is_enabled: boolean;
	created_at: string;
}

// 调度类型
export interface Schedule {
	id: number;
	user_id: number;
	schedule_name: string;
	schedule_type: string;
	cron_expression: string;
	config: string;
	is_enabled: boolean;
	next_run_time: string | null;
	last_run_time: string | null;
	run_count: number;
	created_at: string;
}

// AI引用类型
export interface AICitation {
	id: number;
	content_id: number;
	ai_model: string;
	query_text: string;
	is_cited: boolean;
	citation_position: number;
	citation_text: string;
	sentiment: string;
	tracked_at: string;
}

// 信源评分类型
export interface SourceScore {
	id: number;
	channel_id: number;
	account_id: number;
	score: number;
	score_dimensions: string;
	updated_at: string;
}

// 系统配置类型
export interface SystemConfig {
	id: number;
	config_key: string;
	config_value: string;
	config_type: string;
	description: string;
	is_public: boolean;
	updated_at: string;
}

// 插件类型
export interface Plugin {
	id: number;
	plugin_name: string;
	plugin_type: string;
	description: string;
	version: string;
	author: string;
	is_enabled: boolean;
	installed_at: string;
}

// 知识图谱实体类型
export interface KnowledgeEntity {
	id: number;
	user_id: number;
	entity_name: string;
	entity_type: string;
	entity_data: string;
	authority_links: string;
	content_count: number;
	created_at: string;
	updated_at: string;
}

// Webhook类型
export interface Webhook {
	id: number;
	user_id: number;
	webhook_name: string;
	url: string;
	events: string;
	is_active: boolean;
	last_trigger: string | null;
	created_at: string;
}

// 查询参数
export interface PaginationParams {
	page?: number;
	page_size?: number;
}

export interface ContentQueryParams extends PaginationParams {
	user_id?: number;
	content_type?: string;
	status?: number;
}

export interface AccountQueryParams extends PaginationParams {
	user_id?: number;
	platform?: string;
}

export interface PublishTaskQueryParams extends PaginationParams {
	user_id?: number;
	status?: number;
}
