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

const { Option } = Select;
const { TextArea } = Input;

export default function SchedulePage() {
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

	// 调度类型
	const scheduleTypes = [
		{
			value: "fixed",
			label: "固定时间",
			color: "blue",
			description: "在指定时间执行一次",
		},
		{
			value: "interval",
			label: "间隔循环",
			color: "green",
			description: "按固定间隔重复执行",
		},
		{
			value: "event",
			label: "事件触发",
			color: "purple",
			description: "当特定事件发生时执行",
		},
		{
			value: "heat",
			label: "热力图推荐",
			color: "orange",
			description: "根据AI活跃度推荐时间",
		},
	];

	// 获取调度类型标签
	const getScheduleTypeTag = (type: string) => {
		const typeInfo = scheduleTypes.find((t) => t.value === type);
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
			title: "调度名称",
			dataIndex: "schedule_name",
			key: "schedule_name",
			render: (text: string) => (
				<Space>
					<ScheduleOutlined />
					<span className="font-medium">{text}</span>
				</Space>
			),
		},
		{
			title: "调度类型",
			dataIndex: "schedule_type",
			key: "schedule_type",
			width: 120,
			render: (type: string) => getScheduleTypeTag(type),
		},
		{
			title: "Cron表达式",
			dataIndex: "cron_expression",
			key: "cron_expression",
			width: 150,
			render: (text: string) =>
				text ? (
					<code className="bg-gray-100 px-2 py-1 rounded text-sm">{text}</code>
				) : (
					"-"
				),
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
			title: "下次运行",
			dataIndex: "next_run_time",
			key: "next_run_time",
			width: 180,
			render: (text: string) => (text ? new Date(text).toLocaleString() : "-"),
		},
		{
			title: "运行次数",
			dataIndex: "run_count",
			key: "run_count",
			width: 80,
			render: (count: number) => (
				<Badge count={count} showZero style={{ backgroundColor: "#1890ff" }} />
			),
		},
		{
			title: "操作",
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="编辑">
						<Button
							type="text"
							icon={<EditOutlined />}
							onClick={() => handleEdit(record)}
						/>
					</Tooltip>
					<Tooltip title={record.is_enabled ? "禁用" : "启用"}>
						<Switch
							checked={record.is_enabled}
							size="small"
							checkedChildren={<PlayCircleOutlined />}
							unCheckedChildren={<PauseCircleOutlined />}
							onChange={(checked) => handleToggle(record.id, checked)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个调度吗？"
						onConfirm={() => handleDelete(record.id)}
						okText="确定"
						cancelText="取消"
					>
						<Tooltip title="删除">
							<Button type="text" danger icon={<DeleteOutlined />} />
						</Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	// 创建调度
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

	// 编辑调度
	const handleEdit = (record: any) => {
		setEditingSchedule(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await updateMutation.mutateAsync({
				id: editingSchedule.id,
				data: values,
			});
			message.success("更新成功");
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingSchedule(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || "更新失败");
		}
	};

	// 切换状态
	const handleToggle = async (id: number, enabled: boolean) => {
		try {
			if (enabled) {
				await enableMutation.mutateAsync(id);
			} else {
				await disableMutation.mutateAsync(id);
			}
			message.success(enabled ? "已启用" : "已禁用");
		} catch (error: any) {
			message.error("操作失败");
		}
	};

	// 删除调度
	const handleDelete = async (id: number) => {
		try {
			await deleteMutation.mutateAsync(id);
			message.success("删除成功");
		} catch (error: any) {
			message.error("删除失败");
		}
	};

	// 日历选择
	const onCalendarSelect = (date: Dayjs) => {
		message.info(`选择了日期: ${date.format("YYYY-MM-DD")}`);
	};

	// 统计数据
	const stats = {
		total: schedules.length,
		enabled: schedules.filter((s: any) => s.is_enabled).length,
		fixed: schedules.filter((s: any) => s.schedule_type === "fixed").length,
		interval: schedules.filter((s: any) => s.schedule_type === "interval")
			.length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">调度管理</h1>
				<p className="text-gray-500 mt-1">管理定时发布任务</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}>
					<Card>
						<div className="text-center">
							<div className="text-2xl font-bold text-gray-800">
								{stats.total}
							</div>
							<div className="text-gray-500">总调度数</div>
						</div>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<div className="text-center">
							<div className="text-2xl font-bold text-green-500">
								{stats.enabled}
							</div>
							<div className="text-gray-500">已启用</div>
						</div>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<div className="text-center">
							<div className="text-2xl font-bold text-blue-500">
								{stats.fixed}
							</div>
							<div className="text-gray-500">固定时间</div>
						</div>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<div className="text-center">
							<div className="text-2xl font-bold text-purple-500">
								{stats.interval}
							</div>
							<div className="text-gray-500">间隔循环</div>
						</div>
					</Card>
				</Col>
			</Row>

			{/* 调度类型说明 */}
			<Alert
				message="调度类型说明"
				description={
					<div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-2">
						{scheduleTypes.map((type) => (
							<div key={type.value} className="flex items-start space-x-2">
								<Tag color={type.color}>{type.label}</Tag>
								<span className="text-gray-500 text-sm">
									{type.description}
								</span>
							</div>
						))}
					</div>
				}
				type="info"
				showIcon
				className="mb-4"
			/>

			<Row gutter={[16, 16]}>
				{/* 调度列表 */}
				<Col xs={24} lg={16}>
					<Card
						title={
							<Space>
								<ScheduleOutlined />
								<span>调度列表</span>
							</Space>
						}
						extra={
							<Button
								type="primary"
								icon={<PlusOutlined />}
								onClick={() => setCreateModalVisible(true)}
							>
								创建调度
							</Button>
						}
					>
						<Table
							columns={columns}
							dataSource={schedules}
							rowKey="id"
							loading={isLoading}
							pagination={{
								showSizeChanger: true,
								showQuickJumper: true,
								showTotal: (total) => `共 ${total} 条`,
							}}
						/>
					</Card>
				</Col>

				{/* 日历视图 */}
				<Col xs={24} lg={8}>
					<Card
						title={
							<Space>
								<CalendarOutlined /> 发布日历
							</Space>
						}
					>
						<Calendar fullscreen={false} onSelect={onCalendarSelect} />
					</Card>
				</Col>
			</Row>

			{/* 创建调度弹窗 */}
			<Modal
				title="创建调度"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={600}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="schedule_name"
						label="调度名称"
						rules={[{ required: true, message: "请输入调度名称" }]}
					>
						<Input placeholder="请输入调度名称" />
					</Form.Item>
					<Form.Item
						name="schedule_type"
						label="调度类型"
						rules={[{ required: true, message: "请选择调度类型" }]}
					>
						<Select placeholder="请选择调度类型">
							{scheduleTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<Space>
										<Tag color={type.color}>{type.label}</Tag>
										<span>{type.description}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item
						name="cron_expression"
						label="Cron表达式"
						tooltip="例如: 0 0 9 * * ? 表示每天9点执行"
					>
						<Input placeholder="请输入Cron表达式，例如: 0 0 9 * * ?" />
					</Form.Item>
					<Form.Item name="config" label="配置">
						<TextArea rows={3} placeholder="请输入JSON格式的配置" />
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

			{/* 编辑调度弹窗 */}
			<Modal
				title="编辑调度"
				open={editModalVisible}
				onCancel={() => {
					setEditModalVisible(false);
					setEditingSchedule(null);
				}}
				footer={null}
				width={600}
			>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item
						name="schedule_name"
						label="调度名称"
						rules={[{ required: true, message: "请输入调度名称" }]}
					>
						<Input placeholder="请输入调度名称" />
					</Form.Item>
					<Form.Item
						name="schedule_type"
						label="调度类型"
						rules={[{ required: true, message: "请选择调度类型" }]}
					>
						<Select placeholder="请选择调度类型">
							{scheduleTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<Space>
										<Tag color={type.color}>{type.label}</Tag>
										<span>{type.description}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="cron_expression" label="Cron表达式">
						<Input placeholder="请输入Cron表达式" />
					</Form.Item>
					<Form.Item name="config" label="配置">
						<TextArea rows={3} placeholder="请输入JSON格式的配置" />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								htmlType="submit"
								loading={updateMutation.isPending}
							>
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
