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
import { useIntl } from "react-intl";

const { Option } = Select;
const { TextArea } = Input;

export default function TemplatesPage() {
	const intl = useIntl();
	const queryClient = useQueryClient();

	const templateTypes = [
		{ value: "prompt", label: intl.formatMessage({ id: 'templates.type.promptLabel' }), color: "blue", icon: <ThunderboltOutlined /> },
		{ value: "schema", label: intl.formatMessage({ id: 'templates.type.schemaLabel' }), color: "green", icon: <CodeOutlined /> },
		{ value: "layout", label: intl.formatMessage({ id: 'templates.type.layoutLabel' }), color: "purple", icon: <LayoutOutlined /> },
	];

	const tagColors: Record<string, string> = {
		article: "blue",
		video: "purple",
		image: "orange",
	};

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

	const getTemplateTypeTag = (type: string) => {
		const typeInfo = templateTypes.find((t) => t.value === type);
		return <Tag color={typeInfo?.color || "default"} icon={typeInfo?.icon}>{typeInfo?.label || type}</Tag>;
	};

	const filteredTemplates = templates.filter((template: any) => {
		const matchSearch = template.name.toLowerCase().includes(searchText.toLowerCase()) || template.description.toLowerCase().includes(searchText.toLowerCase());
		const matchTab = activeTab === "all" || template.template_type === activeTab;
		return matchSearch && matchTab;
	});

	const stats = {
		total: templates.length,
		prompt: templates.filter((t: any) => t.template_type === "prompt").length,
		schema: templates.filter((t: any) => t.template_type === "schema").length,
		layout: templates.filter((t: any) => t.template_type === "layout").length,
		public: templates.filter((t: any) => t.is_public).length,
		totalUsage: templates.reduce((sum: number, t: any) => sum + (t.usage_count || 0), 0),
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{
			title: intl.formatMessage({ id: 'templates.column.name' }),
			dataIndex: "name",
			key: "name",
			render: (text: string, record: any) => <a onClick={() => handleShowDetail(record)} className="text-blue-500 font-medium">{text}</a>,
		},
		{ title: intl.formatMessage({ id: 'common.column.type' }), dataIndex: "template_type", key: "template_type", width: 120, render: (type: string) => getTemplateTypeTag(type) },
		{
			title: intl.formatMessage({ id: 'common.column.ratings' }),
			dataIndex: "rating",
			key: "rating",
			width: 150,
			render: (rating: number) => <Space><Rate disabled defaultValue={rating} allowHalf style={{ fontSize: 14 }} /><span className="text-gray-500">{rating}</span></Space>,
		},
		{ title: intl.formatMessage({ id: 'templates.column.usageCount' }), dataIndex: "usage_count", key: "usage_count", width: 100, render: (count: number) => <Space><DownloadOutlined /><span>{count}</span></Space> },
		{ title: intl.formatMessage({ id: 'common.column.author' }), dataIndex: "author", key: "author", width: 120 },
		{ title: intl.formatMessage({ id: 'common.column.tags' }), dataIndex: "tags", key: "tags", render: (tags: string[]) => <Space size={[0, 4]} wrap>{tags.map((tag) => <Tag key={tag} color={tagColors[tag] || "default"}>{tag}</Tag>)}</Space> },
		{
			title: intl.formatMessage({ id: 'common.column.visibility' }),
			dataIndex: "is_public",
			key: "is_public",
			width: 80,
			render: (isPublic: boolean) => isPublic ? <Badge status="success" text={intl.formatMessage({ id: 'common.visibility.public' })} /> : <Badge status="default" text={intl.formatMessage({ id: 'common.visibility.private' })} />,
		},
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.view' })}><Button type="text" icon={<FileTextOutlined />} onClick={() => handleShowDetail(record)} /></Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => handleEdit(record)} /></Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'common.action.share' })}><Button type="text" icon={<ShareAltOutlined />} onClick={() => handleShare(record)} /></Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'templates.confirmDelete' })} onConfirm={() => handleDelete(record.id)} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	const handleShowDetail = (record: any) => {
		setSelectedTemplate(record);
		setDetailModalVisible(true);
	};

	const handleCopy = (record: any) => {
		navigator.clipboard.writeText(record.template_data);
		message.success(intl.formatMessage({ id: 'templates.message.copied' }));
	};

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
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: ["templates"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.createFailed' }));
		}
	};

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
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingTemplate(null);
			queryClient.invalidateQueries({ queryKey: ["templates"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.templates.delete(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
			queryClient.invalidateQueries({ queryKey: ["templates"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const handleShare = (record: any) => {
		navigator.clipboard.writeText(record.template_data);
		message.success(intl.formatMessage({ id: 'templates.message.shared' }));
	};

	const handleUseTemplate = (record: any) => {
		navigator.clipboard.writeText(record.template_data);
		message.success(intl.formatMessage({ id: 'templates.message.copiedForUse' }));
	};

	const popularTemplates = [...templates].sort((a, b) => b.usage_count - a.usage_count).slice(0, 5);
	const topRatedTemplates = [...templates].sort((a, b) => b.rating - a.rating).slice(0, 5);

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'templates.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'templates.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'templates.stat.total' })} value={stats.total} prefix={<BookOutlined />} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'templates.stat.prompt' })} value={stats.prompt} valueStyle={{ color: "#1890ff" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'templates.stat.schema' })} value={stats.schema} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'templates.stat.layout' })} value={stats.layout} valueStyle={{ color: "#722ed1" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'templates.stat.public' })} value={stats.public} valueStyle={{ color: "#fa8c16" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'templates.stat.totalUsage' })} value={stats.totalUsage} prefix={<DownloadOutlined />} /></Card></Col>
			</Row>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={12}>
					<Card title={intl.formatMessage({ id: 'templates.section.popular' })} size="small">
						<List
							dataSource={popularTemplates}
							renderItem={(item, index) => (
								<List.Item>
									<List.Item.Meta
										avatar={<Avatar style={{ backgroundColor: index < 3 ? "#1890ff" : "#d9d9d9" }}>{index + 1}</Avatar>}
										title={<Space><span>{item.name}</span>{getTemplateTypeTag(item.template_type)}</Space>}
										description={<Space><span><DownloadOutlined /> {item.usage_count}</span><span><StarOutlined /> {item.rating}</span></Space>}
									/>
								</List.Item>
							)}
						/>
					</Card>
				</Col>
				<Col xs={24} lg={12}>
					<Card title={intl.formatMessage({ id: 'templates.section.topRated' })} size="small">
						<List
							dataSource={topRatedTemplates}
							renderItem={(item, index) => (
								<List.Item>
									<List.Item.Meta
										avatar={<Avatar style={{ backgroundColor: index < 3 ? "#52c41a" : "#d9d9d9" }}>{index + 1}</Avatar>}
										title={<Space><span>{item.name}</span>{getTemplateTypeTag(item.template_type)}</Space>}
										description={<Space><Rate disabled defaultValue={item.rating} allowHalf style={{ fontSize: 12 }} /><span>{item.rating}</span></Space>}
									/>
								</List.Item>
							)}
						/>
					</Card>
				</Col>
			</Row>

			<Card
				title={<Space><BookOutlined /><span>{intl.formatMessage({ id: 'templates.section.list' })}</span></Space>}
				extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'templates.action.create' })}</Button>}
			>
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={[
						{ key: "all", label: intl.formatMessage({ id: 'templates.tab.all' }, { count: stats.total }) },
						...templateTypes.map((type) => ({
							key: type.value,
							label: <Space>{type.icon}<span>{type.label}</span><Badge count={templates.filter((t: any) => t.template_type === type.value).length} /></Space>,
						})),
					]}
					className="mb-4"
				/>

				<div className="mb-4">
					<Input.Search placeholder={intl.formatMessage({ id: 'templates.placeholder.search' })} allowClear style={{ width: 300 }} onSearch={setSearchText} onChange={(e) => setSearchText(e.target.value)} />
				</div>

				<Table columns={columns} dataSource={filteredTemplates} rowKey="id" loading={isLoading} pagination={{ showSizeChanger: true, showQuickJumper: true, showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }) }} />
			</Card>

			<Modal title={intl.formatMessage({ id: 'templates.modal.create' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={700}>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item name="name" label={intl.formatMessage({ id: 'templates.form.name' })} rules={[{ required: true, message: intl.formatMessage({ id: 'templates.validation.enterName' }) }]}><Input placeholder={intl.formatMessage({ id: 'templates.placeholder.name' })} /></Form.Item>
					<Form.Item name="template_type" label={intl.formatMessage({ id: 'templates.form.type' })} rules={[{ required: true, message: intl.formatMessage({ id: 'templates.validation.selectType' }) }]}>
						<Select placeholder={intl.formatMessage({ id: 'templates.placeholder.type' })}>
							{templateTypes.map((type) => <Option key={type.value} value={type.value}><Space>{type.icon}<span>{type.label}</span></Space></Option>)}
						</Select>
					</Form.Item>
					<Form.Item name="description" label={intl.formatMessage({ id: 'common.column.description' })} rules={[{ required: true, message: intl.formatMessage({ id: 'common.validation.enterDescription' }) }]}><TextArea rows={2} placeholder={intl.formatMessage({ id: 'templates.placeholder.description' })} /></Form.Item>
					<Form.Item name="template_data" label={intl.formatMessage({ id: 'templates.form.content' })} rules={[{ required: true, message: intl.formatMessage({ id: 'templates.validation.enterContent' }) }]}><TextArea rows={6} placeholder='{"prompt": "...", "variables": [...]}' /></Form.Item>
					<Form.Item name="tags" label={intl.formatMessage({ id: 'common.column.tags' })}>
						<Select mode="tags" placeholder={intl.formatMessage({ id: 'common.placeholder.tagsInput' })}>
							{Object.keys(tagColors).map((tag) => <Option key={tag} value={tag}>{tag}</Option>)}
						</Select>
					</Form.Item>
					<Form.Item name="is_public" label={intl.formatMessage({ id: 'common.column.visibility' })} initialValue={true}>
						<Select>
							<Option value={true}>{intl.formatMessage({ id: 'templates.option.publicAll' })}</Option>
							<Option value={false}>{intl.formatMessage({ id: 'templates.option.privateOnly' })}</Option>
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title={intl.formatMessage({ id: 'templates.modal.detail' })} open={detailModalVisible} onCancel={() => setDetailModalVisible(false)} footer={null} width={700}>
				{selectedTemplate && (
					<div>
						<Descriptions column={2} bordered className="mb-4">
							<Descriptions.Item label={intl.formatMessage({ id: 'templates.desc.name' })} span={2}>{selectedTemplate.name}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'templates.desc.type' })}>{getTemplateTypeTag(selectedTemplate.template_type)}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'templates.desc.author' })}>{selectedTemplate.author}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'templates.desc.rating' })}><Rate disabled defaultValue={selectedTemplate.rating} allowHalf /> {selectedTemplate.rating}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'templates.desc.usageCount' })}>{selectedTemplate.usage_count}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'templates.desc.visibility' })}>{selectedTemplate.is_public ? <Badge status="success" text={intl.formatMessage({ id: 'common.visibility.public' })} /> : <Badge status="default" text={intl.formatMessage({ id: 'common.visibility.private' })} />}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'templates.desc.tags' })} span={2}><Space size={[0, 4]} wrap>{selectedTemplate.tags.map((tag: string) => <Tag key={tag} color={tagColors[tag] || "default"}>{tag}</Tag>)}</Space></Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'templates.desc.description' })} span={2}>{selectedTemplate.description}</Descriptions.Item>
						</Descriptions>

						<Card title={intl.formatMessage({ id: 'templates.desc.content' })} size="small">
							<pre className="bg-gray-50 p-4 rounded text-sm overflow-auto max-h-60">{JSON.stringify(JSON.parse(selectedTemplate.template_data), null, 2)}</pre>
						</Card>

						<div className="mt-4 flex justify-end">
							<Space>
								<Button icon={<CopyOutlined />} onClick={() => handleCopy(selectedTemplate)}>{intl.formatMessage({ id: 'templates.action.copy' })}</Button>
								<Button icon={<ShareAltOutlined />} onClick={() => handleShare(selectedTemplate)}>{intl.formatMessage({ id: 'templates.action.share' })}</Button>
								<Button type="primary" icon={<DownloadOutlined />} onClick={() => handleUseTemplate(selectedTemplate)}>{intl.formatMessage({ id: 'templates.action.use' })}</Button>
							</Space>
						</div>
					</div>
				)}
			</Modal>

			<Modal title={intl.formatMessage({ id: 'templates.modal.edit' })} open={editModalVisible} onCancel={() => { setEditModalVisible(false); setEditingTemplate(null); }} footer={null} width={700}>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item name="name" label={intl.formatMessage({ id: 'templates.form.name' })} rules={[{ required: true, message: intl.formatMessage({ id: 'templates.validation.enterName' }) }]}><Input placeholder={intl.formatMessage({ id: 'templates.placeholder.name' })} /></Form.Item>
					<Form.Item name="template_type" label={intl.formatMessage({ id: 'templates.form.type' })} rules={[{ required: true, message: intl.formatMessage({ id: 'templates.validation.selectType' }) }]}>
						<Select placeholder={intl.formatMessage({ id: 'templates.placeholder.type' })}>
							{templateTypes.map((type) => <Option key={type.value} value={type.value}><Space>{type.icon}<span>{type.label}</span></Space></Option>)}
						</Select>
					</Form.Item>
					<Form.Item name="description" label={intl.formatMessage({ id: 'common.column.description' })} rules={[{ required: true, message: intl.formatMessage({ id: 'common.validation.enterDescription' }) }]}><TextArea rows={2} placeholder={intl.formatMessage({ id: 'templates.placeholder.description' })} /></Form.Item>
					<Form.Item name="template_data" label={intl.formatMessage({ id: 'templates.form.content' })} rules={[{ required: true, message: intl.formatMessage({ id: 'templates.validation.enterContent' }) }]}><TextArea rows={6} placeholder='{"prompt": "...", "variables": [...]}' /></Form.Item>
					<Form.Item name="tags" label={intl.formatMessage({ id: 'common.column.tags' })}>
						<Select mode="tags" placeholder={intl.formatMessage({ id: 'common.placeholder.tagsInput' })}>
							{Object.keys(tagColors).map((tag) => <Option key={tag} value={tag}>{tag}</Option>)}
						</Select>
					</Form.Item>
					<Form.Item name="is_public" label={intl.formatMessage({ id: 'common.column.visibility' })}>
						<Select>
							<Option value={true}>{intl.formatMessage({ id: 'templates.option.publicAll' })}</Option>
							<Option value={false}>{intl.formatMessage({ id: 'templates.option.privateOnly' })}</Option>
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.save' })}</Button>
							<Button onClick={() => { setEditModalVisible(false); setEditingTemplate(null); }}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
