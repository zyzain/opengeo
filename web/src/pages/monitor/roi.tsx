"use client";

import { useROIMetrics } from "@/hooks";
import {
	BarChartOutlined,
	DollarOutlined,
	ExportOutlined,
	FallOutlined,
	FilterOutlined,
	LineChartOutlined,
	PieChartOutlined,
	RiseOutlined,
} from "@ant-design/icons";
import {
	Button,
	Card,
	Col,
	DatePicker,
	message,
	Progress,
	Row,
	Select,
	Space,
	Statistic,
	Table,
	Tabs,
	Tag,
} from "antd";
import type dayjs from "dayjs";
import ReactECharts from "echarts-for-react";
import { useState } from "react";
import { useIntl } from "react-intl";

const { RangePicker } = DatePicker;
const { Option } = Select;

export default function ROIPage() {
	const intl = useIntl();
	const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs] | null>(null);
	const [selectedChannel, setSelectedChannel] = useState<string | undefined>(undefined);
	const [activeTab, setActiveTab] = useState("overview");

	const { data: roiData, isLoading } = useROIMetrics({
		start_date: dateRange?.[0]?.toISOString(),
		end_date: dateRange?.[1]?.toISOString(),
		channel_id: selectedChannel,
	});

	const summary = roiData?.summary || {
		total_investment: 0,
		total_return: 0,
		roi_percentage: 0,
		roi_trend: "up",
		roi_trend_value: 0,
	};
	const by_channel = roiData?.by_channel || [];
	const by_content = roiData?.by_content || [];
	const monthly_trend = roiData?.monthly_trend || [];

	const channelColumns = [
		{ title: intl.formatMessage({ id: 'roi.column.channel' }), dataIndex: "channel", key: "channel", render: (text: string) => <span className="font-medium">{text}</span> },
		{ title: intl.formatMessage({ id: 'roi.column.investment' }), dataIndex: "investment", key: "investment", render: (value: number) => `¥${value.toLocaleString()}` },
		{ title: intl.formatMessage({ id: 'roi.column.output' }), dataIndex: "return_value", key: "return_value", render: (value: number) => `¥${value.toLocaleString()}` },
		{
			title: intl.formatMessage({ id: 'roi.chart.roi' }),
			dataIndex: "roi",
			key: "roi",
			render: (roi: number) => (
				<Space>
					<Progress percent={Math.min(roi / 3, 100)} size="small" strokeColor={roi >= 200 ? "#52c41a" : roi >= 100 ? "#1890ff" : "#faad14"} format={() => `${roi}%`} style={{ width: 100 }} />
				</Space>
			),
		},
		{ title: intl.formatMessage({ id: 'common.column.trend' }), dataIndex: "trend", key: "trend", render: (trend: string) => trend === "up" ? <RiseOutlined style={{ color: "#52c41a" }} /> : <FallOutlined style={{ color: "#ff4d4f" }} /> },
	];

	const contentColumns = [
		{ title: intl.formatMessage({ id: 'roi.column.content' }), dataIndex: "content", key: "content", render: (text: string) => <span className="font-medium">{text}</span> },
		{ title: intl.formatMessage({ id: 'roi.column.investment' }), dataIndex: "investment", key: "investment", render: (value: number) => `¥${value.toLocaleString()}` },
		{ title: intl.formatMessage({ id: 'roi.column.output' }), dataIndex: "return_value", key: "return_value", render: (value: number) => `¥${value.toLocaleString()}` },
		{ title: intl.formatMessage({ id: 'roi.chart.roi' }), dataIndex: "roi", key: "roi", render: (roi: number) => <Tag color={roi >= 500 ? "green" : roi >= 300 ? "blue" : roi >= 100 ? "orange" : "red"}>{roi}%</Tag> },
		{ title: intl.formatMessage({ id: 'roi.column.views' }), dataIndex: "views", key: "views", render: (value: number) => value.toLocaleString() },
		{ title: intl.formatMessage({ id: 'roi.column.conversions' }), dataIndex: "conversions", key: "conversions", render: (value: number) => value.toLocaleString() },
		{ title: intl.formatMessage({ id: 'roi.column.conversionRate' }), key: "conversion_rate", render: (_: any, record: any) => <span>{((record.conversions / record.views) * 100).toFixed(2)}%</span> },
	];

	const getTrendChartOption = () => ({
		title: { text: intl.formatMessage({ id: 'roi.chart.monthlyTrend' }), left: "center" },
		tooltip: { trigger: "axis" },
		legend: { bottom: 0, data: [intl.formatMessage({ id: 'roi.chart.investment' }), intl.formatMessage({ id: 'roi.chart.output' }), intl.formatMessage({ id: 'roi.chart.roi' })] },
		xAxis: { type: "category", data: monthly_trend.map((m: any) => m.month) },
		yAxis: [
			{ type: "value", name: intl.formatMessage({ id: 'roi.chart.amount' }), position: "left" },
			{ type: "value", name: "ROI (%)", position: "right" },
		],
		series: [
			{ name: intl.formatMessage({ id: 'roi.chart.investment' }), type: "bar", data: monthly_trend.map((m: any) => m.investment), itemStyle: { color: "#1890ff" } },
			{ name: intl.formatMessage({ id: 'roi.chart.output' }), type: "bar", data: monthly_trend.map((m: any) => m.return), itemStyle: { color: "#52c41a" } },
			{ name: intl.formatMessage({ id: 'roi.chart.roi' }), type: "line", yAxisIndex: 1, data: monthly_trend.map((m: any) => m.roi), itemStyle: { color: "#ff7a45" }, lineStyle: { width: 3 } },
		],
	});

	const getChannelChartOption = () => ({
		title: { text: intl.formatMessage({ id: 'roi.chart.channelROI' }), left: "center" },
		tooltip: { trigger: "item" },
		legend: { bottom: 0 },
		series: [{
			type: "pie",
			radius: ["40%", "70%"],
			avoidLabelOverlap: false,
			itemStyle: { borderRadius: 10, borderColor: "#fff", borderWidth: 2 },
			label: { show: false, position: "center" },
			emphasis: { label: { show: true, fontSize: 20, fontWeight: "bold" } },
			labelLine: { show: false },
			data: by_channel.map((c: any) => ({ value: c.roi, name: c.channel })),
		}],
	});

	const getContentChartOption = () => ({
		title: { text: intl.formatMessage({ id: 'roi.chart.contentROI' }), left: "center" },
		tooltip: { trigger: "axis" },
		xAxis: { type: "value" },
		yAxis: { type: "category", data: by_content.map((c: any) => c.content), inverse: true },
		series: [{
			type: "bar",
			data: by_content.map((c: any) => ({ value: c.roi, itemStyle: { color: c.roi >= 500 ? "#52c41a" : c.roi >= 300 ? "#1890ff" : "#faad14" } })),
		}],
	});

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'nav.monitor.roi' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'roi.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'roi.stat.totalInvestment' })} value={summary.total_investment} prefix="¥" valueStyle={{ color: "#1890ff" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'roi.stat.totalReturn' })} value={summary.total_return} prefix="¥" valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'roi.stat.roi' })} value={summary.roi_percentage} suffix="%" prefix={<RiseOutlined />} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'roi.chart.roi' })} value={summary.roi_trend_value} suffix="%" prefix={summary.roi_trend === "up" ? <RiseOutlined /> : <FallOutlined />} valueStyle={{ color: summary.roi_trend === "up" ? "#52c41a" : "#ff4d4f" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'roi.stat.channels' })} value={by_channel.length} prefix={<PieChartOutlined />} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'roi.stat.contents' })} value={by_content.length} prefix={<BarChartOutlined />} /></Card></Col>
			</Row>

			<Card className="mb-4">
				<Space>
					<RangePicker value={dateRange} onChange={(dates) => setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs])} />
					<Select placeholder={intl.formatMessage({ id: 'roi.column.channel' })} allowClear style={{ width: 150 }} onChange={setSelectedChannel}>
						{by_channel.map((c: any) => <Option key={c.channel} value={c.channel}>{c.channel}</Option>)}
					</Select>
					<Button icon={<FilterOutlined />} onClick={() => message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }))}>{intl.formatMessage({ id: 'common.action.search' })}</Button>
					<Button icon={<ExportOutlined />} onClick={() => {
						const csvContent = [[intl.formatMessage({ id: 'roi.column.channel' }), intl.formatMessage({ id: 'roi.column.investment' }), intl.formatMessage({ id: 'roi.column.output' }), intl.formatMessage({ id: 'roi.chart.roi' }), intl.formatMessage({ id: 'common.column.trend' })].join(","), ...by_channel.map((c: any) => [c.channel, c.investment, c.return_value, `${c.roi}%`, c.trend].join(","))].join("\n");
						const blob = new Blob(["\uFEFF" + csvContent], { type: "text/csv;charset=utf-8;" });
						const url = URL.createObjectURL(blob);
						const a = document.createElement("a");
						a.href = url;
						a.download = `roi_report_${new Date().toISOString().split("T")[0]}.csv`;
						a.click();
						message.success(intl.formatMessage({ id: 'common.message.exportSuccess' }));
					}}>{intl.formatMessage({ id: 'common.action.export' })}</Button>
				</Space>
			</Card>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24}><Card><ReactECharts option={getTrendChartOption()} style={{ height: 350 }} /></Card></Col>
			</Row>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={12}><Card><ReactECharts option={getChannelChartOption()} style={{ height: 300 }} /></Card></Col>
				<Col xs={24} lg={12}><Card><ReactECharts option={getContentChartOption()} style={{ height: 300 }} /></Card></Col>
			</Row>

			<Tabs
				activeKey={activeTab}
				onChange={setActiveTab}
				items={[
					{
						key: "overview",
						label: intl.formatMessage({ id: 'roi.column.channel' }),
						children: (
							<Card>
								<Table
									columns={channelColumns}
									dataSource={by_channel}
									rowKey="channel"
									pagination={false}
									summary={() => (
										<Table.Summary.Row>
											<Table.Summary.Cell index={0}><span className="font-bold">{intl.formatMessage({ id: 'common.summaryTotal' })}</span></Table.Summary.Cell>
											<Table.Summary.Cell index={1}><span className="font-bold">¥{by_channel.reduce((sum: number, c: any) => sum + c.investment, 0).toLocaleString()}</span></Table.Summary.Cell>
											<Table.Summary.Cell index={2}><span className="font-bold">¥{by_channel.reduce((sum: number, c: any) => sum + c.return_value, 0).toLocaleString()}</span></Table.Summary.Cell>
											<Table.Summary.Cell index={3}><span className="font-bold">{Math.round((by_channel.reduce((sum: number, c: any) => sum + c.return_value, 0) / by_channel.reduce((sum: number, c: any) => sum + c.investment, 0)) * 100)}%</span></Table.Summary.Cell>
											<Table.Summary.Cell index={4} />
										</Table.Summary.Row>
									)}
								/>
							</Card>
						),
					},
					{
						key: "content",
						label: intl.formatMessage({ id: 'roi.column.content' }),
						children: (
							<Card>
								<Table columns={contentColumns} dataSource={by_content} rowKey="content" pagination={false} />
							</Card>
						),
					},
				]}
			/>

			<Card title={intl.formatMessage({ id: 'roi.suggestions.title' })} className="mt-4">
				<Row gutter={[16, 16]}>
					<Col xs={24} md={8}>
						<Card size="small" type="inner" title={intl.formatMessage({ id: 'roi.suggestions.channelOptimize' })}>
							<ul className="list-disc pl-4 space-y-2">
								{by_channel.length > 0 && by_channel.sort((a: any, b: any) => b.roi - a.roi)[0] && (
									<li>{intl.formatMessage({ id: 'roi.suggestions.topChannel' }, { channel: by_channel.sort((a: any, b: any) => b.roi - a.roi)[0].channel, roi: by_channel.sort((a: any, b: any) => b.roi - a.roi)[0].roi })}</li>
								)}
								{by_channel.filter((c: any) => c.trend === "down").length > 0 && (
									<li>{intl.formatMessage({ id: 'roi.suggestions.decliningChannel' }, { channels: by_channel.filter((c: any) => c.trend === "down").map((c: any) => c.channel).join(", ") })}</li>
								)}
								{by_channel.length > 0 && by_channel.sort((a: any, b: any) => a.roi - b.roi)[0] && (
									<li>{intl.formatMessage({ id: 'roi.suggestions.lowChannel' }, { channel: by_channel.sort((a: any, b: any) => a.roi - b.roi)[0].channel, roi: by_channel.sort((a: any, b: any) => a.roi - b.roi)[0].roi })}</li>
								)}
							</ul>
						</Card>
					</Col>
					<Col xs={24} md={8}>
						<Card size="small" type="inner" title={intl.formatMessage({ id: 'roi.suggestions.contentOptimize' })}>
							<ul className="list-disc pl-4 space-y-2">
								{by_content.length > 0 && by_content.sort((a: any, b: any) => b.roi - a.roi)[0] && (
									<li>{intl.formatMessage({ id: 'roi.suggestions.topContent' }, { content: by_content.sort((a: any, b: any) => b.roi - a.roi)[0].content, roi: by_content.sort((a: any, b: any) => b.roi - a.roi)[0].roi })}</li>
								)}
								{by_content.filter((c: any) => c.conversions / c.views > 0.005).length > 0 && <li>{intl.formatMessage({ id: 'roi.suggestions.highConversion' })}</li>}
								{by_content.length > 0 && by_content.sort((a: any, b: any) => a.roi - b.roi)[0] && (
									<li>{intl.formatMessage({ id: 'roi.suggestions.optimizeContent' }, { content: by_content.sort((a: any, b: any) => a.roi - b.roi)[0].content })}</li>
								)}
							</ul>
						</Card>
					</Col>
					<Col xs={24} md={8}>
						<Card size="small" type="inner" title={intl.formatMessage({ id: 'roi.suggestions.overallStrategy' })}>
							<ul className="list-disc pl-4 space-y-2">
								{summary.roi_trend === "up" ? (
									<li>{intl.formatMessage({ id: 'roi.suggestions.roiUp' }, { value: summary.roi_trend_value })}</li>
								) : (
									<li>{intl.formatMessage({ id: 'roi.suggestions.roiDown' }, { value: summary.roi_trend_value })}</li>
								)}
								{monthly_trend.length >= 2 && monthly_trend[monthly_trend.length - 1].roi > monthly_trend[monthly_trend.length - 2].roi && (
									<li>{intl.formatMessage({ id: 'roi.suggestions.roiImproved' })}</li>
								)}
								<li>{intl.formatMessage({ id: 'roi.suggestions.continueInvest' })}</li>
							</ul>
						</Card>
					</Col>
				</Row>
			</Card>
		</div>
	);
}
