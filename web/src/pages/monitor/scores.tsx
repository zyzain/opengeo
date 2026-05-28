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
import { useIntl } from "react-intl";

const { Option } = Select;

const platformColors: Record<string, string> = {
	wechat: "green",
	zhihu: "blue",
	douyin: "purple",
	xiaohongshu: "pink",
	weibo: "red",
};

export default function ScoresPage() {
	const intl = useIntl();
	const [selectedPlatform, setSelectedPlatform] = useState<string | undefined>(undefined);

	const platformNames: Record<string, string> = {
		wechat: intl.formatMessage({ id: 'monitor.platform.wechat' }),
		zhihu: intl.formatMessage({ id: 'monitor.platform.zhihu' }),
		douyin: intl.formatMessage({ id: 'monitor.platform.douyin' }),
		xiaohongshu: intl.formatMessage({ id: 'monitor.platform.xiaohongshu' }),
		weibo: intl.formatMessage({ id: 'monitor.platform.weibo' }),
	};

	const { data: scoresData, isLoading, refetch } = useSourceScores();
	const scores = scoresData?.items || [];

	const filteredScores = selectedPlatform
		? scores.filter((s: any) => s.platform === selectedPlatform)
		: scores;

	const stats = {
		total: scores.length,
		avgScore: scores.length > 0 ? (scores.reduce((sum: number, s: any) => sum + (s.score || 0), 0) / scores.length).toFixed(1) : "0",
		highest: scores.length > 0 ? Math.max(...scores.map((s: any) => s.score || 0)) : 0,
		lowest: scores.length > 0 ? Math.min(...scores.map((s: any) => s.score || 0)) : 0,
		improving: scores.filter((s: any) => s.trend === "up").length,
		declining: scores.filter((s: any) => s.trend === "down").length,
	};

	const dimensionLabels: Record<string, string> = {
		recency_speed: intl.formatMessage({ id: 'scores.dimension.recencySpeed' }),
		ranking_stability: intl.formatMessage({ id: 'scores.dimension.rankingStability' }),
		citation_frequency: intl.formatMessage({ id: 'scores.dimension.citationFrequency' }),
		authority_score: intl.formatMessage({ id: 'scores.dimension.authorityScore' }),
		content_quality: intl.formatMessage({ id: 'scores.dimension.contentQuality' }),
	};

	const columns = [
		{
			title: intl.formatMessage({ id: 'scores.column.rank' }),
			key: "rank",
			width: 80,
			render: (_: any, __: any, index: number) => (
				<div className="text-center">
					{index < 3 ? (
						<TrophyOutlined style={{ color: ["#ffd700", "#c0c0c0", "#cd7f32"][index], fontSize: 20 }} />
					) : (
						<span className="text-gray-500">{index + 1}</span>
					)}
				</div>
			),
		},
		{
			title: intl.formatMessage({ id: 'scores.column.channel' }),
			dataIndex: "channel_name",
			key: "channel_name",
			render: (text: string, record: any) => (
				<Space>
					<Tag color={platformColors[record.platform]}>{platformNames[record.platform]}</Tag>
					<span className="font-medium">{text}</span>
				</Space>
			),
		},
		{ title: intl.formatMessage({ id: 'scores.column.account' }), dataIndex: "account_name", key: "account_name" },
		{
			title: intl.formatMessage({ id: 'scores.column.score' }),
			dataIndex: "score",
			key: "score",
			width: 150,
			render: (score: number) => (
				<div>
					<Progress
						percent={score}
						size="small"
						strokeColor={score >= 90 ? "#52c41a" : score >= 80 ? "#1890ff" : score >= 70 ? "#faad14" : "#ff4d4f"}
						format={(percent) => <span className="font-bold">{percent}</span>}
					/>
				</div>
			),
		},
		{
			title: intl.formatMessage({ id: 'common.column.trend' }),
			dataIndex: "trend",
			key: "trend",
			width: 120,
			render: (trend: string, record: any) => (
				<Space>
					{trend === "up" ? <RiseOutlined style={{ color: "#52c41a" }} /> : <FallOutlined style={{ color: "#ff4d4f" }} />}
					<span style={{ color: trend === "up" ? "#52c41a" : "#ff4d4f" }}>{trend === "up" ? "+" : "-"}{record.trend_value}%</span>
				</Space>
			),
		},
		{
			title: intl.formatMessage({ id: 'scores.column.dimensions' }),
			key: "dimensions",
			render: (_: any, record: any) => (
				<Space size={[0, 4]} wrap>
					{Object.entries(record.dimensions).map(([key, value]) => (
						<Tooltip key={key} title={`${dimensionLabels[key]}: ${value}`}>
							<Tag color={(value as number) >= 90 ? "green" : (value as number) >= 80 ? "blue" : "orange"}>
								{dimensionLabels[key].slice(0, 2)}: {value as number}
							</Tag>
						</Tooltip>
					))}
				</Space>
			),
		},
		{ title: intl.formatMessage({ id: 'scores.column.updatedAt' }), dataIndex: "updated_at", key: "updated_at", width: 180, render: (text: string) => new Date(text).toLocaleString() },
	];

	const getRadarOption = () => {
		const indicators = Object.values(dimensionLabels).map((label) => ({ name: label, max: 100 }));
		return {
			title: { text: intl.formatMessage({ id: 'scores.chart.radarTitle' }), left: "center" },
			tooltip: {},
			legend: { bottom: 0, data: scores.slice(0, 3).map((s: any) => s.account_name) },
			radar: { indicator: indicators, shape: "polygon" },
			series: [{
				type: "radar",
				data: scores.slice(0, 3).map((s: any) => ({
					value: Object.values(s.dimensions),
					name: s.account_name,
					areaStyle: { opacity: 0.1 },
				})),
			}],
		};
	};

	const getBarOption = () => ({
		title: { text: intl.formatMessage({ id: 'scores.chart.barTitle' }), left: "center" },
		tooltip: { trigger: "axis" },
		xAxis: { type: "category", data: scores.map((s: any) => s.channel_name), axisLabel: { rotate: 30 } },
		yAxis: {
			type: "value",
			min: scores.length > 0 ? Math.max(0, Math.floor(Math.min(...scores.map((s: any) => s.score || 0)) / 10) * 10 - 10) : 0,
			max: 100,
		},
		series: [{
			type: "bar",
			data: scores.map((s: any) => ({
				value: s.score,
				itemStyle: { color: s.score >= 90 ? "#52c41a" : s.score >= 80 ? "#1890ff" : s.score >= 70 ? "#faad14" : "#ff4d4f" },
			})),
			barWidth: "50%",
		}],
	});

	const dimensionExplanations: Record<string, string> = {
		recency_speed: intl.formatMessage({ id: 'scores.explanation.recencySpeed' }),
		ranking_stability: intl.formatMessage({ id: 'scores.explanation.rankingStability' }),
		citation_frequency: intl.formatMessage({ id: 'scores.explanation.citationFrequency' }),
		authority_score: intl.formatMessage({ id: 'scores.explanation.authorityScore' }),
		content_quality: intl.formatMessage({ id: 'scores.explanation.contentQuality' }),
	};

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'scores.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'scores.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'scores.stat.channelCount' })} value={stats.total} prefix={<StarOutlined />} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'scores.stat.avgScore' })} value={stats.avgScore} suffix={intl.formatMessage({ id: 'common.unit.score' })} valueStyle={{ color: "#1890ff" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'scores.stat.highest' })} value={stats.highest} suffix={intl.formatMessage({ id: 'common.unit.score' })} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'scores.stat.lowest' })} value={stats.lowest} suffix={intl.formatMessage({ id: 'common.unit.score' })} valueStyle={{ color: "#ff4d4f" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'scores.stat.improving' })} value={stats.improving} prefix={<RiseOutlined />} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'scores.stat.declining' })} value={stats.declining} prefix={<FallOutlined />} valueStyle={{ color: "#ff4d4f" }} /></Card></Col>
			</Row>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={12}><Card><ReactECharts option={getRadarOption()} style={{ height: 350 }} /></Card></Col>
				<Col xs={24} lg={12}><Card><ReactECharts option={getBarOption()} style={{ height: 350 }} /></Card></Col>
			</Row>

			<Card
				title={<Space><StarOutlined /><span>{intl.formatMessage({ id: 'scores.section.detail' })}</span></Space>}
				extra={
					<Space>
						<Select placeholder={intl.formatMessage({ id: 'scores.placeholder.filterPlatform' })} allowClear style={{ width: 150 }} onChange={setSelectedPlatform}>
							{Object.entries(platformNames).map(([key, name]) => (
								<Option key={key} value={key}><Tag color={platformColors[key]}>{name}</Tag></Option>
							))}
						</Select>
						<Button icon={<ReloadOutlined />} onClick={() => refetch()}>{intl.formatMessage({ id: 'common.action.refresh' })}</Button>
					</Space>
				}
			>
				<Table columns={columns} dataSource={filteredScores} rowKey="id" loading={isLoading} pagination={false} />

				<div className="mt-4 p-4 bg-gray-50 rounded-lg">
					<h4 className="font-medium mb-2"><InfoCircleOutlined /> {intl.formatMessage({ id: 'scores.section.explanation' })}</h4>
					<Row gutter={[16, 8]}>
						{Object.entries(dimensionLabels).map(([key, label]) => (
							<Col key={key} xs={12} sm={8} lg={4}>
								<div className="text-sm">
									<span className="text-gray-500">{label}：</span>
									<span className="text-gray-700">{dimensionExplanations[key]}</span>
								</div>
							</Col>
						))}
					</Row>
				</div>
			</Card>
		</div>
	);
}
