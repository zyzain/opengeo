"use client";

import {
	useAccounts,
	useCreateAccount,
	useDeleteAccount,
	useUpdateAccount,
} from "@/hooks";
import api from "@/lib/api";
import {
	CheckCircleOutlined,
	DeleteOutlined,
	EditOutlined,
	PlusOutlined,
	ReloadOutlined,
	SearchOutlined,
	TeamOutlined,
} from "@ant-design/icons";
import {
	Badge,
	Button,
	Card,
	Form,
	Input,
	Modal,
	Popconfirm,
	Progress,
	Select,
	Space,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";
import { useIntl } from "react-intl";

const { Option } = Select;

export default function AccountListPage() {
	const intl = useIntl();
	const [searchForm] = Form.useForm();
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [editModalVisible, setEditModalVisible] = useState(false);
	const [editingAccount, setEditingAccount] = useState<any>(null);
	const [createForm] = Form.useForm();
	const [editForm] = Form.useForm();

	const [queryParams, setQueryParams] = useState({
		page: 1,
		page_size: 10,
		platform: undefined,
	});

	const { data, isLoading, refetch } = useAccounts(queryParams);
	const createMutation = useCreateAccount();
	const updateMutation = useUpdateAccount();
	const deleteMutation = useDeleteAccount();

	const accounts = data?.items || [];
	const total = data?.total || 0;

	const platformTypes = [
		{ value: "wechat", label: intl.formatMessage({ id: 'platform.channels.wechat' }), color: "green" },
		{ value: "weibo", label: intl.formatMessage({ id: 'platform.channels.weibo' }), color: "red" },
		{ value: "douyin", label: intl.formatMessage({ id: 'platform.channels.douyin' }), color: "purple" },
		{ value: "xiaohongshu", label: intl.formatMessage({ id: 'platform.channels.xiaohongshu' }), color: "pink" },
		{ value: "zhihu", label: intl.formatMessage({ id: 'platform.channels.zhihu' }), color: "blue" },
		{ value: "toutiao", label: intl.formatMessage({ id: 'platform.channels.toutiao' }), color: "orange" },
	];

	const getPlatformTag = (platform: string) => {
		const platformInfo = platformTypes.find((p) => p.value === platform);
		return <Tag color={platformInfo?.color || "default"}>{platformInfo?.label || platform}</Tag>;
	};

	const getStatusTag = (status: number) => {
		const statusMap: Record<number, { color: string; text: string }> = {
			0: { color: "error", text: intl.formatMessage({ id: 'common.status.disabled' }) },
			1: { color: "success", text: intl.formatMessage({ id: 'common.status.active' }) },
			2: { color: "warning", text: intl.formatMessage({ id: 'account.status.abnormal' }) },
		};
		const config = statusMap[status] || { color: "default", text: intl.formatMessage({ id: 'common.status.unknown' }) };
		return <Badge status={config.color as any} text={config.text} />;
	};

	const getHealthColor = (score: number) => {
		if (score >= 80) return "#52c41a";
		if (score >= 60) return "#faad14";
		return "#ff4d4f";
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{ title: intl.formatMessage({ id: 'account.name' }), dataIndex: "account_name", key: "account_name", ellipsis: true },
		{ title: intl.formatMessage({ id: 'account.platform' }), dataIndex: "platform", key: "platform", width: 120, render: (platform: string) => getPlatformTag(platform) },
		{ title: intl.formatMessage({ id: 'account.accountId' }), dataIndex: "account_id", key: "account_id", width: 150, ellipsis: true },
		{ title: intl.formatMessage({ id: 'account.status' }), dataIndex: "status", key: "status", width: 100, render: (status: number) => getStatusTag(status) },
		{ title: intl.formatMessage({ id: 'account.health' }), dataIndex: "health_score", key: "health_score", width: 150, render: (score: number) => <Progress percent={score} size="small" strokeColor={getHealthColor(score)} format={(percent) => `${percent}%`} /> },
		{ title: intl.formatMessage({ id: 'common.column.createdAt' }), dataIndex: "created_at", key: "created_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}><Button type="text" icon={<EditOutlined />} onClick={() => handleEdit(record)} /></Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'common.check' })}><Button type="text" icon={<CheckCircleOutlined />} onClick={() => handleHealthCheck(record.id)} /></Tooltip>
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
		setEditingAccount(record);
		editForm.setFieldsValue(record);
		setEditModalVisible(true);
	};

	const handleUpdate = async (values: any) => {
		try {
			await updateMutation.mutateAsync({ id: editingAccount.id, data: values });
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			setEditModalVisible(false);
			editForm.resetFields();
			setEditingAccount(null);
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.updateFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await deleteMutation.mutateAsync(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const handleHealthCheck = async (id: number) => {
		try {
			const res = await api.accounts.health(id);
			const healthData = res.data.data;
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			refetch();
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const handleSearch = (values: any) => {
		setQueryParams({ ...queryParams, ...values, page: 1 });
	};

	const handleReset = () => {
		searchForm.resetFields();
		setQueryParams({ page: 1, page_size: 10, platform: undefined });
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.account.list' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'account.page.subtitle' })}</p>
			</div>

			<Card className="mb-4">
				<Form form={searchForm} layout="inline" onFinish={handleSearch}>
					<Form.Item name="platform" label={intl.formatMessage({ id: 'account.platform' })}>
						<Select placeholder={intl.formatMessage({ id: 'account.placeholder.selectPlatform' })} allowClear style={{ width: 150 }}>
							{platformTypes.map((p) => <Option key={p.value} value={p.value}>{p.label}</Option>)}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" icon={<SearchOutlined />} htmlType="submit">{intl.formatMessage({ id: 'common.action.search' })}</Button>
							<Button onClick={handleReset}>{intl.formatMessage({ id: 'common.action.reset' })}</Button>
							<Button icon={<ReloadOutlined />} onClick={() => refetch()}>{intl.formatMessage({ id: 'common.action.refresh' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Card>

			<Card
				title={<Space><TeamOutlined /><span>{intl.formatMessage({ id: 'account.column.accountList' })}</span><Badge count={total} style={{ backgroundColor: "#1890ff" }} /></Space>}
				extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'account.action.addAccount' })}</Button>}
			>
				<Table columns={columns} dataSource={accounts} rowKey="id" loading={isLoading} pagination={{
					current: queryParams.page,
					pageSize: queryParams.page_size,
					total,
					showSizeChanger: true,
					showQuickJumper: true,
					showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }),
					onChange: (page, pageSize) => setQueryParams({ ...queryParams, page, page_size: pageSize }),
				}} />
			</Card>

		<Modal title={intl.formatMessage({ id: 'account.modal.addTitle' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={500}>
			<Form form={createForm} layout="vertical" onFinish={handleCreate}>
				<Form.Item name="platform" label={intl.formatMessage({ id: 'account.platform' })} rules={[{ required: true, message: intl.formatMessage({ id: 'account.validation.selectPlatform' }) }]}>
					<Select placeholder={intl.formatMessage({ id: 'account.placeholder.selectPlatform' })}>{platformTypes.map((p) => <Option key={p.value} value={p.value}>{p.label}</Option>)}</Select>
				</Form.Item>
				<Form.Item name="account_name" label={intl.formatMessage({ id: 'account.name' })} rules={[{ required: true, message: intl.formatMessage({ id: 'account.validation.enterAccountName' }) }]}><Input placeholder={intl.formatMessage({ id: 'account.placeholder.accountName' })} /></Form.Item>
				<Form.Item name="account_id" label={intl.formatMessage({ id: 'account.accountId' })} rules={[{ required: true, message: intl.formatMessage({ id: 'account.validation.enterAccountId' }) }]}><Input placeholder={intl.formatMessage({ id: 'account.placeholder.accountId' })} /></Form.Item>
				<Form.Item name="credentials" label={intl.formatMessage({ id: 'account.form.credentials' })} tooltip={intl.formatMessage({ id: 'account.form.credentialsTip' })}><Input.TextArea rows={3} placeholder={intl.formatMessage({ id: 'account.placeholder.credentials' })} /></Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit" loading={createMutation.isPending}>{intl.formatMessage({ id: 'common.action.create' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

		<Modal title={intl.formatMessage({ id: 'account.modal.editTitle' })} open={editModalVisible} onCancel={() => { setEditModalVisible(false); setEditingAccount(null); }} footer={null} width={500}>
			<Form form={editForm} layout="vertical" onFinish={handleUpdate}>
				<Form.Item name="account_name" label={intl.formatMessage({ id: 'account.name' })} rules={[{ required: true, message: intl.formatMessage({ id: 'account.validation.enterAccountName' }) }]}><Input placeholder={intl.formatMessage({ id: 'account.placeholder.accountName' })} /></Form.Item>
					<Form.Item name="status" label={intl.formatMessage({ id: 'account.status' })}>
						<Select>
							<Option value={1}>{intl.formatMessage({ id: 'common.status.active' })}</Option>
							<Option value={0}>{intl.formatMessage({ id: 'common.status.disabled' })}</Option>
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit" loading={updateMutation.isPending}>{intl.formatMessage({ id: 'common.action.save' })}</Button>
							<Button onClick={() => setEditModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>
		</div>
	);
}
