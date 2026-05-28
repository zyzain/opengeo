"use client";

import { useChannels, useCreateChannel } from "@/hooks";
import api from "@/lib/api";
import {
	CheckCircleOutlined,
	CloseCircleOutlined,
	DeleteOutlined,
	EditOutlined,
	GlobalOutlined,
	PlusOutlined,
	SettingOutlined,
} from "@ant-design/icons";
import { useQueryClient } from "@tanstack/react-query";
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

export default function ChannelsPage() {
	const intl = useIntl();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [createForm] = Form.useForm();
	const queryClient = useQueryClient();

	const { data, isLoading } = useChannels();
	const createMutation = useCreateChannel();

	const channels = data?.items || [];

	const handleToggleEnabled = async (record: any) => {
		try {
			if (record.is_enabled) {
				await api.channels.disable(record.id);
			} else {
				await api.channels.enable(record.id);
			}
			message.success(record.is_enabled ? intl.formatMessage({ id: 'common.status.disabled' }) : intl.formatMessage({ id: 'common.status.enabled' }));
			queryClient.invalidateQueries({ queryKey: ["channels"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.channels.delete(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
			queryClient.invalidateQueries({ queryKey: ["channels"] });
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const handleSettings = (record: any) => {
		message.info(intl.formatMessage({ id: 'settings.config.description' }) + `: ${record.channel_name}`);
	};

	const supportedPlatforms = [
		{ value: "wechat", label: intl.formatMessage({ id: 'platform.channels.wechat' }), color: "green", icon: "🟢" },
		{ value: "weibo", label: intl.formatMessage({ id: 'platform.channels.weibo' }), color: "red", icon: "🔴" },
		{ value: "douyin", label: intl.formatMessage({ id: 'platform.channels.douyin' }), color: "purple", icon: "🟣" },
		{ value: "xiaohongshu", label: intl.formatMessage({ id: 'platform.channels.xiaohongshu' }), color: "pink", icon: "🩷" },
		{ value: "zhihu", label: intl.formatMessage({ id: 'platform.channels.zhihu' }), color: "blue", icon: "🔵" },
		{ value: "toutiao", label: intl.formatMessage({ id: 'platform.channels.toutiao' }), color: "orange", icon: "🟠" },
	];

	const getPlatformTag = (platform: string) => {
		const platformInfo = supportedPlatforms.find((p) => p.value === platform);
		return (
			<Tag color={platformInfo?.color || "default"}>
				{platformInfo?.icon} {platformInfo?.label || platform}
			</Tag>
		);
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{
			title: intl.formatMessage({ id: 'channels.column.channelName' }),
			dataIndex: "channel_name",
			key: "channel_name",
			render: (text: string, record: any) => (
				<Space><GlobalOutlined /><span className="font-medium">{text}</span></Space>
			),
		},
		{ title: intl.formatMessage({ id: 'channels.column.platformType' }), dataIndex: "channel_type", key: "channel_type", width: 150, render: (type: string) => getPlatformTag(type) },
		{
			title: intl.formatMessage({ id: 'common.column.status' }),
			dataIndex: "is_enabled",
			key: "is_enabled",
			width: 100,
			render: (enabled: boolean) => (
				<Badge status={enabled ? "success" : "default"} text={enabled ? intl.formatMessage({ id: 'common.status.enabled' }) : intl.formatMessage({ id: 'common.status.disabled' })} />
			),
		},
		{ title: intl.formatMessage({ id: 'common.column.createdAt' }), dataIndex: "created_at", key: "created_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}>
						<Button type="text" icon={<SettingOutlined />} onClick={() => handleSettings(record)} />
					</Tooltip>
					<Tooltip title={record.is_enabled ? intl.formatMessage({ id: 'common.action.disable' }) : intl.formatMessage({ id: 'common.action.enable' })}>
						<Switch checked={record.is_enabled} size="small" checkedChildren={<CheckCircleOutlined />} unCheckedChildren={<CloseCircleOutlined />} onChange={() => handleToggleEnabled(record)} />
					</Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'common.confirmDelete' })} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })} onConfirm={() => handleDelete(record.id)}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}>
							<Button type="text" danger icon={<DeleteOutlined />} />
						</Tooltip>
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

	const stats = {
		total: channels.length,
		enabled: channels.filter((c: any) => c.is_enabled).length,
		disabled: channels.filter((c: any) => !c.is_enabled).length,
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.publish.channels' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'channels.page.subtitle' })}</p>
			</div>

			<div className="grid grid-cols-3 gap-4 mb-4">
				<Card><div className="text-center"><div className="text-2xl font-bold text-gray-800">{stats.total}</div><div className="text-gray-500">{intl.formatMessage({ id: 'channels.stat.total' })}</div></div></Card>
				<Card><div className="text-center"><div className="text-2xl font-bold text-green-500">{stats.enabled}</div><div className="text-gray-500">{intl.formatMessage({ id: 'common.status.enabled' })}</div></div></Card>
				<Card><div className="text-center"><div className="text-2xl font-bold text-gray-400">{stats.disabled}</div><div className="text-gray-500">{intl.formatMessage({ id: 'common.status.disabled' })}</div></div></Card>
			</div>

			<Card title={intl.formatMessage({ id: 'channels.platforms.supported' })} className="mb-4">
				<div className="flex flex-wrap gap-4">
					{supportedPlatforms.map((platform) => (
						<Card key={platform.value} size="small" hoverable className="w-40">
							<div className="text-center">
								<div className="text-2xl mb-2">{platform.icon}</div>
								<div className="font-medium">{platform.label}</div>
								<Tag color={platform.color} className="mt-2">
									{intl.formatMessage({ id: 'channels.platforms.count' }, { count: channels.filter((c: any) => c.channel_type === platform.value).length })}
								</Tag>
							</div>
						</Card>
					))}
				</div>
			</Card>

			<Card
				title={<Space><GlobalOutlined /><span>{intl.formatMessage({ id: 'channels.card.channelList' })}</span></Space>}
				extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'channels.action.addChannel' })}</Button>}
			>
				<Table
					columns={columns}
					dataSource={channels}
					rowKey="id"
					loading={isLoading}
					pagination={{
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }),
					}}
				/>
			</Card>

		<Modal title={intl.formatMessage({ id: 'channels.modal.addTitle' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={500}>
			<Form form={createForm} layout="vertical" onFinish={handleCreate}>
				<Form.Item name="channel_type" label={intl.formatMessage({ id: 'channels.form.platformType' })} rules={[{ required: true, message: intl.formatMessage({ id: 'channels.validation.selectPlatform' }) }]}>
					<Select placeholder={intl.formatMessage({ id: 'channels.placeholder.selectPlatform' })}>
							{supportedPlatforms.map((platform) => (
								<Option key={platform.value} value={platform.value}>
									<Space><span>{platform.icon}</span><span>{platform.label}</span></Space>
								</Option>
							))}
						</Select>
					</Form.Item>
				<Form.Item name="channel_name" label={intl.formatMessage({ id: 'channels.form.channelName' })} rules={[{ required: true, message: intl.formatMessage({ id: 'channels.validation.enterChannelName' }) }]}>
					<Input placeholder={intl.formatMessage({ id: 'channels.placeholder.channelName' })} />
				</Form.Item>
				<Form.Item name="channel_config" label={intl.formatMessage({ id: 'channels.form.channelConfig' })} rules={[{ required: true, message: intl.formatMessage({ id: 'channels.validation.enterChannelConfig' }) }]}>
					<TextArea rows={4} placeholder={intl.formatMessage({ id: 'channels.placeholder.channelConfig' })} />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit" loading={createMutation.isPending}>{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
