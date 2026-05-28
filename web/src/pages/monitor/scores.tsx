"use client";

import { useSourceScores } from "@/hooks";
import {
	FallOutlined,
	FilterOutlined,
	InfoCircleOutlined,
	ReloadOutlined,
	RiseOutlined,
	StarOutlined,
	TrophyOutlined,
} from "@ant-design/icons";
import {
	Badge,
	Button,
	Card,
	Col,
	Progress,
	Row,
	Select,
	Space,
	Spin,
	Statistic,
	Table,
	Tag,
	Tooltip,
} from "antd";
import ReactECharts from "echarts-for-react";
import { useState } from "react";

const { Option } = Select;

const platformColors: Record<string, string> = {
	wechat: "green",
	zhihu: "blue",
	douyin: "purple",
	xiaohongshu: "pink",
	weibo: "red",
};

const platformNames: Record<string, string> = {
	wechat: "微信",
	zhihu: "知乎",
	douyin: "抖音",
	xiaohongshu: "小红书",
	weibo: "微博",
};

export default function ScoresPage() {
	const [selectedPlatform, setSelectedPlatform] = useState<string | undefined>(
		undefined,
	);

	const { data: scoresData, isLoading, refetch } = useSourceScores();
	const scores = scoresData?.items || [];

	// 过滤数据
	const filteredScores = selectedPlatform
		? scores.filter((s: any) => s.platform === selectedPlatform)
		: scores;

	// 统计数据
	const stats = {
		total: scores.length,
		avgScore:
			scores.length > 0
				? (
						scores.reduce((sum: number, s: any) => sum + (s.score || 0), 0) /
						scores.length
					).toFixed(1)
				: "0",
		highest:
			scores.length > 0 ? Math.max(...scores.map((s: any) => s.score || 0)) : 0,
		lowest:
			scores.length > 0 ? Math.min(...scores.map((s: any) => s.score || 0)) : 0,
		improving: scores.filter((s: any) => s.trend === "up").length,
		declining: scores.filter((s: any) => s.trend === "down").length,
	};

	// 评分维度
	const dimensionLabels: Record<string, string> = {
		recency_speed: "收录速度",
		ranking_stability: "排名稳定性",
		citation_frequency: "引用频次",
		authority_score: "权威性评分",
		content_quality: "内容质量",
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
			title: "渠道",
			dataIndex: "channel_name",
			key: "channel_name",
			render: (text: string, record: any) => (
				<Space>
					<Tag color={platformColors[record.platform]}>
						{platformNames[record.platform]}
					</Tag>
					<span className="font-medium">{text}</span>
				</Space>
			),
		},
		{
			title: "账号",
			dataIndex: "account_name",
			key: "account_name",
		},
		{
			title: "综合评分",
			dataIndex: "score",
			key: "score",
			width: 150,
			render: (score: number) => (
				<div>
					<Progress
						percent={score}
						size="small"
						strokeColor={
							score >= 90
								? "#52c41a"
								: score >= 80
									? "#1890ff"
									: score >= 70
										? "#faad14"
										: "#ff4d4f"
						}
						format={(percent) => <span className="font-bold">{percent}</span>}
					/>
				</div>
			),
		},
		{
			title: "趋势",
			dataIndex: "trend",
			key: "trend",
			width: 120,
			render: (trend: string, record: any) => (
				<Space>
					{trend === "up" ? (
						<RiseOutlined style={{ color: "#52c41a" }} />
					) : (
						<FallOutlined style={{ color: "#ff4d4f" }} />
					)}
					<span style={{ color: trend === "up" ? "#52c41a" : "#ff4d4f" }}>
						{trend === "up" ? "+" : "-"}
						{record.trend_value}%
					</span>
				</Space>
			),
		},
		{
			title: "各维度评分",
			key: "dimensions",
			render: (_: any, record: any) => (
				<Space size={[0, 4]} wrap>
					{Object.entries(record.dimensions).map(([key, value]) => (
						<Tooltip key={key} title={`${dimensionLabels[key]}: ${value}`}>
							<Tag
								color={
									(value as number) >= 90
										? "green"
										: (value as number) >= 80
											? "blue"
											: "orange"
								}
							>
								{dimensionLabels[key].slice(0, 2)}: {value as number}
							</Tag>
						</Tooltip>
					))}
				</Space>
			),
		},
		{
			title: "更新时间",
			dataIndex: "updated_at",
			key: "updated_at",
			width: 180,
			render: (text: string) => new Date(text).toLocaleString(),
		},
	];

	// 雷达图配置
	const getRadarOption = () => {
		const indicators = Object.values(dimensionLabels).map((label) => ({
			name: label,
			max: 100,
		}));

		return {
			title: {
				text: "评分维度对比",
				left: "center",
			},
			tooltip: {},
			legend: {
				bottom: 0,
				data: scores.slice(0, 3).map((s: any) => s.account_name),
			},
			radar: {
				indicator: indicators,
				shape: "polygon",
			},
			series: [
				{
					type: "radar",
					data: scores.slice(0, 3).map((s: any) => ({
						value: Object.values(s.dimensions),
						name: s.account_name,
						areaStyle: { opacity: 0.1 },
					})),
				},
			],
		};
	};

	// 柱状图配置
	const getBarOption = () => ({
		title: {
			text: "渠道评分对比",
			left: "center",
		},
		tooltip: {
			trigger: "axis",
		},
		xAxis: {
			type: "category",
			data: scores.map((s: any) => s.channel_name),
			axisLabel: { rotate: 30 },
		},
		yAxis: {
			type: "value",
			min: scores.length > 0 ? Math.max(0, Math.floor(Math.min(...scores.map((s: any) => s.score || 0)) / 10) * 10 - 10) : 0,
			max: 100,
		},
		series: [
			{
				type: "bar",
				data: scores.map((s: any) => ({
					value: s.score,
					itemStyle: {
						color:
							s.score >= 90
								? "#52c41a"
								: s.score >= 80
									? "#1890ff"
									: s.score >= 70
										? "#faad14"
										: "#ff4d4f",
					},
				})),
				barWidth: "50%",
			},
		],
	});

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">信源评分</h1>
				<p className="text-gray-500 mt-1">
					综合评估各渠道和账号的权威性、收录速度、引用频次等指标
				</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="渠道数量"
							value={stats.total}
							prefix={<StarOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="平均评分"
							value={stats.avgScore}
							suffix="分"
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="最高评分"
							value={stats.highest}
							suffix="分"
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="最低评分"
							value={stats.lowest}
							suffix="分"
							valueStyle={{ color: "#ff4d4f" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="上升渠道"
							value={stats.improving}
							prefix={<RiseOutlined />}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={8} lg={4}>
					<Card>
						<Statistic
							title="下降渠道"
							value={stats.declining}
							prefix={<FallOutlined />}
							valueStyle={{ color: "#ff4d4f" }}
						/>
					</Card>
				</Col>
			</Row>

			{/* 图表 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={12}>
					<Card>
						<ReactECharts option={getRadarOption()} style={{ height: 350 }} />
					</Card>
				</Col>
				<Col xs={24} lg={12}>
					<Card>
						<ReactECharts option={getBarOption()} style={{ height: 350 }} />
					</Card>
				</Col>
			</Row>

			{/* 评分详情 */}
			<Card
				title={
					<Space>
						<StarOutlined />
						<span>评分详情</span>
					</Space>
				}
				extra={
					<Space>
						<Select
							placeholder="筛选平台"
							allowClear
							style={{ width: 150 }}
							onChange={setSelectedPlatform}
						>
							{Object.entries(platformNames).map(([key, name]) => (
								<Option key={key} value={key}>
									<Tag color={platformColors[key]}>{name}</Tag>
								</Option>
							))}
						</Select>
						<Button icon={<ReloadOutlined />} onClick={() => refetch()}>
							刷新
						</Button>
					</Space>
				}
			>
				<Table
					columns={columns}
					dataSource={filteredScores}
					rowKey="id"
					loading={isLoading}
					pagination={false}
				/>

				<div className="mt-4 p-4 bg-gray-50 rounded-lg">
					<h4 className="font-medium mb-2">
						<InfoCircleOutlined /> 评分说明
					</h4>
					<Row gutter={[16, 8]}>
						{Object.entries(dimensionLabels).map(([key, label]) => (
							<Col key={key} xs={12} sm={8} lg={4}>
								<div className="text-sm">
									<span className="text-gray-500">{label}：</span>
									<span className="text-gray-700">
										{key === "recency_speed" && "内容被AI收录的速度"}
										{key === "ranking_stability" && "在AI回答中排名的稳定性"}
										{key === "citation_frequency" && "被AI引用的频率"}
										{key === "authority_score" && "内容来源的权威性"}
										{key === "content_quality" && "内容质量和结构化程度"}
									</span>
								</div>
							</Col>
						))}
					</Row>
				</div>
			</Card>
		</div>
	);
}
