"use client";

import { useChannels, useCreateChannel } from "@/hooks";
import api from "@/lib/api";
import {
	CheckCircleOutlined,
	CloseCircleOutlined,
	DeleteOutlined,
	EditOutlined,
	GlobalOutlined,
	PlusOutlined,
	SettingOutlined,
} from "@ant-design/icons";
import { useQueryClient } from "@tanstack/react-query";
import {
	Badge,
	Button,
	Card,
	Form,
	Input,
	Modal,
	Popconfirm,
	Select,
	Space,
	Switch,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";

const { Option } = Select;
const { TextArea } = Input;

export default function ChannelsPage() {
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [createForm] = Form.useForm();
	const queryClient = useQueryClient();

	const { data, isLoading } = useChannels();
	const createMutation = useCreateChannel();

	const channels = data?.items || [];

	const handleToggleEnabled = async (record: any) => {
		try {
			if (record.is_enabled) {
				await api.channels.disable(record.id);
			} else {
				await api.channels.enable(record.id);
			}
			message.success(record.is_enabled ? "已禁用" : "已启用");
			queryClient.invalidateQueries({ queryKey: ["channels"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "操作失败");
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.channels.delete(id);
			message.success("删除成功");
			queryClient.invalidateQueries({ queryKey: ["channels"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "删除失败");
		}
	};

	const handleSettings = (record: any) => {
		message.info(`渠道设置功能开发中: ${record.channel_name}`);
	};

	// 支持的平台
	const supportedPlatforms = [
		{ value: "wechat", label: "微信公众号", color: "green", icon: "🟢" },
		{ value: "weibo", label: "微博", color: "red", icon: "🔴" },
		{ value: "douyin", label: "抖音", color: "purple", icon: "🟣" },
		{ value: "xiaohongshu", label: "小红书", color: "pink", icon: "🩷" },
		{ value: "zhihu", label: "知乎", color: "blue", icon: "🔵" },
		{ value: "toutiao", label: "今日头条", color: "orange", icon: "🟠" },
	];

	// 获取平台标签
	const getPlatformTag = (platform: string) => {
		const platformInfo = supportedPlatforms.find((p) => p.value === platform);
		return (
			<Tag color={platformInfo?.color || "default"}>
				{platformInfo?.icon} {platformInfo?.label || platform}
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
			title: "渠道名称",
			dataIndex: "channel_name",
			key: "channel_name",
			render: (text: string, record: any) => (
				<Space>
					<GlobalOutlined />
					<span className="font-medium">{text}</span>
				</Space>
			),
		},
		{
			title: "平台类型",
			dataIndex: "channel_type",
			key: "channel_type",
			width: 150,
			render: (type: string) => getPlatformTag(type),
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
			title: "创建时间",
			dataIndex: "created_at",
			key: "created_at",
			width: 180,
			render: (text: string) => new Date(text).toLocaleString(),
		},
		{
			title: "操作",
			key: "action",
			width: 200,
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
							onChange={() => handleToggleEnabled(record)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个渠道吗？"
						okText="确定"
						cancelText="取消"
						onConfirm={() => handleDelete(record.id)}
					>
						<Tooltip title="删除">
							<Button type="text" danger icon={<DeleteOutlined />} />
						</Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	// 创建渠道
	const handleCreate = async (values: any) => {
		try {
			await createMutation.mutateAsync(values);
			message.success("创建成功");
			setCreateModalVisible(false);
			createForm.resetFields();
		} catch (error: any) {
			message.error(error.response?.data?.message || "创建失败");
		}
	};

	// 统计数据
	const stats = {
		total: channels.length,
		enabled: channels.filter((c: any) => c.is_enabled).length,
		disabled: channels.filter((c: any) => !c.is_enabled).length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">渠道管理</h1>
				<p className="text-gray-500 mt-1">管理发布渠道配置</p>
			</div>

			{/* 统计卡片 */}
			<div className="grid grid-cols-3 gap-4 mb-4">
				<Card>
					<div className="text-center">
						<div className="text-2xl font-bold text-gray-800">
							{stats.total}
						</div>
						<div className="text-gray-500">总渠道数</div>
					</div>
				</Card>
				<Card>
					<div className="text-center">
						<div className="text-2xl font-bold text-green-500">
							{stats.enabled}
						</div>
						<div className="text-gray-500">已启用</div>
					</div>
				</Card>
				<Card>
					<div className="text-center">
						<div className="text-2xl font-bold text-gray-400">
							{stats.disabled}
						</div>
						<div className="text-gray-500">已禁用</div>
					</div>
				</Card>
			</div>

			{/* 支持的平台 */}
			<Card title="支持的平台" className="mb-4">
				<div className="flex flex-wrap gap-4">
					{supportedPlatforms.map((platform) => (
						<Card key={platform.value} size="small" hoverable className="w-40">
							<div className="text-center">
								<div className="text-2xl mb-2">{platform.icon}</div>
								<div className="font-medium">{platform.label}</div>
								<Tag color={platform.color} className="mt-2">
									{
										channels.filter(
											(c: any) => c.channel_type === platform.value,
										).length
									}{" "}
									个渠道
								</Tag>
							</div>
						</Card>
					))}
				</div>
			</Card>

			{/* 渠道列表 */}
			<Card
				title={
					<Space>
						<GlobalOutlined />
						<span>渠道列表</span>
					</Space>
				}
				extra={
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => setCreateModalVisible(true)}
					>
						添加渠道
					</Button>
				}
			>
				<Table
					columns={columns}
					dataSource={channels}
					rowKey="id"
					loading={isLoading}
					pagination={{
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
					}}
				/>
			</Card>

			{/* 创建渠道弹窗 */}
			<Modal
				title="添加渠道"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="channel_type"
						label="平台类型"
						rules={[{ required: true, message: "请选择平台类型" }]}
					>
						<Select placeholder="请选择平台类型">
							{supportedPlatforms.map((platform) => (
								<Option key={platform.value} value={platform.value}>
									<Space>
										<span>{platform.icon}</span>
										<span>{platform.label}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item
						name="channel_name"
						label="渠道名称"
						rules={[{ required: true, message: "请输入渠道名称" }]}
					>
						<Input placeholder="请输入渠道名称" />
					</Form.Item>
					<Form.Item
						name="channel_config"
						label="渠道配置"
						rules={[{ required: true, message: "请输入渠道配置" }]}
					>
						<TextArea
							rows={4}
							placeholder='请输入JSON格式的配置，例如：{"app_id": "xxx", "app_secret": "xxx"}'
						/>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								htmlType="submit"
								loading={createMutation.isPending}
							>
								创建
							</Button>
							<Button onClick={() => setCreateModalVisible(false)}>取消</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
