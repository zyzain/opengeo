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
import { useIntl } from "react-intl";

const { Option } = Select;
const { TextArea } = Input;

export default function TenantManagementPage() {
	const intl = useIntl();
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
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			createForm.resetFields();
			setCreateModalVisible(false);
		},
	});

	const updateMutation = useMutation({
		mutationFn: ({ id, data }: { id: number; data: any }) => api.tenants.update(id, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["tenants"] });
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
		},
	});

	const deleteMutation = useMutation({
		mutationFn: api.tenants.delete,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["tenants"] });
			message.success(intl.formatMessage({ id: 'common.message.deleted' }));
		},
	});

	const defaultPlans: Record<string, any> = {
		starter: { label: intl.formatMessage({ id: 'tenant.plan.starter' }), color: "blue", maxUsers: 5, maxStorage: 2, maxContents: 500 },
		professional: { label: intl.formatMessage({ id: 'tenant.plan.professional' }), color: "green", maxUsers: 20, maxStorage: 10, maxContents: 5000 },
		enterprise: { label: intl.formatMessage({ id: 'tenant.plan.enterprise' }), color: "purple", maxUsers: 100, maxStorage: 100, maxContents: 50000 },
	};

	const plansList = plansData?.data?.data?.items || [];
	const plans: Record<string, any> = plansList.length > 0
		? plansList.reduce((acc: any, plan: any) => { acc[plan.id] = plan; return acc; }, {})
		: defaultPlans;

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 60 },
		{
			title: intl.formatMessage({ id: 'tenant.column.name' }),
			dataIndex: "name",
			key: "name",
			render: (text: string, record: any) => (
				<a onClick={() => { setSelectedTenant(record); setDetailModalVisible(true); }} className="text-blue-500 font-medium">{text}</a>
			),
		},
		{
			title: intl.formatMessage({ id: 'tenant.column.domain' }),
			dataIndex: "domain",
			key: "domain",
			render: (t: string) => <code className="bg-gray-100 px-2 py-1 rounded text-sm">{t}</code>,
		},
		{
			title: intl.formatMessage({ id: 'tenant.column.plan' }),
			dataIndex: "plan",
			key: "plan",
			width: 100,
			render: (plan: string) => <Tag color={plans[plan]?.color}>{plans[plan]?.label || plans[plan]?.name}</Tag>,
		},
		{ title: intl.formatMessage({ id: 'tenant.column.userCount' }), dataIndex: "users", key: "users", width: 80, render: (n: number) => <Badge count={n} showZero style={{ backgroundColor: "#1890ff" }} /> },
		{ title: intl.formatMessage({ id: 'tenant.column.contentCount' }), dataIndex: "contents", key: "contents", width: 80 },
		{
			title: intl.formatMessage({ id: 'tenant.column.storage' }),
			key: "storage",
			width: 150,
			render: (_: any, record: any) => (
				<Progress percent={Math.round((record.storage_used / record.storage_limit) * 100)} size="small" strokeColor={record.storage_used / record.storage_limit > 0.8 ? "#ff4d4f" : "#1890ff"} format={() => `${record.storage_used}/${record.storage_limit}GB`} />
			),
		},
		{
			title: intl.formatMessage({ id: 'tenant.column.status' }),
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (s: string) => <Badge status={s === "active" ? "success" : "default"} text={s === "active" ? intl.formatMessage({ id: 'tenant.status.active' }) : intl.formatMessage({ id: 'tenant.status.disabled' })} />,
		},
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.view' })}><Button type="text" icon={<SettingOutlined />} onClick={() => { setSelectedTenant(record); setDetailModalVisible(true); }} /></Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => { setSelectedTenant(record); editForm.setFieldsValue(record); setEditModalVisible(true); }} /></Tooltip>
					{record.id !== 1 && (
						<Popconfirm title={intl.formatMessage({ id: 'tenant.confirmDelete' })} onConfirm={() => deleteMutation.mutate(record.id)}>
							<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
						</Popconfirm>
					)}
				</Space>
			),
		},
	];

	const stats = {
		total: tenants.length,
		active: tenants.filter((t: any) => t.status === "active").length,
		totalUsers: tenants.reduce((sum: number, t: any) => sum + (t.users || 0), 0),
		totalContents: tenants.reduce((sum: number, t: any) => sum + (t.contents || 0), 0),
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'tenant.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'tenant.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'tenant.stat.total' })} value={stats.total} prefix={<ApartmentOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'tenant.stat.active' })} value={stats.active} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'tenant.stat.totalUsers' })} value={stats.totalUsers} prefix={<UserOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'tenant.stat.totalContents' })} value={stats.totalContents} prefix={<DatabaseOutlined />} /></Card></Col>
			</Row>

			<Card title={intl.formatMessage({ id: 'tenant.section.planConfig' })} className="mb-4">
				<Row gutter={[16, 16]}>
					{Object.entries(plans).map(([key, plan]) => (
						<Col xs={24} sm={8} key={key}>
							<div className={`p-4 rounded-lg border ${key === "enterprise" ? "border-purple-300 bg-purple-50" : key === "professional" ? "border-green-300 bg-green-50" : "border-blue-300 bg-blue-50"}`}>
								<div className="flex items-center justify-between mb-2">
									<Tag color={plan.color} className="text-lg">{plan.label || plan.name}</Tag>
									{key === "enterprise" && <CrownOutlined className="text-purple-500 text-xl" />}
								</div>
								<div className="space-y-1 text-sm">
									<div>{intl.formatMessage({ id: 'tenant.plan.maxUsers' })} {plan.max_users ?? plan.maxUsers}</div>
									<div>{intl.formatMessage({ id: 'tenant.plan.storage' })} {plan.max_storage ?? plan.maxStorage}GB</div>
									<div>{intl.formatMessage({ id: 'tenant.plan.maxContents' })} {(plan.max_contents ?? plan.maxContents).toLocaleString()}</div>
								</div>
							</div>
						</Col>
					))}
				</Row>
			</Card>

			<Card
				title={intl.formatMessage({ id: 'tenant.section.list' })}
				extra={
					<Space>
						<Button icon={<ReloadOutlined />} onClick={() => queryClient.invalidateQueries({ queryKey: ["tenants"] })}>{intl.formatMessage({ id: 'common.action.refresh' })}</Button>
						<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'tenant.action.create' })}</Button>
					</Space>
				}
			>
				<Table columns={columns} dataSource={tenants} rowKey="id" pagination={false} loading={isLoading} />
			</Card>

			<Modal title={intl.formatMessage({ id: 'tenant.modal.create' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={500}>
				<Form form={createForm} layout="vertical" onFinish={(values) => createMutation.mutate(values)}>
					<Form.Item name="name" label={intl.formatMessage({ id: 'tenant.form.name' })} rules={[{ required: true }]}><Input placeholder={intl.formatMessage({ id: 'tenant.placeholder.name' })} /></Form.Item>
					<Form.Item name="domain" label={intl.formatMessage({ id: 'tenant.form.domain' })} rules={[{ required: true }]}><Input placeholder="xxx.opengeo.com" /></Form.Item>
					<Form.Item name="plan" label={intl.formatMessage({ id: 'tenant.form.plan' })} rules={[{ required: true }]}>
						<Select placeholder={intl.formatMessage({ id: 'tenant.placeholder.plan' })}>
							{Object.entries(plans).map(([key, plan]) => (
								<Option key={key} value={key}><Tag color={plan.color}>{plan.label || plan.name}</Tag> - {intl.formatMessage({ id: 'tenant.plan.maxUsersLabel' })}{plan.max_users ?? plan.maxUsers}{intl.formatMessage({ id: 'tenant.plan.userLabel' })}{plan.max_storage ?? plan.maxStorage}{intl.formatMessage({ id: 'tenant.plan.storageLabelGB' })}</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="admin_email" label={intl.formatMessage({ id: 'tenant.form.adminEmail' })} rules={[{ required: true }, { type: "email" }]}><Input placeholder={intl.formatMessage({ id: 'tenant.placeholder.adminEmail' })} /></Form.Item>
					<Form.Item name="description" label={intl.formatMessage({ id: 'tenant.form.description' })}><TextArea rows={3} placeholder={intl.formatMessage({ id: 'tenant.placeholder.description' })} /></Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title={intl.formatMessage({ id: 'tenant.modal.detail' })} open={detailModalVisible} onCancel={() => setDetailModalVisible(false)} footer={null} width={600}>
				{selectedTenant && (
					<Descriptions column={2} bordered>
						<Descriptions.Item label={intl.formatMessage({ id: 'tenant.desc.name' })} span={2}>{selectedTenant.name}</Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: 'tenant.desc.domain' })}><code>{selectedTenant.domain}</code></Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: 'tenant.desc.plan' })}><Tag color={plans[selectedTenant.plan]?.color}>{plans[selectedTenant.plan]?.label || plans[selectedTenant.plan]?.name}</Tag></Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: 'tenant.desc.userCount' })}>{selectedTenant.users} / {plans[selectedTenant.plan]?.max_users ?? plans[selectedTenant.plan]?.maxUsers}</Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: 'tenant.desc.contentCount' })}>{selectedTenant.contents} / {plans[selectedTenant.plan]?.max_contents ?? plans[selectedTenant.plan]?.maxContents}</Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: 'tenant.desc.storage' })} span={2}><Progress percent={Math.round((selectedTenant.storage_used / selectedTenant.storage_limit) * 100)} format={() => `${selectedTenant.storage_used}GB / ${selectedTenant.storage_limit}GB`} /></Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: 'tenant.desc.status' })}><Badge status={selectedTenant.status === "active" ? "success" : "default"} text={selectedTenant.status === "active" ? intl.formatMessage({ id: 'tenant.status.active' }) : intl.formatMessage({ id: 'tenant.status.disabled' })} /></Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: 'common.column.createdAt' })}>{new Date(selectedTenant.created_at).toLocaleString()}</Descriptions.Item>
					</Descriptions>
				)}
			</Modal>

			<Modal title={intl.formatMessage({ id: 'tenant.modal.edit' })} open={editModalVisible} onCancel={() => setEditModalVisible(false)} footer={null} width={500}>
				<Form form={editForm} layout="vertical" onFinish={(values) => updateMutation.mutate({ id: selectedTenant.id, data: values })}>
					<Form.Item name="name" label={intl.formatMessage({ id: 'tenant.form.name' })} rules={[{ required: true }]}><Input placeholder={intl.formatMessage({ id: 'tenant.placeholder.name' })} /></Form.Item>
					<Form.Item name="domain" label={intl.formatMessage({ id: 'tenant.form.domain' })} rules={[{ required: true }]}><Input placeholder="xxx.opengeo.com" /></Form.Item>
					<Form.Item name="plan" label={intl.formatMessage({ id: 'tenant.form.plan' })} rules={[{ required: true }]}>
						<Select placeholder={intl.formatMessage({ id: 'tenant.placeholder.plan' })}>
							{Object.entries(plans).map(([key, plan]) => (
								<Option key={key} value={key}><Tag color={plan.color}>{plan.label || plan.name}</Tag> - {intl.formatMessage({ id: 'tenant.plan.maxUsersLabel' })}{plan.max_users ?? plan.maxUsers}{intl.formatMessage({ id: 'tenant.plan.userLabel' })}{plan.max_storage ?? plan.maxStorage}{intl.formatMessage({ id: 'tenant.plan.storageLabelGB' })}</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="status" label={intl.formatMessage({ id: 'tenant.column.status' })} rules={[{ required: true }]}>
						<Select>
							<Option value="active">{intl.formatMessage({ id: 'tenant.status.active' })}</Option>
							<Option value="inactive">{intl.formatMessage({ id: 'tenant.status.disabled' })}</Option>
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.save' })}</Button>
							<Button onClick={() => setEditModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
