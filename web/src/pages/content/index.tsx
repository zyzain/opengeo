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
import { useIntl } from "react-intl";

const { TextArea } = Input;
const { Option } = Select;

export default function ContentPage() {
	const intl = useIntl();
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

	// AI模型列表
	const aiModels = [
		{
			value: "deepseek",
			label: "DeepSeek",
			color: "blue",
			description: intl.formatMessage({ id: 'content.aiModel.deepseek.desc' }),
		},
		{
			value: "kimi",
			label: "Kimi",
			color: "purple",
			description: intl.formatMessage({ id: 'content.aiModel.kimi.desc' }),
		},
		{
			value: "doubao",
			label: intl.formatMessage({ id: 'content.aiModel.doubao' }),
			color: "orange",
			description: intl.formatMessage({ id: 'content.aiModel.doubao.desc' }),
		},
		{
			value: "chatgpt",
			label: "ChatGPT",
			color: "green",
			description: intl.formatMessage({ id: 'content.aiModel.chatgpt.desc' }),
		},
	];

	// 优化类型
	const optimizationTypes = [
		{
			value: "geo_semantic",
			label: intl.formatMessage({ id: 'content.optimization.geoSemantic' }),
			description: intl.formatMessage({ id: 'content.optimization.geoSemantic.desc' }),
		},
		{
			value: "schema_markup",
			label: intl.formatMessage({ id: 'content.optimization.schemaMarkup' }),
			description: intl.formatMessage({ id: 'content.optimization.schemaMarkup.desc' }),
		},
		{
			value: "authority_boost",
			label: intl.formatMessage({ id: 'content.optimization.authorityBoost' }),
			description: intl.formatMessage({ id: 'content.optimization.authorityBoost.desc' }),
		},
		{
			value: "multi_platform",
			label: intl.formatMessage({ id: 'content.optimization.multiPlatform' }),
			description: intl.formatMessage({ id: 'content.optimization.multiPlatform.desc' }),
		},
	];

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
			0: { color: "warning", text: intl.formatMessage({ id: 'content.status.draft' }) },
			1: { color: "success", text: intl.formatMessage({ id: 'content.status.published' }) },
			2: { color: "default", text: intl.formatMessage({ id: 'content.status.archived' }) },
		};
		const config = statusMap[status] || { color: "default", text: intl.formatMessage({ id: 'common.status.unknown' }) };
		return <Tag color={config.color}>{config.text}</Tag>;
	};

	// 内容类型标签
	const getContentTypeTag = (type: string) => {
		const typeMap: Record<string, { color: string; text: string }> = {
			article: { color: "blue", text: intl.formatMessage({ id: 'content.type.article' }) },
			video: { color: "purple", text: intl.formatMessage({ id: 'content.type.video' }) },
			image: { color: "orange", text: intl.formatMessage({ id: 'content.type.image' }) },
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
			title: intl.formatMessage({ id: 'content.title' }),
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
			title: intl.formatMessage({ id: 'common.column.type' }),
			dataIndex: "content_type",
			key: "content_type",
			width: 100,
			render: (type: string) => getContentTypeTag(type),
		},
		{
			title: intl.formatMessage({ id: 'common.column.status' }),
			dataIndex: "status",
			key: "status",
			width: 100,
			render: (status: number) => getStatusTag(status),
		},
		{
			title: intl.formatMessage({ id: 'content.aiScore' }),
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
					{score ? `${score}${intl.formatMessage({ id: 'common.unit.score' })}` : "-"}
				</span>
			),
		},
		{
			title: intl.formatMessage({ id: 'common.column.createdAt' }),
			dataIndex: "created_at",
			key: "created_at",
			width: 180,
			render: (text: string) => new Date(text).toLocaleString(),
		},
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 250,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.view' })}>
						<Button
							type="text"
							icon={<EyeOutlined />}
							onClick={() => navigate(`/content/${record.id}`)}
						/>
					</Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}>
						<Button
							type="text"
							icon={<EditOutlined />}
							onClick={() => handleEdit(record)}
						/>
					</Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'content.optimize' })}>
						<Button
							type="text"
							icon={<ThunderboltOutlined />}
							onClick={() => handleOptimizeClick(record)}
							loading={optimizeMutation.isPending}
						/>
					</Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'content.compliance' })}>
						<Button
							type="text"
							icon={<CheckCircleOutlined />}
							onClick={() => handleComplianceCheck(record.id)}
							loading={complianceMutation.isPending}
						/>
					</Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'content.publish' })}>
						<Button
							type="text"
							icon={<SendOutlined />}
							onClick={() => handlePublish(record.id)}
							disabled={record.status === 1}
						/>
					</Tooltip>
					<Popconfirm
						title={intl.formatMessage({ id: 'content.confirmDelete' })}
						onConfirm={() => handleDelete(record.id)}
						okText={intl.formatMessage({ id: 'common.confirm' })}
						cancelText={intl.formatMessage({ id: 'common.cancel' })}
					>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}>
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
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.createFailed' }));
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
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingContent(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	// 删除内容
	const handleDelete = async (id: number) => {
		try {
			await deleteMutation.mutateAsync(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.deleteFailed' }));
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
			message.success(intl.formatMessage({ id: 'content.optimization.completed' }));
			setOptimizeModalVisible(false);
			optimizeForm.resetFields();
			setOptimizingContent(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'content.optimization.failed' }));
		}
	};

	// 合规检测
	const handleComplianceCheck = async (id: number) => {
		try {
			const res = await complianceMutation.mutateAsync({ id });
			const result = res.data.data;
			if (result.passed) {
				message.success(intl.formatMessage({ id: 'content.compliance.passed' }, { score: result.score }));
			} else {
				message.warning(intl.formatMessage({ id: 'content.compliance.issues' }));
			}
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'content.compliance.failed' }));
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
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'content.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'content.page.subtitle' })}</p>
			</div>

			{/* 搜索表单 */}
			<Card className="mb-4">
				<Form form={searchForm} layout="inline" onFinish={handleSearch}>
					<Form.Item name="content_type" label={intl.formatMessage({ id: 'content.form.contentType' })}>
						<Select placeholder={intl.formatMessage({ id: 'content.select.placeholder' })} allowClear style={{ width: 120 }}>
							<Option value="article">{intl.formatMessage({ id: 'content.type.article' })}</Option>
							<Option value="video">{intl.formatMessage({ id: 'content.type.video' })}</Option>
							<Option value="image">{intl.formatMessage({ id: 'content.type.image' })}</Option>
						</Select>
					</Form.Item>
					<Form.Item name="status" label={intl.formatMessage({ id: 'common.form.status' })}>
						<Select placeholder={intl.formatMessage({ id: 'content.select.placeholder' })} allowClear style={{ width: 120 }}>
							<Option value={0}>{intl.formatMessage({ id: 'content.status.draft' })}</Option>
							<Option value={1}>{intl.formatMessage({ id: 'content.status.published' })}</Option>
							<Option value={2}>{intl.formatMessage({ id: 'content.status.archived' })}</Option>
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								icon={<SearchOutlined />}
								htmlType="submit"
							>
								{intl.formatMessage({ id: 'common.search' })}
							</Button>
							<Button onClick={handleReset}>{intl.formatMessage({ id: 'common.reset' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Card>

			{/* 内容列表 */}
			<Card
				title={intl.formatMessage({ id: 'content.list.title' })}
				extra={
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => setCreateModalVisible(true)}
					>
						{intl.formatMessage({ id: 'content.action.create' })}
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
						showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }),
						onChange: (page, pageSize) =>
							setQueryParams({ ...queryParams, page, page_size: pageSize }),
					}}
				/>
			</Card>

			{/* 创建内容弹窗 */}
			<Modal
				title={intl.formatMessage({ id: 'content.modal.create' })}
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={640}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="title"
						label={intl.formatMessage({ id: 'content.form.title' })}
						rules={[{ required: true, message: intl.formatMessage({ id: 'content.validation.enterTitle' }) }]}
					>
						<Input placeholder={intl.formatMessage({ id: 'content.placeholder.title' })} />
					</Form.Item>
					<Form.Item
						name="content_type"
						label={intl.formatMessage({ id: 'content.form.contentType' })}
						rules={[{ required: true, message: intl.formatMessage({ id: 'content.validation.selectContentType' }) }]}
					>
						<Select placeholder={intl.formatMessage({ id: 'content.select.placeholder' })}>
							<Option value="article">{intl.formatMessage({ id: 'content.type.article' })}</Option>
							<Option value="video">{intl.formatMessage({ id: 'content.type.video' })}</Option>
							<Option value="image">{intl.formatMessage({ id: 'content.type.image' })}</Option>
						</Select>
					</Form.Item>
					<Form.Item
						name="body"
						label={intl.formatMessage({ id: 'content.form.body' })}
						rules={[{ required: true, message: intl.formatMessage({ id: 'content.validation.enterBody' }) }]}
					>
						<TextArea rows={6} placeholder={intl.formatMessage({ id: 'content.placeholder.body' })} />
					</Form.Item>
					<Form.Item name="schema_markup" label={intl.formatMessage({ id: 'content.form.schemaMarkup' })}>
						<TextArea rows={3} placeholder={intl.formatMessage({ id: 'content.placeholder.schemaMarkup' })} />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								htmlType="submit"
								loading={createMutation.isPending}
							>
								{intl.formatMessage({ id: 'common.create' })}
							</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			{/* 编辑内容弹窗 */}
			<Modal
				title={intl.formatMessage({ id: 'content.modal.edit' })}
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
						label={intl.formatMessage({ id: 'content.form.title' })}
						rules={[{ required: true, message: intl.formatMessage({ id: 'content.validation.enterTitle' }) }]}
					>
						<Input placeholder={intl.formatMessage({ id: 'content.placeholder.title' })} />
					</Form.Item>
					<Form.Item
						name="content_type"
						label={intl.formatMessage({ id: 'content.form.contentType' })}
						rules={[{ required: true, message: intl.formatMessage({ id: 'content.validation.selectContentType' }) }]}
					>
						<Select placeholder={intl.formatMessage({ id: 'content.select.placeholder' })}>
							<Option value="article">{intl.formatMessage({ id: 'content.type.article' })}</Option>
							<Option value="video">{intl.formatMessage({ id: 'content.type.video' })}</Option>
							<Option value="image">{intl.formatMessage({ id: 'content.type.image' })}</Option>
						</Select>
					</Form.Item>
					<Form.Item
						name="body"
						label={intl.formatMessage({ id: 'content.form.body' })}
						rules={[{ required: true, message: intl.formatMessage({ id: 'content.validation.enterBody' }) }]}
					>
						<TextArea rows={6} placeholder={intl.formatMessage({ id: 'content.placeholder.body' })} />
					</Form.Item>
					<Form.Item name="schema_markup" label={intl.formatMessage({ id: 'content.form.schemaMarkup' })}>
						<TextArea rows={3} placeholder={intl.formatMessage({ id: 'content.placeholder.schemaMarkup' })} />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								htmlType="submit"
								loading={updateMutation.isPending}
							>
								{intl.formatMessage({ id: 'common.save' })}
							</Button>
							<Button onClick={() => setEditModalVisible(false)}>{intl.formatMessage({ id: 'common.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			{/* AI优化弹窗 */}
			<Modal
				title={
					<Space>
						<ThunderboltOutlined />
						<span>{intl.formatMessage({ id: 'content.modal.optimize' })}</span>
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
						<div className="text-sm text-gray-500">{intl.formatMessage({ id: 'content.placeholder.optimizing' })}</div>
						<div className="font-medium">{optimizingContent.title}</div>
					</div>
				)}
				<Form form={optimizeForm} layout="vertical" onFinish={handleOptimize}>
					<Form.Item
						name="ai_model"
						label={intl.formatMessage({ id: 'content.aiModel.label' })}
						rules={[{ required: true, message: intl.formatMessage({ id: 'content.aiModel.label' }) }]}
					>
						<Select placeholder={intl.formatMessage({ id: 'content.aiModel.label' })}>
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
						label={intl.formatMessage({ id: 'content.optimization.type' })}
						rules={[{ required: true, message: intl.formatMessage({ id: 'content.optimization.type' }) }]}
					>
						<Select placeholder={intl.formatMessage({ id: 'content.optimization.type' })}>
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
								{intl.formatMessage({ id: 'content.optimization.start' })}
							</Button>
							<Button onClick={() => setOptimizeModalVisible(false)}>
								{intl.formatMessage({ id: 'common.cancel' })}
							</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
