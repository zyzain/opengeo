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

const { Option } = Select;
const { TextArea } = Input;

export default function EnvironmentIsolationPage() {
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
			message.success("创建成功");
		},
	});

	const deleteFingerprintMutation = useMutation({
		mutationFn: (id: number) => api.fingerprints.delete(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["fingerprints"] });
			message.success("已删除");
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
			message.success("添加成功");
		},
	});

	const deleteProxyMutation = useMutation({
		mutationFn: (id: number) => api.proxies.delete(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["proxies"] });
			message.success("已删除");
		},
	});

	const checkProxyMutation = useMutation({
		mutationFn: (id: number) => api.proxies.check(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["proxies"] });
			message.success("检查完成");
		},
	});

	const handleAddFingerprint = (values: any) => {
		createFingerprintMutation.mutate(values);
	};

	const handleDeleteFingerprint = (id: number) => {
		deleteFingerprintMutation.mutate(id);
	};

	const handleToggleFingerprint = (id: number) => {
		toggleFingerprintMutation.mutate(id);
	};

	const handleAddProxy = (values: any) => {
		createProxyMutation.mutate(values);
	};

	const handleDeleteProxy = (id: number) => {
		deleteProxyMutation.mutate(id);
	};

	const handleCheckProxy = (id: number) => {
		checkProxyMutation.mutate(id);
	};

	// 指纹表格列
	const fingerprintColumns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 60 },
		{
			title: "配置名称",
			dataIndex: "name",
			key: "name",
			render: (text: string) => <span className="font-medium">{text}</span>,
		},
		{ title: "浏览器", dataIndex: "user_agent", key: "user_agent", width: 120 },
		{ title: "平台", dataIndex: "platform", key: "platform", width: 100 },
		{ title: "分辨率", dataIndex: "screen", key: "screen", width: 120 },
		{ title: "语言", dataIndex: "language", key: "language", width: 80 },
		{ title: "时区", dataIndex: "timezone", key: "timezone", width: 150 },
		{ title: "WebGL", dataIndex: "webgl", key: "webgl", width: 80 },
		{
			title: "状态",
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (status: string) => (
				<Badge
					status={status === "active" ? "success" : "default"}
					text={status === "active" ? "启用" : "禁用"}
				/>
			),
		},
		{
			title: "关联账号",
			dataIndex: "account_count",
			key: "account_count",
			width: 80,
			render: (count: number) => (
				<Badge count={count} showZero style={{ backgroundColor: "#1890ff" }} />
			),
		},
		{
			title: "操作",
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="编辑">
						<Button type="text" icon={<EditOutlined />} onClick={() => message.info("编辑功能开发中")} />
					</Tooltip>
					<Tooltip title={record.status === "active" ? "禁用" : "启用"}>
						<Switch
							size="small"
							checked={record.status === "active"}
							onChange={() => handleToggleFingerprint(record.id)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定删除?"
						onConfirm={() => handleDeleteFingerprint(record.id)}
					>
						<Tooltip title="删除">
							<Button type="text" danger icon={<DeleteOutlined />} />
						</Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	// 代理表格列
	const proxyColumns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 60 },
		{
			title: "IP地址",
			key: "ip",
			render: (_: any, record: any) => (
				<code className="bg-gray-100 px-2 py-1 rounded">
					{record.ip}:{record.port}
				</code>
			),
		},
		{
			title: "协议",
			dataIndex: "protocol",
			key: "protocol",
			width: 80,
			render: (t: string) => <Tag color="blue">{t}</Tag>,
		},
		{ title: "位置", dataIndex: "location", key: "location", width: 120 },
		{
			title: "速度",
			dataIndex: "speed",
			key: "speed",
			width: 100,
			render: (speed: number) => (
				<span
					className={
						speed < 100
							? "text-green-500"
							: speed < 200
								? "text-orange-500"
								: "text-red-500"
					}
				>
					{speed}ms
				</span>
			),
		},
		{
			title: "可用率",
			dataIndex: "uptime",
			key: "uptime",
			width: 120,
			render: (uptime: number) => (
				<Progress
					percent={uptime}
					size="small"
					strokeColor={
						uptime >= 98 ? "#52c41a" : uptime >= 95 ? "#faad14" : "#ff4d4f"
					}
				/>
			),
		},
		{
			title: "状态",
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (status: string) =>
				status === "active" ? (
					<Badge status="success" text="可用" />
				) : (
					<Badge status="error" text="不可用" />
				),
		},
		{
			title: "最后检查",
			dataIndex: "last_check",
			key: "last_check",
			width: 160,
			render: (text: string) => new Date(text).toLocaleString(),
		},
		{
			title: "操作",
			key: "action",
			width: 120,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="检查">
						<Button
							type="text"
							icon={<ReloadOutlined />}
							onClick={() => handleCheckProxy(record.id)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定删除?"
						onConfirm={() => handleDeleteProxy(record.id)}
					>
						<Tooltip title="删除">
							<Button type="text" danger icon={<DeleteOutlined />} />
						</Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	// 统计
	const activeProxies = proxyPool.filter((p) => p.status === "active");
	const stats = {
		totalFingerprints: fingerprints.length,
		activeFingerprints: fingerprints.filter((f) => f.status === "active")
			.length,
		totalProxies: proxyPool.length,
		activeProxies: activeProxies.length,
		avgSpeed: activeProxies.length
			? Math.round(activeProxies.reduce((sum, p) => sum + p.speed, 0) / activeProxies.length)
			: 0,
		avgUptime: activeProxies.length
			? (activeProxies.reduce((sum, p) => sum + p.uptime, 0) / activeProxies.length).toFixed(1)
			: "0.0",
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">环境隔离</h1>
				<p className="text-gray-500 mt-1">
					管理浏览器指纹和代理IP池，防止多账号关联
				</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="指纹配置"
							value={stats.totalFingerprints}
							prefix={<DesktopOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="启用中"
							value={stats.activeFingerprints}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="代理IP"
							value={stats.totalProxies}
							prefix={<GlobalOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="可用代理"
							value={stats.activeProxies}
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
			</Row>

			<Tabs
				activeKey={activeTab}
				onChange={setActiveTab}
				items={[
					{
						key: "fingerprints",
						label: (
							<Space>
								<DesktopOutlined /> 浏览器指纹
							</Space>
						),
						children: (
							<Card
								title="指纹配置列表"
								extra={
									<Button
										type="primary"
										icon={<PlusOutlined />}
										onClick={() => setCreateModalVisible(true)}
									>
										新建指纹
									</Button>
								}
							>
								<Table
									columns={fingerprintColumns}
									dataSource={fingerprints}
									rowKey="id"
									pagination={false}
								/>
							</Card>
						),
					},
					{
						key: "proxies",
						label: (
							<Space>
								<GlobalOutlined /> 代理IP池
							</Space>
						),
						children: (
							<Card
								title="代理IP列表"
								extra={
									<Space>
										<Button icon={<ReloadOutlined />} onClick={() => proxyPool.forEach((p) => checkProxyMutation.mutate(p.id))}>批量检测</Button>
										<Button
											type="primary"
											icon={<PlusOutlined />}
											onClick={() => setProxyModalVisible(true)}
										>
											添加代理
										</Button>
									</Space>
								}
							>
								<Table
									columns={proxyColumns}
									dataSource={proxyPool}
									rowKey="id"
									pagination={false}
								/>
							</Card>
						),
					},
					{
						key: "rules",
						label: (
							<Space>
								<SettingOutlined /> 分配规则
							</Space>
						),
						children: (
							<Card title="环境分配规则">
								<div className="space-y-4">
									<div className="p-4 bg-blue-50 rounded-lg">
										<h4 className="font-medium mb-2">自动分配策略</h4>
										<p className="text-sm text-gray-600 mb-3">
											新账号将自动分配指纹和代理IP，确保每个账号使用独立环境
										</p>
										<Space>
											<Tag color="blue">轮询分配</Tag>
											<Tag color="green">已启用</Tag>
										</Space>
									</div>
									<div className="p-4 bg-green-50 rounded-lg">
										<h4 className="font-medium mb-2">指纹唯一性校验</h4>
										<p className="text-sm text-gray-600 mb-3">
											自动检测指纹重复，确保每个账号使用唯一指纹
										</p>
										<Space>
											<Tag color="green">校验通过率 99.5%</Tag>
										</Space>
									</div>
									<div className="p-4 bg-orange-50 rounded-lg">
										<h4 className="font-medium mb-2">代理IP轮换</h4>
										<p className="text-sm text-gray-600 mb-3">
											按时间间隔自动轮换代理IP，降低被检测风险
										</p>
										<Space>
											<Tag color="orange">轮换间隔: 30分钟</Tag>
											<Tag color="blue">可用率: {stats.avgUptime}%</Tag>
										</Space>
									</div>
								</div>
							</Card>
						),
					},
				]}
			/>

			{/* 新建指纹弹窗 */}
			<Modal
				title="新建指纹配置"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={600}
			>
				<Form
					form={createForm}
					layout="vertical"
					onFinish={handleAddFingerprint}
				>
					<Form.Item name="name" label="配置名称" rules={[{ required: true }]}>
						<Input placeholder="请输入配置名称" />
					</Form.Item>
					<Row gutter={16}>
						<Col span={12}>
							<Form.Item
								name="platform"
								label="操作系统"
								rules={[{ required: true }]}
							>
								<Select placeholder="请选择">
									<Option value="Windows">Windows</Option>
									<Option value="MacOS">MacOS</Option>
									<Option value="Linux">Linux</Option>
								</Select>
							</Form.Item>
						</Col>
						<Col span={12}>
							<Form.Item
								name="browser"
								label="浏览器"
								rules={[{ required: true }]}
							>
								<Select placeholder="请选择">
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
							<Form.Item
								name="screen"
								label="屏幕分辨率"
								rules={[{ required: true }]}
							>
								<Select placeholder="请选择">
									<Option value="1920x1080">1920x1080</Option>
									<Option value="2560x1440">2560x1440</Option>
									<Option value="2560x1600">2560x1600</Option>
									<Option value="3840x2160">3840x2160</Option>
								</Select>
							</Form.Item>
						</Col>
						<Col span={12}>
							<Form.Item
								name="language"
								label="语言"
								rules={[{ required: true }]}
							>
								<Select placeholder="请选择">
									<Option value="zh-CN">中文(简体)</Option>
									<Option value="en-US">English</Option>
									<Option value="ja-JP">日本語</Option>
								</Select>
							</Form.Item>
						</Col>
					</Row>
					<Row gutter={16}>
						<Col span={12}>
							<Form.Item
								name="timezone"
								label="时区"
								rules={[{ required: true }]}
							>
								<Select placeholder="请选择">
									<Option value="Asia/Shanghai">Asia/Shanghai</Option>
									<Option value="America/New_York">America/New_York</Option>
									<Option value="Europe/London">Europe/London</Option>
									<Option value="Asia/Tokyo">Asia/Tokyo</Option>
								</Select>
							</Form.Item>
						</Col>
						<Col span={12}>
							<Form.Item name="webgl" label="WebGL渲染器">
								<Select placeholder="请选择">
									<Option value="Intel">Intel</Option>
									<Option value="NVIDIA">NVIDIA</Option>
									<Option value="AMD">AMD</Option>
									<Option value="Apple">Apple</Option>
								</Select>
							</Form.Item>
						</Col>
					</Row>
					<Form.Item name="user_agent" label="自定义User-Agent">
						<Input placeholder="留空则自动生成" />
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

			{/* 添加代理弹窗 */}
			<Modal
				title="添加代理IP"
				open={proxyModalVisible}
				onCancel={() => setProxyModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form form={proxyForm} layout="vertical" onFinish={handleAddProxy}>
					<Row gutter={16}>
						<Col span={16}>
							<Form.Item name="ip" label="IP地址" rules={[{ required: true }]}>
								<Input placeholder="请输入IP地址" />
							</Form.Item>
						</Col>
						<Col span={8}>
							<Form.Item name="port" label="端口" rules={[{ required: true }]}>
								<InputNumber placeholder="端口" style={{ width: "100%" }} />
							</Form.Item>
						</Col>
					</Row>
					<Form.Item name="protocol" label="协议" rules={[{ required: true }]}>
						<Select placeholder="请选择">
							<Option value="HTTP">HTTP</Option>
							<Option value="HTTPS">HTTPS</Option>
							<Option value="SOCKS5">SOCKS5</Option>
						</Select>
					</Form.Item>
					<Form.Item name="username" label="用户名">
						<Input placeholder="留空表示无需认证" />
					</Form.Item>
					<Form.Item name="password" label="密码">
						<Input.Password placeholder="留空表示无需认证" />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">
								添加
							</Button>
							<Button onClick={() => setProxyModalVisible(false)}>取消</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
