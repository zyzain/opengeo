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
import { useIntl } from "react-intl";

const { Option } = Select;
const { TextArea } = Input;

export default function WebhooksPage() {
	const intl = useIntl();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [historyModalVisible, setHistoryModalVisible] = useState(false);
	const [selectedWebhook, setSelectedWebhook] = useState<any>(null);
	const [webhookHistory, setWebhookHistory] = useState<any[]>([]);
	const [historyLoading, setHistoryLoading] = useState(false);
	const [createForm] = Form.useForm();

	const queryClient = useQueryClient();
	const { data, isLoading } = useWebhooks();
	const webhooks = data?.items || [];

	const eventTypes = [
		{ value: "content.created", label: intl.formatMessage({ id: "webhook.event.contentCreated" }), color: "blue" },
		{ value: "content.published", label: intl.formatMessage({ id: "webhook.event.contentPublished" }), color: "green" },
		{ value: "publish.success", label: intl.formatMessage({ id: "webhook.event.publishSuccess" }), color: "green" },
		{ value: "publish.failed", label: intl.formatMessage({ id: "webhook.event.publishFailed" }), color: "red" },
		{ value: "ai.optimized", label: intl.formatMessage({ id: "webhook.event.aiOptimized" }), color: "purple" },
		{ value: "citation.found", label: intl.formatMessage({ id: "webhook.event.citationFound" }), color: "orange" },
	];

	const getEventTag = (event: string) => {
		const eventInfo = eventTypes.find((e) => e.value === event);
		return <Tag color={eventInfo?.color || "default"}>{eventInfo?.label || event}</Tag>;
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{ title: intl.formatMessage({ id: "webhook.column.name" }), dataIndex: "webhook_name", key: "webhook_name", render: (text: string) => <Space><ApiOutlined /><span className="font-medium">{text}</span></Space> },
		{
			title: "URL",
			dataIndex: "url",
			key: "url",
			ellipsis: true,
			render: (text: string) => <Tooltip title={text}><a href={text} target="_blank" rel="noopener noreferrer" className="text-blue-500"><LinkOutlined /> {text.substring(0, 40)}...</a></Tooltip>,
		},
		{
			title: intl.formatMessage({ id: "webhook.column.events" }),
			dataIndex: "events",
			key: "events",
			render: (events: string) => {
				try {
					const eventList = JSON.parse(events);
					return <Space size={[0, 4]} wrap>{eventList.slice(0, 2).map((event: string) => getEventTag(event))}{eventList.length > 2 && <Tag>+{eventList.length - 2}</Tag>}</Space>;
				} catch { return events; }
			},
		},
		{ title: intl.formatMessage({ id: 'common.column.status' }), dataIndex: "is_active", key: "is_active", width: 100, render: (active: boolean) => <Badge status={active ? "success" : "default"} text={active ? intl.formatMessage({ id: 'common.status.enabled' }) : intl.formatMessage({ id: 'common.status.disabled' })} /> },
	{ title: intl.formatMessage({ id: "webhook.column.failCount" }), dataIndex: "fail_count", key: "fail_count", width: 80, render: (count: number) => count > 0 ? <Badge count={count} style={{ backgroundColor: "#ff4d4f" }} /> : <Badge count={0} style={{ backgroundColor: "#d9d9d9" }} /> },
	{ title: intl.formatMessage({ id: "webhook.column.lastTrigger" }), dataIndex: "last_trigger", key: "last_trigger", width: 180, render: (text: string) => (text ? new Date(text).toLocaleString() : "-") },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
				<Tooltip title={intl.formatMessage({ id: "webhook.action.test" })}><Button type="text" icon={<SendOutlined />} onClick={() => handleTest(record.id)} /></Tooltip>
				<Tooltip title={intl.formatMessage({ id: "webhook.action.history" })}><Button type="text" icon={<HistoryOutlined />} onClick={() => handleShowHistory(record)} /></Tooltip>
					<Tooltip title={record.is_active ? intl.formatMessage({ id: 'common.action.disable' }) : intl.formatMessage({ id: 'common.action.enable' })}>
						<Switch checked={record.is_active} size="small" checkedChildren={<CheckCircleOutlined />} unCheckedChildren={<CloseCircleOutlined />} onChange={(checked) => handleToggleActive(record, checked)} />
					</Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'common.confirmDelete' })} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })} onConfirm={() => handleDelete(record.id)}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	const handleCreate = async (values: any) => {
		try {
			await api.system.createWebhook({ ...values, events: JSON.stringify(values.events) });
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: ["webhooks"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.createFailed' }));
		}
	};

	const handleToggleActive = async (record: any, checked: boolean) => {
		try {
			await api.system.updateWebhook(record.id, { is_active: checked });
			message.success(checked ? intl.formatMessage({ id: 'common.status.enabled' }) : intl.formatMessage({ id: 'common.status.disabled' }));
			queryClient.invalidateQueries({ queryKey: ["webhooks"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.system.deleteWebhook(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
			queryClient.invalidateQueries({ queryKey: ["webhooks"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const handleTest = async (id: number) => {
		try {
			await api.system.testWebhook(id);
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const handleShowHistory = async (record: any) => {
		setSelectedWebhook(record);
		setHistoryModalVisible(true);
		setHistoryLoading(true);
		try {
			const res = await api.webhookHistory(record.id);
			setWebhookHistory(res.data?.data || []);
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.operationFailed' }));
			setWebhookHistory([]);
		} finally {
			setHistoryLoading(false);
		}
	};

	const stats = {
		total: webhooks.length,
		active: webhooks.filter((w: any) => w.is_active).length,
		failed: webhooks.filter((w: any) => w.fail_count > 0).length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.settings.webhooks' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: "webhook.page.subtitle" })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8}><Card><Statistic title={intl.formatMessage({ id: "webhook.stat.total" })} value={stats.total} prefix={<ApiOutlined />} /></Card></Col>
				<Col xs={12} sm={8}><Card><Statistic title={intl.formatMessage({ id: "webhook.stat.active" })} value={stats.active} prefix={<CheckCircleOutlined />} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={8}><Card><Statistic title={intl.formatMessage({ id: "webhook.stat.failed" })} value={stats.failed} prefix={<CloseCircleOutlined />} valueStyle={{ color: "#ff4d4f" }} /></Card></Col>
			</Row>

			<Card title={intl.formatMessage({ id: "webhook.section.eventTypes" })} className="mb-4">
				<div className="flex flex-wrap gap-2">
					{eventTypes.map((event) => <Tag key={event.value} color={event.color}>{event.label}</Tag>)}
				</div>
			</Card>

			<Card
			title={<Space><ApiOutlined /><span>{intl.formatMessage({ id: "webhook.section.list" })}</span></Space>}
			extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: "webhook.action.create" })}</Button>}
			>
				<Table columns={columns} dataSource={webhooks} rowKey="id" loading={isLoading} scroll={{ x: 1500 }} pagination={{ showSizeChanger: true, showQuickJumper: true, showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }) }} />
			</Card>

			<Modal title={intl.formatMessage({ id: "webhook.modal.create" })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={600}>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
				<Form.Item name="webhook_name" label={intl.formatMessage({ id: "webhook.form.name" })} rules={[{ required: true, message: intl.formatMessage({ id: "webhook.validation.enterName" }) }]}><Input placeholder={intl.formatMessage({ id: "webhook.placeholder.name" })} /></Form.Item>
				<Form.Item name="url" label={intl.formatMessage({ id: "webhook.form.url" })} rules={[{ required: true, message: intl.formatMessage({ id: "webhook.validation.enterUrl" }) }, { type: "url", message: intl.formatMessage({ id: "webhook.validation.validUrl" }) }]}><Input placeholder={intl.formatMessage({ id: "webhook.placeholder.url" })} /></Form.Item>
				<Form.Item name="secret" label={intl.formatMessage({ id: "webhook.form.secret" })} tooltip={intl.formatMessage({ id: "webhook.form.secretTooltip" })}><Input.Password placeholder={intl.formatMessage({ id: "webhook.placeholder.secret" })} /></Form.Item>
				<Form.Item name="events" label={intl.formatMessage({ id: "webhook.form.events" })} rules={[{ required: true, message: intl.formatMessage({ id: "webhook.validation.selectEvents" }) }]}>
					<Select mode="multiple" placeholder={intl.formatMessage({ id: "webhook.placeholder.events" })}>
							{eventTypes.map((event) => <Option key={event.value} value={event.value}><Tag color={event.color}>{event.label}</Tag></Option>)}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title={`${intl.formatMessage({ id: "webhook.modal.history" })} - ${selectedWebhook?.webhook_name}`} open={historyModalVisible} onCancel={() => setHistoryModalVisible(false)} footer={null} width={600}>
				<List
					loading={historyLoading}
					dataSource={webhookHistory}
					locale={{ emptyText: intl.formatMessage({ id: "webhook.empty.history" }) }}
					renderItem={(item) => (
						<List.Item>
							<List.Item.Meta
								avatar={<Avatar icon={item.success ? <CheckCircleOutlined /> : <CloseCircleOutlined />} style={{ backgroundColor: item.success ? "#52c41a" : "#ff4d4f" }} />}
								title={<Space>{getEventTag(item.event)}<Tag color={item.success ? "success" : "error"}>HTTP {item.status}</Tag></Space>}
								description={item.triggered_at}
							/>
						</List.Item>
					)}
				/>
			</Modal>
		</div>
	);
}
