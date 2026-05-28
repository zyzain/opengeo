import { api } from "@/lib/api";
import { useAuthStore } from "@/stores";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

// Query Keys
export const queryKeys = {
	users: ["users"] as const,
	user: (id: number) => ["users", id] as const,
	contents: ["contents"] as const,
	content: (id: number) => ["contents", id] as const,
	accounts: ["accounts"] as const,
	account: (id: number) => ["accounts", id] as const,
	accountGroups: ["accountGroups"] as const,
	publishTasks: ["publishTasks"] as const,
	channels: ["channels"] as const,
	schedules: ["schedules"] as const,
	citations: ["citations"] as const,
	scores: ["scores"] as const,
	competitors: ["competitors"] as const,
	roi: ["roi"] as const,
	configs: ["configs"] as const,
	plugins: ["plugins"] as const,
	webhooks: ["webhooks"] as const,
	knowledgeEntities: ["knowledgeEntities"] as const,
	knowledgeEntity: (id: number) => ["knowledgeEntities", id] as const,
};

// 认证 Hooks
export function useLogin() {
	const { setAuth } = useAuthStore();
	const navigate = useNavigate();

	return useMutation({
		mutationFn: api.auth.login,
		onSuccess: (response) => {
			const { token, refresh_token, user_id, username, email } =
				response.data.data;
			setAuth(token, refresh_token, { id: user_id, username, email } as any);
			navigate("/dashboard");
		},
	});
}

export function useRegister() {
	const { setAuth } = useAuthStore();
	const navigate = useNavigate();

	return useMutation({
		mutationFn: api.auth.register,
		onSuccess: (response) => {
			const { token, refresh_token, user_id, username, email } =
				response.data.data;
			setAuth(token, refresh_token, { id: user_id, username, email } as any);
			navigate("/dashboard");
		},
	});
}

export function useLogout() {
	const { logout } = useAuthStore();
	const navigate = useNavigate();

	return () => {
		logout();
		navigate("/auth/login");
	};
}

// 内容 Hooks
export function useContents(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.contents, params],
		queryFn: () => api.contents.list(params),
		select: (response) => response.data.data,
	});
}

export function useContent(id: number) {
	return useQuery({
		queryKey: queryKeys.content(id),
		queryFn: () => api.contents.get(id),
		select: (response) => response.data.data,
		enabled: !!id,
	});
}

export function useCreateContent() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.contents.create,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.contents });
		},
	});
}

export function useUpdateContent() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.contents.update(id, data),
		onSuccess: (_, { id }) => {
			queryClient.invalidateQueries({ queryKey: queryKeys.contents });
			queryClient.invalidateQueries({ queryKey: queryKeys.content(id) });
		},
	});
}

export function useDeleteContent() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.contents.delete,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.contents });
		},
	});
}

export function useOptimizeContent() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.contents.optimize(id, data),
		onSuccess: (_, { id }) => {
			queryClient.invalidateQueries({ queryKey: queryKeys.content(id) });
		},
	});
}

// 账号 Hooks
export function useAccounts(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.accounts, params],
		queryFn: () => api.accounts.list(params),
		select: (response) => response.data.data,
	});
}

export function useAccount(id: number) {
	return useQuery({
		queryKey: queryKeys.account(id),
		queryFn: () => api.accounts.get(id),
		select: (response) => response.data.data,
		enabled: !!id,
	});
}

export function useCreateAccount() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.accounts.create,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.accounts });
		},
	});
}

export function useUpdateAccount() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.accounts.update(id, data),
		onSuccess: (_, { id }) => {
			queryClient.invalidateQueries({ queryKey: queryKeys.accounts });
			queryClient.invalidateQueries({ queryKey: queryKeys.account(id) });
		},
	});
}

export function useDeleteAccount() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.accounts.delete,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.accounts });
		},
	});
}

// 账号分组 Hooks
export function useAccountGroups(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.accountGroups, params],
		queryFn: () => api.accountGroups.list(params),
		select: (response) => response.data.data,
	});
}

export function useCreateAccountGroup() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.accountGroups.create,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.accountGroups });
		},
	});
}

// 发布任务 Hooks
export function usePublishTasks(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.publishTasks, params],
		queryFn: () => api.publishTasks.list(params),
		select: (response) => response.data.data,
	});
}

export function useCreatePublishTask() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.publishTasks.create,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.publishTasks });
		},
	});
}

export function useCancelPublishTask() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.publishTasks.cancel,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.publishTasks });
		},
	});
}

// 渠道 Hooks
export function useChannels(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.channels, params],
		queryFn: () => api.channels.list(params),
		select: (response) => response.data.data,
	});
}

export function useCreateChannel() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.channels.create,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.channels });
		},
	});
}

// 调度 Hooks
export function useSchedules(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.schedules, params],
		queryFn: () => api.schedules.list(params),
		select: (response) => response.data.data,
	});
}

export function useCreateSchedule() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.schedules.create,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules });
		},
	});
}

export function useUpdateSchedule() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.schedules.update(id, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules });
		},
	});
}

export function useDeleteSchedule() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.schedules.delete,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules });
		},
	});
}

// 监测 Hooks
export function useAICitations(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.citations, params],
		queryFn: () => api.monitor.citations(params),
		select: (response) => response.data.data,
	});
}

export function useSourceScores(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.scores, params],
		queryFn: () => api.monitor.scores(params),
		select: (response) => response.data.data,
	});
}

export function useCompetitors(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.competitors, params],
		queryFn: () => api.monitor.competitors(params),
		select: (response) => response.data.data,
	});
}

export function useROIMetrics(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.roi, params],
		queryFn: () => api.monitor.roi(params),
		select: (response) => response.data.data,
	});
}

// 系统 Hooks
export function useSystemConfigs(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.configs, params],
		queryFn: () => api.system.configs(params),
		select: (response) => response.data.data,
	});
}

export function usePlugins(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.plugins, params],
		queryFn: () => api.system.plugins(params),
		select: (response) => response.data.data,
	});
}

export function useWebhooks(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.webhooks, params],
		queryFn: () => api.system.webhooks(params),
		select: (response) => response.data.data,
	});
}

// 知识图谱 Hooks
export function useKnowledgeEntities(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.knowledgeEntities, params],
		queryFn: () => api.knowledge.listEntities(params),
		select: (response) => response.data.data,
	});
}

export function useKnowledgeEntity(id: number) {
	return useQuery({
		queryKey: queryKeys.knowledgeEntity(id),
		queryFn: () => api.knowledge.getEntity(id),
		select: (response) => response.data.data,
		enabled: !!id,
	});
}

export function useCreateKnowledgeEntity() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.knowledge.createEntity,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.knowledgeEntities });
		},
	});
}

export function useUpdateKnowledgeEntity() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.knowledge.updateEntity(id, data),
		onSuccess: (_, { id }) => {
			queryClient.invalidateQueries({ queryKey: queryKeys.knowledgeEntities });
			queryClient.invalidateQueries({
				queryKey: queryKeys.knowledgeEntity(id),
			});
		},
	});
}

export function useDeleteKnowledgeEntity() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.knowledge.deleteEntity,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.knowledgeEntities });
		},
	});
}

// 发布任务补充 Hooks
export function useRetryPublishTask() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.publishTasks.retry,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.publishTasks });
		},
	});
}

// 调度补充 Hooks
export function useEnableSchedule() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.schedules.enable,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules });
		},
	});
}

export function useDisableSchedule() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: api.schedules.disable,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules });
		},
	});
}

// 监测补充 Hooks
export function useGenerateSuggestions() {
	return useMutation({
		mutationFn: (contentId: number) => api.monitor.suggestions(contentId),
	});
}

// 合规检测 Hooks
export function useCheckCompliance() {
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data?: any }) =>
			api.contents.compliance(id, data),
	});
}
