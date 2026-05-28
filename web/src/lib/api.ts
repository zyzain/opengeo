import type { ApiResponse } from "@/types";
import axios, { type AxiosInstance, type AxiosResponse } from "axios";

const apiClient: AxiosInstance = axios.create({
	baseURL: "/api/v1",
	timeout: 30000,
	headers: {
		"Content-Type": "application/json",
	},
});

function getToken(): string | null {
	if (typeof window === "undefined") return null;
	try {
		const raw = localStorage.getItem("auth-storage");
		if (!raw) return null;
		const parsed = JSON.parse(raw);
		return parsed?.state?.token || null;
	} catch {
		return null;
	}
}

function clearAuth() {
	if (typeof window === "undefined") return;
	localStorage.removeItem("auth-storage");
}

apiClient.interceptors.request.use(
	(config) => {
		const token = getToken();
		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => Promise.reject(error),
);

apiClient.interceptors.response.use(
	(response: AxiosResponse<ApiResponse<any>>) => response,
	(error) => {
		if (error.response?.status === 401) {
			const url = error.config?.url || "";
			const isAuthEndpoint = url.startsWith("/auth/");
			if (!isAuthEndpoint) {
				clearAuth();
				if (typeof window !== "undefined") {
					window.location.href = "/auth/login";
				}
			}
		}
		return Promise.reject(error);
	},
);

// API方法
export const api = {
	// 认证
	auth: {
		login: (data: { username: string; password: string }) =>
			apiClient.post<ApiResponse<any>>("/auth/login", data),
		register: (data: { username: string; password: string; email: string }) =>
			apiClient.post<ApiResponse<any>>("/auth/register", data),
		refreshToken: (data: { refresh_token: string }) =>
			apiClient.post<ApiResponse<any>>("/auth/refresh", data),
	},

	// 用户
	users: {
		get: (id: number) => apiClient.get<ApiResponse<any>>(`/users/${id}`),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/users/${id}`, data),
		delete: (id: number) => apiClient.delete<ApiResponse<any>>(`/users/${id}`),
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/users", { params }),
	},

	// 内容
	contents: {
		create: (data: any) => apiClient.post<ApiResponse<any>>("/contents", data),
		get: (id: number) => apiClient.get<ApiResponse<any>>(`/contents/${id}`),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/contents/${id}`, data),
		delete: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/contents/${id}`),
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/contents", { params }),
		optimize: (id: number, data: any) =>
			apiClient.post<ApiResponse<any>>(`/contents/${id}/optimize`, data),
		publish: (id: number, data: any) =>
			apiClient.post<ApiResponse<any>>(`/contents/${id}/publish`, data),
		compliance: (id: number, data?: any) =>
			apiClient.post<ApiResponse<any>>(
				`/contents/${id}/compliance`,
				data || {},
			),
	},

	// 账号
	accounts: {
		create: (data: any) => apiClient.post<ApiResponse<any>>("/accounts", data),
		get: (id: number) => apiClient.get<ApiResponse<any>>(`/accounts/${id}`),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/accounts/${id}`, data),
		delete: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/accounts/${id}`),
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/accounts", { params }),
		health: (id: number) =>
			apiClient.get<ApiResponse<any>>(`/accounts/${id}/health`),
	},

	// 账号分组
	accountGroups: {
		create: (data: any) =>
			apiClient.post<ApiResponse<any>>("/account-groups", data),
		get: (id: number) =>
			apiClient.get<ApiResponse<any>>(`/account-groups/${id}`),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/account-groups/${id}`, data),
		delete: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/account-groups/${id}`),
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/account-groups", { params }),
		addAccount: (groupId: number, accountId: number) =>
			apiClient.post<ApiResponse<any>>(`/account-groups/${groupId}/accounts`, {
				account_id: accountId,
			}),
		removeAccount: (groupId: number, accountId: number) =>
			apiClient.delete<ApiResponse<any>>(
				`/account-groups/${groupId}/accounts/${accountId}`,
			),
	},

	// 发布任务
	publishTasks: {
		create: (data: any) =>
			apiClient.post<ApiResponse<any>>("/publish/tasks", data),
		get: (id: number) =>
			apiClient.get<ApiResponse<any>>(`/publish/tasks/${id}`),
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/publish/tasks", { params }),
		cancel: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/publish/tasks/${id}/cancel`),
		retry: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/publish/tasks/${id}/retry`),
		preview: (data: any) =>
			apiClient.post<ApiResponse<any>>("/publish/preview", data),
		validate: (data: any) =>
			apiClient.post<ApiResponse<any>>("/publish/validate", data),
	},

	// 渠道
	channels: {
		create: (data: any) => apiClient.post<ApiResponse<any>>("/channels", data),
		get: (id: number) => apiClient.get<ApiResponse<any>>(`/channels/${id}`),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/channels/${id}`, data),
		delete: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/channels/${id}`),
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/channels", { params }),
		platforms: () => apiClient.get<ApiResponse<any>>("/channels/platforms"),
		enable: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/channels/${id}/enable`),
		disable: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/channels/${id}/disable`),
	},

	// 调度
	schedules: {
		create: (data: any) => apiClient.post<ApiResponse<any>>("/schedules", data),
		get: (id: number) => apiClient.get<ApiResponse<any>>(`/schedules/${id}`),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/schedules/${id}`, data),
		delete: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/schedules/${id}`),
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/schedules", { params }),
		enable: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/schedules/${id}/enable`),
		disable: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/schedules/${id}/disable`),
	},

	// 监测
	monitor: {
		citations: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/monitor/citations", { params }),
		scores: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/monitor/scores", { params }),
		competitors: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/monitor/competitors", { params }),
		createCompetitor: (data: any) =>
			apiClient.post<ApiResponse<any>>("/monitor/competitors", data),
		deleteCompetitor: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/monitor/competitors/${id}`),
		syncCompetitor: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/monitor/competitors/${id}/sync`),
		ourScore: () =>
			apiClient.get<ApiResponse<any>>("/monitor/competitors/our-score"),
		roi: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/monitor/roi", { params }),
		suggestions: (contentId: number) =>
			apiClient.post<ApiResponse<any>>("/monitor/suggestions/generate", {
				content_id: contentId,
			}),
		applySuggestion: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/monitor/suggestions/${id}/apply`),
		ignoreSuggestion: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/monitor/suggestions/${id}/ignore`),
		listSuggestions: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/monitor/suggestions", { params }),
	},

	// 知识图谱
	knowledge: {
		createEntity: (data: any) =>
			apiClient.post<ApiResponse<any>>("/knowledge/entities", data),
		getEntity: (id: number) =>
			apiClient.get<ApiResponse<any>>(`/knowledge/entities/${id}`),
		updateEntity: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/knowledge/entities/${id}`, data),
		deleteEntity: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/knowledge/entities/${id}`),
		listEntities: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/knowledge/entities", { params }),
		searchEntities: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/knowledge/entities/search", { params }),
	},

	// 系统
	system: {
		configs: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/system/configs", { params }),
		updateConfig: (key: string, value: string) =>
			apiClient.put<ApiResponse<any>>(`/system/configs/${key}`, {
				config_value: value,
			}),
		plugins: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/system/plugins", { params }),
		installPlugin: (data: any) =>
			apiClient.post<ApiResponse<any>>("/system/plugins", data),
		updatePlugin: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/system/plugins/${id}`, data),
		deletePlugin: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/system/plugins/${id}`),
		webhooks: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/system/webhooks", { params }),
		createWebhook: (data: any) =>
			apiClient.post<ApiResponse<any>>("/system/webhooks", data),
		updateWebhook: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/system/webhooks/${id}`, data),
		deleteWebhook: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/system/webhooks/${id}`),
		testWebhook: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/system/webhooks/${id}/test`),
		getPlans: () => apiClient.get<ApiResponse<any>>("/system/plans"),
		getPermissionDefinitions: () =>
			apiClient.get<ApiResponse<any>>("/system/permissions"),
	},

	// 角色
	roles: {
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/roles", { params }),
		create: (data: any) => apiClient.post<ApiResponse<any>>("/roles", data),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/roles/${id}`, data),
		delete: (id: number) => apiClient.delete<ApiResponse<any>>(`/roles/${id}`),
		assign: (data: any) =>
			apiClient.post<ApiResponse<any>>("/roles/assign", data),
		getPermissions: (id: number) =>
			apiClient.get<ApiResponse<any>>(`/roles/${id}/permissions`),
		addPermission: (data: any) =>
			apiClient.post<ApiResponse<any>>("/roles/permissions", data),
		removePermission: (data: any) =>
			apiClient.delete<ApiResponse<any>>("/roles/permissions", { data }),
	},

	// 模板
	templates: {
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/templates", { params }),
		create: (data: any) => apiClient.post<ApiResponse<any>>("/templates", data),
		get: (id: number) => apiClient.get<ApiResponse<any>>(`/templates/${id}`),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/templates/${id}`, data),
		delete: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/templates/${id}`),
	},

	// 指纹管理
	fingerprints: {
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/fingerprints", { params }),
		create: (data: any) =>
			apiClient.post<ApiResponse<any>>("/fingerprints", data),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/fingerprints/${id}`, data),
		delete: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/fingerprints/${id}`),
		toggle: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/fingerprints/${id}/toggle`),
	},

	// 代理IP管理
	proxies: {
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/proxies", { params }),
		create: (data: any) => apiClient.post<ApiResponse<any>>("/proxies", data),
		delete: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/proxies/${id}`),
		check: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/proxies/${id}/check`),
	},

	// 租户管理
	tenants: {
		list: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/tenants", { params }),
		create: (data: any) => apiClient.post<ApiResponse<any>>("/tenants", data),
		get: (id: number) => apiClient.get<ApiResponse<any>>(`/tenants/${id}`),
		update: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/tenants/${id}`, data),
		delete: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/tenants/${id}`),
	},

	// 错峰策略管理
	stagger: {
		listStrategies: (params?: any) =>
			apiClient.get<ApiResponse<any>>("/stagger/strategies", { params }),
		createStrategy: (data: any) =>
			apiClient.post<ApiResponse<any>>("/stagger/strategies", data),
		updateStrategy: (id: number, data: any) =>
			apiClient.put<ApiResponse<any>>(`/stagger/strategies/${id}`, data),
		deleteStrategy: (id: number) =>
			apiClient.delete<ApiResponse<any>>(`/stagger/strategies/${id}`),
		toggleStrategy: (id: number) =>
			apiClient.post<ApiResponse<any>>(`/stagger/strategies/${id}/toggle`),
		getConfig: () => apiClient.get<ApiResponse<any>>("/stagger/config"),
		updateConfig: (data: any) =>
			apiClient.put<ApiResponse<any>>("/stagger/config", data),
	},

	// Webhook历史
	webhookHistory: (id: number) =>
		apiClient.get<ApiResponse<any>>(`/system/webhooks/${id}/history`),

	// 内容去重
	dedup: {
		check: (data: { text: string }) =>
			apiClient.post<ApiResponse<any>>("/publish/dedup/check", data),
	},
};

export default api;
