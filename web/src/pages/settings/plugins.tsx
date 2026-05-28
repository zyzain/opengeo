"use client";

import { queryKeys, usePlugins } from "@/hooks";
import { api } from "@/lib/api";
import {
	ApiOutlined,
	AppstoreOutlined,
	CheckCircleOutlined,
	CloseCircleOutlined,
	CloudOutlined,
	DeleteOutlined,
	EditOutlined,
	PlusOutlined,
	SettingOutlined,
	ThunderboltOutlined,
} from "@ant-design/icons";
import { useQueryClient } from "@tanstack/react-query";
import {
	Badge,
	Button,
	Card,
	Col,
	Form,
	Input,
	Modal,
	Popconfirm,
	Row,
	Select,
	Space,
	Statistic,
	Switch,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";
import { useIntl } from "react-intl";

const { Option } = Select;
const { TextArea } = Input;

export default function PluginsPage() {
	const intl = useIntl();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [createForm] = Form.useForm();
	const queryClient = useQueryClient();

	const { data, isLoading } = usePlugins();
	const plugins = data?.items || [];

	const handleInstall = async (values: any) => {
		try {
			await api.system.installPlugin(values);
			message.success(intl.formatMessage({ id: 'common.message.addSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: queryKeys.plugins });
		} catch (error: any) {
			message.error(error?.response?.data?.message || intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const handleToggleEnabled = async (record: any, checked: boolean) => {
		try {
			await api.system.updatePlugin(record.id, { is_enabled: checked });
			message.success(checked ? intl.formatMessage({ id: 'common.status.enabled' }) : intl.formatMessage({ id: 'common.status.disabled' }));
			queryClient.invalidateQueries({ queryKey: queryKeys.plugins });
		} catch (error: any) {
			message.error(error?.response?.data?.message || intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const handleUninstall = async (id: number) => {
		try {
			await api.system.deletePlugin(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
			queryClient.invalidateQueries({ queryKey: queryKeys.plugins });
		} catch (error: any) {
			message.error(error?.response?.data?.message || intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const handleSettings = (record: any) => {
		Modal.info({
			title: `${record.plugin_name} - ${intl.formatMessage({ id: 'common.action.edit' })}`,
			content: (<div><p>{intl.formatMessage({ id: "plugin.detail.name" })}: {record.plugin_name}</p><p>{intl.formatMessage({ id: "plugin.detail.version" })}: v{record.version}</p><p>{intl.formatMessage({ id: "plugin.detail.author" })}: {record.author}</p><p>{intl.formatMessage({ id: "plugin.detail.description" })}: {record.description || intl.formatMessage({ id: "plugin.detail.noDesc" })}</p></div>),
		});
	};

	const pluginTypes = [
		{ value: "channel", label: intl.formatMessage({ id: "plugin.type.channel" }), color: "blue", icon: <CloudOutlined /> },
		{ value: "ai", label: intl.formatMessage({ id: "plugin.type.ai" }), color: "purple", icon: <ThunderboltOutlined /> },
		{ value: "analyzer", label: intl.formatMessage({ id: "plugin.type.analyzer" }), color: "green", icon: <ApiOutlined /> },
	];

	const getPluginTypeTag = (type: string) => {
		const typeInfo = pluginTypes.find((t) => t.value === type);
		return <Tag color={typeInfo?.color || "default"} icon={typeInfo?.icon}>{typeInfo?.label || type}</Tag>;
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{ title: intl.formatMessage({ id: "plugin.column.name" }), dataIndex: "plugin_name", key: "plugin_name", render: (text: string) => <Space><AppstoreOutlined /><span className="font-medium">{text}</span></Space> },
		{ title: intl.formatMessage({ id: 'common.column.type' }), dataIndex: "plugin_type", key: "plugin_type", width: 120, render: (type: string) => getPluginTypeTag(type) },
		{ title: intl.formatMessage({ id: "plugin.column.version" }), dataIndex: "version", key: "version", width: 100, render: (text: string) => <Tag>v{text}</Tag> },
		{ title: intl.formatMessage({ id: 'common.column.author' }), dataIndex: "author", key: "author", width: 120 },
		{ title: intl.formatMessage({ id: 'common.column.description' }), dataIndex: "description", key: "description", ellipsis: true },
		{ title: intl.formatMessage({ id: 'common.column.status' }), dataIndex: "is_enabled", key: "is_enabled", width: 100, render: (enabled: boolean) => <Badge status={enabled ? "success" : "default"} text={enabled ? intl.formatMessage({ id: 'common.status.enabled' }) : intl.formatMessage({ id: 'common.status.disabled' })} /> },
		{ title: intl.formatMessage({ id: "plugin.column.installedAt" }), dataIndex: "installed_at", key: "installed_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<SettingOutlined />} onClick={() => handleSettings(record)} /></Tooltip>
					<Tooltip title={record.is_enabled ? intl.formatMessage({ id: 'common.action.disable' }) : intl.formatMessage({ id: 'common.action.enable' })}>
						<Switch checked={record.is_enabled} size="small" checkedChildren={<CheckCircleOutlined />} unCheckedChildren={<CloseCircleOutlined />} onChange={(checked) => handleToggleEnabled(record, checked)} />
					</Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'common.confirmDelete' })} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })} onConfirm={() => handleUninstall(record.id)}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	const stats = {
		total: plugins.length,
		enabled: plugins.filter((p: any) => p.is_enabled).length,
		channel: plugins.filter((p: any) => p.plugin_type === "channel").length,
		ai: plugins.filter((p: any) => p.plugin_type === "ai").length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.settings.plugins' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: "plugin.page.subtitle" })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: "plugin.stat.total" })} value={stats.total} prefix={<AppstoreOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'common.status.enabled' })} value={stats.enabled} prefix={<CheckCircleOutlined />} valueStyle={{ color: "#52c41a" }} /></Card></Col>
			<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: "plugin.stat.channel" })} value={stats.channel} prefix={<CloudOutlined />} valueStyle={{ color: "#1890ff" }} /></Card></Col>
			<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: "plugin.stat.ai" })} value={stats.ai} prefix={<ThunderboltOutlined />} valueStyle={{ color: "#722ed1" }} /></Card></Col>
			</Row>

			<Card title={intl.formatMessage({ id: "plugin.section.types" })} className="mb-4">
				<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
					{pluginTypes.map((type) => (
						<Card key={type.value} size="small" hoverable>
							<div className="flex items-center space-x-3">
								<div className="text-2xl">{type.icon}</div>
								<div>
									<div className="font-medium">{type.label}</div>
									<div className="text-gray-500 text-sm">
								{type.value === "channel" && intl.formatMessage({ id: "plugin.type.channel.desc" })}
									{type.value === "ai" && intl.formatMessage({ id: "plugin.type.ai.desc" })}
									{type.value === "analyzer" && intl.formatMessage({ id: "plugin.type.analyzer.desc" })}
									</div>
								</div>
							</div>
						</Card>
					))}
				</div>
			</Card>

			<Card
			title={<Space><AppstoreOutlined /><span>{intl.formatMessage({ id: "plugin.section.list" })}</span></Space>}
			extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: "plugin.action.install" })}</Button>}
			>
				<Table columns={columns} dataSource={plugins} rowKey="id" loading={isLoading} pagination={{ showSizeChanger: true, showQuickJumper: true, showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }) }} />
			</Card>

			<Modal title={intl.formatMessage({ id: "plugin.modal.install" })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={500}>
				<Form form={createForm} layout="vertical" onFinish={handleInstall}>
				<Form.Item name="plugin_name" label={intl.formatMessage({ id: "plugin.form.name" })} rules={[{ required: true, message: intl.formatMessage({ id: "plugin.validation.enterName" }) }]}><Input placeholder={intl.formatMessage({ id: "plugin.placeholder.name" })} /></Form.Item>
				<Form.Item name="plugin_type" label={intl.formatMessage({ id: 'common.column.type' })} rules={[{ required: true, message: intl.formatMessage({ id: "plugin.validation.selectType" }) }]}>
					<Select placeholder={intl.formatMessage({ id: "plugin.placeholder.type" })}>
							{pluginTypes.map((type) => <Option key={type.value} value={type.value}><Space>{type.icon}<span>{type.label}</span></Space></Option>)}
						</Select>
					</Form.Item>
				<Form.Item name="version" label={intl.formatMessage({ id: "plugin.form.version" })}><Input placeholder={intl.formatMessage({ id: "plugin.placeholder.version" })} /></Form.Item>
				<Form.Item name="author" label={intl.formatMessage({ id: 'common.column.author' })}><Input placeholder={intl.formatMessage({ id: "plugin.placeholder.author" })} /></Form.Item>
				<Form.Item name="description" label={intl.formatMessage({ id: 'common.column.description' })}><TextArea rows={3} placeholder={intl.formatMessage({ id: "plugin.placeholder.description" })} /></Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: "plugin.action.install" })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
