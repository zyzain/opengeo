"use client";

import { useAccountGroups, useCreateAccountGroup } from "@/hooks";
import api from "@/lib/api";
import {
	ApartmentOutlined,
	DeleteOutlined,
	EditOutlined,
	FolderOpenOutlined,
	FolderOutlined,
	PlusOutlined,
	TeamOutlined,
} from "@ant-design/icons";
import {
	Button,
	Card,
	Col,
	Empty,
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
import { useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

const { Option } = Select;
const { TextArea } = Input;

export default function AccountGroupsPage() {
	const queryClient = useQueryClient();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [editingGroup, setEditingGroup] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();

	const { data: groupsData, isLoading } = useAccountGroups();
	const createMutation = useCreateAccountGroup();

	const groups = groupsData?.items || [];

	// 分组类型
	const groupTypes = [
		{
			value: "authority",
			label: "权威背书层",
			color: "red",
			description: "官方账号、权威媒体",
		},
		{
			value: "professional",
			label: "专业认证层",
			color: "blue",
			description: "专业账号、KOL",
		},
		{
			value: "ecology",
			label: "生态渗透层",
			color: "green",
			description: "长尾账号、矩阵号",
		},
	];

	// 获取分组类型标签
	const getGroupTypeTag = (type: string) => {
		const typeInfo = groupTypes.find((t) => t.value === type);
		return (
			<Tag color={typeInfo?.color || "default"}>{typeInfo?.label || type}</Tag>
		);
	};

	// 转换为树形结构
	const convertToTreeData = (groups: any[]): any[] => {
		return groups.map((group) => ({
			key: group.id,
			title: (
				<div className="flex items-center space-x-2">
					<span>{group.name}</span>
					{getGroupTypeTag(group.group_type)}
					<span className="text-gray-400 text-xs">
						({group.account_count || 0}个账号)
					</span>
				</div>
			),
			icon: group.children?.length ? (
				<FolderOpenOutlined />
			) : (
				<FolderOutlined />
			),
			children: group.children ? convertToTreeData(group.children) : [],
		}));
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
			title: "分组名称",
			dataIndex: "name",
			key: "name",
			render: (text: string, record: any) => (
				<Space>
					<ApartmentOutlined />
					<span className="font-medium">{text}</span>
				</Space>
			),
		},
		{
			title: "分组类型",
			dataIndex: "group_type",
			key: "group_type",
			width: 150,
			render: (type: string) => getGroupTypeTag(type),
		},
		{
			title: "描述",
			dataIndex: "description",
			key: "description",
			ellipsis: true,
		},
		{
			title: "账号数量",
			dataIndex: "account_count",
			key: "account_count",
			width: 100,
			render: (count: number) => <Tag color="blue">{count || 0}</Tag>,
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
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="编辑">
						<Button
							type="text"
							icon={<EditOutlined />}
							onClick={() => handleEdit(record)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个分组吗？"
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

	// 创建分组
	const handleCreate = async (values: any) => {
		try {
			await createMutation.mutateAsync(values);
			message.success("创建成功");
			setCreateModalVisible(false);
			createForm.resetFields();
		} catch (error: any) {
			message.error(error.response?.data?.message || "创建失败");
		}
	};

	// 编辑分组
	const handleEdit = (record: any) => {
		setEditingGroup(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await api.accountGroups.update(editingGroup.id, values);
			queryClient.invalidateQueries({ queryKey: ['accountGroups'] });
			message.success("更新成功");
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingGroup(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || "更新失败");
		}
	};

	// 删除分组
	const handleDelete = async (id: number) => {
		try {
			await api.accountGroups.delete(id);
			queryClient.invalidateQueries({ queryKey: ['accountGroups'] });
			message.success("删除成功");
		} catch (error: any) {
			message.error(error.response?.data?.message || "删除失败");
		}
	};

	// 统计数据
	const stats: Record<string, number> = {
		total: groups.length,
		authority: groups.filter((g: any) => g.group_type === "authority").length,
		professional: groups.filter((g: any) => g.group_type === "professional")
			.length,
		ecology: groups.filter((g: any) => g.group_type === "ecology").length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">账号分组</h1>
				<p className="text-gray-500 mt-1">按三层架构管理账号分组</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="总分组数"
							value={stats.total}
							prefix={<ApartmentOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="权威背书层"
							value={stats.authority}
							valueStyle={{ color: "#ff4d4f" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="专业认证层"
							value={stats.professional}
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="生态渗透层"
							value={stats.ecology}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
			</Row>

			<Row gutter={[16, 16]}>
				{/* 分组树 */}
				<Col xs={24} lg={8}>
					<Card
						title="分组结构"
						extra={
							<Button
								type="primary"
								size="small"
								icon={<PlusOutlined />}
								onClick={() => setCreateModalVisible(true)}
							>
								新建分组
							</Button>
						}
					>
						{groups.length > 0 ? (
							<Tree
								showIcon
								defaultExpandAll
								treeData={convertToTreeData(groups)}
								className="bg-gray-50 p-4 rounded-lg"
							/>
						) : (
							<Empty description="暂无分组" />
						)}
					</Card>
				</Col>

				{/* 分组列表 */}
				<Col xs={24} lg={16}>
					<Card title="分组列表">
						<Table
							columns={columns}
							dataSource={groups}
							rowKey="id"
							loading={isLoading}
							pagination={false}
						/>
					</Card>
				</Col>
			</Row>

			{/* 创建分组弹窗 */}
			<Modal
				title="新建分组"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="name"
						label="分组名称"
						rules={[{ required: true, message: "请输入分组名称" }]}
					>
						<Input placeholder="请输入分组名称" />
					</Form.Item>
					<Form.Item
						name="group_type"
						label="分组类型"
						rules={[{ required: true, message: "请选择分组类型" }]}
					>
						<Select placeholder="请选择分组类型">
							{groupTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<Space>
										<Tag color={type.color}>{type.label}</Tag>
										<span className="text-gray-400">{type.description}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="parent_id" label="父分组">
						<Select placeholder="请选择父分组（可选）" allowClear>
							{groups.map((group: any) => (
								<Option key={group.id} value={group.id}>
									{group.name}
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="description" label="描述">
						<TextArea rows={3} placeholder="请输入分组描述" />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								htmlType="submit"
								loading={createMutation.isPending}
							>
								创建
							</Button>
							<Button onClick={() => setCreateModalVisible(false)}>取消</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			{/* 编辑分组弹窗 */}
			<Modal
				title="编辑分组"
				open={editModalVisible}
				onCancel={() => {
					setEditModalVisible(false);
					setEditingGroup(null);
				}}
				footer={null}
				width={500}
			>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item
						name="name"
						label="分组名称"
						rules={[{ required: true, message: "请输入分组名称" }]}
					>
						<Input placeholder="请输入分组名称" />
					</Form.Item>
					<Form.Item
						name="group_type"
						label="分组类型"
						rules={[{ required: true, message: "请选择分组类型" }]}
					>
						<Select placeholder="请选择分组类型">
							{groupTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<Space>
										<Tag color={type.color}>{type.label}</Tag>
										<span className="text-gray-400">{type.description}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="description" label="描述">
						<TextArea rows={3} placeholder="请输入分组描述" />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">
								保存
							</Button>
							<Button
								onClick={() => {
									setEditModalVisible(false);
									setEditingGroup(null);
								}}
							>
								取消
							</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
