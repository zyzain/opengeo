"use client";

import api from "@/lib/api";
import {
	BookOutlined,
	CodeOutlined,
	CopyOutlined,
	DeleteOutlined,
	DownloadOutlined,
	EditOutlined,
	FileTextOutlined,
	LayoutOutlined,
	PlusOutlined,
	ShareAltOutlined,
	StarOutlined,
	ThunderboltOutlined,
	UserOutlined,
} from "@ant-design/icons";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
	Avatar,
	Badge,
	Button,
	Card,
	Col,
	Descriptions,
	Form,
	Input,
	List,
	Modal,
	Popconfirm,
	Rate,
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
const { TextArea } = Input;

const templateTypes = [
	{
		value: "prompt",
		label: "Prompt模板",
		color: "blue",
		icon: <ThunderboltOutlined />,
	},
	{
		value: "schema",
		label: "Schema模板",
		color: "green",
		icon: <CodeOutlined />,
	},
	{
		value: "layout",
		label: "布局模板",
		color: "purple",
		icon: <LayoutOutlined />,
	},
];

const tagColors: Record<string, string> = {
	GEO: "blue",
	文章: "green",
	技术: "purple",
	Schema: "orange",
	产品: "cyan",
	发布: "magenta",
	多平台: "red",
	适配: "volcano",
	内容: "geekblue",
	FAQ: "lime",
	问答: "gold",
	合规: "red",
	检测: "orange",
	安全: "purple",
	布局: "blue",
	AI友好: "green",
	结构: "cyan",
};

export default function TemplatesPage() {
	const queryClient = useQueryClient();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [detailModalVisible, setDetailModalVisible] = useState(false);
	const [selectedTemplate, setSelectedTemplate] = useState<any>(null);
	const [editingTemplate, setEditingTemplate] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();
	const [activeTab, setActiveTab] = useState("all");
	const [searchText, setSearchText] = useState("");

	const { data: templatesData, isLoading } = useQuery({
		queryKey: ["templates"],
		queryFn: () => api.templates.list(),
		select: (response: any) => response.data.data,
	});

	const templates = templatesData?.items || templatesData || [];

	// 获取模板类型标签
	const getTemplateTypeTag = (type: string) => {
		const typeInfo = templateTypes.find((t) => t.value === type);
		return (
			<Tag color={typeInfo?.color || "default"} icon={typeInfo?.icon}>
				{typeInfo?.label || type}
			</Tag>
		);
	};

	// 过滤模板
	const filteredTemplates = templates.filter((template: any) => {
		const matchSearch =
			template.name.toLowerCase().includes(searchText.toLowerCase()) ||
			template.description.toLowerCase().includes(searchText.toLowerCase());
		const matchTab =
			activeTab === "all" || template.template_type === activeTab;
		return matchSearch && matchTab;
	});

	// 统计数据
	const stats = {
		total: templates.length,
		prompt: templates.filter((t: any) => t.template_type === "prompt").length,
		schema: templates.filter((t: any) => t.template_type === "schema").length,
		layout: templates.filter((t: any) => t.template_type === "layout").length,
		public: templates.filter((t: any) => t.is_public).length,
		totalUsage: templates.reduce(
			(sum: number, t: any) => sum + (t.usage_count || 0),
			0,
		),
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
			title: "模板名称",
			dataIndex: "name",
			key: "name",
			render: (text: string, record: any) => (
				<a
					onClick={() => handleShowDetail(record)}
					className="text-blue-500 font-medium"
				>
					{text}
				</a>
			),
		},
		{
			title: "类型",
			dataIndex: "template_type",
			key: "template_type",
			width: 120,
			render: (type: string) => getTemplateTypeTag(type),
		},
		{
			title: "评分",
			dataIndex: "rating",
			key: "rating",
			width: 150,
			render: (rating: number) => (
				<Space>
					<Rate
						disabled
						defaultValue={rating}
						allowHalf
						style={{ fontSize: 14 }}
					/>
					<span className="text-gray-500">{rating}</span>
				</Space>
			),
		},
		{
			title: "使用次数",
			dataIndex: "usage_count",
			key: "usage_count",
			width: 100,
			render: (count: number) => (
				<Space>
					<DownloadOutlined />
					<span>{count}</span>
				</Space>
			),
		},
		{
			title: "作者",
			dataIndex: "author",
			key: "author",
			width: 120,
		},
		{
			title: "标签",
			dataIndex: "tags",
			key: "tags",
			render: (tags: string[]) => (
				<Space size={[0, 4]} wrap>
					{tags.map((tag) => (
						<Tag key={tag} color={tagColors[tag] || "default"}>
							{tag}
						</Tag>
					))}
				</Space>
			),
		},
		{
			title: "可见性",
			dataIndex: "is_public",
			key: "is_public",
			width: 80,
			render: (isPublic: boolean) =>
				isPublic ? (
					<Badge status="success" text="公开" />
				) : (
					<Badge status="default" text="私有" />
				),
		},
		{
			title: "操作",
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="查看">
						<Button
							type="text"
							icon={<FileTextOutlined />}
							onClick={() => handleShowDetail(record)}
						/>
					</Tooltip>
					<Tooltip title="编辑">
						<Button
							type="text"
							icon={<EditOutlined />}
							onClick={() => handleEdit(record)}
						/>
					</Tooltip>
					<Tooltip title="分享">
						<Button
							type="text"
							icon={<ShareAltOutlined />}
							onClick={() => handleShare(record)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个模板吗？"
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

	// 显示详情
	const handleShowDetail = (record: any) => {
		setSelectedTemplate(record);
		setDetailModalVisible(true);
	};

	// 复制模板
	const handleCopy = (record: any) => {
		navigator.clipboard.writeText(record.template_data);
		message.success("模板已复制到剪贴板");
	};

	// 创建模板
	const handleCreate = async (values: any) => {
		try {
			const submitData = {
				name: values.name,
				description: values.description,
				template_type: values.template_type,
				content: values.template_data,
				tags: Array.isArray(values.tags) ? values.tags.join(",") : values.tags,
				is_public: values.is_public,
			};
			await api.templates.create(submitData);
			message.success("创建成功");
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: ["templates"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "创建失败");
		}
	};

	// 编辑模板
	const handleEdit = (record: any) => {
		setEditingTemplate(record);
		editForm.setFieldsValue({
			name: record.name,
			description: record.description,
			template_type: record.template_type,
			template_data: record.content || record.template_data,
			tags: record.tags ? (Array.isArray(record.tags) ? record.tags : record.tags.split(",")) : [],
			is_public: record.is_public,
		});
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			const submitData = {
				name: values.name,
				description: values.description,
				template_type: values.template_type,
				content: values.template_data,
				tags: Array.isArray(values.tags) ? values.tags.join(",") : values.tags,
				is_public: values.is_public,
			};
			await api.templates.update(editingTemplate.id, submitData);
			message.success("更新成功");
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingTemplate(null);
			queryClient.invalidateQueries({ queryKey: ["templates"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "更新失败");
		}
	};

	// 删除模板
	const handleDelete = async (id: number) => {
		try {
			await api.templates.delete(id);
			message.success("删除成功");
			queryClient.invalidateQueries({ queryKey: ["templates"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || "删除失败");
		}
	};

	// 分享模板
	const handleShare = (record: any) => {
		navigator.clipboard.writeText(record.template_data);
		message.success("模板内容已复制到剪贴板，可以分享给他人");
	};

	// 使用模板
	const handleUseTemplate = (record: any) => {
		navigator.clipboard.writeText(record.template_data);
		message.success("模板内容已复制，可粘贴使用");
	};

	// 热门模板
	const popularTemplates = [...templates]
		.sort((a, b) => b.usage_count - a.usage_count)
		.slice(0, 5);

	// 高评分模板
	const topRatedTemplates = [...templates]
		.sort((a, b) => b.rating - a.rating)
		.slice(0, 5);

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">Prompt模板市场</h1>
				<p className="text-gray-500 mt-1">共享和发现高质量的GEO优化模板</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="总模板数"
							value={stats.total}
							prefix={<BookOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="Prompt模板"
							value={stats.prompt}
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="Schema模板"
							value={stats.schema}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="布局模板"
							value={stats.layout}
							valueStyle={{ color: "#722ed1" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="公开模板"
							value={stats.public}
							valueStyle={{ color: "#fa8c16" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="总使用次数"
							value={stats.totalUsage}
							prefix={<DownloadOutlined />}
						/>
					</Card>
				</Col>
			</Row>

			<Row gutter={[16, 16]} className="mb-4">
				{/* 热门模板 */}
				<Col xs={24} lg={12}>
					<Card title="热门模板" size="small">
						<List
							dataSource={popularTemplates}
							renderItem={(item, index) => (
								<List.Item>
									<List.Item.Meta
										avatar={
											<Avatar
												style={{
													backgroundColor: index < 3 ? "#1890ff" : "#d9d9d9",
												}}
											>
												{index + 1}
											</Avatar>
										}
										title={
											<Space>
												<span>{item.name}</span>
												{getTemplateTypeTag(item.template_type)}
											</Space>
										}
										description={
											<Space>
												<span>
													<DownloadOutlined /> {item.usage_count}
												</span>
												<span>
													<StarOutlined /> {item.rating}
												</span>
											</Space>
										}
									/>
								</List.Item>
							)}
						/>
					</Card>
				</Col>

				{/* 高评分模板 */}
				<Col xs={24} lg={12}>
					<Card title="高评分模板" size="small">
						<List
							dataSource={topRatedTemplates}
							renderItem={(item, index) => (
								<List.Item>
									<List.Item.Meta
										avatar={
											<Avatar
												style={{
													backgroundColor: index < 3 ? "#52c41a" : "#d9d9d9",
												}}
											>
												{index + 1}
											</Avatar>
										}
										title={
											<Space>
												<span>{item.name}</span>
												{getTemplateTypeTag(item.template_type)}
											</Space>
										}
										description={
											<Space>
												<Rate
													disabled
													defaultValue={item.rating}
													allowHalf
													style={{ fontSize: 12 }}
												/>
												<span>{item.rating}</span>
											</Space>
										}
									/>
								</List.Item>
							)}
						/>
					</Card>
				</Col>
			</Row>

			{/* 模板列表 */}
			<Card
				title={
					<Space>
						<BookOutlined />
						<span>模板列表</span>
					</Space>
				}
				extra={
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => setCreateModalVisible(true)}
					>
						创建模板
					</Button>
				}
			>
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={[
						{ key: "all", label: `全部 (${stats.total})` },
						...templateTypes.map((type) => ({
							key: type.value,
							label: (
								<Space>
									{type.icon}
									<span>{type.label}</span>
									<Badge
										count={
											templates.filter(
												(t: any) => t.template_type === type.value,
											).length
										}
									/>
								</Space>
							),
						})),
					]}
					className="mb-4"
				/>

				<div className="mb-4">
					<Input.Search
						placeholder="搜索模板名称或描述"
						allowClear
						style={{ width: 300 }}
						onSearch={setSearchText}
						onChange={(e) => setSearchText(e.target.value)}
					/>
				</div>

				<Table
					columns={columns}
					dataSource={filteredTemplates}
					rowKey="id"
					loading={isLoading}
					pagination={{
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
					}}
				/>
			</Card>

			{/* 创建模板弹窗 */}
			<Modal
				title="创建模板"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={700}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="name"
						label="模板名称"
						rules={[{ required: true, message: "请输入模板名称" }]}
					>
						<Input placeholder="请输入模板名称" />
					</Form.Item>
					<Form.Item
						name="template_type"
						label="模板类型"
						rules={[{ required: true, message: "请选择模板类型" }]}
					>
						<Select placeholder="请选择模板类型">
							{templateTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<Space>
										{type.icon}
										<span>{type.label}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item
						name="description"
						label="描述"
						rules={[{ required: true, message: "请输入描述" }]}
					>
						<TextArea rows={2} placeholder="请输入模板描述" />
					</Form.Item>
					<Form.Item
						name="template_data"
						label="模板内容 (JSON)"
						rules={[{ required: true, message: "请输入模板内容" }]}
					>
						<TextArea
							rows={6}
							placeholder='{"prompt": "...", "variables": [...]}'
						/>
					</Form.Item>
					<Form.Item name="tags" label="标签">
						<Select mode="tags" placeholder="输入标签后回车">
							{Object.keys(tagColors).map((tag) => (
								<Option key={tag} value={tag}>
									{tag}
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="is_public" label="可见性" initialValue={true}>
						<Select>
							<Option value={true}>公开 - 所有人可见</Option>
							<Option value={false}>私有 - 仅自己可见</Option>
						</Select>
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

			{/* 模板详情弹窗 */}
			<Modal
				title="模板详情"
				open={detailModalVisible}
				onCancel={() => setDetailModalVisible(false)}
				footer={null}
				width={700}
			>
				{selectedTemplate && (
					<div>
						<Descriptions column={2} bordered className="mb-4">
							<Descriptions.Item label="模板名称" span={2}>
								{selectedTemplate.name}
							</Descriptions.Item>
							<Descriptions.Item label="类型">
								{getTemplateTypeTag(selectedTemplate.template_type)}
							</Descriptions.Item>
							<Descriptions.Item label="作者">
								{selectedTemplate.author}
							</Descriptions.Item>
							<Descriptions.Item label="评分">
								<Rate
									disabled
									defaultValue={selectedTemplate.rating}
									allowHalf
								/>{" "}
								{selectedTemplate.rating}
							</Descriptions.Item>
							<Descriptions.Item label="使用次数">
								{selectedTemplate.usage_count}
							</Descriptions.Item>
							<Descriptions.Item label="可见性">
								{selectedTemplate.is_public ? (
									<Badge status="success" text="公开" />
								) : (
									<Badge status="default" text="私有" />
								)}
							</Descriptions.Item>
							<Descriptions.Item label="标签" span={2}>
								<Space size={[0, 4]} wrap>
									{selectedTemplate.tags.map((tag: string) => (
										<Tag key={tag} color={tagColors[tag] || "default"}>
											{tag}
										</Tag>
									))}
								</Space>
							</Descriptions.Item>
							<Descriptions.Item label="描述" span={2}>
								{selectedTemplate.description}
							</Descriptions.Item>
						</Descriptions>

						<Card title="模板内容" size="small">
							<pre className="bg-gray-50 p-4 rounded text-sm overflow-auto max-h-60">
								{JSON.stringify(
									JSON.parse(selectedTemplate.template_data),
									null,
									2,
								)}
							</pre>
						</Card>

						<div className="mt-4 flex justify-end">
							<Space>
								<Button
									icon={<CopyOutlined />}
									onClick={() => handleCopy(selectedTemplate)}
								>
									复制模板
								</Button>
								<Button
									icon={<ShareAltOutlined />}
									onClick={() => handleShare(selectedTemplate)}
								>
									分享模板
								</Button>
								<Button
									type="primary"
									icon={<DownloadOutlined />}
									onClick={() => handleUseTemplate(selectedTemplate)}
								>
									使用模板
								</Button>
							</Space>
						</div>
					</div>
				)}
			</Modal>

			{/* 编辑模板弹窗 */}
			<Modal
				title="编辑模板"
				open={editModalVisible}
				onCancel={() => {
					setEditModalVisible(false);
					setEditingTemplate(null);
				}}
				footer={null}
				width={700}
			>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item
						name="name"
						label="模板名称"
						rules={[{ required: true, message: "请输入模板名称" }]}
					>
						<Input placeholder="请输入模板名称" />
					</Form.Item>
					<Form.Item
						name="template_type"
						label="模板类型"
						rules={[{ required: true, message: "请选择模板类型" }]}
					>
						<Select placeholder="请选择模板类型">
							{templateTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<Space>
										{type.icon}
										<span>{type.label}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item
						name="description"
						label="描述"
						rules={[{ required: true, message: "请输入描述" }]}
					>
						<TextArea rows={2} placeholder="请输入模板描述" />
					</Form.Item>
					<Form.Item
						name="template_data"
						label="模板内容 (JSON)"
						rules={[{ required: true, message: "请输入模板内容" }]}
					>
						<TextArea
							rows={6}
							placeholder='{"prompt": "...", "variables": [...]}'
						/>
					</Form.Item>
					<Form.Item name="tags" label="标签">
						<Select mode="tags" placeholder="输入标签后回车">
							{Object.keys(tagColors).map((tag) => (
								<Option key={tag} value={tag}>
									{tag}
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="is_public" label="可见性">
						<Select>
							<Option value={true}>公开 - 所有人可见</Option>
							<Option value={false}>私有 - 仅自己可见</Option>
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">
								保存
							</Button>
							<Button
								onClick={() => {
									setEditModalVisible(false);
									setEditingTemplate(null);
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
