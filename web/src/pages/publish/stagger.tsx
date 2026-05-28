"use client";

import api from "@/lib/api";
import {
	ClockCircleOutlined,
	InfoCircleOutlined,
	ReloadOutlined,
	SaveOutlined,
	ScheduleOutlined,
	SettingOutlined,
	TeamOutlined,
	ThunderboltOutlined,
} from "@ant-design/icons";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
	Alert,
	Badge,
	Button,
	Card,
	Col,
	Divider,
	Form,
	Input,
	InputNumber,
	Modal,
	Row,
	Select,
	Slider,
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

interface Strategy {
	id: number;
	name: string;
	accounts: number;
	interval: number;
	random_range: number;
	status: "active" | "inactive";
}

export default function StaggeredPublishPage() {
	const [form] = Form.useForm();
	const [strategyForm] = Form.useForm();
	const [editingStrategy, setEditingStrategy] = useState<Strategy | null>(null);
	const [modalVisible, setModalVisible] = useState(false);
	const queryClient = useQueryClient();

	const { data: strategiesData, isLoading: strategiesLoading } = useQuery({
		queryKey: ["stagger-strategies"],
		queryFn: () => api.stagger.listStrategies(),
	});

	const { data: configData, isLoading: configLoading } = useQuery({
		queryKey: ["stagger-config"],
		queryFn: () => api.stagger.getConfig(),
	});

	const strategies = strategiesData?.data?.data?.items || [];
	const config = configData?.data?.data || {};

	const toggleMutation = useMutation({
		mutationFn: (id: number) => api.stagger.toggleStrategy(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["stagger-strategies"] });
			message.success("状态已更新");
		},
		onError: () => message.error("更新状态失败"),
	});

	const createMutation = useMutation({
		mutationFn: (data: any) => api.stagger.createStrategy(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["stagger-strategies"] });
			message.success("策略已创建");
			setModalVisible(false);
		},
		onError: () => message.error("创建策略失败"),
	});

	const updateMutation = useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.stagger.updateStrategy(id, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["stagger-strategies"] });
			message.success("策略已更新");
			setModalVisible(false);
		},
		onError: () => message.error("更新策略失败"),
	});

	const configMutation = useMutation({
		mutationFn: (data: any) => api.stagger.updateConfig(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["stagger-config"] });
			message.success("配置已保存");
		},
		onError: () => message.error("保存配置失败"),
	});

	const handleToggleStatus = (id: number) => {
		toggleMutation.mutate(id);
	};

	const handleEdit = (record: Strategy) => {
		setEditingStrategy(record);
		strategyForm.setFieldsValue(record);
		setModalVisible(true);
	};

	const handleCreate = () => {
		setEditingStrategy(null);
		strategyForm.resetFields();
		setModalVisible(true);
	};

	const handleModalOk = () => {
		strategyForm.validateFields().then((values) => {
			if (editingStrategy) {
				updateMutation.mutate({ id: editingStrategy.id, data: values });
			} else {
				createMutation.mutate(values);
			}
		});
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 60 },
		{
			title: "策略名称",
			dataIndex: "name",
			key: "name",
			render: (t: string) => <span className="font-medium">{t}</span>,
		},
		{
			title: "适用账号数",
			dataIndex: "accounts",
			key: "accounts",
			width: 100,
			render: (n: number) => (
				<Badge count={n} showZero style={{ backgroundColor: "#1890ff" }} />
			),
		},
		{
			title: "基础间隔(分钟)",
			dataIndex: "interval",
			key: "interval",
			width: 130,
			render: (n: number) => <Tag color="blue">{n}分钟</Tag>,
		},
		{
			title: "随机浮动(%)",
			dataIndex: "random_range",
			key: "random_range",
			width: 120,
			render: (n: number) => <Tag color="green">±{n}%</Tag>,
		},
		{
			title: "状态",
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (s: string) => (
				<Badge
					status={s === "active" ? "success" : "default"}
					text={s === "active" ? "启用" : "禁用"}
				/>
			),
		},
		{
			title: "操作",
			key: "action",
			width: 150,
			render: (_: any, record: Strategy) => (
				<Space size="small">
					<Tooltip title="编辑">
						<Button
							type="text"
							icon={<SettingOutlined />}
							onClick={() => handleEdit(record)}
						/>
					</Tooltip>
					<Tooltip title={record.status === "active" ? "禁用" : "启用"}>
						<Switch
							size="small"
							checked={record.status === "active"}
							onChange={() => handleToggleStatus(record.id)}
						/>
					</Tooltip>
				</Space>
			),
		},
	];

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">错峰发布</h1>
				<p className="text-gray-500 mt-1">
					配置矩阵账号的发布间隔和随机延迟，模拟真人操作节奏
				</p>
			</div>

			{/* 功能说明 */}
			<Alert
				message="错峰发布策略"
				description="通过设置账号间的发布间隔和随机延迟，避免批量发布触发平台风控。每个账号将在设定的时间窗口内随机延迟发布，模拟真人操作行为。"
				type="info"
				showIcon
				className="mb-4"
			/>

			{/* 全局开关 */}
			<Card className="mb-4">
				<div className="flex items-center justify-between">
					<div>
						<h3 className="font-medium text-lg">错峰发布功能</h3>
						<p className="text-gray-500">
							启用后，所有批量发布任务将自动应用错峰策略
						</p>
					</div>
					<Switch
						checked={config.enabled !== false}
						onChange={(checked) =>
							configMutation.mutate({ ...config, enabled: checked })
						}
						checkedChildren="启用"
						unCheckedChildren="禁用"
					/>
				</div>
			</Card>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={12}>
					{/* 全局配置 */}
					<Card title="全局配置" loading={configLoading}>
						<Form
							form={form}
							layout="vertical"
							initialValues={config}
							onFinish={(values) => configMutation.mutate(values)}
						>
							<Form.Item
								name="min_interval"
								label="最小间隔(分钟)"
								tooltip="同一账号连续发布的最小间隔"
							>
								<Slider
									min={1}
									max={30}
									marks={{ 1: "1", 5: "5", 10: "10", 15: "15", 30: "30" }}
								/>
							</Form.Item>
							<Form.Item
								name="max_interval"
								label="最大间隔(分钟)"
								tooltip="同一账号连续发布的最大间隔"
							>
								<Slider
									min={5}
									max={60}
									marks={{ 5: "5", 15: "15", 30: "30", 60: "60" }}
								/>
							</Form.Item>
							<Form.Item
								name="random_range"
								label="随机浮动范围(%)"
								tooltip="在基础间隔上随机浮动的百分比"
							>
								<Slider
									min={0}
									max={50}
									marks={{
										0: "0%",
										10: "10%",
										20: "20%",
										30: "30%",
										50: "50%",
									}}
								/>
							</Form.Item>
							<Form.Item
								name="batch_size"
								label="批次大小"
								tooltip="每批同时发布的账号数量"
							>
								<InputNumber min={1} max={50} style={{ width: "100%" }} />
							</Form.Item>
							<Divider />
							<Form.Item
								name="cooldown_after"
								label="冷却触发次数"
								tooltip="连续发布多少次后触发冷却期"
							>
								<InputNumber min={10} max={200} style={{ width: "100%" }} />
							</Form.Item>
							<Form.Item
								name="cooldown_duration"
								label="冷却时长(分钟)"
								tooltip="冷却期内暂停发布"
							>
								<InputNumber min={5} max={120} style={{ width: "100%" }} />
							</Form.Item>
							<Form.Item>
								<Button
									type="primary"
									htmlType="submit"
									icon={<SaveOutlined />}
								>
									保存配置
								</Button>
							</Form.Item>
						</Form>
					</Card>
				</Col>

				<Col xs={24} lg={12}>
					{/* 策略列表 */}
					<Card
						title="策略列表"
						extra={
							<Space>
								<Button
									icon={<ReloadOutlined />}
									onClick={() =>
										queryClient.invalidateQueries({
											queryKey: ["stagger-strategies"],
										})
									}
								>
									刷新
								</Button>
								<Button type="primary" onClick={handleCreate}>
									新建策略
								</Button>
							</Space>
						}
					>
						<Table
							columns={columns}
							dataSource={strategies}
							rowKey="id"
							pagination={false}
							size="small"
							loading={strategiesLoading}
						/>
					</Card>

					{/* 示例计算 */}
					<Card title="发布时间计算示例" className="mt-4">
						<div className="space-y-3">
							<div className="p-3 bg-blue-50 rounded-lg">
								<div className="text-sm text-gray-500">
									场景：10个账号发布同一内容
								</div>
								<div className="font-medium mt-1">
									基础间隔: 5分钟，随机浮动: ±30%
								</div>
								<div className="text-sm text-gray-600 mt-2">
									实际间隔: 3.5~6.5分钟（随机）
								</div>
								<div className="text-sm text-gray-600">
									预计总耗时: 31.5~58.5分钟
								</div>
							</div>
							<div className="p-3 bg-green-50 rounded-lg">
								<div className="text-sm text-gray-500">安全性评估</div>
								<div className="flex items-center mt-1">
									<ThunderboltOutlined className="text-green-500 mr-2" />
									<span className="text-green-600 font-medium">
										低风险 - 符合真人操作节奏
									</span>
								</div>
							</div>
						</div>
					</Card>
				</Col>
			</Row>

			<Modal
				title={editingStrategy ? "编辑策略" : "新建策略"}
				open={modalVisible}
				onOk={handleModalOk}
				onCancel={() => setModalVisible(false)}
				okText="保存"
				cancelText="取消"
				confirmLoading={createMutation.isPending || updateMutation.isPending}
			>
				<Form form={strategyForm} layout="vertical">
					<Form.Item
						name="name"
						label="策略名称"
						rules={[{ required: true, message: "请输入策略名称" }]}
					>
						<Input placeholder="请输入策略名称" />
					</Form.Item>
					<Form.Item
						name="accounts"
						label="适用账号数"
						rules={[{ required: true, message: "请输入账号数" }]}
					>
						<InputNumber
							min={1}
							max={100}
							style={{ width: "100%" }}
							placeholder="请输入账号数"
						/>
					</Form.Item>
					<Form.Item
						name="interval"
						label="基础间隔(分钟)"
						rules={[{ required: true, message: "请输入间隔" }]}
					>
						<InputNumber
							min={1}
							max={60}
							style={{ width: "100%" }}
							placeholder="请输入间隔"
						/>
					</Form.Item>
					<Form.Item
						name="random_range"
						label="随机浮动(%)"
						rules={[{ required: true, message: "请输入随机浮动" }]}
					>
						<InputNumber
							min={0}
							max={100}
							style={{ width: "100%" }}
							placeholder="请输入随机浮动"
						/>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
