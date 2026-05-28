"use client";

import { queryKeys, useCompetitors } from "@/hooks";
import api from "@/lib/api";
import {
	CheckCircleOutlined,
	DeleteOutlined,
	EyeOutlined,
	GlobalOutlined,
	PlusOutlined,
	SyncOutlined,
	TrophyOutlined,
	WarningOutlined,
} from "@ant-design/icons";
import { useQuery, useQueryClient } from "@tanstack/react-query";
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
	Progress,
	Row,
	Space,
	Statistic,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import ReactECharts from "echarts-for-react";
import { useState } from "react";

const { TextArea } = Input;

export default function CompetitorsPage() {
	const queryClient = useQueryClient();
	const { data: competitorsData, isLoading, refetch } = useCompetitors();
	const competitors = competitorsData?.items || [];
	const { data: ourScoreData } = useQuery({
		queryKey: ['ourScore'],
		queryFn: () => api.monitor.ourScore(),
	});
	const [createModalVisible, setCreateModalVisible] = useState(false);
	const [detailModalVisible, setDetailModalVisible] = useState(false);
	const [selectedCompetitor, setSelectedCompetitor] = useState<any>(null);
	const [createForm] = Form.useForm();

	// 统计数据
	const stats = {
		total: competitors.length,
		active: competitors.filter((c: any) => c.status === "active").length,
		avgVisibility:
			competitors.length > 0
				? (
						competitors.reduce(
							(sum: number, c: any) => sum + (c.visibility_score || 0),
							0,
						) / competitors.length
					).toFixed(1)
				: "0",
		totalGaps: competitors.reduce(
			(sum: number, c: any) => sum + (c.content_gap_count || 0),
			0,
		),
		ourScore: ourScoreData?.data?.data?.score || 0,
		ourRank: ourScoreData?.data?.data?.rank || 0,
	};

	// 表格列定义
	const columns = [
		{
			title: "排名",
			key: "rank",
			width: 80,
			render: (_: any, __: any, index: number) => (
				<div className="text-center">
					{index < 3 ? (
						<TrophyOutlined
							style={{
								color: ["#ffd700", "#c0c0c0", "#cd7f32"][index],
								fontSize: 20,
							}}
						/>
					) : (
						<span className="text-gray-500">{index + 1}</span>
					)}
				</div>
			),
		},
		{
			title: "竞品名称",
			dataIndex: "competitor_name",
			key: "competitor_name",
			render: (text: string, record: any) => (
				<div>
					<a
						onClick={() => handleShowDetail(record)}
						className="text-blue-500 font-medium"
					>
						{text}
					</a>
					<div className="text-gray-400 text-xs">
						{record.competitor_domain}
					</div>
				</div>
			),
		},
		{
			title: "可见性评分",
			dataIndex: "visibility_score",
			key: "visibility_score",
			width: 150,
			render: (score: number) => (
				<Progress
					percent={score}
					size="small"
					strokeColor={
						score >= 80 ? "#52c41a" : score >= 70 ? "#1890ff" : "#faad14"
					}
					format={(percent) => <span className="font-bold">{percent}</span>}
				/>
			),
		},
		{
			title: "内容差距",
			dataIndex: "content_gap_count",
			key: "content_gap_count",
			width: 100,
			render: (count: number) => (
				<Badge
					count={count}
					showZero
					style={{ backgroundColor: count > 10 ? "#ff4d4f" : "#1890ff" }}
				/>
			),
		},
		{
			title: "热门查询",
			dataIndex: "top_queries",
			key: "top_queries",
			render: (queries: string[]) => (
				<Space size={[0, 4]} wrap>
					{queries.slice(0, 2).map((query) => (
						<Tag key={query} color="blue">
							{query}
						</Tag>
					))}
					{queries.length > 2 && <Tag>+{queries.length - 2}</Tag>}
				</Space>
			),
		},
		{
			title: "优势",
			dataIndex: "strengths",
			key: "strengths",
			render: (strengths: string[]) => (
				<Space size={[0, 4]} wrap>
					{strengths.slice(0, 1).map((s) => (
						<Tag key={s} color="green" icon={<CheckCircleOutlined />}>
							{s}
						</Tag>
					))}
				</Space>
			),
		},
		{
			title: "劣势",
			dataIndex: "weaknesses",
			key: "weaknesses",
			render: (weaknesses: string[]) => (
				<Space size={[0, 4]} wrap>
					{weaknesses.slice(0, 1).map((w) => (
						<Tag key={w} color="red" icon={<WarningOutlined />}>
							{w}
						</Tag>
					))}
				</Space>
			),
		},
		{
			title: "状态",
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (status: string) => (
				<Badge
					status={status === "active" ? "success" : "default"}
					text={status === "active" ? "监测中" : "已暂停"}
				/>
			),
		},
		{
			title: "操作",
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="查看详情">
						<Button
							type="text"
							icon={<EyeOutlined />}
							onClick={() => handleShowDetail(record)}
						/>
					</Tooltip>
					<Tooltip title="更新数据">
						<Button
							type="text"
							icon={<SyncOutlined />}
							onClick={() => handleSync(record.id)}
						/>
					</Tooltip>
					<Popconfirm
						title="确定要删除这个竞品吗？"
						okText="确定"
						cancelText="取消"
						onConfirm={() => handleDelete(record.id)}
					>
						<Tooltip title="删除">
							<Button type="text" danger icon={<DeleteOutlined />} />
						</Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	// 显示详情
	const handleShowDetail = (record: any) => {
		setSelectedCompetitor(record);
		setDetailModalVisible(true);
	};

	// 创建竞品
	const handleCreate = async (values: any) => {
		try {
			await api.monitor.createCompetitor(values);
			message.success("创建成功");
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: queryKeys.competitors });
		} catch (error: any) {
			message.error("创建失败");
		}
	};

	// 删除竞品
	const handleDelete = async (id: number) => {
		try {
			await api.monitor.deleteCompetitor(id);
			message.success("删除成功");
			queryClient.invalidateQueries({ queryKey: queryKeys.competitors });
		} catch (error: any) {
			message.error("删除失败");
		}
	};

	// 同步/更新竞品数据
	const handleSync = async (id: number) => {
		try {
			await api.monitor.syncCompetitor(id);
			message.success("同步任务已创建");
			queryClient.invalidateQueries({ queryKey: queryKeys.competitors });
		} catch (error: any) {
			message.error("同步失败");
		}
	};

	// 可见性对比图表
	const getVisibilityChartOption = () => ({
		title: {
			text: "可见性评分对比",
			left: "center",
		},
		tooltip: {
			trigger: "axis",
		},
		xAxis: {
			type: "category",
			data: ["我们", ...competitors.map((c: any) => c.competitor_name)],
			axisLabel: { rotate: 30 },
		},
		yAxis: {
			type: "value",
			min: Math.max(0, Math.floor(Math.min(stats.ourScore, ...competitors.map((c: any) => c.visibility_score || 0)) / 10) * 10 - 10),
			max: 100,
		},
		series: [
			{
				type: "bar",
				data: [
					{ value: stats.ourScore, itemStyle: { color: "#52c41a" } },
					...competitors.map((c: any) => ({
						value: c.visibility_score,
						itemStyle: { color: "#1890ff" },
					})),
				],
				barWidth: "50%",
			},
		],
	});

	// 内容差距图表
	const getContentGapChartOption = () => ({
		title: {
			text: "内容差距分析",
			left: "center",
		},
		tooltip: {
			trigger: "item",
		},
		legend: {
			bottom: 0,
		},
		series: [
			{
				type: "pie",
				radius: "50%",
				data: competitors.map((c: any) => ({
					value: c.content_gap_count,
					name: c.competitor_name,
				})),
				emphasis: {
					itemStyle: {
						shadowBlur: 10,
						shadowOffsetX: 0,
						shadowColor: "rgba(0, 0, 0, 0.5)",
					},
				},
			},
		],
	});

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">竞品监测</h1>
				<p className="text-gray-500 mt-1">
					监测竞争对手在AI搜索中的表现，发现内容差距与机会点
				</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="竞品数量"
							value={stats.total}
							prefix={<GlobalOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="监测中"
							value={stats.active}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="我们的评分"
							value={stats.ourScore}
							suffix="分"
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="我们的排名"
							value={stats.ourRank}
							suffix="名"
							prefix={<TrophyOutlined />}
							valueStyle={{ color: "#ffd700" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="竞品平均分"
							value={stats.avgVisibility}
							suffix="分"
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="内容差距数"
							value={stats.totalGaps}
							prefix={<WarningOutlined />}
							valueStyle={{ color: "#ff4d4f" }}
						/>
					</Card>
				</Col>
			</Row>

			{/* 图表 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={14}>
					<Card>
						<ReactECharts
							option={getVisibilityChartOption()}
							style={{ height: 300 }}
						/>
					</Card>
				</Col>
				<Col xs={24} lg={10}>
					<Card>
						<ReactECharts
							option={getContentGapChartOption()}
							style={{ height: 300 }}
						/>
					</Card>
				</Col>
			</Row>

			{/* 竞品列表 */}
			<Card
				title={
					<Space>
						<GlobalOutlined />
						<span>竞品列表</span>
					</Space>
				}
				extra={
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => setCreateModalVisible(true)}
					>
						添加竞品
					</Button>
				}
			>
				<Table
					columns={columns}
					dataSource={competitors}
					rowKey="id"
					pagination={{
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
					}}
				/>
			</Card>

			{/* 创建竞品弹窗 */}
			<Modal
				title="添加竞品"
				open={createModalVisible}
				onCancel={() => setCreateModalVisible(false)}
				footer={null}
				width={500}
			>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item
						name="competitor_name"
						label="竞品名称"
						rules={[{ required: true, message: "请输入竞品名称" }]}
					>
						<Input placeholder="请输入竞品名称" />
					</Form.Item>
					<Form.Item
						name="competitor_domain"
						label="竞品域名"
						rules={[{ required: true, message: "请输入竞品域名" }]}
					>
						<Input placeholder="https://example.com" />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">
								添加
							</Button>
							<Button onClick={() => setCreateModalVisible(false)}>取消</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			{/* 竞品详情弹窗 */}
			<Modal
				title="竞品详情"
				open={detailModalVisible}
				onCancel={() => setDetailModalVisible(false)}
				footer={null}
				width={700}
			>
				{selectedCompetitor && (
					<div>
						<Descriptions column={2} bordered className="mb-4">
							<Descriptions.Item label="竞品名称" span={2}>
								{selectedCompetitor.competitor_name}
							</Descriptions.Item>
							<Descriptions.Item label="域名">
								<a
									href={selectedCompetitor.competitor_domain}
									target="_blank"
									rel="noopener noreferrer"
								>
									{selectedCompetitor.competitor_domain}
								</a>
							</Descriptions.Item>
							<Descriptions.Item label="状态">
								<Badge
									status={
										selectedCompetitor.status === "active"
											? "success"
											: "default"
									}
									text={
										selectedCompetitor.status === "active" ? "监测中" : "已暂停"
									}
								/>
							</Descriptions.Item>
							<Descriptions.Item label="可见性评分">
								<Progress
									percent={selectedCompetitor.visibility_score}
									size="small"
								/>
							</Descriptions.Item>
							<Descriptions.Item label="内容差距数">
								<Badge
									count={selectedCompetitor.content_gap_count}
									showZero
									style={{ backgroundColor: "#ff4d4f" }}
								/>
							</Descriptions.Item>
							<Descriptions.Item label="热门查询" span={2}>
								<Space size={[0, 4]} wrap>
									{selectedCompetitor.top_queries.map((q: string) => (
										<Tag key={q} color="blue">
											{q}
										</Tag>
									))}
								</Space>
							</Descriptions.Item>
							<Descriptions.Item label="优势" span={2}>
								<Space size={[0, 4]} wrap>
									{selectedCompetitor.strengths.map((s: string) => (
										<Tag key={s} color="green" icon={<CheckCircleOutlined />}>
											{s}
										</Tag>
									))}
								</Space>
							</Descriptions.Item>
							<Descriptions.Item label="劣势" span={2}>
								<Space size={[0, 4]} wrap>
									{selectedCompetitor.weaknesses.map((w: string) => (
										<Tag key={w} color="red" icon={<WarningOutlined />}>
											{w}
										</Tag>
									))}
								</Space>
							</Descriptions.Item>
							<Descriptions.Item label="最后检查" span={2}>
								{new Date(selectedCompetitor.last_check_time).toLocaleString()}
							</Descriptions.Item>
						</Descriptions>
					</div>
				)}
			</Modal>
		</div>
	);
}
