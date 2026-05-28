"use client";

import { useAICitations } from "@/hooks";
import {
	CheckCircleOutlined,
	CloseCircleOutlined,
	FallOutlined,
	LineChartOutlined,
	RiseOutlined,
	RobotOutlined,
	SearchOutlined,
} from "@ant-design/icons";
import {
	Badge,
	Button,
	Card,
	Col,
	Form,
	Input,
	Row,
	Select,
	Space,
	Statistic,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import ReactECharts from "echarts-for-react";
import { useState } from "react";

const { Option } = Select;

export default function CitationsPage() {
	const [searchForm] = Form.useForm();
	const [queryParams, setQueryParams] = useState({
		page: 1,
		page_size: 10,
		ai_model: undefined,
	});

	const { data, isLoading } = useAICitations(queryParams);

	const citations = data?.items || [];
	const total = data?.total || 0;

	// AI模型列表
	const aiModels = [
		{ value: "deepseek", label: "DeepSeek", color: "blue" },
		{ value: "kimi", label: "Kimi", color: "purple" },
		{ value: "doubao", label: "豆包", color: "orange" },
		{ value: "chatgpt", label: "ChatGPT", color: "green" },
	];

	// 获取AI模型标签
	const getAIModelTag = (model: string) => {
		const modelInfo = aiModels.find((m) => m.value === model);
		return (
			<Tag color={modelInfo?.color || "default"}>
				<RobotOutlined /> {modelInfo?.label || model}
			</Tag>
		);
	};

	// 情感标签
	const getSentimentTag = (sentiment: string) => {
		const sentimentMap: Record<string, { color: string; text: string }> = {
			positive: { color: "success", text: "正面" },
			neutral: { color: "default", text: "中性" },
			negative: { color: "error", text: "负面" },
		};
		const config = sentimentMap[sentiment] || {
			color: "default",
			text: sentiment,
		};
		return <Tag color={config.color}>{config.text}</Tag>;
	};

	// 统计数据 - 从API数据计算
	const stats = {
		total: total,
		cited: citations.filter((c: any) => c.is_cited).length,
		uncited: citations.filter((c: any) => !c.is_cited).length,
		citationRate:
			citations.length > 0
				? (
						(citations.filter((c: any) => c.is_cited).length /
							citations.length) *
						100
					).toFixed(1)
				: "0",
	};

	// 模型分布数据 - 从API数据计算
	const modelDistribution = aiModels
		.map((model) => ({
			value: citations.filter(
				(c: any) => c.ai_model === model.value && c.is_cited,
			).length,
			name: model.label,
		}))
		.filter((item) => item.value > 0);

	// 趋势数据 - 按日期分组统计
	const trendData = (() => {
		const dateMap: Record<string, number> = {};
		citations.forEach((c: any) => {
			if (c.is_cited && c.tracked_at) {
				const date = new Date(c.tracked_at).toLocaleDateString("zh-CN", {
					month: "short",
					day: "numeric",
				});
				dateMap[date] = (dateMap[date] || 0) + 1;
			}
		});
		const sortedDates = Object.keys(dateMap).slice(-7);
		return {
			dates: sortedDates,
			values: sortedDates.map((d) => dateMap[d] || 0),
		};
	})();

	// 表格列定义
	const columns = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
			width: 80,
		},
		{
			title: "内容ID",
			dataIndex: "content_id",
			key: "content_id",
			width: 80,
		},
		{
			title: "AI模型",
			dataIndex: "ai_model",
			key: "ai_model",
			width: 120,
			render: (model: string) => getAIModelTag(model),
		},
		{
			title: "查询文本",
			dataIndex: "query_text",
			key: "query_text",
			ellipsis: true,
		},
		{
			title: "是否引用",
			dataIndex: "is_cited",
			key: "is_cited",
			width: 100,
			render: (cited: boolean) =>
				cited ? (
					<Badge status="success" text="已引用" />
				) : (
					<Badge status="default" text="未引用" />
				),
		},
		{
			title: "引用位置",
			dataIndex: "citation_position",
			key: "citation_position",
			width: 80,
			render: (position: number) =>
				position > 0 ? <Tag color="blue">第{position}位</Tag> : "-",
		},
		{
			title: "引用内容",
			dataIndex: "citation_text",
			key: "citation_text",
			ellipsis: true,
			render: (text: string) =>
				text ? (
					<Tooltip title={text}>
						<span>{text.substring(0, 50)}...</span>
					</Tooltip>
				) : (
					"-"
				),
		},
		{
			title: "情感",
			dataIndex: "sentiment",
			key: "sentiment",
			width: 80,
			render: (sentiment: string) => getSentimentTag(sentiment),
		},
		{
			title: "追踪时间",
			dataIndex: "tracked_at",
			key: "tracked_at",
			width: 180,
			render: (text: string) => new Date(text).toLocaleString(),
		},
	];

	// 图表配置 - 使用真实数据
	const chartOption = {
		title: {
			text: "AI引用趋势",
			left: "center",
		},
		tooltip: {
			trigger: "axis",
		},
		xAxis: {
			type: "category",
			data: trendData.dates.length > 0 ? trendData.dates : ["暂无数据"],
		},
		yAxis: {
			type: "value",
		},
		series: [
			{
				name: "引用次数",
				type: "line",
				smooth: true,
				data: trendData.values.length > 0 ? trendData.values : [0],
				areaStyle: {
					opacity: 0.3,
				},
				itemStyle: {
					color: "#1890ff",
				},
			},
		],
	};

	// 模型分布图表 - 使用真实数据
	const modelChartOption = {
		title: {
			text: "AI模型引用分布",
			left: "center",
		},
		tooltip: {
			trigger: "item",
		},
		legend: {
			orient: "vertical",
			left: "left",
		},
		series: [
			{
				name: "引用次数",
				type: "pie",
				radius: "50%",
				data:
					modelDistribution.length > 0
						? modelDistribution
						: [{ value: 0, name: "暂无数据" }],
				emphasis: {
					itemStyle: {
						shadowBlur: 10,
						shadowOffsetX: 0,
						shadowColor: "rgba(0, 0, 0, 0.5)",
					},
				},
			},
		],
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">AI引用追踪</h1>
				<p className="text-gray-500 mt-1">监测内容在AI搜索中的引用情况</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="总追踪次数"
							value={stats.total}
							prefix={<LineChartOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="被引用次数"
							value={stats.cited}
							prefix={<CheckCircleOutlined />}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="未引用次数"
							value={stats.uncited}
							prefix={<CloseCircleOutlined />}
							valueStyle={{ color: "#8c8c8c" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="引用率"
							value={stats.citationRate}
							suffix="%"
							prefix={<RiseOutlined />}
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
			</Row>

			{/* 图表 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={14}>
					<Card>
						<ReactECharts option={chartOption} style={{ height: 300 }} />
					</Card>
				</Col>
				<Col xs={24} lg={10}>
					<Card>
						<ReactECharts option={modelChartOption} style={{ height: 300 }} />
					</Card>
				</Col>
			</Row>

			{/* 搜索表单 */}
			<Card className="mb-4">
				<Form
					form={searchForm}
					layout="inline"
					onFinish={(values) =>
						setQueryParams({ ...queryParams, ...values, page: 1 })
					}
				>
					<Form.Item name="ai_model" label="AI模型">
						<Select
							placeholder="请选择AI模型"
							allowClear
							style={{ width: 150 }}
						>
							{aiModels.map((model) => (
								<Option key={model.value} value={model.value}>
									{model.label}
								</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button
								type="primary"
								icon={<SearchOutlined />}
								htmlType="submit"
							>
								搜索
							</Button>
							<Button
								onClick={() => {
									searchForm.resetFields();
									setQueryParams({
										page: 1,
										page_size: 10,
										ai_model: undefined,
									});
								}}
							>
								重置
							</Button>
						</Space>
					</Form.Item>
				</Form>
			</Card>

			{/* 引用列表 */}
			<Card title="AI引用记录">
				<Table
					columns={columns}
					dataSource={citations}
					rowKey="id"
					loading={isLoading}
					scroll={{ x: 1500 }}
					pagination={{
						current: queryParams.page,
						pageSize: queryParams.page_size,
						total,
						showSizeChanger: true,
						showQuickJumper: true,
						showTotal: (total) => `共 ${total} 条`,
						onChange: (page, pageSize) =>
							setQueryParams({ ...queryParams, page, page_size: pageSize }),
					}}
				/>
			</Card>
		</div>
	);
}
