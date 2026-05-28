"use client";

import { useContents } from "@/hooks";
import { api } from "@/lib/api";
import {
	CheckCircleOutlined,
	CopyOutlined,
	ExclamationCircleOutlined,
	EyeOutlined,
	FileTextOutlined,
	ReloadOutlined,
	SearchOutlined,
	WarningOutlined,
} from "@ant-design/icons";
import {
	Alert,
	Badge,
	Button,
	Card,
	Col,
	Descriptions,
	Input,
	Modal,
	Progress,
	Row,
	Space,
	Statistic,
	Table,
	Tag,
	Tooltip,
	message,
} from "antd";
import { useState } from "react";
import { useMemo } from "react";

const { TextArea } = Input;

function getNgrams(text: string, n: number): Set<string> {
	const grams = new Set<string>();
	for (let i = 0; i <= text.length - n; i++) {
		grams.add(text.substring(i, i + n));
	}
	return grams;
}

function computeSimilarity(a: string, b: string): number {
	if (!a || !b) return 0;
	const n = 3;
	const gramsA = getNgrams(a, n);
	const gramsB = getNgrams(b, n);
	let intersection = 0;
	for (const g of gramsA) {
		if (gramsB.has(g)) intersection++;
	}
	const union = gramsA.size + gramsB.size - intersection;
	return union === 0 ? 0 : Math.round((intersection / union) * 100);
}

function findBestMatch(input: string, contentBody: string): string {
	if (!input || !contentBody) return "";
	const inputSentences = input
		.split(/[。！？.!?\n]+/)
		.filter((s) => s.trim().length > 5);
	let best = "";
	let bestScore = 0;
	for (const sentence of inputSentences) {
		const score = computeSimilarity(sentence, contentBody);
		if (score > bestScore) {
			bestScore = score;
			best = sentence.trim();
		}
	}
	return best || input.substring(0, 80);
}

export default function ContentDedupPage() {
	const [inputText, setInputText] = useState("");
	const [checking, setChecking] = useState(false);
	const [checkResult, setCheckResult] = useState<any>(null);
	const [checkHistory, setCheckHistory] = useState<any[]>([]);
	const [detailModalVisible, setDetailModalVisible] = useState(false);
	const [selectedContent, setSelectedContent] = useState<any>(null);

	const { data: contentsData } = useContents({ page: 1, page_size: 100 });
	const contents = contentsData?.items || [];

	const stats = useMemo(() => {
		const high = checkHistory.filter((d) => d.similarity >= 80).length;
		const medium = checkHistory.filter(
			(d) => d.similarity >= 60 && d.similarity < 80,
		).length;
		const low = checkHistory.filter((d) => d.similarity < 60).length;
		const totalDuplicates = checkHistory.length;
		const totalChecked = contents.length;
		const uniqueRate =
			totalChecked === 0
				? 100
				: Math.round(((totalChecked - totalDuplicates) / totalChecked) * 1000) /
					10;
		return {
			totalChecked,
			duplicatesFound: totalDuplicates,
			highRisk: high,
			mediumRisk: medium,
			lowRisk: low,
			uniqueRate,
		};
	}, [checkHistory, contents.length]);

	const handleCheck = async () => {
		if (!inputText.trim()) {
			message.warning("请输入要检测的内容");
			return;
		}
		setChecking(true);
		try {
			const res = await api.dedup.check({ text: inputText });
			const data = res.data?.data || res.data;
			setCheckHistory(data.duplicates || []);
			setCheckResult({
				text_length: inputText.length,
				similarity_score: data.similarity_score ?? 0,
				duplicates: data.duplicates || [],
				suggestions: data.suggestions || [],
			});
			message.success("检测完成");
		} catch (err: any) {
			message.error(err?.message || "检测失败");
		} finally {
			setChecking(false);
		}
	};

	// 相似度颜色
	const getSimilarityColor = (score: number) => {
		if (score >= 80) return "#ff4d4f";
		if (score >= 60) return "#faad14";
		return "#52c41a";
	};

	// 表格列
	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 60 },
		{
			title: "内容标题",
			dataIndex: "title",
			key: "title",
			render: (text: string, record: any) => (
				<a
					onClick={() => {
						setSelectedContent(record);
						setDetailModalVisible(true);
					}}
					className="text-blue-500"
				>
					{text}
				</a>
			),
		},
		{
			title: "相似度",
			dataIndex: "similarity",
			key: "similarity",
			width: 150,
			render: (score: number) => (
				<Progress
					percent={score}
					size="small"
					strokeColor={getSimilarityColor(score)}
					format={(percent) => <span className="font-bold">{percent}%</span>}
				/>
			),
		},
		{ title: "来源", dataIndex: "source", key: "source", width: 100 },
		{
			title: "风险等级",
			dataIndex: "status",
			key: "status",
			width: 100,
			render: (status: string) => {
				const config: Record<string, { color: string; text: string }> = {
					high: { color: "red", text: "高风险" },
					medium: { color: "orange", text: "中风险" },
					low: { color: "green", text: "低风险" },
				};
				const c = config[status] || config.low;
				return <Tag color={c.color}>{c.text}</Tag>;
			},
		},
		{
			title: "操作",
			key: "action",
			width: 120,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title="查看详情">
						<Button
							type="text"
							icon={<EyeOutlined />}
							onClick={() => {
								setSelectedContent(record);
								setDetailModalVisible(true);
							}}
						/>
					</Tooltip>
				</Space>
			),
		},
	];

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">内容去重</h1>
				<p className="text-gray-500 mt-1">
					检测内容相似度，避免跨平台重复内容被判定为低质信源
				</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="已检测内容"
							value={stats.totalChecked}
							prefix={<FileTextOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="发现重复"
							value={stats.duplicatesFound}
							valueStyle={{ color: "#faad14" }}
							prefix={<CopyOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="高风险"
							value={stats.highRisk}
							valueStyle={{ color: "#ff4d4f" }}
							prefix={<WarningOutlined />}
						/>
					</Card>
				</Col>
				<Col xs={12} sm={6}>
					<Card>
						<Statistic
							title="原创率"
							value={stats.uniqueRate}
							suffix="%"
							valueStyle={{ color: "#52c41a" }}
							prefix={<CheckCircleOutlined />}
						/>
					</Card>
				</Col>
			</Row>

			{/* 检测工具 */}
			<Card title="内容去重检测" className="mb-4">
				<div className="space-y-4">
					<TextArea
						rows={6}
						placeholder="粘贴要检测的内容..."
						value={inputText}
						onChange={(e) => setInputText(e.target.value)}
					/>
					<Space>
						<Button
							type="primary"
							icon={<SearchOutlined />}
							onClick={handleCheck}
							loading={checking}
						>
							开始检测
						</Button>
						<Button
							onClick={() => {
								setInputText("");
								setCheckResult(null);
							}}
						>
							清空
						</Button>
						<span className="text-gray-400 text-sm">
							已输入 {inputText.length} 字符
						</span>
					</Space>
				</div>

				{/* 检测结果 */}
				{checkResult && (
					<div className="mt-4 space-y-4">
						<Alert
							message={`检测完成 - 相似度: ${checkResult.similarity_score}%`}
							description={
								checkResult.similarity_score < 30
									? "内容原创度良好"
									: checkResult.similarity_score < 60
										? "存在一定相似度，建议优化"
										: "相似度较高，强烈建议修改"
							}
							type={
								checkResult.similarity_score < 30
									? "success"
									: checkResult.similarity_score < 60
										? "warning"
										: "error"
							}
							showIcon
						/>
						<div className="p-4 bg-gray-50 rounded-lg">
							<h4 className="font-medium mb-2">优化建议</h4>
							<ul className="list-disc pl-4 space-y-1">
								{checkResult.suggestions.map((s: string, i: number) => (
									<li key={i} className="text-sm text-gray-600">
										{s}
									</li>
								))}
							</ul>
						</div>
					</div>
				)}
			</Card>

			{/* 重复内容列表 */}
			<Card title="已发现的重复内容">
				<Table
					columns={columns}
					dataSource={checkResult?.duplicates || []}
					rowKey="id"
					pagination={false}
				/>
			</Card>

			{/* 详情弹窗 */}
			<Modal
				title="重复内容详情"
				open={detailModalVisible}
				onCancel={() => setDetailModalVisible(false)}
				footer={null}
				width={600}
			>
				{selectedContent && (
					<Descriptions column={1} bordered>
						<Descriptions.Item label="内容标题">
							{selectedContent.title}
						</Descriptions.Item>
						<Descriptions.Item label="相似度">
							<Progress
								percent={selectedContent.similarity}
								strokeColor={getSimilarityColor(selectedContent.similarity)}
							/>
						</Descriptions.Item>
						<Descriptions.Item label="来源">
							{selectedContent.source}
						</Descriptions.Item>
						<Descriptions.Item label="匹配文本">
							<div className="p-2 bg-gray-50 rounded text-sm">
								{selectedContent.matched_text}
							</div>
						</Descriptions.Item>
						<Descriptions.Item label="处理建议">
							{selectedContent.status === "high"
								? "建议大幅修改或重新创作"
								: selectedContent.status === "medium"
									? "建议调整段落结构和用词"
									: "可选择性优化"}
						</Descriptions.Item>
					</Descriptions>
				)}
			</Modal>
		</div>
	);
}
