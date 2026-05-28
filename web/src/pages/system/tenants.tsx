"use client";

import api from "@/lib/api";
import {
	ApartmentOutlined,
	CheckCircleOutlined,
	CrownOutlined,
	DatabaseOutlined,
	DeleteOutlined,
	EditOutlined,
	PlusOutlined,
	ReloadOutlined,
	SettingOutlined,
	TeamOutlined,
	UserOutlined,
} from "@ant-design/icons";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
	Badge,
	Button,
	Card,
	Col,
	Descriptions,
	Form,
	Input,
	Modal,
	Popconfirm,
	Progress,
	Row,
	Select,
	Space,
	Statistic,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";

const { Option } = Select;
const { TextArea } = Input;

export default function TenantManagementPage() {
	const queryClient = useQueryClient();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [detailModalVisible, setDetailModalVisible] = useState(false);
	const [selectedTenant, setSelectedTenant] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();

	const { data: tenantsData, isLoading } = useQuery({
		queryKey: ["tenants"],
		queryFn: () => api.tenants.list(),
		select: (response) => response.data.data,
	});

	const tenants = tenantsData?.items || [];

	const { data: plansData } = useQuery({
		queryKey: ["plans"],
		queryFn: () => api.system.getPlans(),
	});

	const createMutation = useMutation({
		mutationFn: api.tenants.create,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["tenants"] });
			message.success("创建成功");
			createForm.resetFields();
			setCreateModalVisible(false);
		},
	});

	const updateMutation = useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) =>
			api.tenants.update(id, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["tenants"] });
			message.success("更新成功");
			setEditModalVisible(false);
		},
	});

	const deleteMutation = useMutation({
		mutationFn: api.tenants.delete,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["tenants"] });
			message.success("已删除");
		},
	});

	const defaultPlans: Record<string, any> = {
		starter: {
			label: "入门版",
			color: "blue",
			maxUsers: 5,
			maxStorage: 2,
			maxContents: 500,
		},
		professional: {
			label: "专业版",
			color: "green",
			maxUsers: 20,
			maxStorage: 10,
			maxContents: 5000,
		},
		enterprise: {
			label: "企业版",
			color: "purple",
			maxUsers: 100,
			maxStorage: 100,
			maxContents: 50000,
		},
	};

	const plansList = plansData?.data?.data?.items || [];
	const plans: Record<string, any> =
		plansList.length > 0
			? plansList.reduce((acc: any, plan: any) => {
					acc[plan.id] = plan;
					return acc;
				}, {})
			: defaultPlans;

	// 表格列
	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 60 },
		{
			title: "租户名称",
			dataIndex: "name",
			key: "name",
			render: (text: string, record: any) => (
				<a
					onClick={() => {
						setSelectedTenant(record);
						setDetailModalVisible(true);
					}}
					className="text-blue-500 font-medium"
				>
					{text}
				</a>
			),
		},
		{
			title: "域名",
			dataIndex: "domain",
			key: "domain",
			render: (t: string) => (
				<code className="bg-gray-100 px-2 py-1 rounded text-sm">{t}</code>
			),
		},
		{
			title: "套餐",
			dataIndex: "plan",
			key: "plan",
			width: 100,
			render: (plan: string) => (
				<Tag color={plans[plan]?.color}>{plans[plan]?.label || plans[plan]?.name}</Tag>
			),
		},
		{
			title: "用户数",
			dataIndex: "users",
			key: "users",
			width: 80,
			render: (n: number) => (
				<Badge count={n} showZero style={{ backgroundColor: "#1890ff" }} />
			),
		},
		{ title: "内容数", dataIndex: "contents", key: "contents", width: 80 },
		{
			title: "存储使用",
			key: "storage",
			width: 150,
			render: (_: any, record: any) => (
				<Progress
					percent={Math.round(
						(record.storage_used / record.storage_limit) * 100,
					)}
					size="small"
					strokeColor={
						record.storage_used / record.storage_limit > 0.8
							? "#ff4d4f"
							: "#1890ff"
					}
					format={() => `${record.storage_used}/${record.storage_limit}GB`}
				/>
			),
		},
		{
			title: "状态",
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (s: string) => (
				<Badge
					status={s === "active" ? "success" : "default"}
					text={s === "active" ? "正常" : "禁用"}
				/>
			),
		},
		{
			title: "操作",
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="详情">
						<Button
							type="text"
							icon={<SettingOutlined />}
							onClick={() => {
								setSelectedTenant(record);
								setDetailModalVisible(true);
							}}
						/>
					</Tooltip>
					<Tooltip title="编辑">
						<Button
							type="text"
							icon={<EditOutlined />}
							onClick={() => {
								setSelectedTenant(record);
								editForm.setFieldsValue(record);
								setEditModalVisible(true);
							}}
						/>
					</Tooltip>
					{record.id !== 1 && (
						<Popconfirm
							title="确定删除?"
							onConfirm={() => deleteMutation.mutate(record.id)}
						>
							<Tooltip title="删除">
								<Button type="text" danger icon={<DeleteOutlined />} />
							</Tooltip>
						</Popconfirm>
					)}
				</Space>
			),
		},
	];

	// 统计
	const stats = {
		total: tenants.length,
		active: tenants.filter((t: any) => t.status === "active").length,
		totalUsers: tenants.reduce((sum: number, t: any) => sum + (t.users || 0), 0),
		totalContents: tenants.reduce((sum: number, t: any) => sum + (t.contents || 0), 0),
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">租户管理</h1>
				<p className="text-gray-500 mt-1">
					管理多租户SaaS模式，支持数据隔离和资源配额配置
				</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="租户总数"
							value={stats.total}
							prefix={<ApartmentOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="活跃租户"
							value={stats.active}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="总用户数"
							value={stats.totalUsers}
							prefix={<UserOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="总内容数"
							value={stats.totalContents}
							prefix={<DatabaseOutlined />}
						/>
					</Card>
				</Col>
			</Row>

			{/* 套餐对比 */}
			<Card title="套餐配置" className="mb-4">
				<Row gutter={[16, 16]}>
					{Object.entries(plans).map(([key, plan]) => (
						<Col xs={24} sm={8} key={key}>
							<div
								className={`p-4 rounded-lg border ${key === "enterprise" ? "border-purple-300 bg-purple-50" : key === "professional" ? "border-green-300 bg-green-50" : "border-blue-300 bg-blue-50"}`}
							>
								<div className="flex items-center justify-between mb-2">
									<Tag color={plan.color} className="text-lg">
										{plan.label || plan.name}
									</Tag>
									{key === "enterprise" && (
										<CrownOutlined className="text-purple-500 text-xl" />
									)}
								</div>
								<div className="space-y-1 text-sm">
									<div>最大用户: {plan.max_users ?? plan.maxUsers}</div>
									<div>存储空间: {plan.max_storage ?? plan.maxStorage}GB</div>
									<div>内容数量: {(plan.max_contents ?? plan.maxContents).toLocaleString()}</div>
								</div>
							</div>
						</Col>
					))}
				</Row>
			</Card>

			{/* 租户列表 */}
			<Card
				title="租户列表"
				extra={
					<Space>
						<Button
							icon={<ReloadOutlined />}
							onClick={() =>
								queryClient.invalidateQueries({ queryKey: ["tenants"] })
							}
						>
							刷新
						</Button>
						<Button
							type="primary"
							icon={<PlusOutlined />}
							onClick={() => setCreateModalVisible(true)}
						>
							新建租户
						</Button>
					</Space>
				}
			>
				<Table
					columns={columns}
					dataSource={tenants}
					rowKey="id"
					pagination={false}
					loading={isLoading}
				/>
			</Card>

			{/* 新建租户弹窗 */}
			<Modal
				title="新建租户"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form
					form={createForm}
					layout="vertical"
					onFinish={(values) => {
						createMutation.mutate(values);
					}}
				>
					<Form.Item name="name" label="租户名称" rules={[{ required: true }]}>
						<Input placeholder="请输入租户名称" />
					</Form.Item>
					<Form.Item name="domain" label="域名" rules={[{ required: true }]}>
						<Input placeholder="xxx.opengeo.com" />
					</Form.Item>
					<Form.Item name="plan" label="套餐" rules={[{ required: true }]}>
						<Select placeholder="请选择套餐">
							{Object.entries(plans).map(([key, plan]) => (
								<Option key={key} value={key}>
									<Tag color={plan.color}>{plan.label || plan.name}</Tag> - 最多
									{plan.max_users ?? plan.maxUsers}用户，{plan.max_storage ?? plan.maxStorage}GB存储
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item
						name="admin_email"
						label="管理员邮箱"
						rules={[{ required: true }, { type: "email" }]}
					>
						<Input placeholder="管理员的邮箱地址" />
					</Form.Item>
					<Form.Item name="description" label="描述">
						<TextArea rows={3} placeholder="租户描述" />
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

			{/* 租户详情弹窗 */}
			<Modal
				title="租户详情"
				open={detailModalVisible}
				onCancel={() => setDetailModalVisible(false)}
				footer={null}
				width={600}
			>
				{selectedTenant && (
					<Descriptions column={2} bordered>
						<Descriptions.Item label="租户名称" span={2}>
							{selectedTenant.name}
						</Descriptions.Item>
						<Descriptions.Item label="域名">
							<code>{selectedTenant.domain}</code>
						</Descriptions.Item>
						<Descriptions.Item label="套餐">
							<Tag color={plans[selectedTenant.plan]?.color}>
								{plans[selectedTenant.plan]?.label || plans[selectedTenant.plan]?.name}
							</Tag>
						</Descriptions.Item>
						<Descriptions.Item label="用户数">
							{selectedTenant.users} / {plans[selectedTenant.plan]?.max_users ?? plans[selectedTenant.plan]?.maxUsers}
						</Descriptions.Item>
						<Descriptions.Item label="内容数">
							{selectedTenant.contents} /{" "}
							{plans[selectedTenant.plan]?.max_contents ?? plans[selectedTenant.plan]?.maxContents}
						</Descriptions.Item>
						<Descriptions.Item label="存储使用" span={2}>
							<Progress
								percent={Math.round(
									(selectedTenant.storage_used / selectedTenant.storage_limit) *
										100,
								)}
								format={() =>
									`${selectedTenant.storage_used}GB / ${selectedTenant.storage_limit}GB`
								}
							/>
						</Descriptions.Item>
						<Descriptions.Item label="状态">
							<Badge
								status={
									selectedTenant.status === "active" ? "success" : "default"
								}
								text={selectedTenant.status === "active" ? "正常" : "禁用"}
							/>
						</Descriptions.Item>
						<Descriptions.Item label="创建时间">
							{new Date(selectedTenant.created_at).toLocaleString()}
						</Descriptions.Item>
					</Descriptions>
				)}
			</Modal>

			{/* 编辑租户弹窗 */}
			<Modal
				title="编辑租户"
				open={editModalVisible}
				onCancel={() => setEditModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form
					form={editForm}
					layout="vertical"
					onFinish={(values) => {
						updateMutation.mutate({ id: selectedTenant.id, data: values });
					}}
				>
					<Form.Item name="name" label="租户名称" rules={[{ required: true }]}>
						<Input placeholder="请输入租户名称" />
					</Form.Item>
					<Form.Item name="domain" label="域名" rules={[{ required: true }]}>
						<Input placeholder="xxx.opengeo.com" />
					</Form.Item>
					<Form.Item name="plan" label="套餐" rules={[{ required: true }]}>
						<Select placeholder="请选择套餐">
							{Object.entries(plans).map(([key, plan]) => (
								<Option key={key} value={key}>
									<Tag color={plan.color}>{plan.label || plan.name}</Tag> - 最多
									{plan.max_users ?? plan.maxUsers}用户，{plan.max_storage ?? plan.maxStorage}GB存储
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="status" label="状态" rules={[{ required: true }]}>
						<Select>
							<Option value="active">正常</Option>
							<Option value="inactive">禁用</Option>
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">
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
