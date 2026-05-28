"use client";

import { useLogout } from "@/hooks";
import { useI18n } from "@/i18n";
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
	const { locale, setLocale } = useI18n();
	const {
		token: { colorBgContainer, borderRadiusLG },
	} = theme.useToken();

	const menuItems = [
		{
			key: "/dashboard",
			icon: <DashboardOutlined />,
			label: "仪表盘",
		},
		{
			key: "/content",
			icon: <FileTextOutlined />,
			label: "内容管理",
			children: [
				{
					key: "content_list",
					label: "内容列表",
				},
				{
					key: "/content/knowledge",
					label: "知识图谱",
				},
				{
					key: "/content/templates",
					label: "Prompt模板",
				},
				{
					key: "/content/dedup",
					label: "内容去重",
				},
			],
		},
		{
			key: "/account",
			icon: <TeamOutlined />,
			label: "账号管理",
			children: [
				{
					key: "/account/list",
					label: "账号列表",
				},
				{
					key: "/account/groups",
					label: "账号分组",
				},
				{
					key: "/account/environment",
					label: "环境隔离",
				},
			],
		},
		{
			key: "/publish",
			icon: <SendOutlined />,
			label: "发布管理",
			children: [
				{
					key: "/publish/tasks",
					label: "发布任务",
				},
				{
					key: "/publish/channels",
					label: "渠道管理",
				},
				{
					key: "/publish/stagger",
					label: "错峰发布",
				},
			],
		},
		{
			key: "/schedule",
			icon: <ScheduleOutlined />,
			label: "调度管理",
		},
		{
			key: "/monitor",
			icon: <LineChartOutlined />,
			label: "监测分析",
			children: [
				{
					key: "/monitor/citations",
					label: "AI引用",
				},
				{
					key: "/monitor/scores",
					label: "信源评分",
				},
				{
					key: "/monitor/competitors",
					label: "竞品监测",
				},
				{
					key: "/monitor/roi",
					label: "ROI分析",
				},
				{
					key: "/monitor/suggestions",
					label: "优化建议",
				},
			],
		},
		{
			key: "/settings",
			icon: <SettingOutlined />,
			label: "系统设置",
			children: [
				{
					key: "/settings/configs",
					label: "系统配置",
				},
				{
					key: "/settings/plugins",
					label: "插件管理",
				},
				{
					key: "/settings/webhooks",
					label: "Webhook",
				},
			],
		},
		{
			key: "/system",
			icon: <SafetyOutlined />,
			label: "权限管理",
			children: [
				{
					key: "/system/users",
					label: "用户管理",
				},
				{
					key: "/system/roles",
					label: "角色管理",
				},
				{
					key: "/system/tenants",
					label: "租户管理",
				},
			],
		},
	];

	const userMenuItems = [
		{
			key: "profile",
			icon: <UserOutlined />,
			label: "个人资料",
			onClick: () => navigate("/settings/profile"),
		},
		{
			type: "divider" as const,
		},
		{
			key: "logout",
			icon: <LogoutOutlined />,
			label: "退出登录",
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
									{user?.username || "用户"}
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
