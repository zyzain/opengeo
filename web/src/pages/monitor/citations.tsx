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
import { useIntl } from "react-intl";

const { Option } = Select;

export default function CitationsPage() {
	const intl = useIntl();
	const [searchForm] = Form.useForm();
	const [queryParams, setQueryParams] = useState({
		page: 1,
		page_size: 10,
		ai_model: undefined,
	});

	const { data, isLoading } = useAICitations(queryParams);

	const citations = data?.items || [];
	const total = data?.total || 0;

	const aiModels = [
		{ value: "deepseek", label: "DeepSeek", color: "blue" },
		{ value: "kimi", label: "Kimi", color: "purple" },
		{ value: "doubao", label: intl.formatMessage({ id: 'content.aiModel.doubao' }), color: "orange" },
		{ value: "chatgpt", label: "ChatGPT", color: "green" },
	];

	const getAIModelTag = (model: string) => {
		const modelInfo = aiModels.find((m) => m.value === model);
		if (modelInfo) {
			return <Tag color={modelInfo.color}>{modelInfo.label}</Tag>;
		}
		return <Tag>{model}</Tag>;
	};

	const getSentimentTag = (sentiment: string) => {
		const sentimentMap: Record<string, { color: string; text: string }> = {
			positive: { color: "success", text: intl.formatMessage({ id: 'citations.sentiment.positive' }) },
			neutral: { color: "default", text: intl.formatMessage({ id: 'citations.sentiment.neutral' }) },
			negative: { color: "error", text: intl.formatMessage({ id: 'citations.sentiment.negative' }) },
		};
		const config = sentimentMap[sentiment] || { color: "default", text: sentiment };
		return <Tag color={config.color}>{config.text}</Tag>;
	};

	const stats = {
		total: total,
		cited: citations.filter((c: any) => c.is_cited).length,
		uncited: citations.filter((c: any) => !c.is_cited).length,
		citationRate: citations.length > 0
			? ((citations.filter((c: any) => c.is_cited).length / citations.length) * 100).toFixed(1)
			: "0",
	};

	const modelDistribution = aiModels
		.map((model) => ({
			value: citations.filter((c: any) => c.ai_model === model.value && c.is_cited).length,
			name: model.label,
		}))
		.filter((item) => item.value > 0);

	const trendData = (() => {
		const dateMap: Record<string, number> = {};
		citations.forEach((c: any) => {
			if (c.is_cited && c.tracked_at) {
				const date = new Date(c.tracked_at).toLocaleDateString("zh-CN", { month: "short", day: "numeric" });
				dateMap[date] = (dateMap[date] || 0) + 1;
			}
		});
		const sortedDates = Object.keys(dateMap).slice(-7);
		return { dates: sortedDates, values: sortedDates.map((d) => dateMap[d] || 0) };
	})();

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 80 },
		{ title: intl.formatMessage({ id: 'publish.column.contentId' }), dataIndex: "content_id", key: "content_id", width: 80 },
		{ title: intl.formatMessage({ id: 'citations.column.aiModel' }), dataIndex: "ai_model", key: "ai_model", width: 120, render: (model: string) => getAIModelTag(model) },
		{ title: intl.formatMessage({ id: 'citations.column.queryText' }), dataIndex: "query_text", key: "query_text", ellipsis: true },
		{
			title: intl.formatMessage({ id: 'citations.column.isCited' }),
			dataIndex: "is_cited",
			key: "is_cited",
			width: 100,
			render: (cited: boolean) => cited ? <Badge status="success" text={intl.formatMessage({ id: 'citations.cited' })} /> : <Badge status="default" text={intl.formatMessage({ id: 'citations.uncited' })} />,
		},
		{
			title: intl.formatMessage({ id: 'citations.column.citationPosition' }),
			dataIndex: "citation_position",
			key: "citation_position",
			width: 80,
			render: (position: number) => position > 0 ? <Tag color="blue">{intl.formatMessage({ id: 'citations.position' }, { pos: position })}</Tag> : "-",
		},
		{
			title: intl.formatMessage({ id: 'citations.column.citationText' }),
			dataIndex: "citation_text",
			key: "citation_text",
			ellipsis: true,
			render: (text: string) => text ? <Tooltip title={text}><span>{text.substring(0, 50)}...</span></Tooltip> : "-",
		},
		{ title: intl.formatMessage({ id: 'citations.column.sentiment' }), dataIndex: "sentiment", key: "sentiment", width: 80, render: (sentiment: string) => getSentimentTag(sentiment) },
		{ title: intl.formatMessage({ id: 'citations.column.trackTime' }), dataIndex: "tracked_at", key: "tracked_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
	];

	const chartOption = {
		title: { text: intl.formatMessage({ id: 'citations.chart.trendTitle' }), left: "center" },
		tooltip: { trigger: "axis" },
		xAxis: { type: "category", data: trendData.dates.length > 0 ? trendData.dates : [intl.formatMessage({ id: 'citations.chart.noData' })] },
		yAxis: { type: "value" },
		series: [{ name: intl.formatMessage({ id: 'citations.chart.trendSeries' }), type: "line", smooth: true, data: trendData.values.length > 0 ? trendData.values : [0], areaStyle: { opacity: 0.3 }, itemStyle: { color: "#1890ff" } }],
	};

	const modelChartOption = {
		title: { text: intl.formatMessage({ id: 'citations.chart.modelDist' }), left: "center" },
		tooltip: { trigger: "item" },
		legend: { orient: "vertical", left: "left" },
		series: [{
			name: intl.formatMessage({ id: 'citations.chart.trendSeries' }),
			type: "pie",
			radius: "50%",
			data: modelDistribution.length > 0 ? modelDistribution : [{ value: 0, name: intl.formatMessage({ id: 'citations.chart.noData' }) }],
			emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: "rgba(0, 0, 0, 0.5)" } },
		}],
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.monitor.citations' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'citations.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}>
					<Card><Statistic title={intl.formatMessage({ id: 'monitor.citations.total' })} value={stats.total} prefix={<LineChartOutlined />} /></Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card><Statistic title={intl.formatMessage({ id: 'monitor.citations.cited' })} value={stats.cited} prefix={<CheckCircleOutlined />} valueStyle={{ color: "#52c41a" }} /></Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card><Statistic title={intl.formatMessage({ id: 'citations.stat.uncited' })} value={stats.uncited} prefix={<CloseCircleOutlined />} valueStyle={{ color: "#8c8c8c" }} /></Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card><Statistic title={intl.formatMessage({ id: 'monitor.citations.rate' })} value={stats.citationRate} suffix="%" prefix={<RiseOutlined />} valueStyle={{ color: "#1890ff" }} /></Card>
				</Col>
			</Row>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={14}><Card><ReactECharts option={chartOption} style={{ height: 300 }} /></Card></Col>
				<Col xs={24} lg={10}><Card><ReactECharts option={modelChartOption} style={{ height: 300 }} /></Card></Col>
			</Row>

			<Card className="mb-4">
				<Form form={searchForm} layout="inline" onFinish={(values) => setQueryParams({ ...queryParams, ...values, page: 1 })}>
				<Form.Item name="ai_model" label={intl.formatMessage({ id: 'citations.form.aiModel' })}>
					<Select placeholder={intl.formatMessage({ id: 'citations.placeholder.selectModel' })} allowClear style={{ width: 150 }}>
							{aiModels.map((model) => (
								<Option key={model.value} value={model.value}>{model.label}</Option>
							))}
						</Select>
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" icon={<SearchOutlined />} htmlType="submit">{intl.formatMessage({ id: 'common.action.search' })}</Button>
							<Button onClick={() => { searchForm.resetFields(); setQueryParams({ page: 1, page_size: 10, ai_model: undefined }); }}>{intl.formatMessage({ id: 'common.action.reset' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Card>

			<Card title={intl.formatMessage({ id: 'monitor.citations.record' })}>
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
						showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }),
						onChange: (page, pageSize) => setQueryParams({ ...queryParams, page, page_size: pageSize }),
					}}
				/>
			</Card>
		</div>
	);
}
