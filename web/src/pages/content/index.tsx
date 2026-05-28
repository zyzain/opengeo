"use client";

import {
	useCheckCompliance,
	useContents,
	useCreateContent,
	useDeleteContent,
	useOptimizeContent,
	useUpdateContent,
} from "@/hooks";
import {
	CheckCircleOutlined,
	DeleteOutlined,
	EditOutlined,
	EyeOutlined,
	PlusOutlined,
	SearchOutlined,
	SendOutlined,
	ThunderboltOutlined,
} from "@ant-design/icons";
import {
	Badge,
	Button,
	Card,
	Form,
	Input,
	Modal,
	Popconfirm,
	Select,
	Space,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

const { TextArea } = Input;
const { Option } = Select;

// AI模型列表
const aiModels = [
	{
		value: "deepseek",
		label: "DeepSeek",
		color: "blue",
		description: "深度求索，擅长中文理解与生成",
	},
	{
		value: "kimi",
		label: "Kimi",
		color: "purple",
		description: "月之暗面，长文本处理能力强",
	},
	{
		value: "doubao",
		label: "豆包",
		color: "orange",
		description: "字节跳动，多模态理解能力",
	},
	{
		value: "chatgpt",
		label: "ChatGPT",
		color: "green",
		description: "OpenAI，全球领先的通用模型",
	},
];

// 优化类型
const optimizationTypes = [
	{
		value: "geo_semantic",
		label: "GEO语义增强",
		description: "优化内容结构，提升AI可读性",
	},
	{
		value: "schema_markup",
		label: "Schema标记优化",
		description: "自动添加结构化数据标记",
	},
	{
		value: "authority_boost",
		label: "权威性提升",
		description: "增强内容信源权重",
	},
	{
		value: "multi_platform",
		label: "多平台适配",
		description: "适配不同平台内容规范",
	},
];

export default function ContentPage() {
	const navigate = useNavigate();
	const [searchForm] = Form.useForm();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [optimizeModalVisible, setOptimizeModalVisible] = useState(false);
	const [editingContent, setEditingContent] = useState<any>(null);
	const [optimizingContent, setOptimizingContent] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();
	const [optimizeForm] = Form.useForm();

	// 查询参数
	const [queryParams, setQueryParams] = useState({
		page: 1,
		page_size: 10,
		content_type: undefined,
		status: undefined,
	});

	const { data, isLoading } = useContents(queryParams);
	const createMutation = useCreateContent();
	const updateMutation = useUpdateContent();
	const deleteMutation = useDeleteContent();
	const optimizeMutation = useOptimizeContent();
	const complianceMutation = useCheckCompliance();

	const contents = data?.items || [];
	const total = data?.total || 0;

	// 内容状态标签
	const getStatusTag = (status: number) => {
		const statusMap: Record<number, { color: string; text: string }> = {
			0: { color: "warning", text: "草稿" },
			1: { color: "success", text: "已发布" },
			2: { color: "default", text: "已归档" },
		};
		const config = statusMap[status] || { color: "default", text: "未知" };
		return <Tag color={config.color}>{config.text}</Tag>;
	};

	// 内容类型标签
	const getContentTypeTag = (type: string) => {
		const typeMap: Record<string, { color: string; text: string }> = {
			article: { color: "blue", text: "文章" },
			video: { color: "purple", text: "视频" },
			image: { color: "orange", text: "图片" },
		};
		const config = typeMap[type] || { color: "default", text: type };
		return <Tag color={config.color}>{config.text}</Tag>;
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
			title: "标题",
			dataIndex: "title",
			key: "title",
			ellipsis: true,
			render: (text: string, record: any) => (
				<a
					onClick={() => navigate(`/content/${record.id}`)}
					className="text-blue-500 hover:text-blue-700"
				>
					{text}
				</a>
			),
		},
		{
			title: "类型",
			dataIndex: "content_type",
			key: "content_type",
			width: 100,
			render: (type: string) => getContentTypeTag(type),
		},
		{
			title: "状态",
			dataIndex: "status",
			key: "status",
			width: 100,
			render: (status: number) => getStatusTag(status),
		},
		{
			title: "AI评分",
			dataIndex: "ai_optimization_score",
			key: "ai_optimization_score",
			width: 100,
			render: (score: number) => (
				<span
					className={
						score >= 80
							? "text-green-500"
							: score >= 60
								? "text-orange-500"
								: "text-red-500"
					}
				>
					{score ? `${score}分` : "-"}
				</span>
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
			width: 250,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="查看">
						<Button
							type="text"
							icon={<EyeOutlined />}
							onClick={() => navigate(`/content/${record.id}`)}
						/>
					</Tooltip>
					<Tooltip title="编辑">
						<Button
							type="text"
							icon={<EditOutlined />}
							onClick={() => handleEdit(record)}
						/>
					</Tooltip>
					<Tooltip title="AI优化">
						<Button
							type="text"
							icon={<ThunderboltOutlined />}
							onClick={() => handleOptimizeClick(record)}
							loading={optimizeMutation.isPending}
						/>
					</Tooltip>
					<Tooltip title="合规检测">
						<Button
							type="text"
							icon={<CheckCircleOutlined />}
							onClick={() => handleComplianceCheck(record.id)}
							loading={complianceMutation.isPending}
						/>
					</Tooltip>
					<Tooltip title="发布">
						<Button
							type="text"
							icon={<SendOutlined />}
							onClick={() => handlePublish(record.id)}
							disabled={record.status === 1}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个内容吗？"
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

	// 创建内容
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

	// 编辑内容
	const handleEdit = (record: any) => {
		setEditingContent(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await updateMutation.mutateAsync({ id: editingContent.id, data: values });
			message.success("更新成功");
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingContent(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || "更新失败");
		}
	};

	// 删除内容
	const handleDelete = async (id: number) => {
		try {
			await deleteMutation.mutateAsync(id);
			message.success("删除成功");
		} catch (error: any) {
			message.error(error.response?.data?.message || "删除失败");
		}
	};

	// AI优化 - 打开弹窗
	const handleOptimizeClick = (record: any) => {
		setOptimizingContent(record);
		optimizeForm.setFieldsValue({
			ai_model: "deepseek",
			optimization_type: "geo_semantic",
		});
		setOptimizeModalVisible(true);
	};

	// AI优化 - 执行
	const handleOptimize = async (values: any) => {
		if (!optimizingContent) return;
		try {
			await optimizeMutation.mutateAsync({
				id: optimizingContent.id,
				data: values,
			});
			message.success("AI优化完成");
			setOptimizeModalVisible(false);
			optimizeForm.resetFields();
			setOptimizingContent(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || "AI优化失败");
		}
	};

	// 合规检测
	const handleComplianceCheck = async (id: number) => {
		try {
			const res = await complianceMutation.mutateAsync({ id });
			const result = res.data.data;
			if (result.passed) {
				message.success(`合规检测通过！评分: ${result.score}`);
			} else {
				message.warning("合规检测发现问题，请查看详情");
			}
		} catch (error: any) {
			message.error(error.response?.data?.message || "合规检测失败");
		}
	};

	// 发布内容
	const handlePublish = (id: number) => {
		navigate(`/publish/tasks?content_id=${id}`);
	};

	// 搜索
	const handleSearch = (values: any) => {
		setQueryParams({ ...queryParams, ...values, page: 1 });
	};

	// 重置搜索
	const handleReset = () => {
		searchForm.resetFields();
		setQueryParams({
			page: 1,
			page_size: 10,
			content_type: undefined,
			status: undefined,
		});
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">内容管理</h1>
				<p className="text-gray-500 mt-1">管理和优化您的内容</p>
			</div>

			{/* 搜索表单 */}
			<Card className="mb-4">
				<Form form={searchForm} layout="inline" onFinish={handleSearch}>
					<Form.Item name="content_type" label="内容类型">
						<Select placeholder="请选择" allowClear style={{ width: 120 }}>
							<Option value="article">文章</Option>
							<Option value="video">视频</Option>
							<Option value="image">图片</Option>
						</Select>
					</Form.Item>
					<Form.Item name="status" label="状态">
						<Select placeholder="请选择" allowClear style={{ width: 120 }}>
							<Option value={0}>草稿</Option>
							<Option value={1}>已发布</Option>
							<Option value={2}>已归档</Option>
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
						</Space>
					</Form.Item>
				</Form>
			</Card>

			{/* 内容列表 */}
			<Card
				title="内容列表"
				extra={
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => setCreateModalVisible(true)}
					>
						创建内容
					</Button>
				}
			>
				<Table
					columns={columns}
					dataSource={contents}
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

			{/* 创建内容弹窗 */}
			<Modal
				title="创建内容"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={640}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="title"
						label="标题"
						rules={[{ required: true, message: "请输入标题" }]}
					>
						<Input placeholder="请输入标题" />
					</Form.Item>
					<Form.Item
						name="content_type"
						label="内容类型"
						rules={[{ required: true, message: "请选择内容类型" }]}
					>
						<Select placeholder="请选择">
							<Option value="article">文章</Option>
							<Option value="video">视频</Option>
							<Option value="image">图片</Option>
						</Select>
					</Form.Item>
					<Form.Item
						name="body"
						label="正文"
						rules={[{ required: true, message: "请输入正文" }]}
					>
						<TextArea rows={6} placeholder="请输入正文" />
					</Form.Item>
					<Form.Item name="schema_markup" label="Schema Markup">
						<TextArea rows={3} placeholder="可选，JSON格式的结构化数据" />
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

			{/* 编辑内容弹窗 */}
			<Modal
				title="编辑内容"
				open={editModalVisible}
				onCancel={() => {
					setEditModalVisible(false);
					setEditingContent(null);
				}}
				footer={null}
				width={640}
			>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item
						name="title"
						label="标题"
						rules={[{ required: true, message: "请输入标题" }]}
					>
						<Input placeholder="请输入标题" />
					</Form.Item>
					<Form.Item
						name="content_type"
						label="内容类型"
						rules={[{ required: true, message: "请选择内容类型" }]}
					>
						<Select placeholder="请选择">
							<Option value="article">文章</Option>
							<Option value="video">视频</Option>
							<Option value="image">图片</Option>
						</Select>
					</Form.Item>
					<Form.Item
						name="body"
						label="正文"
						rules={[{ required: true, message: "请输入正文" }]}
					>
						<TextArea rows={6} placeholder="请输入正文" />
					</Form.Item>
					<Form.Item name="schema_markup" label="Schema Markup">
						<TextArea rows={3} placeholder="可选，JSON格式的结构化数据" />
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

			{/* AI优化弹窗 */}
			<Modal
				title={
					<Space>
						<ThunderboltOutlined />
						<span>AI内容优化</span>
					</Space>
				}
				open={optimizeModalVisible}
				onCancel={() => {
					setOptimizeModalVisible(false);
					setOptimizingContent(null);
				}}
				footer={null}
				width={520}
			>
				{optimizingContent && (
					<div className="mb-4 p-3 bg-blue-50 rounded-lg">
						<div className="text-sm text-gray-500">正在优化内容：</div>
						<div className="font-medium">{optimizingContent.title}</div>
					</div>
				)}
				<Form form={optimizeForm} layout="vertical" onFinish={handleOptimize}>
					<Form.Item
						name="ai_model"
						label="选择AI模型"
						rules={[{ required: true, message: "请选择AI模型" }]}
					>
						<Select placeholder="请选择AI模型">
							{aiModels.map((model) => (
								<Option key={model.value} value={model.value}>
									<Space>
										<Tag color={model.color}>{model.label}</Tag>
										<span className="text-gray-400">{model.description}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item
						name="optimization_type"
						label="优化类型"
						rules={[{ required: true, message: "请选择优化类型" }]}
					>
						<Select placeholder="请选择优化类型">
							{optimizationTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<div>
										<div>{type.label}</div>
										<div className="text-xs text-gray-400">
											{type.description}
										</div>
									</div>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								htmlType="submit"
								loading={optimizeMutation.isPending}
								icon={<ThunderboltOutlined />}
							>
								开始优化
							</Button>
							<Button onClick={() => setOptimizeModalVisible(false)}>
								取消
							</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
