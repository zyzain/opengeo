import api from "@/lib/api";
import { useAuthStore } from "@/stores";
import {
	ApartmentOutlined,
	BookOutlined,
	DeleteOutlined,
	EditOutlined,
	EnvironmentOutlined,
	LinkOutlined,
	NodeIndexOutlined,
	PlusOutlined,
	SearchOutlined,
	ShopOutlined,
	UserOutlined,
} from "@ant-design/icons";
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
	Row,
	Select,
	Space,
	Spin,
	Statistic,
	Table,
	Tag,
	Tooltip,
	Tree,
	message,
} from "antd";
import { useEffect, useState } from "react";

const { Option } = Select;
const { TextArea } = Input;
const { Search } = Input;

const entityTypes = [
	{ value: "brand", label: "品牌", color: "red", icon: <ShopOutlined /> },
	{
		value: "product",
		label: "产品",
		color: "blue",
		icon: <NodeIndexOutlined />,
	},
	{ value: "concept", label: "概念", color: "green", icon: <BookOutlined /> },
	{ value: "person", label: "人物", color: "purple", icon: <UserOutlined /> },
	{
		value: "place",
		label: "地点",
		color: "orange",
		icon: <EnvironmentOutlined />,
	},
];

export default function KnowledgeEntityPage() {
	const user = useAuthStore((s) => s.user);
	const [entities, setEntities] = useState<any[]>([]);
	const [loading, setLoading] = useState(true);
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [editingEntity, setEditingEntity] = useState<any>(null);
	const [detailModalVisible, setDetailModalVisible] = useState(false);
	const [selectedEntity, setSelectedEntity] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();
	const [searchText, setSearchText] = useState("");
	const [filterType, setFilterType] = useState<string | undefined>(undefined);

	const fetchEntities = async () => {
		setLoading(true);
		try {
			const params: any = {};
			if (filterType) params.entity_type = filterType;
			const res = await api.knowledge.listEntities(params);
			setEntities(res.data.data?.items || []);
		} catch {
			message.error("获取实体列表失败");
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		fetchEntities();
	}, [filterType]);

	const getEntityTypeTag = (type: string) => {
		const typeInfo = entityTypes.find((t) => t.value === type);
		return (
			<Tag color={typeInfo?.color || "default"} icon={typeInfo?.icon}>
				{typeInfo?.label || type}
			</Tag>
		);
	};

	const filteredEntities = entities.filter((entity) =>
		entity.entity_name.toLowerCase().includes(searchText.toLowerCase()),
	);

	const stats = {
		total: entities.length,
		brand: entities.filter((e) => e.entity_type === "brand").length,
		product: entities.filter((e) => e.entity_type === "product").length,
		concept: entities.filter((e) => e.entity_type === "concept").length,
		person: entities.filter((e) => e.entity_type === "person").length,
		place: entities.filter((e) => e.entity_type === "place").length,
	};

	const columns = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
			width: 80,
		},
		{
			title: "实体名称",
			dataIndex: "entity_name",
			key: "entity_name",
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
			title: "实体类型",
			dataIndex: "entity_type",
			key: "entity_type",
			width: 120,
			render: (type: string) => getEntityTypeTag(type),
		},
		{
			title: "关联内容数",
			dataIndex: "content_count",
			key: "content_count",
			width: 120,
			render: (count: number) => (
				<Badge count={count} showZero style={{ backgroundColor: "#1890ff" }} />
			),
		},
		{
			title: "权威链接",
			dataIndex: "authority_links",
			key: "authority_links",
			render: (links: string) => {
				try {
					const linkList = JSON.parse(links);
					return linkList.length > 0 ? (
						<Space size={[0, 4]} wrap>
							{linkList.slice(0, 2).map((link: string, index: number) => (
								<Tag key={index} color="blue">
									<LinkOutlined /> {new URL(link).hostname}
								</Tag>
							))}
							{linkList.length > 2 && <Tag>+{linkList.length - 2}</Tag>}
						</Space>
					) : (
						<span className="text-gray-400">暂无</span>
					);
				} catch {
					return <span className="text-gray-400">暂无</span>;
				}
			},
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
					<Tooltip title="查看">
						<Button
							type="text"
							icon={<SearchOutlined />}
							onClick={() => handleShowDetail(record)}
						/>
					</Tooltip>
					<Tooltip title="编辑">
						<Button
							type="text"
							icon={<EditOutlined />}
							onClick={() => handleShowEdit(record)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个实体吗？"
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

	const handleShowDetail = async (record: any) => {
		try {
			const res = await api.knowledge.getEntity(record.id);
			setSelectedEntity(res.data.data);
			setDetailModalVisible(true);
		} catch {
			message.error("获取实体详情失败");
		}
	};

	const handleShowEdit = (record: any) => {
		setEditingEntity(record);
		editForm.setFieldsValue({
			entity_name: record.entity_name,
			entity_type: record.entity_type,
			entity_data: record.entity_data,
			authority_links: record.authority_links,
		});
		setEditModalVisible(true);
	};

	const handleCreate = async (values: any) => {
		if (!user?.id) {
			message.error("用户未登录，无法创建实体");
			return;
		}
		try {
			await api.knowledge.createEntity({
				...values,
				user_id: user.id,
				entity_data: values.entity_data || "{}",
				authority_links: values.authority_links || "[]",
			});
			message.success("创建成功");
			setCreateModalVisible(false);
			createForm.resetFields();
			fetchEntities();
		} catch {
			message.error("创建失败");
		}
	};

	const handleEdit = async (values: any) => {
		if (!editingEntity) return;
		try {
			await api.knowledge.updateEntity(editingEntity.id, values);
			message.success("更新成功");
			setEditModalVisible(false);
			setEditingEntity(null);
			fetchEntities();
		} catch {
			message.error("更新失败");
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.knowledge.deleteEntity(id);
			message.success("删除成功");
			fetchEntities();
		} catch {
			message.error("删除失败");
		}
	};

	const buildTreeData = () => {
		const grouped = entityTypes.map((type) => ({
			key: type.value,
			title: (
				<Space>
					{type.icon}
					<span>{type.label}</span>
					<Badge
						count={entities.filter((e) => e.entity_type === type.value).length}
						style={{ backgroundColor: type.color }}
					/>
				</Space>
			),
			children: entities
				.filter((e) => e.entity_type === type.value)
				.map((entity) => ({
					key: entity.id,
					title: (
						<Space>
							<span>{entity.entity_name}</span>
							<Tag color="blue">{entity.content_count}个内容</Tag>
						</Space>
					),
				})),
		}));
		return grouped.filter((g) => g.children.length > 0);
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">知识图谱</h1>
				<p className="text-gray-500 mt-1">
					管理品牌、产品、概念等实体，提升AI引用权威性
				</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="总实体数"
							value={stats.total}
							prefix={<ApartmentOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="品牌"
							value={stats.brand}
							valueStyle={{ color: "#ff4d4f" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="产品"
							value={stats.product}
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="概念"
							value={stats.concept}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="人物"
							value={stats.person}
							valueStyle={{ color: "#722ed1" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="地点"
							value={stats.place}
							valueStyle={{ color: "#fa8c16" }}
						/>
					</Card>
				</Col>
			</Row>

			<Spin spinning={loading}>
				<Row gutter={[16, 16]}>
					{/* 知识图谱树 */}
					<Col xs={24} lg={8}>
						<Card title="实体分类" className="h-full">
							<Tree
								showIcon
								defaultExpandAll
								treeData={buildTreeData()}
								className="bg-gray-50 p-4 rounded-lg"
							/>
						</Card>
					</Col>

					{/* 实体列表 */}
					<Col xs={24} lg={16}>
						<Card
							title={
								<Space>
									<NodeIndexOutlined />
									<span>实体列表</span>
								</Space>
							}
							extra={
								<Button
									type="primary"
									icon={<PlusOutlined />}
									onClick={() => setCreateModalVisible(true)}
								>
									添加实体
								</Button>
							}
						>
							{/* 搜索和筛选 */}
							<div className="mb-4 flex flex-wrap gap-4">
								<Search
									placeholder="搜索实体名称"
									allowClear
									style={{ width: 250 }}
									onSearch={setSearchText}
									onChange={(e) => setSearchText(e.target.value)}
								/>
								<Select
									placeholder="实体类型"
									allowClear
									style={{ width: 150 }}
									onChange={setFilterType}
								>
									{entityTypes.map((type) => (
										<Option key={type.value} value={type.value}>
											<Space>
												{type.icon}
												<span>{type.label}</span>
											</Space>
										</Option>
									))}
								</Select>
							</div>

							<Table
								columns={columns}
								dataSource={filteredEntities}
								rowKey="id"
								loading={loading}
								pagination={{
									showSizeChanger: true,
									showQuickJumper: true,
									showTotal: (total) => `共 ${total} 条`,
								}}
							/>
						</Card>
					</Col>
				</Row>
			</Spin>

			{/* 创建实体弹窗 */}
			<Modal
				title="添加实体"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={600}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="entity_name"
						label="实体名称"
						rules={[{ required: true, message: "请输入实体名称" }]}
					>
						<Input placeholder="请输入实体名称" />
					</Form.Item>
					<Form.Item
						name="entity_type"
						label="实体类型"
						rules={[{ required: true, message: "请选择实体类型" }]}
					>
						<Select placeholder="请选择实体类型">
							{entityTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<Space>
										{type.icon}
										<span>{type.label}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="entity_data" label="实体数据 (JSON)">
						<TextArea
							rows={3}
							placeholder='{"description": "...", "website": "..."}'
						/>
					</Form.Item>
					<Form.Item name="authority_links" label="权威链接 (JSON数组)">
						<TextArea
							rows={2}
							placeholder='["https://example.com", "https://wikipedia.org/wiki/xxx"]'
						/>
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

			{/* 编辑实体弹窗 */}
			<Modal
				title="编辑实体"
				open={editModalVisible}
				onCancel={() => {
					setEditModalVisible(false);
					setEditingEntity(null);
				}}
				footer={null}
				width={600}
			>
				<Form form={editForm} layout="vertical" onFinish={handleEdit}>
					<Form.Item
						name="entity_name"
						label="实体名称"
						rules={[{ required: true, message: "请输入实体名称" }]}
					>
						<Input placeholder="请输入实体名称" />
					</Form.Item>
					<Form.Item
						name="entity_type"
						label="实体类型"
						rules={[{ required: true, message: "请选择实体类型" }]}
					>
						<Select placeholder="请选择实体类型">
							{entityTypes.map((type) => (
								<Option key={type.value} value={type.value}>
									<Space>
										{type.icon}
										<span>{type.label}</span>
									</Space>
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item name="entity_data" label="实体数据 (JSON)">
						<TextArea
							rows={3}
							placeholder='{"description": "...", "website": "..."}'
						/>
					</Form.Item>
					<Form.Item name="authority_links" label="权威链接 (JSON数组)">
						<TextArea
							rows={2}
							placeholder='["https://example.com", "https://wikipedia.org/wiki/xxx"]'
						/>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">
								保存
							</Button>
							<Button
								onClick={() => {
									setEditModalVisible(false);
									setEditingEntity(null);
								}}
							>
								取消
							</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			{/* 实体详情弹窗 */}
			<Modal
				title="实体详情"
				open={detailModalVisible}
				onCancel={() => setDetailModalVisible(false)}
				footer={null}
				width={600}
			>
				{selectedEntity && (
					<Descriptions column={1} bordered>
						<Descriptions.Item label="实体名称">
							{selectedEntity.entity_name}
						</Descriptions.Item>
						<Descriptions.Item label="实体类型">
							{getEntityTypeTag(selectedEntity.entity_type)}
						</Descriptions.Item>
						<Descriptions.Item label="关联内容数">
							<Badge
								count={selectedEntity.content_count}
								showZero
								style={{ backgroundColor: "#1890ff" }}
							/>
						</Descriptions.Item>
						<Descriptions.Item label="实体数据">
							<pre className="bg-gray-50 p-2 rounded text-sm">
								{JSON.stringify(
									JSON.parse(selectedEntity.entity_data),
									null,
									2,
								)}
							</pre>
						</Descriptions.Item>
						<Descriptions.Item label="权威链接">
							{(() => {
								try {
									const links = JSON.parse(selectedEntity.authority_links);
									return links.length > 0 ? (
										<ul className="list-disc pl-4">
											{links.map((link: string, index: number) => (
												<li key={index}>
													<a
														href={link}
														target="_blank"
														rel="noopener noreferrer"
														className="text-blue-500"
													>
														{link}
													</a>
												</li>
											))}
										</ul>
									) : (
										<span className="text-gray-400">暂无</span>
									);
								} catch {
									return <span className="text-gray-400">暂无</span>;
								}
							})()}
						</Descriptions.Item>
						<Descriptions.Item label="创建时间">
							{new Date(selectedEntity.created_at).toLocaleString()}
						</Descriptions.Item>
					</Descriptions>
				)}
			</Modal>
		</div>
	);
}
