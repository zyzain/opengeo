"use client";

import { useLogout } from "@/hooks";
import { useI18n } from "@/i18n";
import { useIntl } from "react-intl";
import { useAppStore, useAuthStore } from "@/stores";
import {
	BookOutlined,
	BulbOutlined,
	DashboardOutlined,
	DollarOutlined,
	FileTextOutlined,
	GlobalOutlined,
	LineChartOutlined,
	LogoutOutlined,
	MenuFoldOutlined,
	MenuUnfoldOutlined,
	NodeIndexOutlined,
	SafetyOutlined,
	ScheduleOutlined,
	SendOutlined,
	SettingOutlined,
	StarOutlined,
	TeamOutlined,
	UserOutlined,
	TagsOutlined,
	BranchesOutlined,
	HistoryOutlined,
} from "@ant-design/icons";
import { Avatar, Button, Dropdown, Layout, Menu, theme } from "antd";
import { useState } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";

const { Header, Sider, Content } = Layout;

export default function MainLayout() {
	const navigate = useNavigate();
	const location = useLocation();
	const pathname = location.pathname;
	const { user } = useAuthStore();
	const { sidebarCollapsed, toggleSidebar } = useAppStore();
	const logout = useLogout();
	const intl = useIntl();
	const { locale, setLocale } = useI18n();
	const {
		token: { colorBgContainer, borderRadiusLG },
	} = theme.useToken();

	const menuItems = [
		{
			key: "/dashboard",
			icon: <DashboardOutlined />,
			label: intl.formatMessage({ id: 'nav.dashboard' }),
		},
		{
			key: "/brand",
			icon: <StarOutlined />,
			label: intl.formatMessage({ id: 'nav.brand' }),
			children: [
				{
					key: "/brand",
					label: intl.formatMessage({ id: 'nav.brand' }),
				},
				{
					key: "/brand/metadata",
					label: intl.formatMessage({ id: 'brand.tab.metadata' }),
				},
				{
					key: "/brand/glossary",
					label: intl.formatMessage({ id: 'brand.tab.glossary' }),
				},
				{
					key: "/brand/knowledge",
					label: intl.formatMessage({ id: 'nav.content.knowledge' }),
				},
			],
		},
		{
			key: "/content",
			icon: <FileTextOutlined />,
			label: intl.formatMessage({ id: 'nav.content' }),
			children: [
				{
					key: "content_list",
					label: intl.formatMessage({ id: 'nav.content.list' }),
				},
				{
					key: "/content/knowledge",
					label: intl.formatMessage({ id: 'nav.content.knowledge' }),
				},
				{
					key: "/content/templates",
					label: intl.formatMessage({ id: 'nav.content.templates' }),
				},
				{
					key: "/content/dedup",
					label: intl.formatMessage({ id: 'nav.content.dedup' }),
				},
			],
		},
		{
			key: "/account",
			icon: <TeamOutlined />,
			label: intl.formatMessage({ id: 'nav.account' }),
			children: [
				{
					key: "/account/list",
					label: intl.formatMessage({ id: 'nav.account.list' }),
				},
				{
					key: "/account/groups",
					label: intl.formatMessage({ id: 'nav.account.groups' }),
				},
				{
					key: "/account/environment",
					label: intl.formatMessage({ id: 'nav.account.environment' }),
				},
			],
		},
		{
			key: "/publish",
			icon: <SendOutlined />,
			label: intl.formatMessage({ id: 'nav.publish' }),
			children: [
				{
					key: "/publish/tasks",
					label: intl.formatMessage({ id: 'nav.publish.tasks' }),
				},
				{
					key: "/publish/channels",
					label: intl.formatMessage({ id: 'nav.publish.channels' }),
				},
				{
					key: "/publish/stagger",
					label: intl.formatMessage({ id: 'nav.publish.stagger' }),
				},
			],
		},
		{
			key: "/schedule",
			icon: <ScheduleOutlined />,
			label: intl.formatMessage({ id: 'nav.schedule' }),
		},
		{
			key: "/monitor",
			icon: <LineChartOutlined />,
			label: intl.formatMessage({ id: 'nav.monitor' }),
			children: [
				{
					key: "/monitor/citations",
					label: intl.formatMessage({ id: 'nav.monitor.citations' }),
				},
				{
					key: "/monitor/scores",
					label: intl.formatMessage({ id: 'nav.monitor.scores' }),
				},
				{
					key: "/monitor/competitors",
					label: intl.formatMessage({ id: 'nav.monitor.competitors' }),
				},
				{
					key: "/monitor/roi",
					label: intl.formatMessage({ id: 'nav.monitor.roi' }),
				},
				{
					key: "/monitor/suggestions",
					label: intl.formatMessage({ id: 'nav.monitor.suggestions' }),
				},
			],
		},
		{
			key: "/settings",
			icon: <SettingOutlined />,
			label: intl.formatMessage({ id: 'nav.settings' }),
			children: [
				{
					key: "/settings/configs",
					label: intl.formatMessage({ id: 'nav.settings.configs' }),
				},
				{
					key: "/settings/plugins",
					label: intl.formatMessage({ id: 'nav.settings.plugins' }),
				},
				{
					key: "/settings/webhooks",
					label: intl.formatMessage({ id: 'nav.settings.webhooks' }),
				},
			],
		},
		{
			key: "/system",
			icon: <SafetyOutlined />,
			label: intl.formatMessage({ id: 'nav.system' }),
			children: [
				{
					key: "/system/users",
					label: intl.formatMessage({ id: 'nav.system.users' }),
				},
				{
					key: "/system/roles",
					label: intl.formatMessage({ id: 'nav.system.roles' }),
				},
				{
					key: "/system/tenants",
					label: intl.formatMessage({ id: 'nav.system.tenants' }),
				},
			],
		},
	];

	const userMenuItems = [
		{
			key: "profile",
			icon: <UserOutlined />,
			label: intl.formatMessage({ id: 'nav.profile' }),
			onClick: () => navigate("/settings/profile"),
		},
		{
			type: "divider" as const,
		},
		{
			key: "logout",
			icon: <LogoutOutlined />,
			label: intl.formatMessage({ id: 'nav.logout' }),
			onClick: logout,
		},
	];

	const onMenuClick = ({ key }: { key: string }) => {
		if (key === "content_list") {
			navigate("/content");
		} else {
			navigate(key);
		}
	};

	const getSelectedKeys = () => {
		return [pathname];
	};

	const getOpenKeys = () => {
		const pathParts = pathname.split("/");
		if (pathParts.length > 2) {
			return [`/${pathParts[1]}`];
		}
		return [];
	};

	return (
		<Layout className="min-h-screen">
			<Sider
				trigger={null}
				collapsible
				collapsed={sidebarCollapsed}
				width={256}
				style={{
					overflow: "auto",
					height: "100vh",
					position: "fixed",
					left: 0,
					top: 0,
					bottom: 0,
					zIndex: 10,
				}}
			>
				<div className="h-16 flex items-center justify-center border-b border-gray-700">
					<h1
						className={`text-white font-bold ${sidebarCollapsed ? "text-lg" : "text-xl"}`}
					>
						{sidebarCollapsed ? "OG" : "OpenGEO"}
					</h1>
				</div>
				<Menu
					theme="dark"
					mode="inline"
					selectedKeys={getSelectedKeys()}
					defaultOpenKeys={getOpenKeys()}
					items={menuItems}
					onClick={onMenuClick}
				/>
			</Sider>
			<Layout
				style={{
					marginLeft: sidebarCollapsed ? 80 : 256,
					transition: "all 0.2s",
				}}
			>
				<Header
					style={{
						padding: 0,
						background: colorBgContainer,
						boxShadow: "0 1px 4px rgba(0, 0, 0, 0.08)",
						position: "sticky",
						top: 0,
						zIndex: 9,
						display: "flex",
						alignItems: "center",
						justifyContent: "space-between",
					}}
				>
					<Button
						type="text"
						icon={
							sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />
						}
						onClick={toggleSidebar}
						style={{
							fontSize: "16px",
							width: 64,
							height: 64,
						}}
					/>
					<div className="flex items-center mr-6">
						<Button
							type="text"
							size="small"
							className="mr-2"
							onClick={() => setLocale(locale === "zh-CN" ? "en-US" : "zh-CN")}
						>
							{locale === "zh-CN" ? "EN" : "中"}
						</Button>
						<Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
							<div className="flex items-center cursor-pointer hover:bg-gray-50 px-3 py-2 rounded-lg">
								<Avatar icon={<UserOutlined />} className="mr-2" />
								<span className="text-gray-700">
									{user?.username || intl.formatMessage({ id: 'nav.user' })}
								</span>
							</div>
						</Dropdown>
					</div>
				</Header>
				<Content
					style={{
						margin: "24px 16px",
						padding: 24,
						background: "#f5f5f5",
						minHeight: "calc(100vh - 112px)",
						borderRadius: borderRadiusLG,
					}}
				>
					<div className="fade-in">
						<Outlet />
					</div>
				</Content>
			</Layout>
		</Layout>
	);
}
