"use client";

import {
	useCreateSchedule,
	useDeleteSchedule,
	useDisableSchedule,
	useEnableSchedule,
	useSchedules,
	useUpdateSchedule,
} from "@/hooks";
import {
	CalendarOutlined,
	ClockCircleOutlined,
	DeleteOutlined,
	EditOutlined,
	HeatMapOutlined,
	PauseCircleOutlined,
	PlayCircleOutlined,
	PlusOutlined,
	ScheduleOutlined,
} from "@ant-design/icons";
import {
	Alert,
	Badge,
	Button,
	Calendar,
	Card,
	Col,
	Form,
	Input,
	Modal,
	Popconfirm,
	Row,
	Select,
	Space,
	Switch,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import type { Dayjs } from "dayjs";
import { useState } from "react";
import { useIntl } from "react-intl";

const { Option } = Select;
const { TextArea } = Input;

export default function SchedulePage() {
	const intl = useIntl();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [editingSchedule, setEditingSchedule] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();
	const [activeTab, setActiveTab] = useState("list");

	const { data, isLoading } = useSchedules();
	const createMutation = useCreateSchedule();
	const updateMutation = useUpdateSchedule();
	const deleteMutation = useDeleteSchedule();
	const enableMutation = useEnableSchedule();
	const disableMutation = useDisableSchedule();

	const schedules = data?.items || [];

	const scheduleTypes = [
		{ value: "fixed", label: intl.formatMessage({ id: 'schedule.type.fixed' }), color: "blue", description: intl.formatMessage({ id: 'schedule.type.fixedDesc' }) },
		{ value: "interval", label: intl.formatMessage({ id: 'schedule.type.interval' }), color: "green", description: intl.formatMessage({ id: 'schedule.type.intervalDesc' }) },
		{ value: "event", label: intl.formatMessage({ id: 'schedule.type.event' }), color: "purple", description: intl.formatMessage({ id: 'schedule.type.eventDesc' }) },
		{ value: "heat", label: intl.formatMessage({ id: 'schedule.type.heat' }), color: "orange", description: intl.formatMessage({ id: 'schedule.type.heatDesc' }) },
	];

	const getScheduleTypeTag = (type: string) => {
		const typeInfo = scheduleTypes.find((t) => t.value === type);
		return <Tag color={typeInfo?.color || "default"}>{typeInfo?.label || type}</Tag>;
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{ title: intl.formatMessage({ id: 'schedule.name' }), dataIndex: "schedule_name", key: "schedule_name", render: (text: string) => <Space><ScheduleOutlined /><span className="font-medium">{text}</span></Space> },
		{ title: intl.formatMessage({ id: 'schedule.type' }), dataIndex: "schedule_type", key: "schedule_type", width: 120, render: (type: string) => getScheduleTypeTag(type) },
		{ title: intl.formatMessage({ id: 'schedule.cron' }), dataIndex: "cron_expression", key: "cron_expression", width: 150, render: (text: string) => text ? <code className="bg-gray-100 px-2 py-1 rounded text-sm">{text}</code> : "-" },
		{ title: intl.formatMessage({ id: 'schedule.enabled' }), dataIndex: "is_enabled", key: "is_enabled", width: 100, render: (enabled: boolean) => <Badge status={enabled ? "success" : "default"} text={enabled ? intl.formatMessage({ id: 'common.status.enabled' }) : intl.formatMessage({ id: 'common.status.disabled' })} /> },
		{ title: intl.formatMessage({ id: 'schedule.nextRun' }), dataIndex: "next_run_time", key: "next_run_time", width: 180, render: (text: string) => (text ? new Date(text).toLocaleString() : "-") },
		{ title: intl.formatMessage({ id: 'schedule.column.runCount' }), dataIndex: "run_count", key: "run_count", width: 80, render: (count: number) => <Badge count={count} showZero style={{ backgroundColor: "#1890ff" }} /> },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => handleEdit(record)} /></Tooltip>
					<Tooltip title={record.is_enabled ? intl.formatMessage({ id: 'common.action.disable' }) : intl.formatMessage({ id: 'common.action.enable' })}>
						<Switch checked={record.is_enabled} size="small" checkedChildren={<PlayCircleOutlined />} unCheckedChildren={<PauseCircleOutlined />} onChange={(checked) => handleToggle(record.id, checked)} />
					</Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'common.confirmDelete' })} onConfirm={() => handleDelete(record.id)} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	const handleCreate = async (values: any) => {
		try {
			await createMutation.mutateAsync(values);
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.createFailed' }));
		}
	};

	const handleEdit = (record: any) => {
		setEditingSchedule(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await updateMutation.mutateAsync({ id: editingSchedule.id, data: values });
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingSchedule(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	const handleToggle = async (id: number, enabled: boolean) => {
		try {
			if (enabled) { await enableMutation.mutateAsync(id); } else { await disableMutation.mutateAsync(id); }
			message.success(enabled ? intl.formatMessage({ id: 'common.status.enabled' }) : intl.formatMessage({ id: 'common.status.disabled' }));
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await deleteMutation.mutateAsync(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const onCalendarSelect = (date: Dayjs) => {
		message.info(intl.formatMessage({ id: 'common.message.updateSuccess' }));
	};

	const stats = {
		total: schedules.length,
		enabled: schedules.filter((s: any) => s.is_enabled).length,
		fixed: schedules.filter((s: any) => s.schedule_type === "fixed").length,
		interval: schedules.filter((s: any) => s.schedule_type === "interval").length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.schedule' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'schedule.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}><Card><div className="text-center"><div className="text-2xl font-bold text-gray-800">{stats.total}</div><div className="text-gray-500">{intl.formatMessage({ id: 'schedule.stat.total' })}</div></div></Card></Col>
				<Col xs={12} sm={6}><Card><div className="text-center"><div className="text-2xl font-bold text-green-500">{stats.enabled}</div><div className="text-gray-500">{intl.formatMessage({ id: 'common.status.enabled' })}</div></div></Card></Col>
				<Col xs={12} sm={6}><Card><div className="text-center"><div className="text-2xl font-bold text-blue-500">{stats.fixed}</div><div className="text-gray-500">{intl.formatMessage({ id: 'schedule.stat.fixed' })}</div></div></Card></Col>
				<Col xs={12} sm={6}><Card><div className="text-center"><div className="text-2xl font-bold text-purple-500">{stats.interval}</div><div className="text-gray-500">{intl.formatMessage({ id: 'schedule.stat.interval' })}</div></div></Card></Col>
			</Row>

			<Alert
				message={intl.formatMessage({ id: 'schedule.alert.typeDesc' })}
				description={
					<div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-2">
						{scheduleTypes.map((type) => (
							<div key={type.value} className="flex items-start space-x-2">
								<Tag color={type.color}>{type.label}</Tag>
								<span className="text-gray-500 text-sm">{type.description}</span>
							</div>
						))}
					</div>
				}
				type="info"
				showIcon
				className="mb-4"
			/>

			<Row gutter={[16, 16]}>
				<Col xs={24} lg={16}>
					<Card
						title={<Space><ScheduleOutlined /><span>{intl.formatMessage({ id: 'schedule.card.scheduleList' })}</span></Space>}
						extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'schedule.action.createSchedule' })}</Button>}
					>
						<Table columns={columns} dataSource={schedules} rowKey="id" loading={isLoading} pagination={{ showSizeChanger: true, showQuickJumper: true, showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }) }} />
					</Card>
				</Col>

				<Col xs={24} lg={8}>
					<Card title={<Space><CalendarOutlined /> {intl.formatMessage({ id: 'schedule.card.publishCalendar' })}</Space>}>
						<Calendar fullscreen={false} onSelect={onCalendarSelect} />
					</Card>
				</Col>
			</Row>

		<Modal title={intl.formatMessage({ id: 'schedule.modal.createTitle' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={600}>
			<Form form={createForm} layout="vertical" onFinish={handleCreate}>
				<Form.Item name="schedule_name" label={intl.formatMessage({ id: 'schedule.name' })} rules={[{ required: true, message: intl.formatMessage({ id: 'schedule.validation.enterName' }) }]}><Input placeholder={intl.formatMessage({ id: 'schedule.placeholder.name' })} /></Form.Item>
				<Form.Item name="schedule_type" label={intl.formatMessage({ id: 'schedule.type' })} rules={[{ required: true, message: intl.formatMessage({ id: 'schedule.validation.selectType' }) }]}>
					<Select placeholder={intl.formatMessage({ id: 'schedule.placeholder.type' })}>
							{scheduleTypes.map((type) => <Option key={type.value} value={type.value}><Space><Tag color={type.color}>{type.label}</Tag><span>{type.description}</span></Space></Option>)}
						</Select>
					</Form.Item>
				<Form.Item name="cron_expression" label={intl.formatMessage({ id: 'schedule.cron' })} tooltip={intl.formatMessage({ id: 'schedule.placeholder.cron' })}><Input placeholder={intl.formatMessage({ id: 'schedule.placeholder.cron' })} /></Form.Item>
				<Form.Item name="config" label={intl.formatMessage({ id: 'schedule.form.config' })}><TextArea rows={3} placeholder={intl.formatMessage({ id: 'schedule.placeholder.config' })} /></Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit" loading={createMutation.isPending}>{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

		<Modal title={intl.formatMessage({ id: 'schedule.modal.editTitle' })} open={editModalVisible} onCancel={() => { setEditModalVisible(false); setEditingSchedule(null); }} footer={null} width={600}>
			<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
				<Form.Item name="schedule_name" label={intl.formatMessage({ id: 'schedule.name' })} rules={[{ required: true, message: intl.formatMessage({ id: 'schedule.validation.enterName' }) }]}><Input placeholder={intl.formatMessage({ id: 'schedule.placeholder.name' })} /></Form.Item>
				<Form.Item name="schedule_type" label={intl.formatMessage({ id: 'schedule.type' })} rules={[{ required: true, message: intl.formatMessage({ id: 'schedule.validation.selectType' }) }]}>
					<Select placeholder={intl.formatMessage({ id: 'schedule.placeholder.type' })}>
							{scheduleTypes.map((type) => <Option key={type.value} value={type.value}><Space><Tag color={type.color}>{type.label}</Tag><span>{type.description}</span></Space></Option>)}
						</Select>
					</Form.Item>
				<Form.Item name="cron_expression" label={intl.formatMessage({ id: 'schedule.cron' })}><Input placeholder={intl.formatMessage({ id: 'schedule.placeholder.cron' })} /></Form.Item>
				<Form.Item name="config" label={intl.formatMessage({ id: 'schedule.form.config' })}><TextArea rows={3} placeholder={intl.formatMessage({ id: 'schedule.placeholder.config' })} /></Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit" loading={updateMutation.isPending}>{intl.formatMessage({ id: 'common.action.save' })}</Button>
							<Button onClick={() => setEditModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
