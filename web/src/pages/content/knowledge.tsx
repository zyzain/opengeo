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
import { useIntl } from "react-intl";

const { Option } = Select;
const { TextArea } = Input;
const { Search } = Input;

export default function KnowledgeEntityPage() {
	const intl = useIntl();
	const user = useAuthStore((s) => s.user);

	const entityTypes = [
		{ value: "brand", label: intl.formatMessage({ id: 'knowledge.entity.brand' }), color: "red", icon: <ShopOutlined /> },
		{ value: "product", label: intl.formatMessage({ id: 'knowledge.entity.product' }), color: "blue", icon: <NodeIndexOutlined /> },
		{ value: "concept", label: intl.formatMessage({ id: 'knowledge.entity.concept' }), color: "green", icon: <BookOutlined /> },
		{ value: "person", label: intl.formatMessage({ id: 'knowledge.entity.person' }), color: "purple", icon: <UserOutlined /> },
		{ value: "place", label: intl.formatMessage({ id: 'knowledge.entity.place' }), color: "orange", icon: <EnvironmentOutlined /> },
	];
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
			message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		fetchEntities();
	}, [filterType]);

	const getEntityTypeTag = (type: string) => {
		const typeInfo = entityTypes.find((t) => t.value === type);
		return <Tag color={typeInfo?.color || "default"} icon={typeInfo?.icon}>{typeInfo?.label || type}</Tag>;
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
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{
			title: intl.formatMessage({ id: 'knowledge.form.entityName' }),
			dataIndex: "entity_name",
			key: "entity_name",
			render: (text: string, record: any) => <a onClick={() => handleShowDetail(record)} className="text-blue-500 font-medium">{text}</a>,
		},
		{ title: intl.formatMessage({ id: 'knowledge.form.entityType' }), dataIndex: "entity_type", key: "entity_type", width: 120, render: (type: string) => getEntityTypeTag(type) },
		{ title: intl.formatMessage({ id: 'knowledge.column.contentCount' }), dataIndex: "content_count", key: "content_count", width: 120, render: (count: number) => <Badge count={count} showZero style={{ backgroundColor: "#1890ff" }} /> },
		{
			title: intl.formatMessage({ id: 'knowledge.column.authorityLinks' }),
			dataIndex: "authority_links",
			key: "authority_links",
			render: (links: string) => {
				try {
					const linkList = JSON.parse(links);
					return linkList.length > 0 ? (
						<Space size={[0, 4]} wrap>
							{linkList.slice(0, 2).map((link: string, index: number) => <Tag key={index} color="blue"><LinkOutlined /> {new URL(link).hostname}</Tag>)}
							{linkList.length > 2 && <Tag>+{linkList.length - 2}</Tag>}
						</Space>
					) : <span className="text-gray-400">{intl.formatMessage({ id: 'knowledge.none' })}</span>;
				} catch { return <span className="text-gray-400">{intl.formatMessage({ id: 'knowledge.none' })}</span>; }
			},
		},
		{ title: intl.formatMessage({ id: 'common.column.createdAt' }), dataIndex: "created_at", key: "created_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.view' })}><Button type="text" icon={<SearchOutlined />} onClick={() => handleShowDetail(record)} /></Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => handleShowEdit(record)} /></Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'common.confirmDelete' })} onConfirm={() => handleDelete(record.id)} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
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
			message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
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
			message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
			return;
		}
		try {
			await api.knowledge.createEntity({
				...values,
				user_id: user.id,
				entity_data: values.entity_data || "{}",
				authority_links: values.authority_links || "[]",
			});
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
			fetchEntities();
		} catch {
			message.error(intl.formatMessage({ id: 'common.message.createFailed' }));
		}
	};

	const handleEdit = async (values: any) => {
		if (!editingEntity) return;
		try {
			await api.knowledge.updateEntity(editingEntity.id, values);
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
			setEditingEntity(null);
			fetchEntities();
		} catch {
			message.error(intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.knowledge.deleteEntity(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
			fetchEntities();
		} catch {
			message.error(intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const buildTreeData = () => {
		const grouped = entityTypes.map((type) => ({
			key: type.value,
			title: <Space>{type.icon}<span>{type.label}</span><Badge count={entities.filter((e) => e.entity_type === type.value).length} style={{ backgroundColor: type.color }} /></Space>,
			children: entities.filter((e) => e.entity_type === type.value).map((entity) => ({
				key: entity.id,
				title: <Space><span>{entity.entity_name}</span><Tag color="blue">{intl.formatMessage({ id: 'knowledge.tree.contentCount' }, { count: entity.content_count })}</Tag></Space>,
			})),
		}));
		return grouped.filter((g) => g.children.length > 0);
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'knowledge.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'knowledge.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'knowledge.stat.total' })} value={stats.total} prefix={<ApartmentOutlined />} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'knowledge.stat.brand' })} value={stats.brand} valueStyle={{ color: "#ff4d4f" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'knowledge.stat.product' })} value={stats.product} valueStyle={{ color: "#1890ff" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'knowledge.stat.concept' })} value={stats.concept} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'knowledge.stat.person' })} value={stats.person} valueStyle={{ color: "#722ed1" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'knowledge.stat.place' })} value={stats.place} valueStyle={{ color: "#fa8c16" }} /></Card></Col>
			</Row>

			<Spin spinning={loading}>
				<Row gutter={[16, 16]}>
					<Col xs={24} lg={8}>
						<Card title={intl.formatMessage({ id: 'knowledge.card.entityClassify' })} className="h-full">
							<Tree showIcon defaultExpandAll treeData={buildTreeData()} className="bg-gray-50 p-4 rounded-lg" />
						</Card>
					</Col>

					<Col xs={24} lg={16}>
						<Card
							title={<Space><NodeIndexOutlined /><span>{intl.formatMessage({ id: 'knowledge.card.entityList' })}</span></Space>}
							extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'knowledge.action.addEntity' })}</Button>}
						>
							<div className="mb-4 flex flex-wrap gap-4">
								<Search placeholder={intl.formatMessage({ id: 'knowledge.placeholder.entityName' })} allowClear style={{ width: 250 }} onSearch={setSearchText} onChange={(e) => setSearchText(e.target.value)} />
								<Select placeholder={intl.formatMessage({ id: 'knowledge.placeholder.entityType' })} allowClear style={{ width: 150 }} onChange={setFilterType}>
									{entityTypes.map((type) => <Option key={type.value} value={type.value}><Space>{type.icon}<span>{type.label}</span></Space></Option>)}
								</Select>
							</div>

							<Table columns={columns} dataSource={filteredEntities} rowKey="id" loading={loading} pagination={{ showSizeChanger: true, showQuickJumper: true, showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }) }} />
						</Card>
					</Col>
				</Row>
			</Spin>

			<Modal title={intl.formatMessage({ id: 'knowledge.modal.addEntity' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={600}>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item name="entity_name" label={intl.formatMessage({ id: 'knowledge.form.entityName' })} rules={[{ required: true, message: intl.formatMessage({ id: 'knowledge.validation.enterEntityName' }) }]}><Input placeholder={intl.formatMessage({ id: 'knowledge.placeholder.entityName' })} /></Form.Item>
					<Form.Item name="entity_type" label={intl.formatMessage({ id: 'knowledge.form.entityType' })} rules={[{ required: true, message: intl.formatMessage({ id: 'knowledge.validation.selectEntityType' }) }]}>
						<Select placeholder={intl.formatMessage({ id: 'knowledge.placeholder.entityType' })}>
							{entityTypes.map((type) => <Option key={type.value} value={type.value}><Space>{type.icon}<span>{type.label}</span></Space></Option>)}
						</Select>
					</Form.Item>
				<Form.Item name="entity_data" label={intl.formatMessage({ id: 'knowledge.form.entityData' })}><TextArea rows={3} placeholder='{"description": "...", "website": "..."}' /></Form.Item>
				<Form.Item name="authority_links" label={intl.formatMessage({ id: 'knowledge.form.authorityLinks' })}><TextArea rows={2} placeholder='["https://example.com", "https://wikipedia.org/wiki/xxx"]' /></Form.Item>
				<Form.Item>
					<Space>
						<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

		<Modal title={intl.formatMessage({ id: 'knowledge.modal.editEntity' })} open={editModalVisible} onCancel={() => { setEditModalVisible(false); setEditingEntity(null); }} footer={null} width={600}>
			<Form form={editForm} layout="vertical" onFinish={handleEdit}>
					<Form.Item name="entity_name" label={intl.formatMessage({ id: 'knowledge.form.entityName' })} rules={[{ required: true, message: intl.formatMessage({ id: 'knowledge.validation.enterEntityName' }) }]}><Input placeholder={intl.formatMessage({ id: 'knowledge.placeholder.entityName' })} /></Form.Item>
					<Form.Item name="entity_type" label={intl.formatMessage({ id: 'knowledge.form.entityType' })} rules={[{ required: true, message: intl.formatMessage({ id: 'knowledge.validation.selectEntityType' }) }]}>
						<Select placeholder={intl.formatMessage({ id: 'knowledge.placeholder.entityType' })}>
							{entityTypes.map((type) => <Option key={type.value} value={type.value}><Space>{type.icon}<span>{type.label}</span></Space></Option>)}
						</Select>
					</Form.Item>
				<Form.Item name="entity_data" label={intl.formatMessage({ id: 'knowledge.form.entityData' })}><TextArea rows={3} placeholder='{"description": "...", "website": "..."}' /></Form.Item>
				<Form.Item name="authority_links" label={intl.formatMessage({ id: 'knowledge.form.authorityLinks' })}><TextArea rows={2} placeholder='["https://example.com", "https://wikipedia.org/wiki/xxx"]' /></Form.Item>
				<Form.Item>
					<Space>
						<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.save' })}</Button>
							<Button onClick={() => { setEditModalVisible(false); setEditingEntity(null); }}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title={intl.formatMessage({ id: 'knowledge.modal.entityDetail' })} open={detailModalVisible} onCancel={() => setDetailModalVisible(false)} footer={null} width={600}>
				{selectedEntity && (
					<Descriptions column={1} bordered>
						<Descriptions.Item label={intl.formatMessage({ id: 'knowledge.form.entityName' })}>{selectedEntity.entity_name}</Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: 'knowledge.form.entityType' })}>{getEntityTypeTag(selectedEntity.entity_type)}</Descriptions.Item>
					<Descriptions.Item label={intl.formatMessage({ id: 'knowledge.desc.contentCount' })}><Badge count={selectedEntity.content_count} showZero style={{ backgroundColor: "#1890ff" }} /></Descriptions.Item>
					<Descriptions.Item label={intl.formatMessage({ id: 'knowledge.desc.entityData' })}><pre className="bg-gray-50 p-2 rounded text-sm">{JSON.stringify(JSON.parse(selectedEntity.entity_data), null, 2)}</pre></Descriptions.Item>
					<Descriptions.Item label={intl.formatMessage({ id: 'knowledge.desc.authorityLinks' })}>
							{(() => {
								try {
									const links = JSON.parse(selectedEntity.authority_links);
									return links.length > 0 ? (
										<ul className="list-disc pl-4">
											{links.map((link: string, index: number) => <li key={index}><a href={link} target="_blank" rel="noopener noreferrer" className="text-blue-500">{link}</a></li>)}
										</ul>
								) : <span className="text-gray-400">{intl.formatMessage({ id: 'knowledge.none' })}</span>;
							} catch { return <span className="text-gray-400">{intl.formatMessage({ id: 'knowledge.none' })}</span>; }
							})()}
						</Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: 'common.column.createdAt' })}>{new Date(selectedEntity.created_at).toLocaleString()}</Descriptions.Item>
					</Descriptions>
				)}
			</Modal>
		</div>
	);
}
