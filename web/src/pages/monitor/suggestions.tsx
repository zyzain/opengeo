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

const suggestionTypes = [
	{
		value: "content",
		label: "内容优化",
		color: "blue",
		icon: <BulbOutlined />,
	},
	{
		value: "structure",
		label: "结构优化",
		color: "green",
		icon: <ThunderboltOutlined />,
	},
	{
		value: "authority",
		label: "权威性提升",
		color: "purple",
		icon: <StarOutlined />,
	},
];

const priorityConfig: Record<number, { label: string; color: string }> = {
	2: { label: "高", color: "red" },
	1: { label: "中", color: "orange" },
	0: { label: "低", color: "blue" },
};

const statusConfig: Record<
	number,
	{ label: string; color: string; icon: React.ReactNode }
> = {
	0: { label: "待处理", color: "processing", icon: <ClockCircleOutlined /> },
	1: { label: "已应用", color: "success", icon: <CheckCircleOutlined /> },
	2: { label: "已忽略", color: "default", icon: <ExclamationCircleOutlined /> },
};

export default function SuggestionsPage() {
	const queryClient = useQueryClient();
	const generateMutation = useGenerateSuggestions();

	const {
		data: suggestionsData,
		isLoading,
		refetch,
	} = useQuery({
		queryKey: ["suggestions"],
		queryFn: () => api.monitor.listSuggestions(),
		select: (response) => response.data.data,
	});

	const suggestions = suggestionsData?.items || [];
	const [detailModalVisible, setDetailModalVisible] = useState(false);
	const [selectedSuggestion, setSelectedSuggestion] = useState<any>(null);
	const [activeTab, setActiveTab] = useState("all");
	const [filterType, setFilterType] = useState<string | undefined>(undefined);
	const [filterPriority, setFilterPriority] = useState<number | undefined>(
		undefined,
	);
	const [generateContentId, setGenerateContentId] = useState<number>(0);

	const handleGenerateSuggestions = async () => {
		try {
			await generateMutation.mutateAsync(generateContentId);
			message.success("建议生成成功");
			refetch();
		} catch (error: any) {
			message.error("生成失败");
		}
	};

	// 过滤建议
	const filteredSuggestions = suggestions.filter((s: any) => {
		const matchTab =
			activeTab === "all" ||
			(activeTab === "pending" && s.status === 0) ||
			(activeTab === "applied" && s.status === 1) ||
			(activeTab === "ignored" && s.status === 2);
		const matchType = !filterType || s.suggestion_type === filterType;
		const matchPriority =
			filterPriority === undefined || s.priority === filterPriority;
		return matchTab && matchType && matchPriority;
	});

	// 统计数据
	const stats = {
		total: suggestions.length,
		pending: suggestions.filter((s: any) => s.status === 0).length,
		applied: suggestions.filter((s: any) => s.status === 1).length,
		ignored: suggestions.filter((s: any) => s.status === 2).length,
		highPriority: suggestions.filter(
			(s: any) => s.priority === 2 && s.status === 0,
		).length,
	};

	// 显示详情
	const handleShowDetail = (record: any) => {
		setSelectedSuggestion(record);
		setDetailModalVisible(true);
	};

	// 应用建议
	const handleApply = async (id: number) => {
		try {
			await api.monitor.applySuggestion(id);
			message.success("建议已应用");
			refetch();
		} catch (error: any) {
			message.error("操作失败");
		}
	};

	// 忽略建议
	const handleIgnore = async (id: number) => {
		try {
			await api.monitor.ignoreSuggestion(id);
			message.success("建议已忽略");
			refetch();
		} catch (error: any) {
			message.error("操作失败");
		}
	};

	// 获取类型标签
	const getTypeTag = (type: string) => {
		const config = suggestionTypes.find((t) => t.value === type);
		return (
			<Tag color={config?.color} icon={config?.icon}>
				{config?.label}
			</Tag>
		);
	};

	// 获取优先级标签
	const getPriorityTag = (priority: number) => {
		const config = priorityConfig[priority];
		return <Tag color={config?.color}>{config?.label}优先级</Tag>;
	};

	// 获取状态标签
	const getStatusTag = (status: number) => {
		const config = statusConfig[status];
		return <Badge status={config?.color as any} text={config?.label} />;
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
			title: "内容",
			dataIndex: "content_title",
			key: "content_title",
			render: (text: string) => <span className="font-medium">{text}</span>,
		},
		{
			title: "建议类型",
			dataIndex: "suggestion_type",
			key: "suggestion_type",
			width: 120,
			render: (type: string) => getTypeTag(type),
		},
		{
			title: "建议标题",
			key: "title",
			render: (_: any, record: any) => (
				<a onClick={() => handleShowDetail(record)} className="text-blue-500">
					{record.suggestion_data.title}
				</a>
			),
		},
		{
			title: "优先级",
			dataIndex: "priority",
			key: "priority",
			width: 100,
			render: (priority: number) => getPriorityTag(priority),
		},
		{
			title: "影响",
			key: "impact",
			width: 80,
			render: (_: any, record: any) => (
				<Tag
					color={
						record.suggestion_data.impact === "high"
							? "red"
							: record.suggestion_data.impact === "medium"
								? "orange"
								: "blue"
					}
				>
					{record.suggestion_data.impact === "high"
						? "高"
						: record.suggestion_data.impact === "medium"
							? "中"
							: "低"}
				</Tag>
			),
		},
		{
			title: "工作量",
			key: "effort",
			width: 80,
			render: (_: any, record: any) => (
				<Tag
					color={
						record.suggestion_data.effort === "low"
							? "green"
							: record.suggestion_data.effort === "medium"
								? "orange"
								: "red"
					}
				>
					{record.suggestion_data.effort === "low"
						? "低"
						: record.suggestion_data.effort === "medium"
							? "中"
							: "高"}
				</Tag>
			),
		},
		{
			title: "状态",
			dataIndex: "status",
			key: "status",
			width: 100,
			render: (status: number) => getStatusTag(status),
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
			width: 200,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="查看详情">
						<Button
							type="text"
							icon={<BulbOutlined />}
							onClick={() => handleShowDetail(record)}
						/>
					</Tooltip>
					{record.status === 0 && (
						<>
							<Tooltip title="应用">
								<Button
									type="text"
									icon={<LikeOutlined />}
									onClick={() => handleApply(record.id)}
								/>
							</Tooltip>
							<Tooltip title="忽略">
								<Button
									type="text"
									icon={<DislikeOutlined />}
									onClick={() => handleIgnore(record.id)}
								/>
							</Tooltip>
						</>
					)}
				</Space>
			),
		},
	];

	// 高优先级建议列表
	const highPrioritySuggestions = suggestions.filter(
		(s: any) => s.priority === 2 && s.status === 0,
	);

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">优化建议</h1>
				<p className="text-gray-500 mt-1">
					基于AI分析生成的内容优化建议，提升GEO效果
				</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="总建议数"
							value={stats.total}
							prefix={<BulbOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="待处理"
							value={stats.pending}
							prefix={<ClockCircleOutlined />}
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="已应用"
							value={stats.applied}
							prefix={<CheckCircleOutlined />}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="已忽略"
							value={stats.ignored}
							prefix={<ExclamationCircleOutlined />}
							valueStyle={{ color: "#8c8c8c" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="高优先级"
							value={stats.highPriority}
							prefix={<ExclamationCircleOutlined />}
							valueStyle={{ color: "#ff4d4f" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="完成率"
							value={
								stats.total > 0
									? Math.round((stats.applied / stats.total) * 100)
									: 0
							}
							suffix="%"
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
			</Row>

			{/* 高优先级建议 */}
			{highPrioritySuggestions.length > 0 && (
				<Card title="高优先级建议" className="mb-4" type="inner">
					<List
						dataSource={highPrioritySuggestions}
						renderItem={(item: any) => (
							<List.Item
								actions={[
									<Button
										type="primary"
										size="small"
										icon={<LikeOutlined />}
										onClick={() => handleApply(item.id)}
									>
										应用
									</Button>,
									<Button
										size="small"
										icon={<DislikeOutlined />}
										onClick={() => handleIgnore(item.id)}
									>
										忽略
									</Button>,
								]}
							>
								<List.Item.Meta
									avatar={
										<Avatar
											icon={<ExclamationCircleOutlined />}
											style={{ backgroundColor: "#ff4d4f" }}
										/>
									}
									title={
										<Space>
											<span>{item.suggestion_data.title}</span>
											{getTypeTag(item.suggestion_type)}
										</Space>
									}
									description={
										<div>
											<div>{item.content_title}</div>
											<div className="text-gray-400 text-sm">
												{item.suggestion_data.description}
											</div>
										</div>
									}
								/>
							</List.Item>
						)}
					/>
				</Card>
			)}

			{/* 建议列表 */}
			<Card
				title={
					<Space>
						<BulbOutlined />
						<span>建议列表</span>
					</Space>
				}
				extra={
					<Space>
						<InputNumber
							placeholder="内容ID"
							min={0}
							style={{ width: 100 }}
							value={generateContentId}
							onChange={(v) => setGenerateContentId(v || 0)}
						/>
						<Button
							type="primary"
							icon={<BulbOutlined />}
							loading={generateMutation.isPending}
							onClick={handleGenerateSuggestions}
						>
							生成建议
						</Button>
						<Select
							placeholder="建议类型"
							allowClear
							style={{ width: 120 }}
							onChange={setFilterType}
							options={suggestionTypes.map((type) => ({
								value: type.value,
								label: type.label,
							}))}
						/>
						<Select
							placeholder="优先级"
							allowClear
							style={{ width: 100 }}
							onChange={setFilterPriority}
							options={[
								{ value: 2, label: "高" },
								{ value: 1, label: "中" },
								{ value: 0, label: "低" },
							]}
						/>
						<Button icon={<ReloadOutlined />} onClick={() => refetch()}>
							刷新
						</Button>
					</Space>
				}
			>
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={[
						{ key: "all", label: `全部 (${stats.total})` },
						{ key: "pending", label: `待处理 (${stats.pending})` },
						{ key: "applied", label: `已应用 (${stats.applied})` },
						{ key: "ignored", label: `已忽略 (${stats.ignored})` },
					]}
					className="mb-4"
				/>

				<Table
					columns={columns}
					dataSource={filteredSuggestions}
					rowKey="id"
					pagination={{
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
					}}
				/>
			</Card>

			{/* 建议详情弹窗 */}
			<Modal
				title="优化建议详情"
				open={detailModalVisible}
				onCancel={() => setDetailModalVisible(false)}
				footer={null}
				width={600}
			>
				{selectedSuggestion && (
					<div>
						<Descriptions column={1} bordered>
							<Descriptions.Item label="内容">
								{selectedSuggestion.content_title}
							</Descriptions.Item>
							<Descriptions.Item label="建议类型">
								{getTypeTag(selectedSuggestion.suggestion_type)}
							</Descriptions.Item>
							<Descriptions.Item label="优先级">
								{getPriorityTag(selectedSuggestion.priority)}
							</Descriptions.Item>
							<Descriptions.Item label="状态">
								{getStatusTag(selectedSuggestion.status)}
							</Descriptions.Item>
							<Descriptions.Item label="影响程度">
								<Tag
									color={
										selectedSuggestion.suggestion_data.impact === "high"
											? "red"
											: selectedSuggestion.suggestion_data.impact === "medium"
												? "orange"
												: "blue"
									}
								>
									{selectedSuggestion.suggestion_data.impact === "high"
										? "高"
										: selectedSuggestion.suggestion_data.impact === "medium"
											? "中"
											: "低"}
								</Tag>
							</Descriptions.Item>
							<Descriptions.Item label="工作量">
								<Tag
									color={
										selectedSuggestion.suggestion_data.effort === "low"
											? "green"
											: selectedSuggestion.suggestion_data.effort === "medium"
												? "orange"
												: "red"
									}
								>
									{selectedSuggestion.suggestion_data.effort === "low"
										? "低"
										: selectedSuggestion.suggestion_data.effort === "medium"
											? "中"
											: "高"}
								</Tag>
							</Descriptions.Item>
							<Descriptions.Item label="建议标题">
								{selectedSuggestion.suggestion_data.title}
							</Descriptions.Item>
							<Descriptions.Item label="详细描述">
								{selectedSuggestion.suggestion_data.description}
							</Descriptions.Item>
							<Descriptions.Item label="示例">
								<ul className="list-disc pl-4">
									{selectedSuggestion.suggestion_data.examples.map(
										(example: string, index: number) => (
											<li key={index}>{example}</li>
										),
									)}
								</ul>
							</Descriptions.Item>
							<Descriptions.Item label="创建时间">
								{new Date(selectedSuggestion.created_at).toLocaleString()}
							</Descriptions.Item>
						</Descriptions>

						{selectedSuggestion.status === 0 && (
							<div className="mt-4 flex justify-end space-x-4">
								<Button
									icon={<DislikeOutlined />}
									onClick={() => {
										handleIgnore(selectedSuggestion.id);
										setDetailModalVisible(false);
									}}
								>
									忽略
								</Button>
								<Button
									type="primary"
									icon={<LikeOutlined />}
									onClick={() => {
										handleApply(selectedSuggestion.id);
										setDetailModalVisible(false);
									}}
								>
									应用建议
								</Button>
							</div>
						)}
					</div>
				)}
			</Modal>
		</div>
	);
}
