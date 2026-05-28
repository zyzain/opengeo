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

const { RangePicker } = DatePicker;
const { Option } = Select;

export default function ROIPage() {
	const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs] | null>(
		null,
	);
	const [selectedChannel, setSelectedChannel] = useState<string | undefined>(
		undefined,
	);
	const [activeTab, setActiveTab] = useState("overview");

	const { data: roiData, isLoading } = useROIMetrics({
		start_date: dateRange?.[0]?.toISOString(),
		end_date: dateRange?.[1]?.toISOString(),
		channel_id: selectedChannel,
	});

	// 使用API数据或默认值
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

	// 渠道表格列
	const channelColumns = [
		{
			title: "渠道",
			dataIndex: "channel",
			key: "channel",
			render: (text: string) => <span className="font-medium">{text}</span>,
		},
		{
			title: "投入",
			dataIndex: "investment",
			key: "investment",
			render: (value: number) => `¥${value.toLocaleString()}`,
		},
		{
			title: "产出",
			dataIndex: "return_value",
			key: "return_value",
			render: (value: number) => `¥${value.toLocaleString()}`,
		},
		{
			title: "ROI",
			dataIndex: "roi",
			key: "roi",
			render: (roi: number) => (
				<Space>
					<Progress
						percent={Math.min(roi / 3, 100)}
						size="small"
						strokeColor={
							roi >= 200 ? "#52c41a" : roi >= 100 ? "#1890ff" : "#faad14"
						}
						format={() => `${roi}%`}
						style={{ width: 100 }}
					/>
				</Space>
			),
		},
		{
			title: "趋势",
			dataIndex: "trend",
			key: "trend",
			render: (trend: string) =>
				trend === "up" ? (
					<RiseOutlined style={{ color: "#52c41a" }} />
				) : (
					<FallOutlined style={{ color: "#ff4d4f" }} />
				),
		},
	];

	// 内容表格列
	const contentColumns = [
		{
			title: "内容",
			dataIndex: "content",
			key: "content",
			render: (text: string) => <span className="font-medium">{text}</span>,
		},
		{
			title: "投入",
			dataIndex: "investment",
			key: "investment",
			render: (value: number) => `¥${value.toLocaleString()}`,
		},
		{
			title: "产出",
			dataIndex: "return_value",
			key: "return_value",
			render: (value: number) => `¥${value.toLocaleString()}`,
		},
		{
			title: "ROI",
			dataIndex: "roi",
			key: "roi",
			render: (roi: number) => (
				<Tag
					color={
						roi >= 500
							? "green"
							: roi >= 300
								? "blue"
								: roi >= 100
									? "orange"
									: "red"
					}
				>
					{roi}%
				</Tag>
			),
		},
		{
			title: "浏览量",
			dataIndex: "views",
			key: "views",
			render: (value: number) => value.toLocaleString(),
		},
		{
			title: "转化数",
			dataIndex: "conversions",
			key: "conversions",
			render: (value: number) => value.toLocaleString(),
		},
		{
			title: "转化率",
			key: "conversion_rate",
			render: (_: any, record: any) => (
				<span>{((record.conversions / record.views) * 100).toFixed(2)}%</span>
			),
		},
	];

	// ROI趋势图表
	const getTrendChartOption = () => ({
		title: {
			text: "月度ROI趋势",
			left: "center",
		},
		tooltip: {
			trigger: "axis",
		},
		legend: {
			bottom: 0,
			data: ["投入", "产出", "ROI"],
		},
		xAxis: {
			type: "category",
			data: monthly_trend.map((m: any) => m.month),
		},
		yAxis: [
			{
				type: "value",
				name: "金额 (¥)",
				position: "left",
			},
			{
				type: "value",
				name: "ROI (%)",
				position: "right",
			},
		],
		series: [
			{
				name: "投入",
				type: "bar",
				data: monthly_trend.map((m: any) => m.investment),
				itemStyle: { color: "#1890ff" },
			},
			{
				name: "产出",
				type: "bar",
				data: monthly_trend.map((m: any) => m.return),
				itemStyle: { color: "#52c41a" },
			},
			{
				name: "ROI",
				type: "line",
				yAxisIndex: 1,
				data: monthly_trend.map((m: any) => m.roi),
				itemStyle: { color: "#ff7a45" },
				lineStyle: { width: 3 },
			},
		],
	});

	// 渠道ROI图表
	const getChannelChartOption = () => ({
		title: {
			text: "渠道ROI对比",
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
				radius: ["40%", "70%"],
				avoidLabelOverlap: false,
				itemStyle: {
					borderRadius: 10,
					borderColor: "#fff",
					borderWidth: 2,
				},
				label: {
					show: false,
					position: "center",
				},
				emphasis: {
					label: {
						show: true,
						fontSize: 20,
						fontWeight: "bold",
					},
				},
				labelLine: {
					show: false,
				},
				data: by_channel.map((c: any) => ({
					value: c.roi,
					name: c.channel,
				})),
			},
		],
	});

	// 内容ROI图表
	const getContentChartOption = () => ({
		title: {
			text: "内容ROI排名",
			left: "center",
		},
		tooltip: {
			trigger: "axis",
		},
		xAxis: {
			type: "value",
		},
		yAxis: {
			type: "category",
			data: by_content.map((c: any) => c.content),
			inverse: true,
		},
		series: [
			{
				type: "bar",
				data: by_content.map((c: any) => ({
					value: c.roi,
					itemStyle: {
						color:
							c.roi >= 500 ? "#52c41a" : c.roi >= 300 ? "#1890ff" : "#faad14",
					},
				})),
			},
		],
	});

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">ROI分析</h1>
				<p className="text-gray-500 mt-1">
					关联发布内容与转化数据，量化GEO投入产出比
				</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="总投入"
							value={summary.total_investment}
							prefix="¥"
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="总产出"
							value={summary.total_return}
							prefix="¥"
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="ROI"
							value={summary.roi_percentage}
							suffix="%"
							prefix={<RiseOutlined />}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="ROI趋势"
							value={summary.roi_trend_value}
							suffix="%"
							prefix={
								summary.roi_trend === "up" ? <RiseOutlined /> : <FallOutlined />
							}
							valueStyle={{
								color: summary.roi_trend === "up" ? "#52c41a" : "#ff4d4f",
							}}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="渠道数量"
							value={by_channel.length}
							prefix={<PieChartOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="内容数量"
							value={by_content.length}
							prefix={<BarChartOutlined />}
						/>
					</Card>
				</Col>
			</Row>

			{/* 筛选条件 */}
			<Card className="mb-4">
				<Space>
					<RangePicker
						value={dateRange}
						onChange={(dates) =>
							setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs])
						}
					/>
					<Select
						placeholder="选择渠道"
						allowClear
						style={{ width: 150 }}
						onChange={setSelectedChannel}
					>
						{by_channel.map((c: any) => (
							<Option key={c.channel} value={c.channel}>
								{c.channel}
							</Option>
						))}
					</Select>
					<Button
						icon={<FilterOutlined />}
						onClick={() => message.success("筛选已应用")}
					>
						筛选
					</Button>
					<Button
						icon={<ExportOutlined />}
						onClick={() => {
							const csvContent = [
								["渠道", "投入", "产出", "ROI", "趋势"].join(","),
								...by_channel.map((c: any) =>
									[
										c.channel,
										c.investment,
										c.return_value,
										`${c.roi}%`,
										c.trend,
									].join(","),
								),
							].join("\n");
							const blob = new Blob(["\uFEFF" + csvContent], {
								type: "text/csv;charset=utf-8;",
							});
							const url = URL.createObjectURL(blob);
							const a = document.createElement("a");
							a.href = url;
							a.download = `roi_report_${new Date().toISOString().split("T")[0]}.csv`;
							a.click();
							message.success("导出成功");
						}}
					>
						导出
					</Button>
				</Space>
			</Card>

			{/* 图表 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24}>
					<Card>
						<ReactECharts
							option={getTrendChartOption()}
							style={{ height: 350 }}
						/>
					</Card>
				</Col>
			</Row>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={12}>
					<Card>
						<ReactECharts
							option={getChannelChartOption()}
							style={{ height: 300 }}
						/>
					</Card>
				</Col>
				<Col xs={24} lg={12}>
					<Card>
						<ReactECharts
							option={getContentChartOption()}
							style={{ height: 300 }}
						/>
					</Card>
				</Col>
			</Row>

			{/* 详细数据 */}
			<Tabs
				activeKey={activeTab}
				onChange={setActiveTab}
				items={[
					{
						key: "overview",
						label: "渠道分析",
						children: (
							<Card>
								<Table
									columns={channelColumns}
									dataSource={by_channel}
									rowKey="channel"
									pagination={false}
									summary={() => (
										<Table.Summary.Row>
											<Table.Summary.Cell index={0}>
												<span className="font-bold">合计</span>
											</Table.Summary.Cell>
											<Table.Summary.Cell index={1}>
												<span className="font-bold">
													¥
													{by_channel
														.reduce(
															(sum: number, c: any) => sum + c.investment,
															0,
														)
														.toLocaleString()}
												</span>
											</Table.Summary.Cell>
											<Table.Summary.Cell index={2}>
												<span className="font-bold">
													¥
													{by_channel
														.reduce(
															(sum: number, c: any) => sum + c.return_value,
															0,
														)
														.toLocaleString()}
												</span>
											</Table.Summary.Cell>
											<Table.Summary.Cell index={3}>
												<span className="font-bold">
													{Math.round(
														(by_channel.reduce(
															(sum: number, c: any) => sum + c.return_value,
															0,
														) /
															by_channel.reduce(
																(sum: number, c: any) => sum + c.investment,
																0,
															)) *
															100,
													)}
													%
												</span>
											</Table.Summary.Cell>
											<Table.Summary.Cell index={4} />
										</Table.Summary.Row>
									)}
								/>
							</Card>
						),
					},
					{
						key: "content",
						label: "内容分析",
						children: (
							<Card>
								<Table
									columns={contentColumns}
									dataSource={by_content}
									rowKey="content"
									pagination={false}
								/>
							</Card>
						),
					},
				]}
			/>

			{/* 优化建议 - 基于数据动态生成 */}
			<Card title="优化建议" className="mt-4">
				<Row gutter={[16, 16]}>
					<Col xs={24} md={8}>
						<Card size="small" type="inner" title="渠道优化">
							<ul className="list-disc pl-4 space-y-2">
								{by_channel.length > 0 &&
									by_channel.sort((a: any, b: any) => b.roi - a.roi)[0] && (
										<li>
											{
												by_channel.sort((a: any, b: any) => b.roi - a.roi)[0]
													.channel
											}
											ROI最高(
											{
												by_channel.sort((a: any, b: any) => b.roi - a.roi)[0]
													.roi
											}
											%)，建议增加投入
										</li>
									)}
								{by_channel.filter((c: any) => c.trend === "down").length >
									0 && (
									<li>
										{by_channel
											.filter((c: any) => c.trend === "down")
											.map((c: any) => c.channel)
											.join("、")}
										ROI下降，需要优化内容策略
									</li>
								)}
								{by_channel.length > 0 &&
									by_channel.sort((a: any, b: any) => a.roi - b.roi)[0] && (
										<li>
											{
												by_channel.sort((a: any, b: any) => a.roi - b.roi)[0]
													.channel
											}
											ROI较低(
											{
												by_channel.sort((a: any, b: any) => a.roi - b.roi)[0]
													.roi
											}
											%)，考虑减少投入或调整方向
										</li>
									)}
							</ul>
						</Card>
					</Col>
					<Col xs={24} md={8}>
						<Card size="small" type="inner" title="内容优化">
							<ul className="list-disc pl-4 space-y-2">
								{by_content.length > 0 &&
									by_content.sort((a: any, b: any) => b.roi - a.roi)[0] && (
										<li>
											《
											{
												by_content.sort((a: any, b: any) => b.roi - a.roi)[0]
													.content
											}
											》ROI最高(
											{
												by_content.sort((a: any, b: any) => b.roi - a.roi)[0]
													.roi
											}
											%)，建议复制模式
										</li>
									)}
								{by_content.filter((c: any) => c.conversions / c.views > 0.005)
									.length > 0 && <li>高转化率内容应增加产出</li>}
								{by_content.length > 0 &&
									by_content.sort((a: any, b: any) => a.roi - b.roi)[0] && (
										<li>
											优化《
											{
												by_content.sort((a: any, b: any) => a.roi - b.roi)[0]
													.content
											}
											》的转化路径
										</li>
									)}
							</ul>
						</Card>
					</Col>
					<Col xs={24} md={8}>
						<Card size="small" type="inner" title="整体策略">
							<ul className="list-disc pl-4 space-y-2">
								{summary.roi_trend === "up" ? (
									<li>
										总ROI持续上升(+{summary.roi_trend_value}%)，策略方向正确
									</li>
								) : (
									<li>总ROI下降({summary.roi_trend_value}%)，需要调整策略</li>
								)}
								{monthly_trend.length >= 2 &&
									monthly_trend[monthly_trend.length - 1].roi >
										monthly_trend[monthly_trend.length - 2].roi && (
										<li>本月ROI较上月提升，归因于内容优化</li>
									)}
								<li>建议继续投入GEO优化和AI适配</li>
							</ul>
						</Card>
					</Col>
				</Row>
			</Card>
		</div>
	);
}
