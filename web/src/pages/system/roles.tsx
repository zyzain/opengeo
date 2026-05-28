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

const { Option } = Select;
const { TextArea } = Input;

// 权限定义（硬编码回退数据）
const fallbackPermissionGroups = [
	{
		title: "内容管理",
		permissions: [
			{ id: "content:create", label: "创建内容" },
			{ id: "content:read", label: "查看内容" },
			{ id: "content:update", label: "编辑内容" },
			{ id: "content:delete", label: "删除内容" },
			{ id: "content:publish", label: "发布内容" },
			{ id: "content:optimize", label: "AI优化内容" },
		],
	},
	{
		title: "账号管理",
		permissions: [
			{ id: "account:create", label: "创建账号" },
			{ id: "account:read", label: "查看账号" },
			{ id: "account:update", label: "编辑账号" },
			{ id: "account:delete", label: "删除账号" },
		],
	},
	{
		title: "发布管理",
		permissions: [
			{ id: "publish:create", label: "创建发布任务" },
			{ id: "publish:read", label: "查看发布任务" },
			{ id: "publish:cancel", label: "取消发布任务" },
			{ id: "publish:retry", label: "重试发布任务" },
		],
	},
	{
		title: "监测管理",
		permissions: [
			{ id: "monitor:read", label: "查看监测数据" },
			{ id: "monitor:configure", label: "配置监测规则" },
		],
	},
	{
		title: "系统管理",
		permissions: [
			{ id: "system:config", label: "系统配置" },
			{ id: "system:user", label: "用户管理" },
			{ id: "system:role", label: "角色管理" },
			{ id: "system:plugin", label: "插件管理" },
		],
	},
];

export default function RoleManagementPage() {
	const queryClient = useQueryClient();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [permModalVisible, setPermModalVisible] = useState(false);
	const [editingRole, setEditingRole] = useState<any>(null);
	const [selectedRole, setSelectedRole] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();

	// 获取权限定义
	const { data: permDefsData } = useQuery({
		queryKey: ["permissionDefinitions"],
		queryFn: () => api.system.getPermissionDefinitions(),
	});

	const permissionGroups =
		permDefsData?.data?.data?.groups || fallbackPermissionGroups;

	// 查询角色列表
	const { data: rolesData, isLoading } = useQuery({
		queryKey: ["roles"],
		queryFn: () => api.roles.list(),
		select: (response) => response.data.data,
	});

	const roles = rolesData?.items || [];

	// 表格列定义
	const columns = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
			width: 80,
		},
		{
			title: "角色名称",
			dataIndex: "name",
			key: "name",
			render: (text: string, record: any) => (
				<Space>
					<SafetyOutlined />
					<span className="font-medium">{text}</span>
					{record.is_system && <Tag color="blue">系统角色</Tag>}
				</Space>
			),
		},
		{
			title: "描述",
			dataIndex: "description",
			key: "description",
			ellipsis: true,
		},
		{
			title: "用户数",
			dataIndex: "user_count",
			key: "user_count",
			width: 100,
			render: (count: number) => (
				<Badge
					count={count || 0}
					showZero
					style={{ backgroundColor: "#1890ff" }}
				/>
			),
		},
		{
			title: "创建时间",
			dataIndex: "created_at",
			key: "created_at",
			width: 180,
			render: (text: string) => new Date(text).toLocaleString(),
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
					<Tooltip title="权限配置">
						<Button
							type="text"
							icon={<KeyOutlined />}
							onClick={() => handlePermissions(record)}
						/>
					</Tooltip>
					{!record.is_system && (
						<Popconfirm
							title="确定要删除这个角色吗？"
							onConfirm={() => handleDelete(record.id)}
							okText="确定"
							cancelText="取消"
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

	// 创建角色
	const handleCreate = async (values: any) => {
		try {
			await api.roles.create(values);
			message.success("创建成功");
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: ["roles"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "创建失败");
		}
	};

	// 编辑角色
	const handleEdit = (record: any) => {
		setEditingRole(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await api.roles.update(editingRole.id, values);
			message.success("更新成功");
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingRole(null);
			queryClient.invalidateQueries({ queryKey: ["roles"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "更新失败");
		}
	};

	// 删除角色
	const handleDelete = async (id: number) => {
		try {
			await api.roles.delete(id);
			message.success("删除成功");
			queryClient.invalidateQueries({ queryKey: ["roles"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "删除失败");
		}
	};

	// 权限配置
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
				...added.map((perm: string) =>
					api.roles.addPermission({
						role_id: selectedRole.id,
						permission: perm,
					}),
				),
				...removed.map((perm: string) =>
					api.roles.removePermission({
						role_id: selectedRole.id,
						permission: perm,
					}),
				),
			]);
			message.success("权限更新成功");
			queryClient.invalidateQueries({
				queryKey: ["rolePermissions", selectedRole.id],
			});
		} catch (error: any) {
			message.error("权限更新失败");
		}
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">角色管理</h1>
				<p className="text-gray-500 mt-1">管理系统角色和权限配置</p>
			</div>

			{/* 预设角色说明 */}
			<Card className="mb-4">
				<Row gutter={[16, 16]}>
					<Col xs={24} sm={8}>
						<div className="p-3 bg-red-50 rounded-lg">
							<div className="font-medium text-red-700">管理员</div>
							<div className="text-sm text-gray-500">拥有所有系统权限</div>
						</div>
					</Col>
					<Col xs={24} sm={8}>
						<div className="p-3 bg-blue-50 rounded-lg">
							<div className="font-medium text-blue-700">运营人员</div>
							<div className="text-sm text-gray-500">内容和发布管理权限</div>
						</div>
					</Col>
					<Col xs={24} sm={8}>
						<div className="p-3 bg-green-50 rounded-lg">
							<div className="font-medium text-green-700">查看者</div>
							<div className="text-sm text-gray-500">
								仅查看权限，无修改权限
							</div>
						</div>
					</Col>
				</Row>
			</Card>

			{/* 角色列表 */}
			<Card
				title={
					<Space>
						<SafetyOutlined />
						<span>角色列表</span>
					</Space>
				}
				extra={
					<Space>
						<Button
							icon={<ReloadOutlined />}
							onClick={() =>
								queryClient.invalidateQueries({ queryKey: ["roles"] })
							}
						>
							刷新
						</Button>
						<Button
							type="primary"
							icon={<PlusOutlined />}
							onClick={() => setCreateModalVisible(true)}
						>
							创建角色
						</Button>
					</Space>
				}
			>
				<Table
					columns={columns}
					dataSource={roles}
					rowKey="id"
					loading={isLoading}
					pagination={false}
				/>
			</Card>

			{/* 创建角色弹窗 */}
			<Modal
				title="创建角色"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="name"
						label="角色名称"
						rules={[{ required: true, message: "请输入角色名称" }]}
					>
						<Input placeholder="请输入角色名称" />
					</Form.Item>
					<Form.Item name="description" label="描述">
						<TextArea rows={3} placeholder="请输入角色描述" />
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

			{/* 编辑角色弹窗 */}
			<Modal
				title="编辑角色"
				open={editModalVisible}
				onCancel={() => {
					setEditModalVisible(false);
					setEditingRole(null);
				}}
				footer={null}
				width={500}
			>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item
						name="name"
						label="角色名称"
						rules={[{ required: true, message: "请输入角色名称" }]}
					>
						<Input placeholder="请输入角色名称" />
					</Form.Item>
					<Form.Item name="description" label="描述">
						<TextArea rows={3} placeholder="请输入角色描述" />
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

			{/* 权限配置弹窗 */}
			<Modal
				title={`权限配置 - ${selectedRole?.name}`}
				open={permModalVisible}
				onCancel={() => setPermModalVisible(false)}
				footer={null}
				width={600}
			>
				<div className="space-y-4">
					{permissionGroups.map((group: any) => (
						<div key={group.title}>
							<Divider orientation="left">{group.title}</Divider>
							<Checkbox.Group
								value={group.permissions
									.filter((p: any) => currentPermissions.includes(p.id))
									.map((p: any) => p.id)}
								onChange={handlePermissionChange}
							>
								<div className="grid grid-cols-2 gap-2">
									{group.permissions.map((perm: any) => (
										<Checkbox key={perm.id} value={perm.id}>
											{perm.label}
										</Checkbox>
									))}
								</div>
							</Checkbox.Group>
						</div>
					))}
				</div>
			</Modal>
		</div>
	);
}
