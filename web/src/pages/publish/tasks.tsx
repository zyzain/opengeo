"use client";

import {
	useCancelPublishTask,
	useChannels,
	useContents,
	useCreatePublishTask,
	usePublishTasks,
	useRetryPublishTask,
} from "@/hooks";
import api from "@/lib/api";
import {
	CheckCircleOutlined,
	ClockCircleOutlined,
	CloseCircleOutlined,
	ExclamationCircleOutlined,
	EyeOutlined,
	PlusOutlined,
	ReloadOutlined,
	SearchOutlined,
	SendOutlined,
	StopOutlined,
} from "@ant-design/icons";
import {
	Badge,
	Button,
	Card,
	DatePicker,
	Form,
	Input,
	Modal,
	Popconfirm,
	Select,
	Space,
	Steps,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import dayjs from "dayjs";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

const { Option } = Select;

export default function PublishTasksPage() {
	const navigate = useNavigate();
	const [searchForm] = Form.useForm();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [previewModalVisible, setPreviewModalVisible] = useState(false);
	const [previewData, setPreviewData] = useState<any>(null);
	const [createForm] = Form.useForm();

	// 查询参数
	const [queryParams, setQueryParams] = useState({
		page: 1,
		page_size: 10,
		status: undefined,
	});

	const { data, isLoading } = usePublishTasks(queryParams);
	const { data: contentsData } = useContents({ page: 1, page_size: 100 });
	const { data: channelsData } = useChannels({ page: 1, page_size: 100 });
	const createMutation = useCreatePublishTask();
	const cancelMutation = useCancelPublishTask();
	const retryMutation = useRetryPublishTask();

	const tasks = data?.items || [];
	const total = data?.total || 0;
	const contents = contentsData?.items || [];
	const channels = channelsData?.items || [];

	// 任务状态
	const taskStatuses = [
		{
			value: 0,
			label: "待发布",
			color: "processing",
			icon: <ClockCircleOutlined />,
		},
		{ value: 1, label: "发布中", color: "processing", icon: <SendOutlined /> },
		{
			value: 2,
			label: "已发布",
			color: "success",
			icon: <CheckCircleOutlined />,
		},
		{ value: 3, label: "失败", color: "error", icon: <CloseCircleOutlined /> },
		{ value: 4, label: "已取消", color: "default", icon: <StopOutlined /> },
	];

	// 任务优先级
	const taskPriorities = [
		{ value: 2, label: "高", color: "red" },
		{ value: 1, label: "中", color: "orange" },
		{ value: 0, label: "低", color: "blue" },
	];

	// 获取状态标签
	const getStatusTag = (status: number) => {
		const statusInfo = taskStatuses.find((s) => s.value === status);
		return (
			<Tag icon={statusInfo?.icon} color={statusInfo?.color}>
				{statusInfo?.label || "未知"}
			</Tag>
		);
	};

	// 获取优先级标签
	const getPriorityTag = (priority: number) => {
		const priorityInfo = taskPriorities.find((p) => p.value === priority);
		return <Tag color={priorityInfo?.color}>{priorityInfo?.label || "中"}</Tag>;
	};

	// 表格列定义
	const columns = [
		{
			title: "任务ID",
			dataIndex: "id",
			key: "id",
			width: 80,
		},
		{
			title: "内容ID",
			dataIndex: "content_id",
			key: "content_id",
			width: 80,
			render: (id: number) => (
				<a onClick={() => navigate(`/content/${id}`)} className="text-blue-500">
					{id}
				</a>
			),
		},
		{
			title: "渠道ID",
			dataIndex: "channel_id",
			key: "channel_id",
			width: 80,
		},
		{
			title: "状态",
			dataIndex: "status",
			key: "status",
			width: 100,
			render: (status: number) => getStatusTag(status),
		},
		{
			title: "优先级",
			dataIndex: "priority",
			key: "priority",
			width: 80,
			render: (priority: number) => getPriorityTag(priority),
		},
		{
			title: "计划时间",
			dataIndex: "scheduled_time",
			key: "scheduled_time",
			width: 180,
			render: (text: string) => (text ? new Date(text).toLocaleString() : "-"),
		},
		{
			title: "发布时间",
			dataIndex: "published_time",
			key: "published_time",
			width: 180,
			render: (text: string) => (text ? new Date(text).toLocaleString() : "-"),
		},
		{
			title: "重试次数",
			dataIndex: "retry_count",
			key: "retry_count",
			width: 80,
			render: (count: number) => (
				<Badge
					count={count}
					showZero
					style={{ backgroundColor: count > 0 ? "#faad14" : "#d9d9d9" }}
				/>
			),
		},
		{
			title: "错误信息",
			dataIndex: "error_message",
			key: "error_message",
			ellipsis: true,
			render: (text: string) =>
				text ? (
					<Tooltip title={text}>
						<span className="text-red-500">{text}</span>
					</Tooltip>
				) : (
					"-"
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
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="查看">
						<Button
							type="text"
							icon={<EyeOutlined />}
							onClick={() => handlePreview(record)}
						/>
					</Tooltip>
					{record.status === 0 && (
						<Popconfirm
							title="确定要取消这个任务吗？"
							onConfirm={() => handleCancel(record.id)}
							okText="确定"
							cancelText="取消"
						>
							<Tooltip title="取消">
								<Button type="text" danger icon={<StopOutlined />} />
							</Tooltip>
						</Popconfirm>
					)}
					{record.status === 3 && (
						<Popconfirm
							title="确定要重试这个任务吗？"
							onConfirm={() => handleRetry(record.id)}
							okText="确定"
							cancelText="取消"
						>
							<Tooltip title="重试">
								<Button
									type="text"
									icon={<ReloadOutlined />}
									loading={retryMutation.isPending}
								/>
							</Tooltip>
						</Popconfirm>
					)}
				</Space>
			),
		},
	];

	// 创建任务
	const handleCreate = async (values: any) => {
		try {
			const submitData = {
				...values,
				scheduled_time: values.scheduled_time
					? values.scheduled_time.toISOString()
					: null,
			};
			await createMutation.mutateAsync(submitData);
			message.success("创建成功");
			setCreateModalVisible(false);
			createForm.resetFields();
		} catch (error: any) {
			message.error(error.response?.data?.message || "创建失败");
		}
	};

	// 取消任务
	const handleCancel = async (id: number) => {
		try {
			await cancelMutation.mutateAsync(id);
			message.success("取消成功");
		} catch (error: any) {
			message.error(error.response?.data?.message || "取消失败");
		}
	};

	// 重试任务
	const handleRetry = async (id: number) => {
		try {
			await retryMutation.mutateAsync(id);
			message.success("重试任务已提交");
		} catch (error: any) {
			message.error(error.response?.data?.message || "重试失败");
		}
	};

	// 预览发布
	const handlePreview = async (record: any) => {
		try {
			const res = await api.publishTasks.get(record.id);
			setPreviewData({
				...res.data.data,
				content_title: `内容 #${record.content_id}`,
				channel_name: `渠道 #${record.channel_id}`,
			});
			setPreviewModalVisible(true);
		} catch (error: any) {
			message.error("获取预览数据失败");
		}
	};

	// 搜索
	const handleSearch = (values: any) => {
		setQueryParams({ ...queryParams, ...values, page: 1 });
	};

	// 重置搜索
	const handleReset = () => {
		searchForm.resetFields();
		setQueryParams({ page: 1, page_size: 10, status: undefined });
	};

	// 统计数据
	const stats = {
		total: total,
		pending: tasks.filter((t: any) => t.status === 0).length,
		publishing: tasks.filter((t: any) => t.status === 1).length,
		success: tasks.filter((t: any) => t.status === 2).length,
		failed: tasks.filter((t: any) => t.status === 3).length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">发布任务</h1>
				<p className="text-gray-500 mt-1">管理内容发布任务</p>
			</div>

			{/* 统计卡片 */}
			<div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-4">
				<Card size="small">
					<div className="text-center">
						<div className="text-2xl font-bold text-gray-800">
							{stats.total}
						</div>
						<div className="text-gray-500">总任务数</div>
					</div>
				</Card>
				<Card size="small">
					<div className="text-center">
						<div className="text-2xl font-bold text-blue-500">
							{stats.pending}
						</div>
						<div className="text-gray-500">待发布</div>
					</div>
				</Card>
				<Card size="small">
					<div className="text-center">
						<div className="text-2xl font-bold text-orange-500">
							{stats.publishing}
						</div>
						<div className="text-gray-500">发布中</div>
					</div>
				</Card>
				<Card size="small">
					<div className="text-center">
						<div className="text-2xl font-bold text-green-500">
							{stats.success}
						</div>
						<div className="text-gray-500">已发布</div>
					</div>
				</Card>
				<Card size="small">
					<div className="text-center">
						<div className="text-2xl font-bold text-red-500">
							{stats.failed}
						</div>
						<div className="text-gray-500">失败</div>
					</div>
				</Card>
			</div>

			{/* 搜索表单 */}
			<Card className="mb-4">
				<Form form={searchForm} layout="inline" onFinish={handleSearch}>
					<Form.Item name="status" label="状态">
						<Select placeholder="请选择状态" allowClear style={{ width: 120 }}>
							{taskStatuses.map((s) => (
								<Option key={s.value} value={s.value}>
									{s.label}
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								icon={<SearchOutlined />}
								htmlType="submit"
							>
								搜索
							</Button>
							<Button onClick={handleReset}>重置</Button>
						</Space>
					</Form.Item>
				</Form>
			</Card>

			{/* 任务列表 */}
			<Card
				title={
					<Space>
						<SendOutlined />
						<span>发布任务列表</span>
					</Space>
				}
				extra={
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => setCreateModalVisible(true)}
					>
						创建任务
					</Button>
				}
			>
				<Table
					columns={columns}
					dataSource={tasks}
					rowKey="id"
					loading={isLoading}
					scroll={{ x: 1500 }}
					pagination={{
						current: queryParams.page,
						pageSize: queryParams.page_size,
						total,
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
						onChange: (page, pageSize) =>
							setQueryParams({ ...queryParams, page, page_size: pageSize }),
					}}
				/>
			</Card>

			{/* 创建任务弹窗 */}
			<Modal
				title="创建发布任务"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={600}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="content_id"
						label="选择内容"
						rules={[{ required: true, message: "请选择内容" }]}
					>
						<Select
							placeholder="请选择内容"
							showSearch
							optionFilterProp="label"
						>
							{contents.map((content: any) => (
								<Select.Option
									key={content.id}
									value={content.id}
									label={content.title}
								>
									{content.title}
								</Select.Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item
						name="channel_id"
						label="选择渠道"
						rules={[{ required: true, message: "请选择渠道" }]}
					>
						<Select placeholder="请选择渠道">
							{channels.map((channel: any) => (
								<Select.Option key={channel.id} value={channel.id}>
									{channel.channel_name} ({channel.channel_type})
								</Select.Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="scheduled_time" label="计划发布时间">
						<DatePicker
							showTime
							format="YYYY-MM-DD HH:mm:ss"
							placeholder="留空表示立即发布"
							style={{ width: "100%" }}
						/>
					</Form.Item>
					<Form.Item name="priority" label="优先级" initialValue={1}>
						<Select placeholder="请选择优先级">
							{taskPriorities.map((p) => (
								<Select.Option key={p.value} value={p.value}>
									<Tag color={p.color}>{p.label}</Tag>
								</Select.Option>
							))}
						</Select>
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

			{/* 预览弹窗 */}
			<Modal
				title="发布预览"
				open={previewModalVisible}
				onCancel={() => setPreviewModalVisible(false)}
				footer={null}
				width={700}
			>
				{previewData && (
					<div className="space-y-4">
						<div className="p-4 bg-blue-50 rounded-lg">
							<h3 className="font-medium mb-2">任务信息</h3>
							<div className="grid grid-cols-2 gap-2 text-sm">
								<div>任务ID: {previewData.id}</div>
								<div>状态: {getStatusTag(previewData.status)}</div>
								<div>内容: {previewData.content_title}</div>
								<div>渠道: {previewData.channel_name}</div>
								<div>计划时间: {previewData.scheduled_time || "立即发布"}</div>
								<div>重试次数: {previewData.retry_count || 0}</div>
							</div>
						</div>
						<div className="p-4 bg-green-50 rounded-lg">
							<h3 className="font-medium mb-2">校验结果</h3>
							<div className="space-y-2">
								{(previewData?.validation || [
									{ name: '内容完整性校验', passed: previewData?.status === 2 },
									{ name: '渠道配置校验', passed: previewData?.status === 2 },
									{ name: 'Schema标记校验', passed: previewData?.status === 2 },
									{ name: '图片ALT标签校验', passed: previewData?.status === 2 },
								]).map((item: { name: string; passed: boolean }, index: number) => (
									<div key={index} className="flex items-center">
										{item.passed ? (
											<CheckCircleOutlined className="text-green-500 mr-2" />
										) : (
											<CloseCircleOutlined className="text-red-500 mr-2" />
										)}
										<span>{item.name}{item.passed ? '通过' : '未通过'}</span>
									</div>
								))}
							</div>
						</div>
					</div>
				)}
			</Modal>
		</div>
	);
}
