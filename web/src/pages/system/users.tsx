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

const { Option } = Select;

export default function UserManagementPage() {
	const queryClient = useQueryClient();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [roleModalVisible, setRoleModalVisible] = useState(false);
	const [editingUser, setEditingUser] = useState<any>(null);
	const [selectedUser, setSelectedUser] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();
	const [roleForm] = Form.useForm();

	// 查询用户列表
	const { data: usersData, isLoading } = useQuery({
		queryKey: ["users"],
		queryFn: () => api.users.list(),
		select: (response) => response.data.data,
	});

	// 查询角色列表
	const { data: rolesData } = useQuery({
		queryKey: ["roles"],
		queryFn: () => api.roles.list(),
		select: (response) => response.data.data,
	});

	const users = usersData?.items || [];
	const roles = rolesData?.items || [];

	// 用户状态
	const userStatuses = [
		{ value: 1, label: "正常", color: "success" },
		{ value: 0, label: "禁用", color: "error" },
		{ value: 2, label: "待审核", color: "warning" },
	];

	// 获取状态标签
	const getStatusTag = (status: number) => {
		const statusInfo = userStatuses.find((s) => s.value === status);
		return (
			<Badge
				status={statusInfo?.color as any}
				text={statusInfo?.label || "未知"}
			/>
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
			title: "用户名",
			dataIndex: "username",
			key: "username",
			render: (text: string) => (
				<Space>
					<UserOutlined />
					<span className="font-medium">{text}</span>
				</Space>
			),
		},
		{
			title: "邮箱",
			dataIndex: "email",
			key: "email",
			ellipsis: true,
		},
		{
			title: "状态",
			dataIndex: "status",
			key: "status",
			width: 100,
			render: (status: number) => getStatusTag(status),
		},
		{
			title: "创建时间",
			dataIndex: "created_at",
			key: "created_at",
			width: 180,
			render: (text: string) => new Date(text).toLocaleString(),
		},
		{
			title: "最后登录",
			dataIndex: "last_login_at",
			key: "last_login_at",
			width: 180,
			render: (text: string) => (text ? new Date(text).toLocaleString() : "-"),
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
					<Tooltip title="分配角色">
						<Button
							type="text"
							icon={<KeyOutlined />}
							onClick={() => handleAssignRole(record)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个用户吗？"
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

	// 创建用户
	const handleCreate = async (values: any) => {
		try {
			await api.auth.register(values);
			message.success("创建成功");
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: ["users"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "创建失败");
		}
	};

	// 编辑用户
	const handleEdit = (record: any) => {
		setEditingUser(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await api.users.update(editingUser.id, values);
			message.success("更新成功");
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingUser(null);
			queryClient.invalidateQueries({ queryKey: ["users"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "更新失败");
		}
	};

	// 删除用户
	const handleDelete = async (id: number) => {
		try {
			await api.users.delete(id);
			message.success("删除成功");
			queryClient.invalidateQueries({ queryKey: ["users"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "删除失败");
		}
	};

	// 分配角色
	const handleAssignRole = (record: any) => {
		setSelectedUser(record);
		roleForm.resetFields();
		setRoleModalVisible(true);
	};

	const handleRoleSubmit = async (values: any) => {
		try {
			await api.roles.assign({
				user_id: selectedUser.id,
				role_id: values.role_id,
			});
			message.success("角色分配成功");
			setRoleModalVisible(false);
			roleForm.resetFields();
		} catch (error: any) {
			message.error(error.response?.data?.message || "角色分配失败");
		}
	};

	// 统计数据
	const stats = {
		total: users.length,
		active: users.filter((u: any) => u.status === 1).length,
		disabled: users.filter((u: any) => u.status === 0).length,
		pending: users.filter((u: any) => u.status === 2).length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">用户管理</h1>
				<p className="text-gray-500 mt-1">管理系统用户账号和权限</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="总用户数"
							value={stats.total}
							prefix={<TeamOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="正常用户"
							value={stats.active}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="禁用用户"
							value={stats.disabled}
							valueStyle={{ color: "#ff4d4f" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="待审核"
							value={stats.pending}
							valueStyle={{ color: "#faad14" }}
						/>
					</Card>
				</Col>
			</Row>

			{/* 用户列表 */}
			<Card
				title={
					<Space>
						<UserOutlined />
						<span>用户列表</span>
					</Space>
				}
				extra={
					<Space>
						<Button
							icon={<ReloadOutlined />}
							onClick={() =>
								queryClient.invalidateQueries({ queryKey: ["users"] })
							}
						>
							刷新
						</Button>
						<Button
							type="primary"
							icon={<PlusOutlined />}
							onClick={() => setCreateModalVisible(true)}
						>
							添加用户
						</Button>
					</Space>
				}
			>
				<Table
					columns={columns}
					dataSource={users}
					rowKey="id"
					loading={isLoading}
					pagination={{
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
					}}
				/>
			</Card>

			{/* 创建用户弹窗 */}
			<Modal
				title="添加用户"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="username"
						label="用户名"
						rules={[
							{ required: true, message: "请输入用户名" },
							{ min: 3, message: "至少3个字符" },
						]}
					>
						<Input placeholder="请输入用户名" />
					</Form.Item>
					<Form.Item
						name="email"
						label="邮箱"
						rules={[
							{ required: true, message: "请输入邮箱" },
							{ type: "email", message: "请输入有效邮箱" },
						]}
					>
						<Input placeholder="请输入邮箱" />
					</Form.Item>
					<Form.Item
						name="password"
						label="密码"
						rules={[
							{ required: true, message: "请输入密码" },
							{ min: 8, message: "密码至少8个字符" },
						]}
					>
						<Input.Password placeholder="请输入密码" />
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

			{/* 编辑用户弹窗 */}
			<Modal
				title="编辑用户"
				open={editModalVisible}
				onCancel={() => {
					setEditModalVisible(false);
					setEditingUser(null);
				}}
				footer={null}
				width={500}
			>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item
						name="email"
						label="邮箱"
						rules={[{ type: "email", message: "请输入有效邮箱" }]}
					>
						<Input placeholder="请输入邮箱" />
					</Form.Item>
					<Form.Item name="status" label="状态">
						<Select>
							{userStatuses.map((s) => (
								<Option key={s.value} value={s.value}>
									{s.label}
								</Option>
							))}
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

			{/* 分配角色弹窗 */}
			<Modal
				title={`分配角色 - ${selectedUser?.username}`}
				open={roleModalVisible}
				onCancel={() => setRoleModalVisible(false)}
				footer={null}
				width={400}
			>
				<Form form={roleForm} layout="vertical" onFinish={handleRoleSubmit}>
					<Form.Item
						name="role_id"
						label="选择角色"
						rules={[{ required: true, message: "请选择角色" }]}
					>
						<Select placeholder="请选择角色">
							{(roles || []).map((role: any) => (
								<Option key={role.id} value={role.id}>
									{role.name}
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">
								分配
							</Button>
							<Button onClick={() => setRoleModalVisible(false)}>取消</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
