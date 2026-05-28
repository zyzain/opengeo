"use client";

import { queryKeys, usePlugins } from "@/hooks";
import { api } from "@/lib/api";
import {
	ApiOutlined,
	AppstoreOutlined,
	CheckCircleOutlined,
	CloseCircleOutlined,
	CloudOutlined,
	DeleteOutlined,
	EditOutlined,
	PlusOutlined,
	SettingOutlined,
	ThunderboltOutlined,
} from "@ant-design/icons";
import { useQueryClient } from "@tanstack/react-query";
import {
	Badge,
	Button,
	Card,
	Col,
	Form,
	Input,
	Modal,
	Popconfirm,
	Row,
	Select,
	Space,
	Statistic,
	Switch,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";

const { Option } = Select;
const { TextArea } = Input;

export default function PluginsPage() {
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [createForm] = Form.useForm();
	const queryClient = useQueryClient();

	const { data, isLoading } = usePlugins();
	const plugins = data?.items || [];

	const handleInstall = async (values: any) => {
		try {
			await api.system.installPlugin(values);
			message.success("插件安装成功");
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: queryKeys.plugins });
		} catch (error: any) {
			message.error(error?.response?.data?.message || "安装失败");
		}
	};

	const handleToggleEnabled = async (record: any, checked: boolean) => {
		try {
			await api.system.updatePlugin(record.id, { is_enabled: checked });
			message.success(checked ? "插件已启用" : "插件已禁用");
			queryClient.invalidateQueries({ queryKey: queryKeys.plugins });
		} catch (error: any) {
			message.error(error?.response?.data?.message || "操作失败");
		}
	};

	const handleUninstall = async (id: number) => {
		try {
			await api.system.deletePlugin(id);
			message.success("插件已卸载");
			queryClient.invalidateQueries({ queryKey: queryKeys.plugins });
		} catch (error: any) {
			message.error(error?.response?.data?.message || "卸载失败");
		}
	};

	const handleSettings = (record: any) => {
		Modal.info({
			title: `${record.plugin_name} - 设置`,
			content: (
				<div>
					<p>插件名称: {record.plugin_name}</p>
					<p>版本: v{record.version}</p>
					<p>作者: {record.author}</p>
					<p>描述: {record.description || "无"}</p>
				</div>
			),
		});
	};

	// 插件类型
	const pluginTypes = [
		{
			value: "channel",
			label: "渠道插件",
			color: "blue",
			icon: <CloudOutlined />,
		},
		{
			value: "ai",
			label: "AI插件",
			color: "purple",
			icon: <ThunderboltOutlined />,
		},
		{
			value: "analyzer",
			label: "分析插件",
			color: "green",
			icon: <ApiOutlined />,
		},
	];

	// 获取插件类型标签
	const getPluginTypeTag = (type: string) => {
		const typeInfo = pluginTypes.find((t) => t.value === type);
		return (
			<Tag color={typeInfo?.color || "default"} icon={typeInfo?.icon}>
				{typeInfo?.label || type}
			</Tag>
		);
	};

	// 表格列定义
	const columns = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
			width: 80,
		},
		{
			title: "插件名称",
			dataIndex: "plugin_name",
			key: "plugin_name",
			render: (text: string) => (
				<Space>
					<AppstoreOutlined />
					<span className="font-medium">{text}</span>
				</Space>
			),
		},
		{
			title: "类型",
			dataIndex: "plugin_type",
			key: "plugin_type",
			width: 120,
			render: (type: string) => getPluginTypeTag(type),
		},
		{
			title: "版本",
			dataIndex: "version",
			key: "version",
			width: 100,
			render: (text: string) => <Tag>v{text}</Tag>,
		},
		{
			title: "作者",
			dataIndex: "author",
			key: "author",
			width: 120,
		},
		{
			title: "描述",
			dataIndex: "description",
			key: "description",
			ellipsis: true,
		},
		{
			title: "状态",
			dataIndex: "is_enabled",
			key: "is_enabled",
			width: 100,
			render: (enabled: boolean) => (
				<Badge
					status={enabled ? "success" : "default"}
					text={enabled ? "已启用" : "已禁用"}
				/>
			),
		},
		{
			title: "安装时间",
			dataIndex: "installed_at",
			key: "installed_at",
			width: 180,
			render: (text: string) => new Date(text).toLocaleString(),
		},
		{
			title: "操作",
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="设置">
						<Button
							type="text"
							icon={<SettingOutlined />}
							onClick={() => handleSettings(record)}
						/>
					</Tooltip>
					<Tooltip title={record.is_enabled ? "禁用" : "启用"}>
						<Switch
							checked={record.is_enabled}
							size="small"
							checkedChildren={<CheckCircleOutlined />}
							unCheckedChildren={<CloseCircleOutlined />}
							onChange={(checked) => handleToggleEnabled(record, checked)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要卸载这个插件吗？"
						okText="确定"
						cancelText="取消"
						onConfirm={() => handleUninstall(record.id)}
					>
						<Tooltip title="卸载">
							<Button type="text" danger icon={<DeleteOutlined />} />
						</Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	// 统计数据
	const stats = {
		total: plugins.length,
		enabled: plugins.filter((p: any) => p.is_enabled).length,
		channel: plugins.filter((p: any) => p.plugin_type === "channel").length,
		ai: plugins.filter((p: any) => p.plugin_type === "ai").length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">插件管理</h1>
				<p className="text-gray-500 mt-1">管理系统扩展插件</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="总插件数"
							value={stats.total}
							prefix={<AppstoreOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="已启用"
							value={stats.enabled}
							prefix={<CheckCircleOutlined />}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="渠道插件"
							value={stats.channel}
							prefix={<CloudOutlined />}
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="AI插件"
							value={stats.ai}
							prefix={<ThunderboltOutlined />}
							valueStyle={{ color: "#722ed1" }}
						/>
					</Card>
				</Col>
			</Row>

			{/* 插件类型说明 */}
			<Card title="插件类型" className="mb-4">
				<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
					{pluginTypes.map((type) => (
						<Card key={type.value} size="small" hoverable>
							<div className="flex items-center space-x-3">
								<div className="text-2xl">{type.icon}</div>
								<div>
									<div className="font-medium">{type.label}</div>
									<div className="text-gray-500 text-sm">
										{type.value === "channel" && "支持多平台内容发布"}
										{type.value === "ai" && "AI模型集成与优化"}
										{type.value === "analyzer" && "数据分析与报告"}
									</div>
								</div>
							</div>
						</Card>
					))}
				</div>
			</Card>

			{/* 插件列表 */}
			<Card
				title={
					<Space>
						<AppstoreOutlined />
						<span>插件列表</span>
					</Space>
				}
				extra={
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => setCreateModalVisible(true)}
					>
						安装插件
					</Button>
				}
			>
				<Table
					columns={columns}
					dataSource={plugins}
					rowKey="id"
					loading={isLoading}
					pagination={{
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
					}}
				/>
			</Card>

			{/* 安装插件弹窗 */}
			<Modal
				title="安装插件"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form form={createForm} layout="vertical" onFinish={handleInstall}>
					<Form.Item
						name="plugin_name"
						label="插件名称"
						rules={[{ required: true, message: "请输入插件名称" }]}
					>
						<Input placeholder="请输入插件名称" />
					</Form.Item>
					<Form.Item
						name="plugin_type"
						label="插件类型"
						rules={[{ required: true, message: "请选择插件类型" }]}
					>
						<Select placeholder="请选择插件类型">
							{pluginTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<Space>
										{type.icon}
										<span>{type.label}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="version" label="版本">
						<Input placeholder="请输入版本号" />
					</Form.Item>
					<Form.Item name="author" label="作者">
						<Input placeholder="请输入作者" />
					</Form.Item>
					<Form.Item name="description" label="描述">
						<TextArea rows={3} placeholder="请输入插件描述" />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">
								安装
							</Button>
							<Button onClick={() => setCreateModalVisible(false)}>取消</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
