"use client";

import { useSystemConfigs } from "@/hooks";
import api from "@/lib/api";
import {
	CopyOutlined,
	DeleteOutlined,
	EditOutlined,
	GlobalOutlined,
	LockOutlined,
	PlusOutlined,
	ReloadOutlined,
	SaveOutlined,
	SettingOutlined,
} from "@ant-design/icons";
import {
	Badge,
	Button,
	Card,
	Descriptions,
	Form,
	Input,
	Modal,
	Popconfirm,
	Select,
	Space,
	Switch,
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

export default function ConfigsPage() {
	const intl = useIntl();
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [editingConfig, setEditingConfig] = useState<any>(null);
	const [editForm] = Form.useForm();

	const { data, isLoading, refetch } = useSystemConfigs();
	const configs = data?.items || [];

	const configTypes = [
		{ value: "string", label: intl.formatMessage({ id: "config.type.string" }), color: "blue" },
		{ value: "number", label: intl.formatMessage({ id: "config.type.number" }), color: "green" },
		{ value: "json", label: "JSON", color: "purple" },
		{ value: "boolean", label: intl.formatMessage({ id: "config.type.boolean" }), color: "orange" },
	];

	const getConfigTypeTag = (type: string) => {
		const typeInfo = configTypes.find((t) => t.value === type);
		return <Tag color={typeInfo?.color || "default"}>{typeInfo?.label || type}</Tag>;
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{ title: intl.formatMessage({ id: 'settings.config.key' }), dataIndex: "config_key", key: "config_key", render: (text: string) => <code className="bg-gray-100 px-2 py-1 rounded text-sm">{text}</code> },
		{
			title: intl.formatMessage({ id: 'settings.config.value' }),
			dataIndex: "config_value",
			key: "config_value",
			ellipsis: true,
			render: (text: string, record: any) => {
				if (record.config_type === "boolean") return <Switch checked={text === "true"} disabled size="small" />;
				if (record.config_type === "json") return <Tooltip title={<pre>{JSON.stringify(JSON.parse(text), null, 2)}</pre>}><span className="text-blue-500 cursor-pointer">{text.substring(0, 50)}...</span></Tooltip>;
				return text;
			},
		},
		{ title: intl.formatMessage({ id: 'settings.config.type' }), dataIndex: "config_type", key: "config_type", width: 100, render: (type: string) => getConfigTypeTag(type) },
		{ title: intl.formatMessage({ id: 'settings.config.description' }), dataIndex: "description", key: "description", ellipsis: true },
		{ title: intl.formatMessage({ id: "config.column.public" }), dataIndex: "is_public", key: "is_public", width: 80, render: (isPublic: boolean) => isPublic ? <Badge status="success" text={intl.formatMessage({ id: "config.yes" })} /> : <Badge status="default" text={intl.formatMessage({ id: "config.no" })} /> },
		{ title: intl.formatMessage({ id: 'common.column.updatedAt' }), dataIndex: "updated_at", key: "updated_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => handleEdit(record)} /></Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'common.action.share' })}><Button type="text" icon={<CopyOutlined />} onClick={() => handleCopy(record)} /></Tooltip>
				</Space>
			),
		},
	];

	const handleEdit = (record: any) => {
		setEditingConfig(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await api.system.updateConfig(values.config_key, values.config_value);
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingConfig(null);
			refetch();
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	const handleCopy = (record: any) => {
		navigator.clipboard.writeText(`${record.config_key}=${record.config_value}`);
		message.success(intl.formatMessage({ id: 'common.message.exportSuccess' }));
	};

	const configCategories = [
		{ title: intl.formatMessage({ id: "config.tab.system" }), icon: <SettingOutlined />, configs: configs.filter((c: any) => c.config_key.startsWith("system.")) },
		{ title: intl.formatMessage({ id: "config.tab.publish" }), icon: <GlobalOutlined />, configs: configs.filter((c: any) => c.config_key.startsWith("publish.")) },
		{ title: intl.formatMessage({ id: "config.tab.ai" }), icon: <LockOutlined />, configs: configs.filter((c: any) => c.config_key.startsWith("ai.")) },
		{ title: intl.formatMessage({ id: "config.tab.monitor" }), icon: <SettingOutlined />, configs: configs.filter((c: any) => c.config_key.startsWith("monitor.")) },
	];

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.settings.configs' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: "config.page.subtitle" })}</p>
			</div>

			<Tabs
				defaultActiveKey="all"
				items={[
					{
						key: "all",
						label: intl.formatMessage({ id: "config.tab.all" }),
						children: (
							<Card extra={<Space><Button icon={<ReloadOutlined />} onClick={() => refetch()}>{intl.formatMessage({ id: 'common.action.refresh' })}</Button></Space>}>
								<Table columns={columns} dataSource={configs} rowKey="id" loading={isLoading} pagination={{ showSizeChanger: true, showQuickJumper: true, showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }) }} />
							</Card>
						),
					},
					...configCategories.map((category) => ({
						key: category.title,
						label: <Space>{category.icon}<span>{category.title}</span><Badge count={category.configs.length} style={{ backgroundColor: "#1890ff" }} /></Space>,
						children: <Card><Table columns={columns} dataSource={category.configs} rowKey="id" pagination={false} /></Card>,
					})),
				]}
			/>

			<Modal title={intl.formatMessage({ id: "config.modal.edit" })} open={editModalVisible} onCancel={() => { setEditModalVisible(false); setEditingConfig(null); }} footer={null} width={500}>
				<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
					<Form.Item name="config_key" label={intl.formatMessage({ id: 'settings.config.key' })}><Input disabled /></Form.Item>
				<Form.Item name="config_value" label={intl.formatMessage({ id: 'settings.config.value' })} rules={[{ required: true, message: intl.formatMessage({ id: "config.validation.enterValue" }) }]}>
					{editingConfig?.config_type === "json" ? <TextArea rows={4} placeholder={intl.formatMessage({ id: "config.placeholder.value" })} /> : editingConfig?.config_type === "boolean" ? <Select><Option value="true">{intl.formatMessage({ id: "config.yes" })}</Option><Option value="false">{intl.formatMessage({ id: "config.no" })}</Option></Select> : <Input placeholder={intl.formatMessage({ id: "config.placeholder.value" })} />}
					</Form.Item>
					<Form.Item name="description" label={intl.formatMessage({ id: 'settings.config.description' })}><TextArea rows={2} placeholder={intl.formatMessage({ id: "config.placeholder.description" })} /></Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit" icon={<SaveOutlined />}>{intl.formatMessage({ id: 'common.action.save' })}</Button>
							<Button onClick={() => setEditModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
