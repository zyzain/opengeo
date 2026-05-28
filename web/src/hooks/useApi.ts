import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useAuthStore } from "@/stores";
import { useNavigate } from "react-router-dom";

export const queryKeys = {
	contents: ["contents"] as const,
	content: (id: number) => ["contents", id] as const,
	accounts: ["accounts"] as const,
	account: (id: number) => ["accounts", id] as const,
	accountGroups: ["accountGroups"] as const,
	publishTasks: ["publishTasks"] as const,
	publishTask: (id: number) => ["publishTasks", id] as const,
	channels: ["channels"] as const,
	channel: (id: number) => ["channels", id] as const,
	schedules: ["schedules"] as const,
	schedule: (id: number) => ["schedules", id] as const,
	citations: ["citations"] as const,
	scores: ["scores"] as const,
	competitors: ["competitors"] as const,
	competitor: (id: number) => ["competitors", id] as const,
	roi: ["roi"] as const,
	suggestions: ["suggestions"] as const,
	users: ["users"] as const,
	user: (id: number) => ["users", id] as const,
	roles: ["roles"] as const,
	role: (id: number) => ["roles", id] as const,
	tenants: ["tenants"] as const,
	tenant: (id: number) => ["tenants", id] as const,
	systemConfigs: ["systemConfigs"] as const,
	plugins: ["plugins"] as const,
	webhooks: ["webhooks"] as const,
	templates: ["templates"] as const,
	template: (id: number) => ["templates", id] as const,
	fingerprints: ["fingerprints"] as const,
	proxies: ["proxies"] as const,
	knowledgeEntities: ["knowledgeEntities"] as const,
	knowledgeEntity: (id: number) => ["knowledgeEntities", id] as const,
	staggerStrategies: ["staggerStrategies"] as const,
	dedup: ["dedup"] as const,
	plans: ["plans"] as const,
	permissionDefinitions: ["permissionDefinitions"] as const,
};

// ==================== Auth ====================

export function useLogin() {
	const { setAuth } = useAuthStore();
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: (data: { username: string; password: string }) =>
			api.auth.login(data),
		onSuccess: (res) => {
			const { token, refresh_token, user_id, username, email } = res.data.data;
			setAuth(token, refresh_token, { id: user_id, username, email });
			queryClient.clear();
			navigate("/dashboard");
		},
	});
}

export function useRegister() {
	const { setAuth } = useAuthStore();
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: (data: { username: string; password: string; email: string }) =>
			api.auth.register(data),
		onSuccess: (res) => {
			const { token, refresh_token, user_id, username, email } = res.data.data;
			setAuth(token, refresh_token, { id: user_id, username, email });
			queryClient.clear();
			navigate("/dashboard");
		},
	});
}

export function useLogout() {
	const { logout } = useAuthStore();
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	return () => {
		logout();
		queryClient.clear();
		navigate("/auth/login");
	};
}

// ==================== Contents ====================

export function useContents(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.contents, params],
		queryFn: () => api.contents.list(params).then((r) => r.data.data),
	});
}

export function useContent(id: number) {
	return useQuery({
		queryKey: queryKeys.content(id),
		queryFn: () => api.contents.get(id).then((r) => r.data.data),
		enabled: !!id,
	});
}

export function useCreateContent() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.contents.create(data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.contents }),
	});
}

export function useUpdateContent() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.contents.update(id, data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.contents }),
	});
}

export function useDeleteContent() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.contents.delete(id),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.contents }),
	});
}

export function useOptimizeContent() {
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.contents.optimize(id, data),
	});
}

export function useCheckCompliance() {
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data?: any }) =>
			api.contents.compliance(id, data),
	});
}

// ==================== Accounts ====================

export function useAccounts(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.accounts, params],
		queryFn: () => api.accounts.list(params).then((r) => r.data.data),
	});
}

export function useAccount(id: number) {
	return useQuery({
		queryKey: queryKeys.account(id),
		queryFn: () => api.accounts.get(id).then((r) => r.data.data),
		enabled: !!id,
	});
}

export function useCreateAccount() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.accounts.create(data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.accounts }),
	});
}

export function useUpdateAccount() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.accounts.update(id, data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.accounts }),
	});
}

export function useDeleteAccount() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.accounts.delete(id),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.accounts }),
	});
}

// ==================== Account Groups ====================

export function useAccountGroups(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.accountGroups, params],
		queryFn: () => api.accountGroups.list(params).then((r) => r.data.data),
	});
}

export function useCreateAccountGroup() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.accountGroups.create(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.accountGroups }),
	});
}

export function useDeleteAccountGroup() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.accountGroups.delete(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.accountGroups }),
	});
}

export function useAddAccountToGroup() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ groupId, accountId }: { groupId: number; accountId: number }) =>
			api.accountGroups.addAccount(groupId, accountId),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.accountGroups }),
	});
}

export function useRemoveAccountFromGroup() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({
			groupId,
			accountId,
		}: { groupId: number; accountId: number }) =>
			api.accountGroups.removeAccount(groupId, accountId),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.accountGroups }),
	});
}

// ==================== Publish Tasks ====================

export function usePublishTasks(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.publishTasks, params],
		queryFn: () => api.publishTasks.list(params).then((r) => r.data.data),
	});
}

export function useCreatePublishTask() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.publishTasks.create(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.publishTasks }),
	});
}

export function useCancelPublishTask() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.publishTasks.cancel(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.publishTasks }),
	});
}

export function useRetryPublishTask() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.publishTasks.retry(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.publishTasks }),
	});
}

export function usePublishPreview() {
	return useMutation({
		mutationFn: (data: any) => api.publishTasks.preview(data),
	});
}

// ==================== Channels ====================

export function useChannels(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.channels, params],
		queryFn: () => api.channels.list(params).then((r) => r.data.data),
	});
}

export function useCreateChannel() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.channels.create(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.channels }),
	});
}

export function useDeleteChannel() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.channels.delete(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.channels }),
	});
}

// ==================== Schedules ====================

export function useSchedules(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.schedules, params],
		queryFn: () => api.schedules.list(params).then((r) => r.data.data),
	});
}

export function useCreateSchedule() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.schedules.create(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules }),
	});
}

export function useUpdateSchedule() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.schedules.update(id, data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules }),
	});
}

export function useDeleteSchedule() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.schedules.delete(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules }),
	});
}

export function useEnableSchedule() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.schedules.enable(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules }),
	});
}

export function useDisableSchedule() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.schedules.disable(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.schedules }),
	});
}

// ==================== Monitor ====================

export function useAICitations(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.citations, params],
		queryFn: () => api.monitor.citations(params).then((r) => r.data.data),
	});
}

export function useSourceScores(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.scores, params],
		queryFn: () => api.monitor.scores(params).then((r) => r.data.data),
	});
}

export function useCompetitors(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.competitors, params],
		queryFn: () => api.monitor.competitors(params).then((r) => r.data.data),
	});
}

export function useCreateCompetitor() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.monitor.createCompetitor(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.competitors }),
	});
}

export function useDeleteCompetitor() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.monitor.deleteCompetitor(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.competitors }),
	});
}

export function useSyncCompetitor() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.monitor.syncCompetitor(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.competitors }),
	});
}

export function useROIMetrics(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.roi, params],
		queryFn: () => api.monitor.roi(params).then((r) => r.data.data),
	});
}

export function useGenerateSuggestions() {
	return useMutation({
		mutationFn: (contentId: number) => api.monitor.suggestions(contentId),
	});
}

export function useApplySuggestion() {
	return useMutation({
		mutationFn: (id: number) => api.monitor.applySuggestion(id),
	});
}

export function useIgnoreSuggestion() {
	return useMutation({
		mutationFn: (id: number) => api.monitor.ignoreSuggestion(id),
	});
}

export function useListSuggestions(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.suggestions, params],
		queryFn: () => api.monitor.listSuggestions(params).then((r) => r.data.data),
	});
}

// ==================== Users ====================

export function useUsers(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.users, params],
		queryFn: () => api.users.list(params).then((r) => r.data.data),
	});
}

export function useUser(id: number) {
	return useQuery({
		queryKey: queryKeys.user(id),
		queryFn: () => api.users.get(id).then((r) => r.data.data),
		enabled: !!id,
	});
}

export function useUpdateUser() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.users.update(id, data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.users }),
	});
}

export function useDeleteUser() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.users.delete(id),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.users }),
	});
}

// ==================== Roles ====================

export function useRoles(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.roles, params],
		queryFn: () => api.roles.list(params).then((r) => r.data.data),
	});
}

export function useCreateRole() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.roles.create(data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.roles }),
	});
}

export function useUpdateRole() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.roles.update(id, data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.roles }),
	});
}

export function useDeleteRole() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.roles.delete(id),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.roles }),
	});
}

export function useRolePermissions(roleId: number) {
	return useQuery({
		queryKey: [...queryKeys.role(roleId), "permissions"],
		queryFn: () => api.roles.getPermissions(roleId).then((r) => r.data.data),
		enabled: !!roleId,
	});
}

export function useAddRolePermission() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.roles.addPermission(data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.roles }),
	});
}

export function useRemoveRolePermission() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.roles.removePermission(data),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeys.roles }),
	});
}

// ==================== Tenants ====================

export function useTenants(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.tenants, params],
		queryFn: () => api.tenants.list(params).then((r) => r.data.data),
	});
}

export function useCreateTenant() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.tenants.create(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.tenants }),
	});
}

export function useUpdateTenant() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.tenants.update(id, data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.tenants }),
	});
}

export function useDeleteTenant() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.tenants.delete(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.tenants }),
	});
}

// ==================== System ====================

export function useSystemConfigs(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.systemConfigs, params],
		queryFn: () => api.system.configs(params).then((r) => r.data.data),
	});
}

export function useUpdateSystemConfig() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ key, value }: { key: string; value: string }) =>
			api.system.updateConfig(key, value),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.systemConfigs }),
	});
}

export function usePlugins(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.plugins, params],
		queryFn: () => api.system.plugins(params).then((r) => r.data.data),
	});
}

export function useInstallPlugin() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.system.installPlugin(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.plugins }),
	});
}

export function useUpdatePlugin() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.system.updatePlugin(id, data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.plugins }),
	});
}

export function useDeletePlugin() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.system.deletePlugin(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.plugins }),
	});
}

export function useWebhooks(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.webhooks, params],
		queryFn: () => api.system.webhooks(params).then((r) => r.data.data),
	});
}

export function useCreateWebhook() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.system.createWebhook(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.webhooks }),
	});
}

export function useUpdateWebhook() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.system.updateWebhook(id, data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.webhooks }),
	});
}

export function useDeleteWebhook() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.system.deleteWebhook(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.webhooks }),
	});
}

export function useTestWebhook() {
	return useMutation({
		mutationFn: (id: number) => api.system.testWebhook(id),
	});
}

export function useWebhookHistory(id: number) {
	return useQuery({
		queryKey: [...queryKeys.webhooks, id, "history"],
		queryFn: () => api.webhookHistory(id).then((r) => r.data.data),
		enabled: !!id,
	});
}

export function usePlans() {
	return useQuery({
		queryKey: queryKeys.plans,
		queryFn: () => api.system.getPlans().then((r) => r.data.data),
	});
}

export function usePermissionDefinitions() {
	return useQuery({
		queryKey: queryKeys.permissionDefinitions,
		queryFn: () => api.system.getPermissionDefinitions().then((r) => r.data.data),
	});
}

// ==================== Templates ====================

export function useTemplates(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.templates, params],
		queryFn: () => api.templates.list(params).then((r) => r.data.data),
	});
}

export function useCreateTemplate() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.templates.create(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.templates }),
	});
}

export function useUpdateTemplate() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.templates.update(id, data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.templates }),
	});
}

export function useDeleteTemplate() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.templates.delete(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.templates }),
	});
}

// ==================== Fingerprints ====================

export function useFingerprints(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.fingerprints, params],
		queryFn: () => api.fingerprints.list(params).then((r) => r.data.data),
	});
}

export function useCreateFingerprint() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.fingerprints.create(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.fingerprints }),
	});
}

export function useDeleteFingerprint() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.fingerprints.delete(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.fingerprints }),
	});
}

export function useToggleFingerprint() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.fingerprints.toggle(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.fingerprints }),
	});
}

// ==================== Proxies ====================

export function useProxies(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.proxies, params],
		queryFn: () => api.proxies.list(params).then((r) => r.data.data),
	});
}

export function useCreateProxy() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.proxies.create(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.proxies }),
	});
}

export function useDeleteProxy() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.proxies.delete(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.proxies }),
	});
}

export function useCheckProxy() {
	return useMutation({
		mutationFn: (id: number) => api.proxies.check(id),
	});
}

// ==================== Stagger ====================

export function useStaggerStrategies(params?: any) {
	return useQuery({
		queryKey: [...queryKeys.staggerStrategies, params],
		queryFn: () => api.stagger.listStrategies(params).then((r) => r.data.data),
	});
}

export function useCreateStaggerStrategy() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (data: any) => api.stagger.createStrategy(data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.staggerStrategies }),
	});
}

export function useUpdateStaggerStrategy() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.stagger.updateStrategy(id, data),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.staggerStrategies }),
	});
}

export function useDeleteStaggerStrategy() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.stagger.deleteStrategy(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.staggerStrategies }),
	});
}

export function useToggleStaggerStrategy() {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: (id: number) => api.stagger.toggleStrategy(id),
		onSuccess: () =>
			queryClient.invalidateQueries({ queryKey: queryKeys.staggerStrategies }),
	});
}

// ==================== Dedup ====================

export function useCheckDedup() {
	return useMutation({
		mutationFn: (data: { text: string }) => api.dedup.check(data),
	});
}
