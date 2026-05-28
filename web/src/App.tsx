import AuthGuard from "@/components/layout/AuthGuard";
import MainLayout from "@/components/layout/MainLayout";
import { I18nProvider } from "@/i18n";
import HomePage from "@/pages";
import AccountEnvironmentPage from "@/pages/account/environment";
import AccountGroupsPage from "@/pages/account/groups";
import AccountListPage from "@/pages/account/list";
import LoginPage from "@/pages/auth/login";
import BrandListPage from "@/pages/brand";
import BrandDetailPage from "@/pages/brand/detail";
import BrandMetadataPage from "@/pages/brand/metadata";
import BrandGlossaryPage from "@/pages/brand/glossary";
import BrandKnowledgePage from "@/pages/brand/knowledge";
import ContentListPage from "@/pages/content";
import ContentDedupPage from "@/pages/content/dedup";
import ContentKnowledgePage from "@/pages/content/knowledge";
import ContentTemplatesPage from "@/pages/content/templates";
import DashboardPage from "@/pages/dashboard";
import MonitorCitationsPage from "@/pages/monitor/citations";
import MonitorCompetitorsPage from "@/pages/monitor/competitors";
import MonitorRoiPage from "@/pages/monitor/roi";
import MonitorScoresPage from "@/pages/monitor/scores";
import MonitorSuggestionsPage from "@/pages/monitor/suggestions";
import PublishChannelsPage from "@/pages/publish/channels";
import PublishStaggerPage from "@/pages/publish/stagger";
import PublishTasksPage from "@/pages/publish/tasks";
import SchedulePage from "@/pages/schedule";
import SettingsConfigsPage from "@/pages/settings/configs";
import SettingsPluginsPage from "@/pages/settings/plugins";
import SettingsWebhooksPage from "@/pages/settings/webhooks";
import SystemRolesPage from "@/pages/system/roles";
import SystemTenantsPage from "@/pages/system/tenants";
import SystemUsersPage from "@/pages/system/users";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { App as AntApp, ConfigProvider } from "antd";
import { useState } from "react";
import {
	Navigate,
	Outlet,
	RouterProvider,
	createBrowserRouter,
} from "react-router-dom";

const router = createBrowserRouter([
	{
		path: "/",
		element: <Outlet />,
		children: [
			{ index: true, element: <HomePage /> },
			{ path: "auth/login", element: <LoginPage /> },
			{
				element: <AuthGuard />,
				children: [
					{
						element: <MainLayout />,
						children: [
							{ path: "dashboard", element: <DashboardPage /> },
							{ path: "brand", element: <BrandListPage /> },
							{ path: "brand/:id", element: <BrandDetailPage /> },
							{ path: "brand/metadata", element: <BrandMetadataPage /> },
							{ path: "brand/glossary", element: <BrandGlossaryPage /> },
							{ path: "brand/knowledge", element: <BrandKnowledgePage /> },
							{ path: "content", element: <ContentListPage /> },
							{ path: "content/knowledge", element: <ContentKnowledgePage /> },
							{ path: "content/templates", element: <ContentTemplatesPage /> },
							{ path: "content/dedup", element: <ContentDedupPage /> },
							{ path: "account/list", element: <AccountListPage /> },
							{ path: "account/groups", element: <AccountGroupsPage /> },
							{
								path: "account/environment",
								element: <AccountEnvironmentPage />,
							},
							{ path: "publish/tasks", element: <PublishTasksPage /> },
							{ path: "publish/channels", element: <PublishChannelsPage /> },
							{ path: "publish/stagger", element: <PublishStaggerPage /> },
							{ path: "schedule", element: <SchedulePage /> },
							{ path: "monitor/citations", element: <MonitorCitationsPage /> },
							{ path: "monitor/scores", element: <MonitorScoresPage /> },
							{
								path: "monitor/competitors",
								element: <MonitorCompetitorsPage />,
							},
							{ path: "monitor/roi", element: <MonitorRoiPage /> },
							{
								path: "monitor/suggestions",
								element: <MonitorSuggestionsPage />,
							},
							{ path: "settings/configs", element: <SettingsConfigsPage /> },
							{ path: "settings/plugins", element: <SettingsPluginsPage /> },
							{ path: "settings/webhooks", element: <SettingsWebhooksPage /> },
							{ path: "system/users", element: <SystemUsersPage /> },
							{ path: "system/roles", element: <SystemRolesPage /> },
							{ path: "system/tenants", element: <SystemTenantsPage /> },
						],
					},
				],
			},
			{ path: "*", element: <Navigate to="/" replace /> },
		],
	},
]);

export default function App() {
	const [queryClient] = useState(
		() =>
			new QueryClient({
				defaultOptions: {
					queries: {
						staleTime: 60 * 1000,
						retry: 1,
						refetchOnWindowFocus: false,
					},
				},
			}),
	);

	return (
		<ConfigProvider>
			<AntApp>
				<I18nProvider>
					<QueryClientProvider client={queryClient}>
						<RouterProvider router={router} />
					</QueryClientProvider>
				</I18nProvider>
			</AntApp>
		</ConfigProvider>
	);
}
