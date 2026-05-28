"use client";

import {
	useAccounts,
	useCreateAccount,
	useDeleteAccount,
	useUpdateAccount,
} from "@/hooks";
import api from "@/lib/api";
import {
	CheckCircleOutlined,
	DeleteOutlined,
	EditOutlined,
	PlusOutlined,
	ReloadOutlined,
	SearchOutlined,
	TeamOutlined,
} from "@ant-design/icons";
import {
	Badge,
	Button,
	Card,
	Form,
	Input,
	Modal,
	Popconfirm,
	Progress,
	Select,
	Space,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";

const { Option } = Select;

export default function AccountListPage() {
	const [searchForm] = Form.useForm();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [editingAccount, setEditingAccount] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();

	// 查询参数
	const [queryParams, setQueryParams] = useState({
		page: 1,
		page_size: 10,
		platform: undefined,
	});

	const { data, isLoading, refetch } = useAccounts(queryParams);
	const createMutation = useCreateAccount();
	const updateMutation = useUpdateAccount();
	const deleteMutation = useDeleteAccount();

	const accounts = data?.items || [];
	const total = data?.total || 0;

	// 平台类型
	const platformTypes = [
		{ value: "wechat", label: "微信公众号", color: "green" },
		{ value: "weibo", label: "微博", color: "red" },
		{ value: "douyin", label: "抖音", color: "purple" },
		{ value: "xiaohongshu", label: "小红书", color: "pink" },
		{ value: "zhihu", label: "知乎", color: "blue" },
		{ value: "toutiao", label: "今日头条", color: "orange" },
	];

	// 获取平台标签
	const getPlatformTag = (platform: string) => {
		const platformInfo = platformTypes.find((p) => p.value === platform);
		return (
			<Tag color={platformInfo?.color || "default"}>
				{platformInfo?.label || platform}
			</Tag>
		);
	};

	// 状态标签
	const getStatusTag = (status: number) => {
		const statusMap: Record<number, { color: string; text: string }> = {
			0: { color: "error", text: "禁用" },
			1: { color: "success", text: "正常" },
			2: { color: "warning", text: "异常" },
		};
		const config = statusMap[status] || { color: "default", text: "未知" };
		return <Badge status={config.color as any} text={config.text} />;
	};

	// 健康度颜色
	const getHealthColor = (score: number) => {
		if (score >= 80) return "#52c41a";
		if (score >= 60) return "#faad14";
		return "#ff4d4f";
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
			title: "账号名称",
			dataIndex: "account_name",
			key: "account_name",
			ellipsis: true,
		},
		{
			title: "平台",
			dataIndex: "platform",
			key: "platform",
			width: 120,
			render: (platform: string) => getPlatformTag(platform),
		},
		{
			title: "账号ID",
			dataIndex: "account_id",
			key: "account_id",
			width: 150,
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
			title: "健康度",
			dataIndex: "health_score",
			key: "health_score",
			width: 150,
			render: (score: number) => (
				<Progress
					percent={score}
					size="small"
					strokeColor={getHealthColor(score)}
					format={(percent) => `${percent}%`}
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
					<Tooltip title="健康检查">
						<Button
							type="text"
							icon={<CheckCircleOutlined />}
							onClick={() => handleHealthCheck(record.id)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个账号吗？"
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

	// 创建账号
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

	// 编辑账号
	const handleEdit = (record: any) => {
		setEditingAccount(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await updateMutation.mutateAsync({ id: editingAccount.id, data: values });
			message.success("更新成功");
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingAccount(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || "更新失败");
		}
	};

	// 删除账号
	const handleDelete = async (id: number) => {
		try {
			await deleteMutation.mutateAsync(id);
			message.success("删除成功");
		} catch (error: any) {
			message.error(error.response?.data?.message || "删除失败");
		}
	};

	// 健康检查
	const handleHealthCheck = async (id: number) => {
		try {
			const res = await api.accounts.health(id);
			const healthData = res.data.data;
			message.success(
				`健康检查完成 - 评分: ${healthData.health_score}, 状态: ${healthData.status}`,
			);
			refetch();
		} catch (error: any) {
			message.error("健康检查失败");
		}
	};

	// 搜索
	const handleSearch = (values: any) => {
		setQueryParams({ ...queryParams, ...values, page: 1 });
	};

	// 重置搜索
	const handleReset = () => {
		searchForm.resetFields();
		setQueryParams({ page: 1, page_size: 10, platform: undefined });
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">账号管理</h1>
				<p className="text-gray-500 mt-1">管理您的多平台账号</p>
			</div>

			{/* 搜索表单 */}
			<Card className="mb-4">
				<Form form={searchForm} layout="inline" onFinish={handleSearch}>
					<Form.Item name="platform" label="平台">
						<Select placeholder="请选择平台" allowClear style={{ width: 150 }}>
							{platformTypes.map((p) => (
								<Option key={p.value} value={p.value}>
									{p.label}
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								icon={<SearchOutlined />}
								htmlType="submit"
							>
								搜索
							</Button>
							<Button onClick={handleReset}>重置</Button>
							<Button icon={<ReloadOutlined />} onClick={() => refetch()}>
								刷新
							</Button>
						</Space>
					</Form.Item>
				</Form>
			</Card>

			{/* 账号列表 */}
			<Card
				title={
					<Space>
						<TeamOutlined />
						<span>账号列表</span>
						<Badge count={total} style={{ backgroundColor: "#1890ff" }} />
					</Space>
				}
				extra={
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => setCreateModalVisible(true)}
					>
						添加账号
					</Button>
				}
			>
				<Table
					columns={columns}
					dataSource={accounts}
					rowKey="id"
					loading={isLoading}
					pagination={{
						current: queryParams.page,
						pageSize: queryParams.page_size,
						total,
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
						onChange: (page, pageSize) =>
							setQueryParams({ ...queryParams, page, page_size: pageSize }),
					}}
				/>
			</Card>

			{/* 创建账号弹窗 */}
			<Modal
				title="添加账号"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="platform"
						label="平台"
						rules={[{ required: true, message: "请选择平台" }]}
					>
						<Select placeholder="请选择平台">
							{platformTypes.map((p) => (
								<Option key={p.value} value={p.value}>
									{p.label}
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item
						name="account_name"
						label="账号名称"
						rules={[{ required: true, message: "请输入账号名称" }]}
					>
						<Input placeholder="请输入账号名称" />
					</Form.Item>
					<Form.Item
						name="account_id"
						label="平台账号ID"
						rules={[{ required: true, message: "请输入平台账号ID" }]}
					>
						<Input placeholder="请输入平台账号ID" />
					</Form.Item>
					<Form.Item
						name="credentials"
						label="凭证信息"
						tooltip="Cookie、Token 或 API Key，加密存储"
					>
						<Input.TextArea
							rows={3}
							placeholder="请输入凭证信息（Cookie/Token/API Key）"
						/>
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

			{/* 编辑账号弹窗 */}
			<Modal
				title="编辑账号"
				open={editModalVisible}
				onCancel={() => {
					setEditModalVisible(false);
					setEditingAccount(null);
				}}
				footer={null}
				width={500}
			>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item
						name="account_name"
						label="账号名称"
						rules={[{ required: true, message: "请输入账号名称" }]}
					>
						<Input placeholder="请输入账号名称" />
					</Form.Item>
					<Form.Item name="status" label="状态">
						<Select>
							<Option value={1}>正常</Option>
							<Option value={0}>禁用</Option>
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								htmlType="submit"
								loading={updateMutation.isPending}
							>
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
