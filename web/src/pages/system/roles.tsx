"use client";

import api from "@/lib/api";
import {
	CheckCircleOutlined,
	DeleteOutlined,
	EditOutlined,
	KeyOutlined,
	PlusOutlined,
	ReloadOutlined,
	SafetyOutlined,
} from "@ant-design/icons";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
	Badge,
	Button,
	Card,
	Checkbox,
	Col,
	Divider,
	Form,
	Input,
	Modal,
	Popconfirm,
	Row,
	Select,
	Space,
	Statistic,
	Table,
	Tag,
	Tooltip,
	Tree,
	message,
} from "antd";
import { useState } from "react";
import { useIntl } from "react-intl";

const { Option } = Select;
const { TextArea } = Input;

export default function RoleManagementPage() {
	const intl = useIntl();
	const queryClient = useQueryClient();

	const fallbackPermissionGroups = [
		{ title: intl.formatMessage({ id: 'roles.permission.contentManagement' }), permissions: [
			{ id: "content:create", label: intl.formatMessage({ id: 'roles.permission.createContent' }) },
			{ id: "content:read", label: intl.formatMessage({ id: 'roles.permission.viewContent' }) },
			{ id: "content:update", label: intl.formatMessage({ id: 'roles.permission.editContent' }) },
			{ id: "content:delete", label: intl.formatMessage({ id: 'roles.permission.deleteContent' }) },
			{ id: "content:publish", label: intl.formatMessage({ id: 'roles.permission.publishContent' }) },
			{ id: "content:optimize", label: intl.formatMessage({ id: 'roles.permission.aiOptimizeContent' }) },
		]},
		{ title: intl.formatMessage({ id: 'roles.permission.accountManagement' }), permissions: [
			{ id: "account:create", label: intl.formatMessage({ id: 'roles.permission.createAccount' }) },
			{ id: "account:read", label: intl.formatMessage({ id: 'roles.permission.viewAccount' }) },
			{ id: "account:update", label: intl.formatMessage({ id: 'roles.permission.editAccount' }) },
			{ id: "account:delete", label: intl.formatMessage({ id: 'roles.permission.deleteAccount' }) },
		]},
		{ title: intl.formatMessage({ id: 'roles.permission.publishManagement' }), permissions: [
			{ id: "publish:create", label: intl.formatMessage({ id: 'roles.permission.createPublishTask' }) },
			{ id: "publish:read", label: intl.formatMessage({ id: 'roles.permission.viewPublishTask' }) },
			{ id: "publish:cancel", label: intl.formatMessage({ id: 'roles.permission.cancelPublishTask' }) },
			{ id: "publish:retry", label: intl.formatMessage({ id: 'roles.permission.retryPublishTask' }) },
		]},
		{ title: intl.formatMessage({ id: 'roles.permission.monitorManagement' }), permissions: [
			{ id: "monitor:read", label: intl.formatMessage({ id: 'roles.permission.viewMonitorData' }) },
			{ id: "monitor:configure", label: intl.formatMessage({ id: 'roles.permission.configMonitorRules' }) },
		]},
		{ title: intl.formatMessage({ id: 'roles.permission.systemManagement' }), permissions: [
			{ id: "system:config", label: intl.formatMessage({ id: 'roles.permission.systemConfig' }) },
			{ id: "system:user", label: intl.formatMessage({ id: 'roles.permission.userManagement' }) },
			{ id: "system:role", label: intl.formatMessage({ id: 'roles.permission.roleManagement' }) },
			{ id: "system:plugin", label: intl.formatMessage({ id: 'roles.permission.pluginManagement' }) },
		]},
	];
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [permModalVisible, setPermModalVisible] = useState(false);
	const [editingRole, setEditingRole] = useState<any>(null);
	const [selectedRole, setSelectedRole] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();

	const { data: permDefsData } = useQuery({
		queryKey: ["permissionDefinitions"],
		queryFn: () => api.system.getPermissionDefinitions(),
	});

	const permissionGroups = permDefsData?.data?.data?.groups || fallbackPermissionGroups;

	const { data: rolesData, isLoading } = useQuery({
		queryKey: ["roles"],
		queryFn: () => api.roles.list(),
		select: (response) => response.data.data,
	});

	const roles = rolesData?.items || [];

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{
			title: intl.formatMessage({ id: 'roles.column.name' }),
			dataIndex: "name",
			key: "name",
			render: (text: string, record: any) => (
				<Space>
					<SafetyOutlined />
					<span className="font-medium">{text}</span>
					{record.is_system && <Tag color="blue">{intl.formatMessage({ id: 'roles.tag.systemRole' })}</Tag>}
				</Space>
			),
		},
		{ title: intl.formatMessage({ id: 'role.description' }), dataIndex: "description", key: "description", ellipsis: true },
		{ title: intl.formatMessage({ id: 'roles.column.userCount' }), dataIndex: "user_count", key: "user_count", width: 100, render: (count: number) => <Badge count={count || 0} showZero style={{ backgroundColor: "#1890ff" }} /> },
		{ title: intl.formatMessage({ id: 'roles.column.createdAt' }), dataIndex: "created_at", key: "created_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => handleEdit(record)} /></Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'roles.action.permissions' })}><Button type="text" icon={<KeyOutlined />} onClick={() => handlePermissions(record)} /></Tooltip>
					{!record.is_system && (
						<Popconfirm title={intl.formatMessage({ id: 'roles.confirmDelete' })} onConfirm={() => handleDelete(record.id)} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })}>
							<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
						</Popconfirm>
					)}
				</Space>
			),
		},
	];

	const handleCreate = async (values: any) => {
		try {
			await api.roles.create(values);
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: ["roles"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.createFailed' }));
		}
	};

	const handleEdit = (record: any) => {
		setEditingRole(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await api.roles.update(editingRole.id, values);
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingRole(null);
			queryClient.invalidateQueries({ queryKey: ["roles"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.roles.delete(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
			queryClient.invalidateQueries({ queryKey: ["roles"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const handlePermissions = (record: any) => {
		setSelectedRole(record);
		setPermModalVisible(true);
	};

	const { data: permissionsData } = useQuery({
		queryKey: ["rolePermissions", selectedRole?.id],
		queryFn: () => api.roles.getPermissions(selectedRole.id),
		select: (response) => response.data.data,
		enabled: !!selectedRole?.id && permModalVisible,
	});

	const currentPermissions = permissionsData?.permissions || [];

	const handlePermissionChange = async (checkedValues: string[]) => {
		if (!selectedRole) return;
		const oldPerms = currentPermissions;
		const newPerms = checkedValues;
		const added = newPerms.filter((p: string) => !oldPerms.includes(p));
		const removed = oldPerms.filter((p: string) => !newPerms.includes(p));
		try {
			await Promise.all([
				...added.map((perm: string) => api.roles.addPermission({ role_id: selectedRole.id, permission: perm })),
				...removed.map((perm: string) => api.roles.removePermission({ role_id: selectedRole.id, permission: perm })),
			]);
			message.success(intl.formatMessage({ id: 'roles.message.permissionUpdateSuccess' }));
			queryClient.invalidateQueries({ queryKey: ["rolePermissions", selectedRole.id] });
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'roles.message.permissionUpdateFailed' }));
		}
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'roles.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'roles.page.subtitle' })}</p>
			</div>

			<Card className="mb-4">
				<Row gutter={[16, 16]}>
					<Col xs={24} sm={8}>
						<div className="p-3 bg-red-50 rounded-lg">
							<div className="font-medium text-red-700">{intl.formatMessage({ id: 'roles.preset.admin' })}</div>
							<div className="text-sm text-gray-500">{intl.formatMessage({ id: 'roles.preset.adminDesc' })}</div>
						</div>
					</Col>
					<Col xs={24} sm={8}>
						<div className="p-3 bg-blue-50 rounded-lg">
							<div className="font-medium text-blue-700">{intl.formatMessage({ id: 'roles.preset.operator' })}</div>
							<div className="text-sm text-gray-500">{intl.formatMessage({ id: 'roles.preset.operatorDesc' })}</div>
						</div>
					</Col>
					<Col xs={24} sm={8}>
						<div className="p-3 bg-green-50 rounded-lg">
							<div className="font-medium text-green-700">{intl.formatMessage({ id: 'roles.preset.viewer' })}</div>
							<div className="text-sm text-gray-500">{intl.formatMessage({ id: 'roles.preset.viewerDesc' })}</div>
						</div>
					</Col>
				</Row>
			</Card>

			<Card
				title={<Space><SafetyOutlined /><span>{intl.formatMessage({ id: 'roles.section.list' })}</span></Space>}
				extra={
					<Space>
						<Button icon={<ReloadOutlined />} onClick={() => queryClient.invalidateQueries({ queryKey: ["roles"] })}>{intl.formatMessage({ id: 'common.action.refresh' })}</Button>
						<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'roles.action.create' })}</Button>
					</Space>
				}
			>
				<Table columns={columns} dataSource={roles} rowKey="id" loading={isLoading} pagination={false} />
			</Card>

			<Modal title={intl.formatMessage({ id: 'roles.modal.create' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={500}>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item name="name" label={intl.formatMessage({ id: 'roles.form.name' })} rules={[{ required: true, message: intl.formatMessage({ id: 'roles.validation.enterName' }) }]}>
						<Input placeholder={intl.formatMessage({ id: 'roles.placeholder.name' })} />
					</Form.Item>
					<Form.Item name="description" label={intl.formatMessage({ id: 'roles.form.description' })}>
						<TextArea rows={3} placeholder={intl.formatMessage({ id: 'roles.placeholder.description' })} />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title={intl.formatMessage({ id: 'roles.modal.edit' })} open={editModalVisible} onCancel={() => { setEditModalVisible(false); setEditingRole(null); }} footer={null} width={500}>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item name="name" label={intl.formatMessage({ id: 'roles.form.name' })} rules={[{ required: true, message: intl.formatMessage({ id: 'roles.validation.enterName' }) }]}>
						<Input placeholder={intl.formatMessage({ id: 'roles.placeholder.name' })} />
					</Form.Item>
					<Form.Item name="description" label={intl.formatMessage({ id: 'roles.form.description' })}>
						<TextArea rows={3} placeholder={intl.formatMessage({ id: 'roles.placeholder.description' })} />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.save' })}</Button>
							<Button onClick={() => setEditModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title={`${intl.formatMessage({ id: 'roles.modal.permissionsTitle' })} ${selectedRole?.name}`} open={permModalVisible} onCancel={() => setPermModalVisible(false)} footer={null} width={600}>
				<div className="space-y-4">
					{permissionGroups.map((group: any) => (
						<div key={group.title}>
							<Divider orientation="left">{group.title}</Divider>
							<Checkbox.Group
								value={group.permissions.filter((p: any) => currentPermissions.includes(p.id)).map((p: any) => p.id)}
								onChange={handlePermissionChange}
							>
								<div className="grid grid-cols-2 gap-2">
									{group.permissions.map((perm: any) => <Checkbox key={perm.id} value={perm.id}>{perm.label}</Checkbox>)}
								</div>
							</Checkbox.Group>
						</div>
					))}
				</div>
			</Modal>
		</div>
	);
}
