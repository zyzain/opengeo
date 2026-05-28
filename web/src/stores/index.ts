import type { User } from "@/types";
import { create } from "zustand";
import { persist } from "zustand/middleware";

// 认证状态
interface AuthState {
	token: string | null;
	refreshToken: string | null;
	user: User | null;
	isAuthenticated: boolean;
	setAuth: (token: string, refreshToken: string, user: User) => void;
	logout: () => void;
	updateUser: (user: Partial<User>) => void;
}

export const useAuthStore = create<AuthState>()(
	persist(
		(set) => ({
			token: null,
			refreshToken: null,
			user: null,
			isAuthenticated: false,
			setAuth: (token, refreshToken, user) =>
				set({
					token,
					refreshToken,
					user,
					isAuthenticated: true,
				}),
			logout: () =>
				set({
					token: null,
					refreshToken: null,
					user: null,
					isAuthenticated: false,
				}),
			updateUser: (userData) =>
				set((state) => ({
					user: state.user ? { ...state.user, ...userData } : null,
				})),
		}),
		{
			name: "auth-storage",
		},
	),
);

// 应用状态
interface AppState {
	sidebarCollapsed: boolean;
	toggleSidebar: () => void;
	setSidebarCollapsed: (collapsed: boolean) => void;
}

export const useAppStore = create<AppState>((set) => ({
	sidebarCollapsed: false,
	toggleSidebar: () =>
		set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
	setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
}));

// 内容编辑状态
interface ContentEditorState {
	currentContent: any | null;
	isDirty: boolean;
	setCurrentContent: (content: any) => void;
	setIsDirty: (dirty: boolean) => void;
	reset: () => void;
}

export const useContentEditorStore = create<ContentEditorState>((set) => ({
	currentContent: null,
	isDirty: false,
	setCurrentContent: (content) =>
		set({ currentContent: content, isDirty: false }),
	setIsDirty: (dirty) => set({ isDirty: dirty }),
	reset: () => set({ currentContent: null, isDirty: false }),
}));

// 发布状态
interface PublishState {
	selectedContentId: number | null;
	selectedChannelIds: number[];
	scheduledTime: string | null;
	setSelectedContentId: (id: number | null) => void;
	setSelectedChannelIds: (ids: number[]) => void;
	setScheduledTime: (time: string | null) => void;
	reset: () => void;
}

export const usePublishStore = create<PublishState>((set) => ({
	selectedContentId: null,
	selectedChannelIds: [],
	scheduledTime: null,
	setSelectedContentId: (id) => set({ selectedContentId: id }),
	setSelectedChannelIds: (ids) => set({ selectedChannelIds: ids }),
	setScheduledTime: (time) => set({ scheduledTime: time }),
	reset: () =>
		set({
			selectedContentId: null,
			selectedChannelIds: [],
			scheduledTime: null,
		}),
}));
