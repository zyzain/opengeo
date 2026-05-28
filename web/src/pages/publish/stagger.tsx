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
import { useIntl } from "react-intl";

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
	const intl = useIntl();
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
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
		},
		onError: () => message.error(intl.formatMessage({ id: 'common.message.updateFailed' })),
	});

	const createMutation = useMutation({
		mutationFn: (data: any) => api.stagger.createStrategy(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["stagger-strategies"] });
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setModalVisible(false);
		},
		onError: () => message.error(intl.formatMessage({ id: 'common.message.createFailed' })),
	});

	const updateMutation = useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) => api.stagger.updateStrategy(id, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["stagger-strategies"] });
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setModalVisible(false);
		},
		onError: () => message.error(intl.formatMessage({ id: 'common.message.updateFailed' })),
	});

	const configMutation = useMutation({
		mutationFn: (data: any) => api.stagger.updateConfig(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["stagger-config"] });
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
		},
		onError: () => message.error(intl.formatMessage({ id: 'common.message.updateFailed' })),
	});

	const handleToggleStatus = (id: number) => { toggleMutation.mutate(id); };

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
		{ title: intl.formatMessage({ id: 'stagger.strategy.name' }), dataIndex: "name", key: "name", render: (t: string) => <span className="font-medium">{t}</span> },
		{ title: intl.formatMessage({ id: 'stagger.strategy.accounts' }), dataIndex: "accounts", key: "accounts", width: 100, render: (n: number) => <Badge count={n} showZero style={{ backgroundColor: "#1890ff" }} /> },
		{ title: intl.formatMessage({ id: 'stagger.strategy.interval' }), dataIndex: "interval", key: "interval", width: 130, render: (n: number) => <Tag color="blue">{intl.formatMessage({ id: 'stagger.strategy.intervalUnit' }, { n })}</Tag> },
		{ title: intl.formatMessage({ id: 'stagger.strategy.randomRange' }), dataIndex: "random_range", key: "random_range", width: 120, render: (n: number) => <Tag color="green">±{n}%</Tag> },
		{
			title: intl.formatMessage({ id: 'common.column.status' }),
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (s: string) => <Badge status={s === "active" ? "success" : "default"} text={s === "active" ? intl.formatMessage({ id: 'common.action.enable' }) : intl.formatMessage({ id: 'common.action.disable' })} />,
		},
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 150,
			render: (_: any, record: Strategy) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}>
						<Button type="text" icon={<SettingOutlined />} onClick={() => handleEdit(record)} />
					</Tooltip>
					<Tooltip title={record.status === "active" ? intl.formatMessage({ id: 'common.action.disable' }) : intl.formatMessage({ id: 'common.action.enable' })}>
						<Switch size="small" checked={record.status === "active"} onChange={() => handleToggleStatus(record.id)} />
					</Tooltip>
				</Space>
			),
		},
	];

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.publish.stagger' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'stagger.page.subtitle' })}</p>
			</div>

			<Alert message={intl.formatMessage({ id: 'stagger.alert.title' })} description={intl.formatMessage({ id: 'stagger.alert.description' })} type="info" showIcon className="mb-4" />

			<Card className="mb-4">
				<div className="flex items-center justify-between">
					<div>
						<h3 className="font-medium text-lg">{intl.formatMessage({ id: 'stagger.feature.title' })}</h3>
						<p className="text-gray-500">{intl.formatMessage({ id: 'stagger.feature.description' })}</p>
					</div>
					<Switch checked={config.enabled !== false} onChange={(checked) => configMutation.mutate({ ...config, enabled: checked })} checkedChildren={intl.formatMessage({ id: 'common.action.enable' })} unCheckedChildren={intl.formatMessage({ id: 'common.action.disable' })} />
				</div>
			</Card>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={12}>
				<Card title={intl.formatMessage({ id: 'stagger.config.global' })} loading={configLoading}>
					<Form form={form} layout="vertical" initialValues={config} onFinish={(values) => configMutation.mutate(values)}>
						<Form.Item name="min_interval" label={intl.formatMessage({ id: 'stagger.config.minInterval' })} tooltip={intl.formatMessage({ id: 'stagger.config.minIntervalTip' })}>
							<Slider min={1} max={30} marks={{ 1: "1", 5: "5", 10: "10", 15: "15", 30: "30" }} />
						</Form.Item>
						<Form.Item name="max_interval" label={intl.formatMessage({ id: 'stagger.config.maxInterval' })} tooltip={intl.formatMessage({ id: 'stagger.config.maxIntervalTip' })}>
							<Slider min={5} max={60} marks={{ 5: "5", 15: "15", 30: "30", 60: "60" }} />
						</Form.Item>
						<Form.Item name="random_range" label={intl.formatMessage({ id: 'stagger.config.randomRange' })} tooltip={intl.formatMessage({ id: 'stagger.config.randomRangeTip' })}>
							<Slider min={0} max={50} marks={{ 0: "0%", 10: "10%", 20: "20%", 30: "30%", 50: "50%" }} />
						</Form.Item>
						<Form.Item name="batch_size" label={intl.formatMessage({ id: 'stagger.config.batchSize' })} tooltip={intl.formatMessage({ id: 'stagger.config.batchSizeTip' })}>
							<InputNumber min={1} max={50} style={{ width: "100%" }} />
						</Form.Item>
						<Divider />
						<Form.Item name="cooldown_after" label={intl.formatMessage({ id: 'stagger.config.cooldownAfter' })} tooltip={intl.formatMessage({ id: 'stagger.config.cooldownAfterTip' })}>
							<InputNumber min={10} max={200} style={{ width: "100%" }} />
						</Form.Item>
						<Form.Item name="cooldown_duration" label={intl.formatMessage({ id: 'stagger.config.cooldownDuration' })} tooltip={intl.formatMessage({ id: 'stagger.config.cooldownDurationTip' })}>
							<InputNumber min={5} max={120} style={{ width: "100%" }} />
						</Form.Item>
						<Form.Item>
							<Button type="primary" htmlType="submit" icon={<SaveOutlined />}>{intl.formatMessage({ id: 'stagger.config.save' })}</Button>
						</Form.Item>
					</Form>
				</Card>
				</Col>

				<Col xs={24} lg={12}>
					<Card
						title={intl.formatMessage({ id: 'stagger.strategy.list' })}
						extra={
							<Space>
								<Button icon={<ReloadOutlined />} onClick={() => queryClient.invalidateQueries({ queryKey: ["stagger-strategies"] })}>{intl.formatMessage({ id: 'common.action.refresh' })}</Button>
								<Button type="primary" onClick={handleCreate}>{intl.formatMessage({ id: 'stagger.strategy.create' })}</Button>
							</Space>
						}
					>
						<Table columns={columns} dataSource={strategies} rowKey="id" pagination={false} size="small" loading={strategiesLoading} />
					</Card>

				<Card title={intl.formatMessage({ id: 'stagger.example.title' })} className="mt-4">
					<div className="space-y-3">
						<div className="p-3 bg-blue-50 rounded-lg">
							<div className="text-sm text-gray-500">{intl.formatMessage({ id: 'stagger.example.scenario' })}</div>
							<div className="font-medium mt-1">{intl.formatMessage({ id: 'stagger.example.params' })}</div>
							<div className="text-sm text-gray-600 mt-2">{intl.formatMessage({ id: 'stagger.example.actualInterval' })}</div>
							<div className="text-sm text-gray-600">{intl.formatMessage({ id: 'stagger.example.totalTime' })}</div>
						</div>
						<div className="p-3 bg-green-50 rounded-lg">
							<div className="text-sm text-gray-500">{intl.formatMessage({ id: 'stagger.example.safety' })}</div>
							<div className="flex items-center mt-1">
								<ThunderboltOutlined className="text-green-500 mr-2" />
								<span className="text-green-600 font-medium">{intl.formatMessage({ id: 'stagger.example.lowRisk' })}</span>
							</div>
						</div>
					</div>
				</Card>
				</Col>
			</Row>

			<Modal
				title={editingStrategy ? intl.formatMessage({ id: 'stagger.modal.editStrategy' }) : intl.formatMessage({ id: 'stagger.modal.createStrategy' })}
				open={modalVisible}
				onOk={handleModalOk}
				onCancel={() => setModalVisible(false)}
				okText={intl.formatMessage({ id: 'common.action.save' })}
				cancelText={intl.formatMessage({ id: 'common.action.cancel' })}
				confirmLoading={createMutation.isPending || updateMutation.isPending}
			>
				<Form form={strategyForm} layout="vertical">
				<Form.Item name="name" label={intl.formatMessage({ id: 'stagger.form.strategyName' })} rules={[{ required: true, message: intl.formatMessage({ id: 'stagger.validation.enterName' }) }]}>
					<Input placeholder={intl.formatMessage({ id: 'stagger.placeholder.name' })} />
				</Form.Item>
				<Form.Item name="accounts" label={intl.formatMessage({ id: 'stagger.form.accounts' })} rules={[{ required: true, message: intl.formatMessage({ id: 'stagger.validation.enterAccounts' }) }]}>
					<InputNumber min={1} max={100} style={{ width: "100%" }} placeholder={intl.formatMessage({ id: 'stagger.placeholder.accounts' })} />
				</Form.Item>
				<Form.Item name="interval" label={intl.formatMessage({ id: 'stagger.form.interval' })} rules={[{ required: true, message: intl.formatMessage({ id: 'stagger.validation.enterInterval' }) }]}>
					<InputNumber min={1} max={60} style={{ width: "100%" }} placeholder={intl.formatMessage({ id: 'stagger.placeholder.interval' })} />
				</Form.Item>
				<Form.Item name="random_range" label={intl.formatMessage({ id: 'stagger.form.randomRange' })} rules={[{ required: true, message: intl.formatMessage({ id: 'stagger.validation.enterRandomRange' }) }]}>
					<InputNumber min={0} max={100} style={{ width: "100%" }} placeholder={intl.formatMessage({ id: 'stagger.placeholder.randomRange' })} />
				</Form.Item>
			</Form>
			</Modal>
		</div>
	);
}
