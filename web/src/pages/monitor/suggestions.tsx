"use client";

import { useGenerateSuggestions } from "@/hooks";
import api from "@/lib/api";
import {
	BulbOutlined,
	CheckCircleOutlined,
	ClockCircleOutlined,
	DislikeOutlined,
	ExclamationCircleOutlined,
	FilterOutlined,
	LikeOutlined,
	ReloadOutlined,
	StarOutlined,
	ThunderboltOutlined,
} from "@ant-design/icons";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import {
	Avatar,
	Badge,
	Button,
	Card,
	Col,
	Descriptions,
	InputNumber,
	List,
	Modal,
	Progress,
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

const priorityConfig: Record<number, { color: string }> = {
	2: { color: "red" },
	1: { color: "orange" },
	0: { color: "blue" },
};

const statusConfig: Record<number, { color: string; icon: React.ReactNode }> = {
	0: { color: "processing", icon: <ClockCircleOutlined /> },
	1: { color: "success", icon: <CheckCircleOutlined /> },
	2: { color: "default", icon: <ExclamationCircleOutlined /> },
};

export default function SuggestionsPage() {
	const intl = useIntl();

	const suggestionTypes = [
		{ value: "content", label: intl.formatMessage({ id: "suggestions.type.content" }), color: "blue", icon: <BulbOutlined /> },
		{ value: "structure", label: intl.formatMessage({ id: "suggestions.type.structure" }), color: "green", icon: <ThunderboltOutlined /> },
		{ value: "authority", label: intl.formatMessage({ id: "suggestions.type.authority" }), color: "purple", icon: <StarOutlined /> },
	];

	const priorityLabels: Record<number, string> = {
		2: intl.formatMessage({ id: "suggestions.priority.high" }),
		1: intl.formatMessage({ id: "suggestions.priority.medium" }),
		0: intl.formatMessage({ id: "suggestions.priority.low" }),
	};

	const statusLabels: Record<number, string> = {
		0: intl.formatMessage({ id: "suggestions.status.pending" }),
		1: intl.formatMessage({ id: "suggestions.status.applied" }),
		2: intl.formatMessage({ id: "suggestions.status.ignored" }),
	};
	const queryClient = useQueryClient();
	const generateMutation = useGenerateSuggestions();

	const { data: suggestionsData, isLoading, refetch } = useQuery({
		queryKey: ["suggestions"],
		queryFn: () => api.monitor.listSuggestions(),
		select: (response) => response.data.data,
	});

	const suggestions = suggestionsData?.items || [];
	const [detailModalVisible, setDetailModalVisible] = useState(false);
	const [selectedSuggestion, setSelectedSuggestion] = useState<any>(null);
	const [activeTab, setActiveTab] = useState("all");
	const [filterType, setFilterType] = useState<string | undefined>(undefined);
	const [filterPriority, setFilterPriority] = useState<number | undefined>(undefined);
	const [generateContentId, setGenerateContentId] = useState<number>(0);

	const handleGenerateSuggestions = async () => {
		try {
			await generateMutation.mutateAsync(generateContentId);
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			refetch();
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const filteredSuggestions = suggestions.filter((s: any) => {
		const matchTab = activeTab === "all" || (activeTab === "pending" && s.status === 0) || (activeTab === "applied" && s.status === 1) || (activeTab === "ignored" && s.status === 2);
		const matchType = !filterType || s.suggestion_type === filterType;
		const matchPriority = filterPriority === undefined || s.priority === filterPriority;
		return matchTab && matchType && matchPriority;
	});

	const stats = {
		total: suggestions.length,
		pending: suggestions.filter((s: any) => s.status === 0).length,
		applied: suggestions.filter((s: any) => s.status === 1).length,
		ignored: suggestions.filter((s: any) => s.status === 2).length,
		highPriority: suggestions.filter((s: any) => s.priority === 2 && s.status === 0).length,
	};

	const handleShowDetail = (record: any) => {
		setSelectedSuggestion(record);
		setDetailModalVisible(true);
	};

	const handleApply = async (id: number) => {
		try {
			await api.monitor.applySuggestion(id);
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			refetch();
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const handleIgnore = async (id: number) => {
		try {
			await api.monitor.ignoreSuggestion(id);
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
			refetch();
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
		}
	};

	const getTypeTag = (type: string) => {
		const config = suggestionTypes.find((t) => t.value === type);
		return <Tag color={config?.color} icon={config?.icon}>{config?.label}</Tag>;
	};

	const getPriorityTag = (priority: number) => {
		const config = priorityConfig[priority];
		return <Tag color={config?.color}>{priorityLabels[priority]}</Tag>;
	};

	const getStatusTag = (status: number) => {
		const config = statusConfig[status];
		return <Badge status={config?.color as any} text={statusLabels[status]} />;
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
	{ title: intl.formatMessage({ id: "suggestions.column.content" }), dataIndex: "content_title", key: "content_title", render: (text: string) => <span className="font-medium">{text}</span> },
	{ title: intl.formatMessage({ id: "suggestions.column.type" }), dataIndex: "suggestion_type", key: "suggestion_type", width: 120, render: (type: string) => getTypeTag(type) },
	{ title: intl.formatMessage({ id: "suggestions.column.content" }), key: "title", render: (_: any, record: any) => <a onClick={() => handleShowDetail(record)} className="text-blue-500">{record.suggestion_data.title}</a> },
		{ title: intl.formatMessage({ id: 'common.column.priority' }), dataIndex: "priority", key: "priority", width: 100, render: (priority: number) => getPriorityTag(priority) },
		{
			title: intl.formatMessage({ id: "suggestions.column.priority" }),
			key: "impact",
			width: 80,
			render: (_: any, record: any) => (
				<Tag color={record.suggestion_data.impact === "high" ? "red" : record.suggestion_data.impact === "medium" ? "orange" : "blue"}>
					{record.suggestion_data.impact === "high" ? intl.formatMessage({ id: 'common.priority.high' }) : record.suggestion_data.impact === "medium" ? intl.formatMessage({ id: 'common.priority.medium' }) : intl.formatMessage({ id: 'common.priority.low' })}
				</Tag>
			),
		},
		{
			title: intl.formatMessage({ id: "suggestions.column.priority" }),
			key: "effort",
			width: 80,
			render: (_: any, record: any) => (
				<Tag color={record.suggestion_data.effort === "low" ? "green" : record.suggestion_data.effort === "medium" ? "orange" : "red"}>
					{record.suggestion_data.effort === "low" ? intl.formatMessage({ id: 'common.priority.low' }) : record.suggestion_data.effort === "medium" ? intl.formatMessage({ id: 'common.priority.medium' }) : intl.formatMessage({ id: 'common.priority.high' })}
				</Tag>
			),
		},
		{ title: intl.formatMessage({ id: 'common.column.status' }), dataIndex: "status", key: "status", width: 100, render: (status: number) => getStatusTag(status) },
		{ title: intl.formatMessage({ id: 'common.column.createdAt' }), dataIndex: "created_at", key: "created_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.viewDetail' })}><Button type="text" icon={<BulbOutlined />} onClick={() => handleShowDetail(record)} /></Tooltip>
					{record.status === 0 && (
						<>
							<Tooltip title={intl.formatMessage({ id: 'common.action.enable' })}><Button type="text" icon={<LikeOutlined />} onClick={() => handleApply(record.id)} /></Tooltip>
							<Tooltip title={intl.formatMessage({ id: 'common.action.disable' })}><Button type="text" icon={<DislikeOutlined />} onClick={() => handleIgnore(record.id)} /></Tooltip>
						</>
					)}
				</Space>
			),
		},
	];

	const highPrioritySuggestions = suggestions.filter((s: any) => s.priority === 2 && s.status === 0);

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.monitor.suggestions' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: "suggestions.page.subtitle" })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
			<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: "suggestions.column.content" })} value={stats.total} prefix={<BulbOutlined />} /></Card></Col>
			<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: "suggestions.status.pending" })} value={stats.pending} prefix={<ClockCircleOutlined />} valueStyle={{ color: "#1890ff" }} /></Card></Col>
			<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: "suggestions.status.applied" })} value={stats.applied} prefix={<CheckCircleOutlined />} valueStyle={{ color: "#52c41a" }} /></Card></Col>
			<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: "suggestions.status.ignored" })} value={stats.ignored} prefix={<ExclamationCircleOutlined />} valueStyle={{ color: "#8c8c8c" }} /></Card></Col>
			<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: "suggestions.priority.high" })} value={stats.highPriority} prefix={<ExclamationCircleOutlined />} valueStyle={{ color: "#ff4d4f" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: "common.status.active" })} value={stats.total > 0 ? Math.round((stats.applied / stats.total) * 100) : 0} suffix="%" valueStyle={{ color: "#52c41a" }} /></Card></Col>
			</Row>

			{highPrioritySuggestions.length > 0 && (
				<Card title={intl.formatMessage({ id: "suggestions.priority.high" })} className="mb-4" type="inner">
					<List
						dataSource={highPrioritySuggestions}
						renderItem={(item: any) => (
							<List.Item
								actions={[
								<Button type="primary" size="small" icon={<LikeOutlined />} onClick={() => handleApply(item.id)}>{intl.formatMessage({ id: "suggestions.action.apply" })}</Button>,
								<Button size="small" icon={<DislikeOutlined />} onClick={() => handleIgnore(item.id)}>{intl.formatMessage({ id: "suggestions.action.ignore" })}</Button>,
								]}
							>
								<List.Item.Meta
									avatar={<Avatar icon={<ExclamationCircleOutlined />} style={{ backgroundColor: "#ff4d4f" }} />}
									title={<Space><span>{item.suggestion_data.title}</span>{getTypeTag(item.suggestion_type)}</Space>}
									description={<div><div>{item.content_title}</div><div className="text-gray-400 text-sm">{item.suggestion_data.description}</div></div>}
								/>
							</List.Item>
						)}
					/>
				</Card>
			)}

			<Card
				title={<Space><BulbOutlined /><span>{intl.formatMessage({ id: "suggestions.column.content" })}</span></Space>}
				extra={
					<Space>
					<InputNumber placeholder="Content ID" min={0} style={{ width: 100 }} value={generateContentId} onChange={(v) => setGenerateContentId(v || 0)} />
					<Button type="primary" icon={<BulbOutlined />} loading={generateMutation.isPending} onClick={handleGenerateSuggestions}>{intl.formatMessage({ id: "suggestions.action.apply" })}</Button>
					<Select placeholder={intl.formatMessage({ id: "suggestions.column.type" })} allowClear style={{ width: 120 }} onChange={setFilterType} options={suggestionTypes.map((type) => ({ value: type.value, label: type.label }))} />
						<Select placeholder={intl.formatMessage({ id: 'common.column.priority' })} allowClear style={{ width: 100 }} onChange={setFilterPriority} options={[{ value: 2, label: intl.formatMessage({ id: 'common.priority.high' }) }, { value: 1, label: intl.formatMessage({ id: 'common.priority.medium' }) }, { value: 0, label: intl.formatMessage({ id: 'common.priority.low' }) }]} />
						<Button icon={<ReloadOutlined />} onClick={() => refetch()}>{intl.formatMessage({ id: 'common.action.refresh' })}</Button>
					</Space>
				}
			>
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={[
					{ key: "all", label: `${intl.formatMessage({ id: "config.tab.all" })} (${stats.total})` },
					{ key: "pending", label: `${intl.formatMessage({ id: "suggestions.status.pending" })} (${stats.pending})` },
					{ key: "applied", label: `${intl.formatMessage({ id: "suggestions.status.applied" })} (${stats.applied})` },
					{ key: "ignored", label: `${intl.formatMessage({ id: "suggestions.status.ignored" })} (${stats.ignored})` },
					]}
					className="mb-4"
				/>

				<Table columns={columns} dataSource={filteredSuggestions} rowKey="id" pagination={{ showSizeChanger: true, showQuickJumper: true, showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }) }} />
			</Card>

			<Modal title={intl.formatMessage({ id: "suggestions.page.title" })} open={detailModalVisible} onCancel={() => setDetailModalVisible(false)} footer={null} width={600}>
				{selectedSuggestion && (
					<div>
						<Descriptions column={1} bordered>
						<Descriptions.Item label={intl.formatMessage({ id: "suggestions.column.content" })}>{selectedSuggestion.content_title}</Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: "suggestions.column.type" })}>{getTypeTag(selectedSuggestion.suggestion_type)}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'common.column.priority' })}>{getPriorityTag(selectedSuggestion.priority)}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'common.column.status' })}>{getStatusTag(selectedSuggestion.status)}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: "suggestions.column.priority" })}>
								<Tag color={selectedSuggestion.suggestion_data.impact === "high" ? "red" : selectedSuggestion.suggestion_data.impact === "medium" ? "orange" : "blue"}>
									{selectedSuggestion.suggestion_data.impact === "high" ? intl.formatMessage({ id: 'common.priority.high' }) : selectedSuggestion.suggestion_data.impact === "medium" ? intl.formatMessage({ id: 'common.priority.medium' }) : intl.formatMessage({ id: 'common.priority.low' })}
								</Tag>
							</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: "suggestions.column.priority" })}>
								<Tag color={selectedSuggestion.suggestion_data.effort === "low" ? "green" : selectedSuggestion.suggestion_data.effort === "medium" ? "orange" : "red"}>
									{selectedSuggestion.suggestion_data.effort === "low" ? intl.formatMessage({ id: 'common.priority.low' }) : selectedSuggestion.suggestion_data.effort === "medium" ? intl.formatMessage({ id: 'common.priority.medium' }) : intl.formatMessage({ id: 'common.priority.high' })}
								</Tag>
							</Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: "suggestions.column.content" })}>{selectedSuggestion.suggestion_data.title}</Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: "suggestions.column.content" })}>{selectedSuggestion.suggestion_data.description}</Descriptions.Item>
						<Descriptions.Item label={intl.formatMessage({ id: "suggestions.column.content" })}>
								<ul className="list-disc pl-4">
									{selectedSuggestion.suggestion_data.examples.map((example: string, index: number) => <li key={index}>{example}</li>)}
								</ul>
							</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'common.column.createdAt' })}>{new Date(selectedSuggestion.created_at).toLocaleString()}</Descriptions.Item>
						</Descriptions>

						{selectedSuggestion.status === 0 && (
							<div className="mt-4 flex justify-end space-x-4">
							<Button icon={<DislikeOutlined />} onClick={() => { handleIgnore(selectedSuggestion.id); setDetailModalVisible(false); }}>{intl.formatMessage({ id: "suggestions.action.ignore" })}</Button>
							<Button type="primary" icon={<LikeOutlined />} onClick={() => { handleApply(selectedSuggestion.id); setDetailModalVisible(false); }}>{intl.formatMessage({ id: "suggestions.action.apply" })}</Button>
							</div>
						)}
					</div>
				)}
			</Modal>
		</div>
	);
}
