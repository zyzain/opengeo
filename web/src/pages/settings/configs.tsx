"use client";

import { useSystemConfigs } from "@/hooks";
import api from "@/lib/api";
import {
	CopyOutlined,
	DeleteOutlined,
	EditOutlined,
	GlobalOutlined,
	LockOutlined,
	PlusOutlined,
	ReloadOutlined,
	SaveOutlined,
	SettingOutlined,
} from "@ant-design/icons";
import {
	Badge,
	Button,
	Card,
	Descriptions,
	Form,
	Input,
	Modal,
	Popconfirm,
	Select,
	Space,
	Switch,
	Table,
	Tabs,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";

const { Option } = Select;
const { TextArea } = Input;

export default function ConfigsPage() {
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [editingConfig, setEditingConfig] = useState<any>(null);
	const [editForm] = Form.useForm();

	const { data, isLoading, refetch } = useSystemConfigs();
	const configs = data?.items || [];

	// 配置类型
	const configTypes = [
		{ value: "string", label: "字符串", color: "blue" },
		{ value: "number", label: "数字", color: "green" },
		{ value: "json", label: "JSON", color: "purple" },
		{ value: "boolean", label: "布尔值", color: "orange" },
	];

	// 获取配置类型标签
	const getConfigTypeTag = (type: string) => {
		const typeInfo = configTypes.find((t) => t.value === type);
		return (
			<Tag color={typeInfo?.color || "default"}>{typeInfo?.label || type}</Tag>
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
			title: "配置键",
			dataIndex: "config_key",
			key: "config_key",
			render: (text: string) => (
				<code className="bg-gray-100 px-2 py-1 rounded text-sm">{text}</code>
			),
		},
		{
			title: "配置值",
			dataIndex: "config_value",
			key: "config_value",
			ellipsis: true,
			render: (text: string, record: any) => {
				if (record.config_type === "boolean") {
					return <Switch checked={text === "true"} disabled size="small" />;
				}
				if (record.config_type === "json") {
					return (
						<Tooltip
							title={<pre>{JSON.stringify(JSON.parse(text), null, 2)}</pre>}
						>
							<span className="text-blue-500 cursor-pointer">
								{text.substring(0, 50)}...
							</span>
						</Tooltip>
					);
				}
				return text;
			},
		},
		{
			title: "类型",
			dataIndex: "config_type",
			key: "config_type",
			width: 100,
			render: (type: string) => getConfigTypeTag(type),
		},
		{
			title: "描述",
			dataIndex: "description",
			key: "description",
			ellipsis: true,
		},
		{
			title: "公开",
			dataIndex: "is_public",
			key: "is_public",
			width: 80,
			render: (isPublic: boolean) =>
				isPublic ? (
					<Badge status="success" text="是" />
				) : (
					<Badge status="default" text="否" />
				),
		},
		{
			title: "更新时间",
			dataIndex: "updated_at",
			key: "updated_at",
			width: 180,
			render: (text: string) => new Date(text).toLocaleString(),
		},
		{
			title: "操作",
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="编辑">
						<Button
							type="text"
							icon={<EditOutlined />}
							onClick={() => handleEdit(record)}
						/>
					</Tooltip>
					<Tooltip title="复制">
						<Button
							type="text"
							icon={<CopyOutlined />}
							onClick={() => handleCopy(record)}
						/>
					</Tooltip>
				</Space>
			),
		},
	];

	// 编辑配置
	const handleEdit = (record: any) => {
		setEditingConfig(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await api.system.updateConfig(values.config_key, values.config_value);
			message.success("更新成功");
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingConfig(null);
			refetch();
		} catch (error: any) {
			message.error("更新失败");
		}
	};

	// 复制配置
	const handleCopy = (record: any) => {
		navigator.clipboard.writeText(
			`${record.config_key}=${record.config_value}`,
		);
		message.success("已复制到剪贴板");
	};

	// 配置分类
	const configCategories = [
		{
			title: "系统配置",
			icon: <SettingOutlined />,
			configs: configs.filter((c: any) => c.config_key.startsWith("system.")),
		},
		{
			title: "发布配置",
			icon: <GlobalOutlined />,
			configs: configs.filter((c: any) => c.config_key.startsWith("publish.")),
		},
		{
			title: "AI配置",
			icon: <LockOutlined />,
			configs: configs.filter((c: any) => c.config_key.startsWith("ai.")),
		},
		{
			title: "监控配置",
			icon: <SettingOutlined />,
			configs: configs.filter((c: any) => c.config_key.startsWith("monitor.")),
		},
	];

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">系统配置</h1>
				<p className="text-gray-500 mt-1">管理系统全局配置参数</p>
			</div>

			{/* 配置分类 */}
			<Tabs
				defaultActiveKey="all"
				items={[
					{
						key: "all",
						label: "所有配置",
						children: (
							<Card
								extra={
									<Space>
										<Button icon={<ReloadOutlined />} onClick={() => refetch()}>
											刷新
										</Button>
									</Space>
								}
							>
								<Table
									columns={columns}
									dataSource={configs}
									rowKey="id"
									loading={isLoading}
									pagination={{
										showSizeChanger: true,
										showQuickJumper: true,
										showTotal: (total) => `共 ${total} 条`,
									}}
								/>
							</Card>
						),
					},
					...configCategories.map((category) => ({
						key: category.title,
						label: (
							<Space>
								{category.icon}
								<span>{category.title}</span>
								<Badge
									count={category.configs.length}
									style={{ backgroundColor: "#1890ff" }}
								/>
							</Space>
						),
						children: (
							<Card>
								<Table
									columns={columns}
									dataSource={category.configs}
									rowKey="id"
									pagination={false}
								/>
							</Card>
						),
					})),
				]}
			/>

			{/* 编辑配置弹窗 */}
			<Modal
				title="编辑配置"
				open={editModalVisible}
				onCancel={() => {
					setEditModalVisible(false);
					setEditingConfig(null);
				}}
				footer={null}
				width={500}
			>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item name="config_key" label="配置键">
						<Input disabled />
					</Form.Item>
					<Form.Item
						name="config_value"
						label="配置值"
						rules={[{ required: true, message: "请输入配置值" }]}
					>
						{editingConfig?.config_type === "json" ? (
							<TextArea rows={4} placeholder="请输入JSON格式的配置值" />
						) : editingConfig?.config_type === "boolean" ? (
							<Select>
								<Option value="true">是</Option>
								<Option value="false">否</Option>
							</Select>
						) : (
							<Input placeholder="请输入配置值" />
						)}
					</Form.Item>
					<Form.Item name="description" label="描述">
						<TextArea rows={2} placeholder="请输入配置描述" />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit" icon={<SaveOutlined />}>
								保存
							</Button>
							<Button onClick={() => setEditModalVisible(false)}>取消</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
