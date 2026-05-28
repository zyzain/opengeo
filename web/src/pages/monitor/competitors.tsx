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
import { useIntl } from "react-intl";

const { TextArea } = Input;

export default function CompetitorsPage() {
	const intl = useIntl();
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

	const stats = {
		total: competitors.length,
		active: competitors.filter((c: any) => c.status === "active").length,
		avgVisibility: competitors.length > 0 ? (competitors.reduce((sum: number, c: any) => sum + (c.visibility_score || 0), 0) / competitors.length).toFixed(1) : "0",
		totalGaps: competitors.reduce((sum: number, c: any) => sum + (c.content_gap_count || 0), 0),
		ourScore: ourScoreData?.data?.data?.score || 0,
		ourRank: ourScoreData?.data?.data?.rank || 0,
	};

	const columns = [
		{
			title: intl.formatMessage({ id: 'competitor.column.rank' }),
			key: "rank",
			width: 80,
			render: (_: any, __: any, index: number) => (
				<div className="text-center">
					{index < 3 ? <TrophyOutlined style={{ color: ["#ffd700", "#c0c0c0", "#cd7f32"][index], fontSize: 20 }} /> : <span className="text-gray-500">{index + 1}</span>}
				</div>
			),
		},
		{
			title: intl.formatMessage({ id: 'competitor.column.name' }),
			dataIndex: "competitor_name",
			key: "competitor_name",
			render: (text: string, record: any) => (
				<div>
					<a onClick={() => handleShowDetail(record)} className="text-blue-500 font-medium">{text}</a>
					<div className="text-gray-400 text-xs">{record.competitor_domain}</div>
				</div>
			),
		},
		{
			title: intl.formatMessage({ id: 'competitor.column.visibilityScore' }),
			dataIndex: "visibility_score",
			key: "visibility_score",
			width: 150,
			render: (score: number) => <Progress percent={score} size="small" strokeColor={score >= 80 ? "#52c41a" : score >= 70 ? "#1890ff" : "#faad14"} format={(percent) => <span className="font-bold">{percent}</span>} />,
		},
		{
			title: intl.formatMessage({ id: 'competitor.column.contentGap' }),
			dataIndex: "content_gap_count",
			key: "content_gap_count",
			width: 100,
			render: (count: number) => <Badge count={count} showZero style={{ backgroundColor: count > 10 ? "#ff4d4f" : "#1890ff" }} />,
		},
		{
			title: intl.formatMessage({ id: 'competitor.column.topQueries' }),
			dataIndex: "top_queries",
			key: "top_queries",
			render: (queries: string[]) => (
				<Space size={[0, 4]} wrap>
					{queries.slice(0, 2).map((query) => <Tag key={query} color="blue">{query}</Tag>)}
					{queries.length > 2 && <Tag>+{queries.length - 2}</Tag>}
				</Space>
			),
		},
		{
			title: intl.formatMessage({ id: 'competitor.column.strengths' }),
			dataIndex: "strengths",
			key: "strengths",
			render: (strengths: string[]) => (
				<Space size={[0, 4]} wrap>
					{strengths.slice(0, 1).map((s) => <Tag key={s} color="green" icon={<CheckCircleOutlined />}>{s}</Tag>)}
				</Space>
			),
		},
		{
			title: intl.formatMessage({ id: 'competitor.column.weaknesses' }),
			dataIndex: "weaknesses",
			key: "weaknesses",
			render: (weaknesses: string[]) => (
				<Space size={[0, 4]} wrap>
					{weaknesses.slice(0, 1).map((w) => <Tag key={w} color="red" icon={<WarningOutlined />}>{w}</Tag>)}
				</Space>
			),
		},
		{
			title: intl.formatMessage({ id: 'competitor.column.status' }),
			dataIndex: "status",
			key: "status",
			width: 80,
			render: (status: string) => <Badge status={status === "active" ? "success" : "default"} text={status === "active" ? intl.formatMessage({ id: 'competitor.status.monitoring' }) : intl.formatMessage({ id: 'competitor.status.paused' })} />,
		},
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 150,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.viewDetail' })}><Button type="text" icon={<EyeOutlined />} onClick={() => handleShowDetail(record)} /></Tooltip>
					<Tooltip title={intl.formatMessage({ id: 'competitor.action.sync' })}><Button type="text" icon={<SyncOutlined />} onClick={() => handleSync(record.id)} /></Tooltip>
					<Popconfirm title={intl.formatMessage({ id: 'competitor.confirmDelete' })} okText={intl.formatMessage({ id: 'common.action.confirm' })} cancelText={intl.formatMessage({ id: 'common.action.cancel' })} onConfirm={() => handleDelete(record.id)}>
						<Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}><Button type="text" danger icon={<DeleteOutlined />} /></Tooltip>
					</Popconfirm>
				</Space>
			),
		},
	];

	const handleShowDetail = (record: any) => {
		setSelectedCompetitor(record);
		setDetailModalVisible(true);
	};

	const handleCreate = async (values: any) => {
		try {
			await api.monitor.createCompetitor(values);
			message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
			setCreateModalVisible(false);
			createForm.resetFields();
			queryClient.invalidateQueries({ queryKey: queryKeys.competitors });
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'common.message.createFailed' }));
		}
	};

	const handleDelete = async (id: number) => {
		try {
			await api.monitor.deleteCompetitor(id);
			message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
			queryClient.invalidateQueries({ queryKey: queryKeys.competitors });
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'common.message.deleteFailed' }));
		}
	};

	const handleSync = async (id: number) => {
		try {
			await api.monitor.syncCompetitor(id);
			message.success(intl.formatMessage({ id: 'competitor.message.syncCreated' }));
			queryClient.invalidateQueries({ queryKey: queryKeys.competitors });
		} catch (error: any) {
			message.error(intl.formatMessage({ id: 'competitor.message.syncFailed' }));
		}
	};

	const getVisibilityChartOption = () => ({
		title: { text: intl.formatMessage({ id: 'competitor.chart.visibilityTitle' }), left: "center" },
		tooltip: { trigger: "axis" },
		xAxis: { type: "category", data: [intl.formatMessage({ id: 'competitor.chart.us' }), ...competitors.map((c: any) => c.competitor_name)], axisLabel: { rotate: 30 } },
		yAxis: { type: "value", min: Math.max(0, Math.floor(Math.min(stats.ourScore, ...competitors.map((c: any) => c.visibility_score || 0)) / 10) * 10 - 10), max: 100 },
		series: [{
			type: "bar",
			data: [
				{ value: stats.ourScore, itemStyle: { color: "#52c41a" } },
				...competitors.map((c: any) => ({ value: c.visibility_score, itemStyle: { color: "#1890ff" } })),
			],
			barWidth: "50%",
		}],
	});

	const getContentGapChartOption = () => ({
		title: { text: intl.formatMessage({ id: 'competitor.chart.contentGapTitle' }), left: "center" },
		tooltip: { trigger: "item" },
		legend: { bottom: 0 },
		series: [{
			type: "pie",
			radius: "50%",
			data: competitors.map((c: any) => ({ value: c.content_gap_count, name: c.competitor_name })),
			emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: "rgba(0, 0, 0, 0.5)" } },
		}],
	});

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'competitor.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'competitor.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'competitor.stat.total' })} value={stats.total} prefix={<GlobalOutlined />} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'competitor.stat.monitoring' })} value={stats.active} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'competitor.stat.ourScore' })} value={stats.ourScore} suffix={intl.formatMessage({ id: 'common.unit.score' })} valueStyle={{ color: "#52c41a" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'competitor.stat.ourRank' })} value={stats.ourRank} suffix={intl.formatMessage({ id: 'common.unit.rank' })} prefix={<TrophyOutlined />} valueStyle={{ color: "#ffd700" }} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'competitor.stat.avgScore' })} value={stats.avgVisibility} suffix={intl.formatMessage({ id: 'common.unit.score' })} /></Card></Col>
				<Col xs={12} sm={8} lg={4}><Card><Statistic title={intl.formatMessage({ id: 'competitor.stat.totalGaps' })} value={stats.totalGaps} prefix={<WarningOutlined />} valueStyle={{ color: "#ff4d4f" }} /></Card></Col>
			</Row>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={24} lg={14}><Card><ReactECharts option={getVisibilityChartOption()} style={{ height: 300 }} /></Card></Col>
				<Col xs={24} lg={10}><Card><ReactECharts option={getContentGapChartOption()} style={{ height: 300 }} /></Card></Col>
			</Row>

			<Card
				title={<Space><GlobalOutlined /><span>{intl.formatMessage({ id: 'competitor.section.list' })}</span></Space>}
				extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>{intl.formatMessage({ id: 'competitor.action.add' })}</Button>}
			>
				<Table columns={columns} dataSource={competitors} rowKey="id" pagination={{ showSizeChanger: true, showQuickJumper: true, showTotal: (total) => intl.formatMessage({ id: 'common.pagination.totalShort' }, { total }) }} />
			</Card>

			<Modal title={intl.formatMessage({ id: 'competitor.modal.add' })} open={createModalVisible} onCancel={() => setCreateModalVisible(false)} footer={null} width={500}>
				<Form form={createForm} layout="vertical" onFinish={handleCreate}>
					<Form.Item name="competitor_name" label={intl.formatMessage({ id: 'competitor.form.name' })} rules={[{ required: true, message: intl.formatMessage({ id: 'competitor.validation.enterName' }) }]}>
						<Input placeholder={intl.formatMessage({ id: 'competitor.placeholder.name' })} />
					</Form.Item>
					<Form.Item name="competitor_domain" label={intl.formatMessage({ id: 'competitor.form.domain' })} rules={[{ required: true, message: intl.formatMessage({ id: 'competitor.validation.enterDomain' }) }]}>
						<Input placeholder="https://example.com" />
					</Form.Item>
					<Form.Item>
						<Space>
							<Button type="primary" htmlType="submit">{intl.formatMessage({ id: 'common.action.add' })}</Button>
							<Button onClick={() => setCreateModalVisible(false)}>{intl.formatMessage({ id: 'common.action.cancel' })}</Button>
						</Space>
					</Form.Item>
				</Form>
			</Modal>

			<Modal title={intl.formatMessage({ id: 'competitor.modal.detail' })} open={detailModalVisible} onCancel={() => setDetailModalVisible(false)} footer={null} width={700}>
				{selectedCompetitor && (
					<div>
						<Descriptions column={2} bordered className="mb-4">
							<Descriptions.Item label={intl.formatMessage({ id: 'competitor.desc.name' })} span={2}>{selectedCompetitor.competitor_name}</Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'competitor.desc.domain' })}><a href={selectedCompetitor.competitor_domain} target="_blank" rel="noopener noreferrer">{selectedCompetitor.competitor_domain}</a></Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'competitor.desc.status' })}><Badge status={selectedCompetitor.status === "active" ? "success" : "default"} text={selectedCompetitor.status === "active" ? intl.formatMessage({ id: 'competitor.status.monitoring' }) : intl.formatMessage({ id: 'competitor.status.paused' })} /></Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'competitor.desc.visibilityScore' })}><Progress percent={selectedCompetitor.visibility_score} size="small" /></Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'competitor.desc.contentGap' })}><Badge count={selectedCompetitor.content_gap_count} showZero style={{ backgroundColor: "#ff4d4f" }} /></Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'competitor.desc.topQueries' })} span={2}><Space size={[0, 4]} wrap>{selectedCompetitor.top_queries.map((q: string) => <Tag key={q} color="blue">{q}</Tag>)}</Space></Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'competitor.desc.strengths' })} span={2}><Space size={[0, 4]} wrap>{selectedCompetitor.strengths.map((s: string) => <Tag key={s} color="green" icon={<CheckCircleOutlined />}>{s}</Tag>)}</Space></Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'competitor.desc.weaknesses' })} span={2}><Space size={[0, 4]} wrap>{selectedCompetitor.weaknesses.map((w: string) => <Tag key={w} color="red" icon={<WarningOutlined />}>{w}</Tag>)}</Space></Descriptions.Item>
							<Descriptions.Item label={intl.formatMessage({ id: 'competitor.desc.lastCheck' })} span={2}>{new Date(selectedCompetitor.last_check_time).toLocaleString()}</Descriptions.Item>
						</Descriptions>
					</div>
				)}
			</Modal>
		</div>
	);
}
