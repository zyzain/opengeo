"use client";

import api from "@/lib/api";
import {
	CheckCircleOutlined,
	CloseCircleOutlined,
	DeleteOutlined,
	DesktopOutlined,
	EditOutlined,
	EnvironmentOutlined,
	GlobalOutlined,
	PlusOutlined,
	ReloadOutlined,
	SettingOutlined,
} from "@ant-design/icons";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
	Badge,
	Button,
	Card,
	Col,
	Form,
	Input,
	InputNumber,
	Modal,
	Popconfirm,
	Progress,
	Row,
	Select,
	Space,
	Statistic,
	Switch,
	Table,
	Tabs,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";
import { useIntl } from "react-intl";

const { Option } = Select;
const { TextArea } = Input;

export default function EnvironmentIsolationPage() {
	const intl = useIntl();
	const queryClient = useQueryClient();
	const [activeTab, setActiveTab] = useState("fingerprints");
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [proxyModalVisible, setProxyModalVisible] = useState(false);
	const [createForm] = Form.useForm();
	const [proxyForm] = Form.useForm();

	interface Fingerprint {
		id: number;
		name: string;
		user_agent: string;
		platform: string;
		screen: string;
		language: string;
		timezone: string;
		webgl: string;
		canvas: string;
		status: string;
		account_count: number;
	}

	interface Proxy {
		id: number;
		ip: string;
		port: number;
		protocol: string;
		location: string;
		speed: number;
		uptime: number;
		status: string;
		last_check: string;
	}

	const { data: fingerprintsData } = useQuery({
		queryKey: ["fingerprints"],
		queryFn: () => api.fingerprints.list(),
	});
	const fingerprints: Fingerprint[] = fingerprintsData?.data?.data?.items ?? [];

	const { data: proxiesData } = useQuery({
		queryKey: ["proxies"],
		queryFn: () => api.proxies.list(),
	});
	const proxyPool: Proxy[] = proxiesData?.data?.data?.items ?? [];

	const createFingerprintMutation = useMutation({
		mutationFn: (data: any) => api.fingerprints.create(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["fingerprints"] });
			createForm.resetFields();
			setCreateModalVisible(false);
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
		},
	});

	const deleteFingerprintMutation = useMutation({
		mutationFn: (id: number) => api.fingerprints.delete(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["fingerprints"] });
			message.success(intl.formatMessage({ id: 'common.message.deleted' }));
		},
	});

	const toggleFingerprintMutation = useMutation({
		mutationFn: (id: number) => api.fingerprints.toggle(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["fingerprints"] });
		},
	});

	const createProxyMutation = useMutation({
		mutationFn: (data: any) => api.proxies.create(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["proxies"] });
			proxyForm.resetFields();
			setProxyModalVisible(false);
			message.success(intl.formatMessage({ id: 'common.message.addSuccess' }));
		},
	});

	const deleteProxyMutation = useMutation({
		mutationFn: (id: number) => api.proxies.delete(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["proxies"] });
			message.success(intl.formatMessage({ id: 'common.message.deleted' }));
		},
	});

	const checkProxyMutation = useMutation({
		mutationFn: (id: number) => api.proxies.check(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["proxies"] });
			message.success(intl.formatMessage({ id: 'environment.message.checkComplete' }));
		},
	});

	const handleAddFingerprint = (values: any) => { createFingerprintMutation.mutate(values); };
	const handleDeleteFingerprint = (id: number) => { deleteFingerprintMutation.mutate(id); };
	const handleToggleFingerprint = (id: number) => { toggleFingerprintMutation.mutate(id); };
	const handleAddProxy = (values: any) => { createProxyMutation.mutate(values); };
	const handleDeleteProxy = (id: number) => { deleteProxyMutation.mutate(id); };
	const handleCheckProxy = (id: number) => { checkProxyMutation.mutate(id); };

	const fingerprintColumns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 60 },
		{ title: intl.formatMessage({ id: 'environment.column.name' }), dataIndex: "name", key: "name", render: (text: string) => <span className="font-medium">{text}</span> },
		{ title: intl.formatMessage({ id: 'environment.column.browser' }), dataIndex: "user_agent", key: "user_agent", width: 120 },
		{ title: intl.formatMessage({ id: 'environment.column.platform' }), dataIndex: "platform", key: "platform", width: 100 },
		{ title: intl.formatMessage({ id: 'environment.column.screen' }), dataIndex: "screen", key: "screen", width: 120 },
		{ title: intl.formatMessage({ id: 'environment.column.language' }), dataIndex: "language", key: "language", width: 80 },
		{ title: intl.formatMessage({ id: 'environment.column.timezone' }), dataIndex: "timezone", key: "timezone", width: 150 },
		{ title: "WebGL", dataIndex: "webgl", key: "webgl", width: 80 },
		{
			title: intl.formatMessage({ id: 'environment.column.status' }),
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (status: string) => <Badge status={status === "active" ? "success" : "default"} text={status === "active" ? intl.formatMessage({ id: 'common.status.enabled' }) : intl.formatMessage({ id: 'common.status.disabled' })} />,
		},
		{ title: intl.formatMessage({ id: 'environment.column.accountCount' }), dataIndex: "account_count", key: "account_count", width: 80, render: (count: number) => <Badge count={count} showZero style={{ backgroundColor: "#1890ff" }} /> },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => message.info(intl.formatMessage({ id: 'environment.message.editInDev' }))} /></Tooltip>
					<Tooltip title={record.status === "active" ? intl.formatMessage({ id: 'common.action.disable' }) : intl.formatMessage({ id: 'common.action.enable' })}>
						<Switch size="small" checked={record.status === "active"} onChange={() => handleToggleFingerprint(record.id)} />
					</Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'environment.confirmDelete' })} onConfirm={() => handleDeleteFingerprint(record.id)}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	const proxyColumns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 60 },
		{ title: intl.formatMessage({ id: 'environment.column.ip' }), key: "ip", render: (_: any, record: any) => <code className="bg-gray-100 px-2 py-1 rounded">{record.ip}:{record.port}</code> },
		{ title: intl.formatMessage({ id: 'environment.column.protocol' }), dataIndex: "protocol", key: "protocol", width: 80, render: (t: string) => <Tag color="blue">{t}</Tag> },
		{ title: intl.formatMessage({ id: 'environment.column.location' }), dataIndex: "location", key: "location", width: 120 },
		{ title: intl.formatMessage({ id: 'environment.column.speed' }), dataIndex: "speed", key: "speed", width: 100, render: (speed: number) => <span className={speed < 100 ? "text-green-500" : speed < 200 ? "text-orange-500" : "text-red-500"}>{speed}ms</span> },
		{ title: intl.formatMessage({ id: 'environment.column.uptime' }), dataIndex: "uptime", key: "uptime", width: 120, render: (uptime: number) => <Progress percent={uptime} size="small" strokeColor={uptime >= 98 ? "#52c41a" : uptime >= 95 ? "#faad14" : "#ff4d4f"} /> },
		{
			title: intl.formatMessage({ id: 'environment.column.status' }),
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (status: string) => status === "active" ? <Badge status="success" text={intl.formatMessage({ id: 'environment.status.available' })} /> : <Badge status="error" text={intl.formatMessage({ id: 'environment.status.unavailable' })} />,
		},
		{ title: intl.formatMessage({ id: 'environment.column.lastCheck' }), dataIndex: "last_check", key: "last_check", width: 160, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 120,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.check' })}><Button type="text" icon={<ReloadOutlined />} onClick={() => handleCheckProxy(record.id)} /></Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'environment.confirmDelete' })} onConfirm={() => handleDeleteProxy(record.id)}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	const activeProxies = proxyPool.filter((p) => p.status === "active");
	const stats = {
		totalFingerprints: fingerprints.length,
		activeFingerprints: fingerprints.filter((f) => f.status === "active").length,
		totalProxies: proxyPool.length,
		activeProxies: activeProxies.length,
		avgSpeed: activeProxies.length ? Math.round(activeProxies.reduce((sum, p) => sum + p.speed, 0) / activeProxies.length) : 0,
		avgUptime: activeProxies.length ? (activeProxies.reduce((sum, p) => sum + p.uptime, 0) / activeProxies.length).toFixed(1) : "0.0",
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'environment.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'environment.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'environment.stat.fingerprints' })} value={stats.totalFingerprints} prefix={<DesktopOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'environment.stat.activeFingerprints' })} value={stats.activeFingerprints} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'environment.stat.proxies' })} value={stats.totalProxies} prefix={<GlobalOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'environment.stat.activeProxies' })} value={stats.activeProxies} valueStyle={{ color: "#1890ff" }} /></Card></Col>
			</Row>

			<Tabs
				activeKey={activeTab}
				onChange={setActiveTab}
				items={[
					{
						key: "fingerprints",
						label: <Space><DesktopOutlined /> {intl.formatMessage({ id: 'environment.tab.fingerprints' })}</Space>,
						children: (
							<Card
								title={intl.formatMessage({ id: 'environment.section.fingerprintList' })}
								extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'environment.action.createFingerprint' })}</Button>}
							>
								<Table columns={fingerprintColumns} dataSource={fingerprints} rowKey="id" pagination={false} />
							</Card>
						),
					},
					{
						key: "proxies",
						label: <Space><GlobalOutlined /> {intl.formatMessage({ id: 'environment.tab.proxies' })}</Space>,
						children: (
							<Card
								title={intl.formatMessage({ id: 'environment.section.proxyList' })}
								extra={
									<Space>
										<Button icon={<ReloadOutlined />} onClick={() => proxyPool.forEach((p) => checkProxyMutation.mutate(p.id))}>{intl.formatMessage({ id: 'environment.action.batchCheck' })}</Button>
										<Button type="primary" icon={<PlusOutlined />} onClick={() => setProxyModalVisible(true)}>{intl.formatMessage({ id: 'environment.action.addProxy' })}</Button>
									</Space>
								}
							>
								<Table columns={proxyColumns} dataSource={proxyPool} rowKey="id" pagination={false} />
							</Card>
						),
					},
					{
						key: "rules",
						label: <Space><SettingOutlined /> {intl.formatMessage({ id: 'environment.tab.rules' })}</Space>,
						children: (
							<Card title={intl.formatMessage({ id: 'environment.section.rules' })}>
								<div className="space-y-4">
									<div className="p-4 bg-blue-50 rounded-lg">
										<h4 className="font-medium mb-2">{intl.formatMessage({ id: 'environment.rule.autoAssign' })}</h4>
										<p className="text-sm text-gray-600 mb-3">{intl.formatMessage({ id: 'environment.rule.autoAssignDesc' })}</p>
										<Space>
											<Tag color="blue">{intl.formatMessage({ id: 'environment.rule.roundRobin' })}</Tag>
											<Tag color="green">{intl.formatMessage({ id: 'common.status.enabled' })}</Tag>
										</Space>
									</div>
									<div className="p-4 bg-green-50 rounded-lg">
										<h4 className="font-medium mb-2">{intl.formatMessage({ id: 'environment.rule.uniqueCheck' })}</h4>
										<p className="text-sm text-gray-600 mb-3">{intl.formatMessage({ id: 'environment.rule.uniqueCheckDesc' })}</p>
										<Space><Tag color="green">{intl.formatMessage({ id: 'environment.rule.checkRate' })}</Tag></Space>
									</div>
									<div className="p-4 bg-orange-50 rounded-lg">
										<h4 className="font-medium mb-2">{intl.formatMessage({ id: 'environment.rule.proxyRotation' })}</h4>
										<p className="text-sm text-gray-600 mb-3">{intl.formatMessage({ id: 'environment.rule.proxyRotationDesc' })}</p>
										<Space>
											<Tag color="orange">{intl.formatMessage({ id: 'environment.rule.rotationInterval' })}</Tag>
											<Tag color="blue">{intl.formatMessage({ id: 'environment.rule.uptime' })} {stats.avgUptime}%</Tag>
										</Space>
									</div>
								</div>
							</Card>
						),
					},
				]}
			/>

			<Modal title={intl.formatMessage({ id: 'environment.modal.createFingerprint' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={600}>
				<Form form={createForm} layout="vertical" onFinish={handleAddFingerprint}>
					<Form.Item name="name" label={intl.formatMessage({ id: 'environment.form.name' })} rules={[{ required: true }]}><Input placeholder={intl.formatMessage({ id: 'environment.placeholder.name' })} /></Form.Item>
					<Row gutter={16}>
						<Col span={12}>
							<Form.Item name="platform" label={intl.formatMessage({ id: 'environment.form.os' })} rules={[{ required: true }]}>
								<Select placeholder={intl.formatMessage({ id: 'common.placeholder.select' })}>
									<Option value="Windows">Windows</Option>
									<Option value="MacOS">MacOS</Option>
									<Option value="Linux">Linux</Option>
								</Select>
							</Form.Item>
						</Col>
						<Col span={12}>
							<Form.Item name="browser" label={intl.formatMessage({ id: 'environment.form.browser' })} rules={[{ required: true }]}>
								<Select placeholder={intl.formatMessage({ id: 'common.placeholder.select' })}>
									<Option value="Chrome">Chrome</Option>
									<Option value="Firefox">Firefox</Option>
									<Option value="Safari">Safari</Option>
									<Option value="Edge">Edge</Option>
								</Select>
							</Form.Item>
						</Col>
					</Row>
					<Row gutter={16}>
						<Col span={12}>
							<Form.Item name="screen" label={intl.formatMessage({ id: 'environment.form.screen' })} rules={[{ required: true }]}>
								<Select placeholder={intl.formatMessage({ id: 'common.placeholder.select' })}>
									<Option value="1920x1080">1920x1080</Option>
									<Option value="2560x1440">2560x1440</Option>
									<Option value="2560x1600">2560x1600</Option>
									<Option value="3840x2160">3840x2160</Option>
								</Select>
							</Form.Item>
						</Col>
						<Col span={12}>
							<Form.Item name="language" label={intl.formatMessage({ id: 'environment.form.language' })} rules={[{ required: true }]}>
								<Select placeholder={intl.formatMessage({ id: 'common.placeholder.select' })}>
									<Option value="zh-CN">{intl.formatMessage({ id: 'environment.language.zhCN' })}</Option>
									<Option value="en-US">English</Option>
									<Option value="ja-JP">日本語</Option>
								</Select>
							</Form.Item>
						</Col>
					</Row>
					<Row gutter={16}>
						<Col span={12}>
							<Form.Item name="timezone" label={intl.formatMessage({ id: 'environment.form.timezone' })} rules={[{ required: true }]}>
								<Select placeholder={intl.formatMessage({ id: 'common.placeholder.select' })}>
									<Option value="Asia/Shanghai">Asia/Shanghai</Option>
									<Option value="America/New_York">America/New_York</Option>
									<Option value="Europe/London">Europe/London</Option>
									<Option value="Asia/Tokyo">Asia/Tokyo</Option>
								</Select>
							</Form.Item>
						</Col>
						<Col span={12}>
							<Form.Item name="webgl" label={intl.formatMessage({ id: 'environment.form.webgl' })}>
								<Select placeholder={intl.formatMessage({ id: 'common.placeholder.select' })}>
									<Option value="Intel">Intel</Option>
									<Option value="NVIDIA">NVIDIA</Option>
									<Option value="AMD">AMD</Option>
									<Option value="Apple">Apple</Option>
								</Select>
							</Form.Item>
						</Col>
					</Row>
					<Form.Item name="user_agent" label={intl.formatMessage({ id: 'environment.form.userAgent' })}><Input placeholder={intl.formatMessage({ id: 'environment.placeholder.userAgent' })} /></Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title={intl.formatMessage({ id: 'environment.modal.addProxy' })} open={proxyModalVisible} onCancel={() => setProxyModalVisible(false)} footer={null} width={500}>
				<Form form={proxyForm} layout="vertical" onFinish={handleAddProxy}>
					<Row gutter={16}>
						<Col span={16}>
							<Form.Item name="ip" label={intl.formatMessage({ id: 'environment.form.ip' })} rules={[{ required: true }]}><Input placeholder={intl.formatMessage({ id: 'environment.placeholder.ip' })} /></Form.Item>
						</Col>
						<Col span={8}>
							<Form.Item name="port" label={intl.formatMessage({ id: 'environment.form.port' })} rules={[{ required: true }]}><InputNumber placeholder={intl.formatMessage({ id: 'environment.form.port' })} style={{ width: "100%" }} /></Form.Item>
						</Col>
					</Row>
					<Form.Item name="protocol" label={intl.formatMessage({ id: 'environment.form.protocol' })} rules={[{ required: true }]}>
						<Select placeholder={intl.formatMessage({ id: 'common.placeholder.select' })}>
							<Option value="HTTP">HTTP</Option>
							<Option value="HTTPS">HTTPS</Option>
							<Option value="SOCKS5">SOCKS5</Option>
						</Select>
					</Form.Item>
					<Form.Item name="username" label={intl.formatMessage({ id: 'environment.form.username' })}><Input placeholder={intl.formatMessage({ id: 'environment.placeholder.noAuth' })} /></Form.Item>
					<Form.Item name="password" label={intl.formatMessage({ id: 'environment.form.password' })}><Input.Password placeholder={intl.formatMessage({ id: 'environment.placeholder.noAuth' })} /></Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.add' })}</Button>
							<Button onClick={() => setProxyModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
