"use client";

import api from "@/lib/api";
import {
	DeleteOutlined,
	EditOutlined,
	KeyOutlined,
	PlusOutlined,
	ReloadOutlined,
	SafetyOutlined,
	SearchOutlined,
	TeamOutlined,
	UserOutlined,
} from "@ant-design/icons";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
	Badge,
	Button,
	Card,
	Col,
	Form,
	Input,
	Modal,
	Popconfirm,
	Row,
	Select,
	Space,
	Statistic,
	Table,
	Tabs,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";
import { useIntl } from "react-intl";

const { Option } = Select;

export default function UserManagementPage() {
	const intl = useIntl();
	const queryClient = useQueryClient();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [roleModalVisible, setRoleModalVisible] = useState(false);
	const [editingUser, setEditingUser] = useState<any>(null);
	const [selectedUser, setSelectedUser] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();
	const [roleForm] = Form.useForm();

	const { data: usersData, isLoading } = useQuery({
		queryKey: ["users"],
		queryFn: () => api.users.list(),
		select: (response) => response.data.data,
	});

	const { data: rolesData } = useQuery({
		queryKey: ["roles"],
		queryFn: () => api.roles.list(),
		select: (response) => response.data.data,
	});

	const users = usersData?.items || [];
	const roles = rolesData?.items || [];

	const userStatuses = [
		{ value: 1, label: intl.formatMessage({ id: 'common.status.active' }), color: "success" },
		{ value: 0, label: intl.formatMessage({ id: 'common.status.disabled' }), color: "error" },
		{ value: 2, label: intl.formatMessage({ id: 'users.status.pending' }), color: "warning" },
	];

	const getStatusTag = (status: number) => {
		const statusInfo = userStatuses.find((s) => s.value === status);
		return <Badge status={statusInfo?.color as any} text={statusInfo?.label || intl.formatMessage({ id: 'common.status.unknown' })} />;
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{ title: intl.formatMessage({ id: 'user.username' }), dataIndex: "username", key: "username", render: (text: string) => <Space><UserOutlined /><span className="font-medium">{text}</span></Space> },
		{ title: intl.formatMessage({ id: 'user.email' }), dataIndex: "email", key: "email", ellipsis: true },
		{ title: intl.formatMessage({ id: 'user.status' }), dataIndex: "status", key: "status", width: 100, render: (status: number) => getStatusTag(status) },
		{ title: intl.formatMessage({ id: 'common.column.createdAt' }), dataIndex: "created_at", key: "created_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{ title: intl.formatMessage({ id: 'users.column.lastLogin' }), dataIndex: "last_login_at", key: "last_login_at", width: 180, render: (text: string) => (text ? new Date(text).toLocaleString() : "-") },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => handleEdit(record)} /></Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'users.column.action.assignRole' })}><Button type="text" icon={<KeyOutlined />} onClick={() => handleAssignRole(record)} /></Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'common.confirmDelete' })} onConfirm={() => handleDelete(record.id)} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	const handleCreate = async (values: any) => {
		try {
			await api.auth.register(values);
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: ["users"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.createFailed' }));
		}
	};

	const handleEdit = (record: any) => {
		setEditingUser(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await api.users.update(editingUser.id, values);
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingUser(null);
			queryClient.invalidateQueries({ queryKey: ["users"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.users.delete(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
			queryClient.invalidateQueries({ queryKey: ["users"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const handleAssignRole = (record: any) => {
		setSelectedUser(record);
		roleForm.resetFields();
		setRoleModalVisible(true);
	};

	const handleRoleSubmit = async (values: any) => {
		try {
			await api.roles.assign({ user_id: selectedUser.id, role_id: values.role_id });
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setRoleModalVisible(false);
			roleForm.resetFields();
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	const stats = {
		total: users.length,
		active: users.filter((u: any) => u.status === 1).length,
		disabled: users.filter((u: any) => u.status === 0).length,
		pending: users.filter((u: any) => u.status === 2).length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.system.users' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'users.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'users.stat.total' })} value={stats.total} prefix={<TeamOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'users.stat.active' })} value={stats.active} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'users.stat.disabled' })} value={stats.disabled} valueStyle={{ color: "#ff4d4f" }} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'users.stat.pending' })} value={stats.pending} valueStyle={{ color: "#faad14" }} /></Card></Col>
			</Row>

			<Card
				title={<Space><UserOutlined /><span>{intl.formatMessage({ id: 'users.card.userList' })}</span></Space>}
				extra={
					<Space>
						<Button icon={<ReloadOutlined />} onClick={() => queryClient.invalidateQueries({ queryKey: ["users"] })}>{intl.formatMessage({ id: 'common.action.refresh' })}</Button>
						<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'users.action.addUser' })}</Button>
					</Space>
				}
			>
				<Table columns={columns} dataSource={users} rowKey="id" loading={isLoading} pagination={{ showSizeChanger: true, showQuickJumper: true, showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }) }} />
			</Card>

		<Modal title={intl.formatMessage({ id: 'users.modal.addTitle' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={500}>
			<Form form={createForm} layout="vertical" onFinish={handleCreate}>
				<Form.Item name="username" label={intl.formatMessage({ id: 'user.username' })} rules={[{ required: true, message: intl.formatMessage({ id: 'users.validation.enterUsername' }) }, { min: 3, message: intl.formatMessage({ id: 'users.validation.usernameMin' }) }]}>
					<Input placeholder={intl.formatMessage({ id: 'users.placeholder.username' })} />
				</Form.Item>
				<Form.Item name="email" label={intl.formatMessage({ id: 'user.email' })} rules={[{ required: true, message: intl.formatMessage({ id: 'users.validation.enterEmail' }) }, { type: "email", message: intl.formatMessage({ id: 'users.validation.validEmail' }) }]}>
					<Input placeholder={intl.formatMessage({ id: 'users.placeholder.email' })} />
				</Form.Item>
				<Form.Item name="password" label={intl.formatMessage({ id: 'users.form.password' })} rules={[{ required: true, message: intl.formatMessage({ id: 'users.validation.enterPassword' }) }, { min: 8, message: intl.formatMessage({ id: 'users.validation.passwordMin' }) }]}>
					<Input.Password placeholder={intl.formatMessage({ id: 'users.placeholder.password' })} />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

		<Modal title={intl.formatMessage({ id: 'users.modal.editTitle' })} open={editModalVisible} onCancel={() => { setEditModalVisible(false); setEditingUser(null); }} footer={null} width={500}>
			<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
				<Form.Item name="email" label={intl.formatMessage({ id: 'user.email' })} rules={[{ type: "email", message: intl.formatMessage({ id: 'users.validation.validEmail' }) }]}>
					<Input placeholder={intl.formatMessage({ id: 'users.placeholder.email' })} />
					</Form.Item>
					<Form.Item name="status" label={intl.formatMessage({ id: 'user.status' })}>
						<Select>{userStatuses.map((s) => <Option key={s.value} value={s.value}>{s.label}</Option>)}</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.save' })}</Button>
							<Button onClick={() => setEditModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

		<Modal title={intl.formatMessage({ id: 'users.modal.assignRole' }, { username: selectedUser?.username })} open={roleModalVisible} onCancel={() => setRoleModalVisible(false)} footer={null} width={400}>
			<Form form={roleForm} layout="vertical" onFinish={handleRoleSubmit}>
				<Form.Item name="role_id" label={intl.formatMessage({ id: 'users.form.selectRole' })} rules={[{ required: true, message: intl.formatMessage({ id: 'users.validation.selectRole' }) }]}>
					<Select placeholder={intl.formatMessage({ id: 'users.placeholder.selectRole' })}>
							{(roles || []).map((role: any) => <Option key={role.id} value={role.id}>{role.name}</Option>)}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'users.form.assign' })}</Button>
							<Button onClick={() => setRoleModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
