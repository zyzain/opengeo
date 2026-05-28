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
import { useIntl } from "react-intl";

const { Option } = Select;
const { TextArea } = Input;

export default function AccountGroupsPage() {
	const intl = useIntl();
	const queryClient = useQueryClient();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [editingGroup, setEditingGroup] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();

	const { data: groupsData, isLoading } = useAccountGroups();
	const createMutation = useCreateAccountGroup();

	const groups = groupsData?.items || [];

	const groupTypes = [
		{ value: "authority", label: intl.formatMessage({ id: 'groups.type.authority' }), color: "red", description: intl.formatMessage({ id: 'groups.type.authorityDesc' }) },
		{ value: "professional", label: intl.formatMessage({ id: 'groups.type.professional' }), color: "blue", description: intl.formatMessage({ id: 'groups.type.professionalDesc' }) },
		{ value: "ecology", label: intl.formatMessage({ id: 'groups.type.ecology' }), color: "green", description: intl.formatMessage({ id: 'groups.type.ecologyDesc' }) },
	];

	const getGroupTypeTag = (type: string) => {
		const typeInfo = groupTypes.find((t) => t.value === type);
		return <Tag color={typeInfo?.color || "default"}>{typeInfo?.label || type}</Tag>;
	};

	const convertToTreeData = (groups: any[]): any[] => {
		return groups.map((group) => ({
			key: group.id,
			title: (
				<div className="flex items-center space-x-2">
					<span>{group.name}</span>
					{getGroupTypeTag(group.group_type)}
					<span className="text-gray-400 text-xs">({intl.formatMessage({ id: 'groups.tree.accountCount' }, { count: group.account_count || 0 })})</span>
				</div>
			),
			icon: group.children?.length ? <FolderOpenOutlined /> : <FolderOutlined />,
			children: group.children ? convertToTreeData(group.children) : [],
		}));
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{ title: intl.formatMessage({ id: 'groups.column.groupName' }), dataIndex: "name", key: "name", render: (text: string, record: any) => <Space><ApartmentOutlined /><span className="font-medium">{text}</span></Space> },
		{ title: intl.formatMessage({ id: 'groups.column.groupType' }), dataIndex: "group_type", key: "group_type", width: 150, render: (type: string) => getGroupTypeTag(type) },
		{ title: intl.formatMessage({ id: 'common.column.description' }), dataIndex: "description", key: "description", ellipsis: true },
		{ title: intl.formatMessage({ id: 'groups.column.accountCount' }), dataIndex: "account_count", key: "account_count", width: 100, render: (count: number) => <Tag color="blue">{count || 0}</Tag> },
		{ title: intl.formatMessage({ id: 'common.column.createdAt' }), dataIndex: "created_at", key: "created_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => handleEdit(record)} /></Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'common.confirmDelete' })} onConfirm={() => handleDelete(record.id)} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

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

	const handleEdit = (record: any) => {
		setEditingGroup(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await api.accountGroups.update(editingGroup.id, values);
			queryClient.invalidateQueries({ queryKey: ['accountGroups'] });
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingGroup(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.accountGroups.delete(id);
			queryClient.invalidateQueries({ queryKey: ['accountGroups'] });
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const stats: Record<string, number> = {
		total: groups.length,
		authority: groups.filter((g: any) => g.group_type === "authority").length,
		professional: groups.filter((g: any) => g.group_type === "professional").length,
		ecology: groups.filter((g: any) => g.group_type === "ecology").length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.account.groups' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'groups.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'groups.stat.total' })} value={stats.total} prefix={<ApartmentOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'groups.stat.authority' })} value={stats.authority} valueStyle={{ color: "#ff4d4f" }} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'groups.stat.professional' })} value={stats.professional} valueStyle={{ color: "#1890ff" }} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'groups.stat.ecology' })} value={stats.ecology} valueStyle={{ color: "#52c41a" }} /></Card></Col>
			</Row>

			<Row gutter={[16, 16]}>
				<Col xs={24} lg={8}>
				<Card
					title={intl.formatMessage({ id: 'groups.card.structure' })}
					extra={<Button type="primary" size="small" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'groups.action.createGroup' })}</Button>}
				>
					{groups.length > 0 ? <Tree showIcon defaultExpandAll treeData={convertToTreeData(groups)} className="bg-gray-50 p-4 rounded-lg" /> : <Empty description={intl.formatMessage({ id: 'groups.tree.empty' })} />}
					</Card>
				</Col>

				<Col xs={24} lg={16}>
					<Card title={intl.formatMessage({ id: 'groups.card.list' })}>
						<Table columns={columns} dataSource={groups} rowKey="id" loading={isLoading} pagination={false} />
					</Card>
				</Col>
			</Row>

		<Modal title={intl.formatMessage({ id: 'groups.modal.createTitle' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={500}>
			<Form form={createForm} layout="vertical" onFinish={handleCreate}>
				<Form.Item name="name" label={intl.formatMessage({ id: 'groups.form.groupName' })} rules={[{ required: true, message: intl.formatMessage({ id: 'groups.validation.enterGroupName' }) }]}><Input placeholder={intl.formatMessage({ id: 'groups.placeholder.groupName' })} /></Form.Item>
				<Form.Item name="group_type" label={intl.formatMessage({ id: 'groups.form.groupType' })} rules={[{ required: true, message: intl.formatMessage({ id: 'groups.validation.selectGroupType' }) }]}>
					<Select placeholder={intl.formatMessage({ id: 'groups.placeholder.groupType' })}>
							{groupTypes.map((type) => <Option key={type.value} value={type.value}><Space><Tag color={type.color}>{type.label}</Tag><span className="text-gray-400">{type.description}</span></Space></Option>)}
						</Select>
					</Form.Item>
				<Form.Item name="parent_id" label={intl.formatMessage({ id: 'groups.form.parentGroup' })}>
					<Select placeholder={intl.formatMessage({ id: 'groups.placeholder.parentGroup' })} allowClear>
							{groups.map((group: any) => <Option key={group.id} value={group.id}>{group.name}</Option>)}
						</Select>
					</Form.Item>
				<Form.Item name="description" label={intl.formatMessage({ id: 'common.column.description' })}><TextArea rows={3} placeholder={intl.formatMessage({ id: 'groups.placeholder.description' })} /></Form.Item>
				<Form.Item>
					<Space>
						<Button type="primary" htmlType="submit" loading={createMutation.isPending}>{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

		<Modal title={intl.formatMessage({ id: 'groups.modal.editTitle' })} open={editModalVisible} onCancel={() => { setEditModalVisible(false); setEditingGroup(null); }} footer={null} width={500}>
			<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
				<Form.Item name="name" label={intl.formatMessage({ id: 'groups.form.groupName' })} rules={[{ required: true, message: intl.formatMessage({ id: 'groups.validation.enterGroupName' }) }]}><Input placeholder={intl.formatMessage({ id: 'groups.placeholder.groupName' })} /></Form.Item>
				<Form.Item name="group_type" label={intl.formatMessage({ id: 'groups.form.groupType' })} rules={[{ required: true, message: intl.formatMessage({ id: 'groups.validation.selectGroupType' }) }]}>
					<Select placeholder={intl.formatMessage({ id: 'groups.placeholder.groupType' })}>
							{groupTypes.map((type) => <Option key={type.value} value={type.value}><Space><Tag color={type.color}>{type.label}</Tag><span className="text-gray-400">{type.description}</span></Space></Option>)}
						</Select>
					</Form.Item>
				<Form.Item name="description" label={intl.formatMessage({ id: 'common.column.description' })}><TextArea rows={3} placeholder={intl.formatMessage({ id: 'groups.placeholder.description' })} /></Form.Item>
				<Form.Item>
					<Space>
						<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.save' })}</Button>
							<Button onClick={() => { setEditModalVisible(false); setEditingGroup(null); }}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
