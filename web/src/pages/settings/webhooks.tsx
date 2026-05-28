"use client";

import { useWebhooks } from "@/hooks";
import { api } from "@/lib/api";
import {
	ApiOutlined,
	CheckCircleOutlined,
	CloseCircleOutlined,
	DeleteOutlined,
	EditOutlined,
	HistoryOutlined,
	LinkOutlined,
	PlusOutlined,
	SendOutlined,
	ThunderboltOutlined,
} from "@ant-design/icons";
import { useQueryClient } from "@tanstack/react-query";
import {
	Avatar,
	Badge,
	Button,
	Card,
	Col,
	Form,
	Input,
	List,
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

export default function WebhooksPage() {
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [historyModalVisible, setHistoryModalVisible] = useState(false);
	const [selectedWebhook, setSelectedWebhook] = useState<any>(null);
	const [webhookHistory, setWebhookHistory] = useState<any[]>([]);
	const [historyLoading, setHistoryLoading] = useState(false);
	const [createForm] = Form.useForm();

	const queryClient = useQueryClient();
	const { data, isLoading } = useWebhooks();
	const webhooks = data?.items || [];

	// 事件类型
	const eventTypes = [
		{ value: "content.created", label: "内容创建", color: "blue" },
		{ value: "content.published", label: "内容发布", color: "green" },
		{ value: "publish.success", label: "发布成功", color: "green" },
		{ value: "publish.failed", label: "发布失败", color: "red" },
		{ value: "ai.optimized", label: "AI优化完成", color: "purple" },
		{ value: "citation.found", label: "AI引用发现", color: "orange" },
	];

	// 获取事件标签
	const getEventTag = (event: string) => {
		const eventInfo = eventTypes.find((e) => e.value === event);
		return (
			<Tag color={eventInfo?.color || "default"}>
				{eventInfo?.label || event}
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
			title: "Webhook名称",
			dataIndex: "webhook_name",
			key: "webhook_name",
			render: (text: string) => (
				<Space>
					<ApiOutlined />
					<span className="font-medium">{text}</span>
				</Space>
			),
		},
		{
			title: "URL",
			dataIndex: "url",
			key: "url",
			ellipsis: true,
			render: (text: string) => (
				<Tooltip title={text}>
					<a
						href={text}
						target="_blank"
						rel="noopener noreferrer"
						className="text-blue-500"
					>
						<LinkOutlined /> {text.substring(0, 40)}...
					</a>
				</Tooltip>
			),
		},
		{
			title: "订阅事件",
			dataIndex: "events",
			key: "events",
			render: (events: string) => {
				try {
					const eventList = JSON.parse(events);
					return (
						<Space size={[0, 4]} wrap>
							{eventList.slice(0, 2).map((event: string) => getEventTag(event))}
							{eventList.length > 2 && <Tag>+{eventList.length - 2}</Tag>}
						</Space>
					);
				} catch {
					return events;
				}
			},
		},
		{
			title: "状态",
			dataIndex: "is_active",
			key: "is_active",
			width: 100,
			render: (active: boolean) => (
				<Badge
					status={active ? "success" : "default"}
					text={active ? "已启用" : "已禁用"}
				/>
			),
		},
		{
			title: "失败次数",
			dataIndex: "fail_count",
			key: "fail_count",
			width: 80,
			render: (count: number) =>
				count > 0 ? (
					<Badge count={count} style={{ backgroundColor: "#ff4d4f" }} />
				) : (
					<Badge count={0} style={{ backgroundColor: "#d9d9d9" }} />
				),
		},
		{
			title: "最后触发",
			dataIndex: "last_trigger",
			key: "last_trigger",
			width: 180,
			render: (text: string) => (text ? new Date(text).toLocaleString() : "-"),
		},
		{
			title: "操作",
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="测试">
						<Button
							type="text"
							icon={<SendOutlined />}
							onClick={() => handleTest(record.id)}
						/>
					</Tooltip>
					<Tooltip title="历史">
						<Button
							type="text"
							icon={<HistoryOutlined />}
							onClick={() => handleShowHistory(record)}
						/>
					</Tooltip>
					<Tooltip title={record.is_active ? "禁用" : "启用"}>
						<Switch
							checked={record.is_active}
							size="small"
							checkedChildren={<CheckCircleOutlined />}
							unCheckedChildren={<CloseCircleOutlined />}
							onChange={(checked) => handleToggleActive(record, checked)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个Webhook吗？"
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

	// 创建Webhook
	const handleCreate = async (values: any) => {
		try {
			await api.system.createWebhook({
				...values,
				events: JSON.stringify(values.events),
			});
			message.success("创建成功");
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: ["webhooks"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "创建失败");
		}
	};

	// 启用/禁用Webhook
	const handleToggleActive = async (record: any, checked: boolean) => {
		try {
			await api.system.updateWebhook(record.id, { is_active: checked });
			message.success(checked ? "已启用" : "已禁用");
			queryClient.invalidateQueries({ queryKey: ["webhooks"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "操作失败");
		}
	};

	// 删除Webhook
	const handleDelete = async (id: number) => {
		try {
			await api.system.deleteWebhook(id);
			message.success("删除成功");
			queryClient.invalidateQueries({ queryKey: ["webhooks"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "删除失败");
		}
	};

	// 测试Webhook
	const handleTest = async (id: number) => {
		try {
			await api.system.testWebhook(id);
			message.success("测试请求已发送");
		} catch (error: any) {
			message.error(error.response?.data?.message || "测试失败");
		}
	};

	// 显示历史
	const handleShowHistory = async (record: any) => {
		setSelectedWebhook(record);
		setHistoryModalVisible(true);
		setHistoryLoading(true);
		try {
			const res = await api.webhookHistory(record.id);
			setWebhookHistory(res.data?.data || []);
		} catch (error: any) {
			message.error(error.response?.data?.message || "获取历史记录失败");
			setWebhookHistory([]);
		} finally {
			setHistoryLoading(false);
		}
	};

	// 统计数据
	const stats = {
		total: webhooks.length,
		active: webhooks.filter((w: any) => w.is_active).length,
		failed: webhooks.filter((w: any) => w.fail_count > 0).length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">Webhook管理</h1>
				<p className="text-gray-500 mt-1">管理事件通知Webhook</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8}>
					<Card>
						<Statistic
							title="总Webhook数"
							value={stats.total}
							prefix={<ApiOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8}>
					<Card>
						<Statistic
							title="已启用"
							value={stats.active}
							prefix={<CheckCircleOutlined />}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8}>
					<Card>
						<Statistic
							title="有失败记录"
							value={stats.failed}
							prefix={<CloseCircleOutlined />}
							valueStyle={{ color: "#ff4d4f" }}
						/>
					</Card>
				</Col>
			</Row>

			{/* 事件类型说明 */}
			<Card title="支持的事件类型" className="mb-4">
				<div className="flex flex-wrap gap-2">
					{eventTypes.map((event) => (
						<Tag key={event.value} color={event.color}>
							{event.label}
						</Tag>
					))}
				</div>
			</Card>

			{/* Webhook列表 */}
			<Card
				title={
					<Space>
						<ApiOutlined />
						<span>Webhook列表</span>
					</Space>
				}
				extra={
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => setCreateModalVisible(true)}
					>
						创建Webhook
					</Button>
				}
			>
				<Table
					columns={columns}
					dataSource={webhooks}
					rowKey="id"
					loading={isLoading}
					scroll={{ x: 1500 }}
					pagination={{
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
					}}
				/>
			</Card>

			{/* 创建Webhook弹窗 */}
			<Modal
				title="创建Webhook"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={600}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="webhook_name"
						label="Webhook名称"
						rules={[{ required: true, message: "请输入名称" }]}
					>
						<Input placeholder="请输入Webhook名称" />
					</Form.Item>
					<Form.Item
						name="url"
						label="回调URL"
						rules={[
							{ required: true, message: "请输入URL" },
							{ type: "url", message: "请输入有效的URL" },
						]}
					>
						<Input placeholder="https://example.com/webhook" />
					</Form.Item>
					<Form.Item name="secret" label="签名密钥" tooltip="用于验证请求来源">
						<Input.Password placeholder="请输入签名密钥（可选）" />
					</Form.Item>
					<Form.Item
						name="events"
						label="订阅事件"
						rules={[{ required: true, message: "请选择至少一个事件" }]}
					>
						<Select mode="multiple" placeholder="请选择要订阅的事件">
							{eventTypes.map((event) => (
								<Option key={event.value} value={event.value}>
									<Tag color={event.color}>{event.label}</Tag>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">
								创建
							</Button>
							<Button onClick={() => setCreateModalVisible(false)}>取消</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			{/* 历史记录弹窗 */}
			<Modal
				title={`Webhook触发历史 - ${selectedWebhook?.webhook_name}`}
				open={historyModalVisible}
				onCancel={() => setHistoryModalVisible(false)}
				footer={null}
				width={600}
			>
        <List
          loading={historyLoading}
          dataSource={webhookHistory}
          locale={{ emptyText: '暂无历史记录' }}
          renderItem={(item) => (
						<List.Item>
							<List.Item.Meta
								avatar={
									<Avatar
										icon={
											item.success ? (
												<CheckCircleOutlined />
											) : (
												<CloseCircleOutlined />
											)
										}
										style={{
											backgroundColor: item.success ? "#52c41a" : "#ff4d4f",
										}}
									/>
								}
								title={
									<Space>
										{getEventTag(item.event)}
										<Tag color={item.success ? "success" : "error"}>
											HTTP {item.status}
										</Tag>
									</Space>
								}
								description={item.triggered_at}
							/>
						</List.Item>
					)}
				/>
			</Modal>
		</div>
	);
}
