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
import { useIntl } from "react-intl";

const { Option } = Select;

export default function PublishTasksPage() {
	const intl = useIntl();
	const navigate = useNavigate();
	const [searchForm] = Form.useForm();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [previewModalVisible, setPreviewModalVisible] = useState(false);
	const [previewData, setPreviewData] = useState<any>(null);
	const [createForm] = Form.useForm();

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

	const taskStatuses = [
		{ value: 0, label: intl.formatMessage({ id: 'publish.status.pending' }), color: "processing", icon: <ClockCircleOutlined /> },
		{ value: 1, label: intl.formatMessage({ id: 'publish.status.publishing' }), color: "processing", icon: <SendOutlined /> },
		{ value: 2, label: intl.formatMessage({ id: 'publish.status.published' }), color: "success", icon: <CheckCircleOutlined /> },
		{ value: 3, label: intl.formatMessage({ id: 'publish.status.failed' }), color: "error", icon: <CloseCircleOutlined /> },
		{ value: 4, label: intl.formatMessage({ id: 'publish.status.cancelled' }), color: "default", icon: <StopOutlined /> },
	];

	const taskPriorities = [
		{ value: 2, label: intl.formatMessage({ id: 'common.priority.high' }), color: "red" },
		{ value: 1, label: intl.formatMessage({ id: 'common.priority.medium' }), color: "orange" },
		{ value: 0, label: intl.formatMessage({ id: 'common.priority.low' }), color: "blue" },
	];

	const getStatusTag = (status: number) => {
		const statusInfo = taskStatuses.find((s) => s.value === status);
		return (
			<Tag icon={statusInfo?.icon} color={statusInfo?.color}>
				{statusInfo?.label || intl.formatMessage({ id: 'common.status.unknown' })}
			</Tag>
		);
	};

	const getPriorityTag = (priority: number) => {
		const priorityInfo = taskPriorities.find((p) => p.value === priority);
		return <Tag color={priorityInfo?.color}>{priorityInfo?.label || intl.formatMessage({ id: 'common.priority.medium' })}</Tag>;
	};

	const columns = [
		{ title: intl.formatMessage({ id: 'publish.column.taskId' }), dataIndex: "id", key: "id", width: 80 },
		{
			title: intl.formatMessage({ id: 'publish.column.contentId' }),
			dataIndex: "content_id",
			key: "content_id",
			width: 80,
			render: (id: number) => (
				<a onClick={() => navigate(`/content/${id}`)} className="text-blue-500">{id}</a>
			),
		},
		{ title: intl.formatMessage({ id: 'publish.column.channelId' }), dataIndex: "channel_id", key: "channel_id", width: 80 },
		{ title: intl.formatMessage({ id: 'publish.column.status' }), dataIndex: "status", key: "status", width: 100, render: (status: number) => getStatusTag(status) },
		{ title: intl.formatMessage({ id: 'publish.column.priority' }), dataIndex: "priority", key: "priority", width: 80, render: (priority: number) => getPriorityTag(priority) },
		{ title: intl.formatMessage({ id: 'publish.column.scheduledTime' }), dataIndex: "scheduled_time", key: "scheduled_time", width: 180, render: (text: string) => (text ? new Date(text).toLocaleString() : "-") },
		{ title: intl.formatMessage({ id: 'publish.column.publishedTime' }), dataIndex: "published_time", key: "published_time", width: 180, render: (text: string) => (text ? new Date(text).toLocaleString() : "-") },
		{ title: intl.formatMessage({ id: 'publish.column.retryCount' }), dataIndex: "retry_count", key: "retry_count", width: 80, render: (count: number) => <Badge count={count} showZero style={{ backgroundColor: count > 0 ? "#faad14" : "#d9d9d9" }} /> },
		{
			title: intl.formatMessage({ id: 'publish.column.errorMessage' }),
			dataIndex: "error_message",
			key: "error_message",
			ellipsis: true,
			render: (text: string) => text ? <Tooltip title={text}><span className="text-red-500">{text}</span></Tooltip> : "-",
		},
		{ title: intl.formatMessage({ id: 'publish.column.createdAt' }), dataIndex: "created_at", key: "created_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.view' })}>
						<Button type="text" icon={<EyeOutlined />} onClick={() => handlePreview(record)} />
					</Tooltip>
					{record.status === 0 && (
						<Popconfirm
							title={intl.formatMessage({ id: 'publish.confirmCancel' })}
							onConfirm={() => handleCancel(record.id)}
							okText={intl.formatMessage({ id: 'common.action.confirm' })}
							cancelText={intl.formatMessage({ id: 'common.action.cancel' })}
						>
							<Tooltip title={intl.formatMessage({ id: 'common.action.cancel' })}>
								<Button type="text" danger icon={<StopOutlined />} />
							</Tooltip>
						</Popconfirm>
					)}
					{record.status === 3 && (
						<Popconfirm
							title={intl.formatMessage({ id: 'publish.confirmRetry' })}
							onConfirm={() => handleRetry(record.id)}
							okText={intl.formatMessage({ id: 'common.action.confirm' })}
							cancelText={intl.formatMessage({ id: 'common.action.cancel' })}
						>
							<Tooltip title={intl.formatMessage({ id: 'publish.action.retry' })}>
								<Button type="text" icon={<ReloadOutlined />} loading={retryMutation.isPending} />
							</Tooltip>
						</Popconfirm>
					)}
				</Space>
			),
		},
	];

	const handleCreate = async (values: any) => {
		try {
			const submitData = {
				...values,
				scheduled_time: values.scheduled_time ? values.scheduled_time.toISOString() : null,
			};
			await createMutation.mutateAsync(submitData);
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.createFailed' }));
		}
	};

	const handleCancel = async (id: number) => {
		try {
			await cancelMutation.mutateAsync(id);
			message.success(intl.formatMessage({ id: 'publish.message.cancelSuccess' }));
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'publish.message.cancelFailed' }));
		}
	};

	const handleRetry = async (id: number) => {
		try {
			await retryMutation.mutateAsync(id);
			message.success(intl.formatMessage({ id: 'publish.message.retrySubmitted' }));
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'publish.message.retryFailed' }));
		}
	};

	const handlePreview = async (record: any) => {
		try {
			const res = await api.publishTasks.get(record.id);
			setPreviewData({
				...res.data.data,
				content_title: intl.formatMessage({ id: 'publish.preview.contentLabel' }, { id: record.content_id }),
				channel_name: intl.formatMessage({ id: 'publish.preview.channelLabel' }, { id: record.channel_id }),
			});
			setPreviewModalVisible(true);
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'publish.message.previewFailed' }));
		}
	};

	const handleSearch = (values: any) => {
		setQueryParams({ ...queryParams, ...values, page: 1 });
	};

	const handleReset = () => {
		searchForm.resetFields();
		setQueryParams({ page: 1, page_size: 10, status: undefined });
	};

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
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'publish.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'publish.page.subtitle' })}</p>
			</div>

			<div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-4">
				<Card size="small"><div className="text-center"><div className="text-2xl font-bold text-gray-800">{stats.total}</div><div className="text-gray-500">{intl.formatMessage({ id: 'publish.stat.total' })}</div></div></Card>
				<Card size="small"><div className="text-center"><div className="text-2xl font-bold text-blue-500">{stats.pending}</div><div className="text-gray-500">{intl.formatMessage({ id: 'publish.stat.pending' })}</div></div></Card>
				<Card size="small"><div className="text-center"><div className="text-2xl font-bold text-orange-500">{stats.publishing}</div><div className="text-gray-500">{intl.formatMessage({ id: 'publish.stat.publishing' })}</div></div></Card>
				<Card size="small"><div className="text-center"><div className="text-2xl font-bold text-green-500">{stats.success}</div><div className="text-gray-500">{intl.formatMessage({ id: 'publish.stat.published' })}</div></div></Card>
				<Card size="small"><div className="text-center"><div className="text-2xl font-bold text-red-500">{stats.failed}</div><div className="text-gray-500">{intl.formatMessage({ id: 'publish.stat.failed' })}</div></div></Card>
			</div>

			<Card className="mb-4">
				<Form form={searchForm} layout="inline" onFinish={handleSearch}>
					<Form.Item name="status" label={intl.formatMessage({ id: 'common.form.status' })}>
						<Select placeholder={intl.formatMessage({ id: 'publish.placeholder.status' })} allowClear style={{ width: 120 }}>
							{taskStatuses.map((s) => (
								<Option key={s.value} value={s.value}>{s.label}</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" icon={<SearchOutlined />} htmlType="submit">{intl.formatMessage({ id: 'common.action.search' })}</Button>
							<Button onClick={handleReset}>{intl.formatMessage({ id: 'common.action.reset' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Card>

			<Card
				title={<Space><SendOutlined /><span>{intl.formatMessage({ id: 'publish.section.list' })}</span></Space>}
				extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'publish.action.create' })}</Button>}
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
						showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }),
						onChange: (page, pageSize) => setQueryParams({ ...queryParams, page, page_size: pageSize }),
					}}
				/>
			</Card>

			<Modal title={intl.formatMessage({ id: 'publish.modal.create' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={600}>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item name="content_id" label={intl.formatMessage({ id: 'publish.form.selectContent' })} rules={[{ required: true, message: intl.formatMessage({ id: 'publish.validation.selectContent' }) }]}>
						<Select placeholder={intl.formatMessage({ id: 'publish.placeholder.content' })} showSearch optionFilterProp="label">
							{contents.map((content: any) => (
								<Select.Option key={content.id} value={content.id} label={content.title}>{content.title}</Select.Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="channel_id" label={intl.formatMessage({ id: 'publish.form.selectChannel' })} rules={[{ required: true, message: intl.formatMessage({ id: 'publish.validation.selectChannel' }) }]}>
						<Select placeholder={intl.formatMessage({ id: 'publish.placeholder.channel' })}>
							{channels.map((channel: any) => (
								<Select.Option key={channel.id} value={channel.id}>{channel.channel_name} ({channel.channel_type})</Select.Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="scheduled_time" label={intl.formatMessage({ id: 'publish.form.scheduledTime' })}>
						<DatePicker showTime format="YYYY-MM-DD HH:mm:ss" placeholder={intl.formatMessage({ id: 'publish.placeholder.immediate' })} style={{ width: "100%" }} />
					</Form.Item>
					<Form.Item name="priority" label={intl.formatMessage({ id: 'publish.form.priority' })} initialValue={1}>
						<Select placeholder={intl.formatMessage({ id: 'publish.placeholder.priority' })}>
							{taskPriorities.map((p) => (
								<Select.Option key={p.value} value={p.value}><Tag color={p.color}>{p.label}</Tag></Select.Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit" loading={createMutation.isPending}>{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title={intl.formatMessage({ id: 'publish.modal.preview' })} open={previewModalVisible} onCancel={() => setPreviewModalVisible(false)} footer={null} width={700}>
				{previewData && (
					<div className="space-y-4">
						<div className="p-4 bg-blue-50 rounded-lg">
							<h3 className="font-medium mb-2">{intl.formatMessage({ id: 'publish.preview.taskInfo' })}</h3>
							<div className="grid grid-cols-2 gap-2 text-sm">
								<div>{intl.formatMessage({ id: 'publish.preview.taskId' })} {previewData.id}</div>
								<div>{intl.formatMessage({ id: 'publish.preview.status' })} {getStatusTag(previewData.status)}</div>
								<div>{intl.formatMessage({ id: 'publish.preview.content' })} {previewData.content_title}</div>
								<div>{intl.formatMessage({ id: 'publish.preview.channel' })} {previewData.channel_name}</div>
								<div>{intl.formatMessage({ id: 'publish.preview.scheduledTime' })} {previewData.scheduled_time || intl.formatMessage({ id: 'publish.preview.immediate' })}</div>
								<div>{intl.formatMessage({ id: 'publish.preview.retryCount' })} {previewData.retry_count || 0}</div>
							</div>
						</div>
						<div className="p-4 bg-green-50 rounded-lg">
							<h3 className="font-medium mb-2">{intl.formatMessage({ id: 'publish.preview.validationResult' })}</h3>
							<div className="space-y-2">
								{(previewData?.validation || [
									{ name: intl.formatMessage({ id: 'publish.validation.contentIntegrity' }), passed: previewData?.status === 2 },
									{ name: intl.formatMessage({ id: 'publish.validation.channelConfig' }), passed: previewData?.status === 2 },
									{ name: intl.formatMessage({ id: 'publish.validation.schemaMarkup' }), passed: previewData?.status === 2 },
									{ name: intl.formatMessage({ id: 'publish.validation.imageAlt' }), passed: previewData?.status === 2 },
								]).map((item: { name: string; passed: boolean }, index: number) => (
									<div key={index} className="flex items-center">
										{item.passed ? <CheckCircleOutlined className="text-green-500 mr-2" /> : <CloseCircleOutlined className="text-red-500 mr-2" />}
										<span>{item.name}{item.passed ? intl.formatMessage({ id: 'common.status.passed' }) : intl.formatMessage({ id: 'common.status.failed' })}</span>
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
