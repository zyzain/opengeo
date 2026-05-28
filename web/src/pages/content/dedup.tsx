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
import { useIntl } from "react-intl";

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
	const inputSentences = input.split(/[。！？.!?\n]+/).filter((s) => s.trim().length > 5);
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
	const intl = useIntl();
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
		const medium = checkHistory.filter((d) => d.similarity >= 60 && d.similarity < 80).length;
		const low = checkHistory.filter((d) => d.similarity < 60).length;
		const totalDuplicates = checkHistory.length;
		const totalChecked = contents.length;
		const uniqueRate = totalChecked === 0 ? 100 : Math.round(((totalChecked - totalDuplicates) / totalChecked) * 1000) / 10;
		return { totalChecked, duplicatesFound: totalDuplicates, highRisk: high, mediumRisk: medium, lowRisk: low, uniqueRate };
	}, [checkHistory, contents.length]);

	const handleCheck = async () => {
		if (!inputText.trim()) {
			message.warning(intl.formatMessage({ id: 'common.message.operationFailed' }));
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
			message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
		} catch (err: any) {
			message.error(err?.message || intl.formatMessage({ id: 'common.message.operationFailed' }));
		} finally {
			setChecking(false);
		}
	};

	const getSimilarityColor = (score: number) => {
		if (score >= 80) return "#ff4d4f";
		if (score >= 60) return "#faad14";
		return "#52c41a";
	};

	const columns = [
		{ title: "ID", dataIndex: "id", key: "id", width: 60 },
		{
			title: intl.formatMessage({ id: 'dedup.column.contentTitle' }),
			dataIndex: "title",
			key: "title",
			render: (text: string, record: any) => <a onClick={() => { setSelectedContent(record); setDetailModalVisible(true); }} className="text-blue-500">{text}</a>,
		},
		{
			title: intl.formatMessage({ id: 'dedup.column.similarity' }),
			dataIndex: "similarity",
			key: "similarity",
			width: 150,
			render: (score: number) => <Progress percent={score} size="small" strokeColor={getSimilarityColor(score)} format={(percent) => <span className="font-bold">{percent}%</span>} />,
		},
		{ title: intl.formatMessage({ id: 'dedup.column.source' }), dataIndex: "source", key: "source", width: 100 },
		{
			title: intl.formatMessage({ id: 'dedup.column.riskLevel' }),
			dataIndex: "status",
			key: "status",
			width: 100,
			render: (status: string) => {
				const config: Record<string, { color: string; text: string }> = {
					high: { color: "red", text: intl.formatMessage({ id: 'dedup.risk.high' }) },
					medium: { color: "orange", text: intl.formatMessage({ id: 'dedup.risk.medium' }) },
					low: { color: "green", text: intl.formatMessage({ id: 'dedup.risk.low' }) },
				};
				const c = config[status] || config.low;
				return <Tag color={c.color}>{c.text}</Tag>;
			},
		},
		{
			title: intl.formatMessage({ id: 'common.column.action' }),
			key: "action",
			width: 120,
			render: (_: any, record: any) => (
				<Space size="small">
					<Tooltip title={intl.formatMessage({ id: 'common.action.viewDetail' })}><Button type="text" icon={<EyeOutlined />} onClick={() => { setSelectedContent(record); setDetailModalVisible(true); }} /></Tooltip>
				</Space>
			),
		},
	];

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'dedup.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'dedup.page.subtitle' })}</p>
			</div>

			<Row gutter={[16, 16]} className="mb-4">
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'dedup.stat.totalChecked' })} value={stats.totalChecked} prefix={<FileTextOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'dedup.stat.duplicatesFound' })} value={stats.duplicatesFound} valueStyle={{ color: "#faad14" }} prefix={<CopyOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'dedup.stat.highRisk' })} value={stats.highRisk} valueStyle={{ color: "#ff4d4f" }} prefix={<WarningOutlined />} /></Card></Col>
				<Col xs={12} sm={6}><Card><Statistic title={intl.formatMessage({ id: 'dedup.stat.uniqueRate' })} value={stats.uniqueRate} suffix="%" valueStyle={{ color: "#52c41a" }} prefix={<CheckCircleOutlined />} /></Card></Col>
			</Row>

			<Card title={intl.formatMessage({ id: 'dedup.card.title' })} className="mb-4">
				<div className="space-y-4">
					<TextArea rows={6} placeholder={intl.formatMessage({ id: 'dedup.placeholder.input' })} value={inputText} onChange={(e) => setInputText(e.target.value)} />
					<Space>
						<Button type="primary" icon={<SearchOutlined />} onClick={handleCheck} loading={checking}>{intl.formatMessage({ id: 'dedup.action.startCheck' })}</Button>
						<Button onClick={() => { setInputText(""); setCheckResult(null); }}>{intl.formatMessage({ id: 'dedup.action.clear' })}</Button>
						<span className="text-gray-400 text-sm">{intl.formatMessage({ id: 'dedup.charCount' }, { count: inputText.length })}</span>
					</Space>
				</div>

				{checkResult && (
					<div className="mt-4 space-y-4">
					<Alert
						message={intl.formatMessage({ id: 'dedup.result.complete' }, { score: checkResult.similarity_score })}
						description={checkResult.similarity_score < 30 ? intl.formatMessage({ id: 'dedup.result.goodOriginal' }) : checkResult.similarity_score < 60 ? intl.formatMessage({ id: 'dedup.result.someSimilarity' }) : intl.formatMessage({ id: 'dedup.result.highSimilarity' })}
							type={checkResult.similarity_score < 30 ? "success" : checkResult.similarity_score < 60 ? "warning" : "error"}
							showIcon
						/>
						<div className="p-4 bg-gray-50 rounded-lg">
							<h4 className="font-medium mb-2">{intl.formatMessage({ id: 'dedup.suggestions.title' })}</h4>
							<ul className="list-disc pl-4 space-y-1">
								{checkResult.suggestions.map((s: string, i: number) => <li key={i} className="text-sm text-gray-600">{s}</li>)}
							</ul>
						</div>
					</div>
				)}
			</Card>

			<Card title={intl.formatMessage({ id: 'dedup.card.duplicates' })}>
				<Table columns={columns} dataSource={checkResult?.duplicates || []} rowKey="id" pagination={false} />
			</Card>

		<Modal title={intl.formatMessage({ id: 'dedup.modal.detailTitle' })} open={detailModalVisible} onCancel={() => setDetailModalVisible(false)} footer={null} width={600}>
			{selectedContent && (
				<Descriptions column={1} bordered>
					<Descriptions.Item label={intl.formatMessage({ id: 'dedup.desc.contentTitle' })}>{selectedContent.title}</Descriptions.Item>
					<Descriptions.Item label={intl.formatMessage({ id: 'dedup.desc.similarity' })}><Progress percent={selectedContent.similarity} strokeColor={getSimilarityColor(selectedContent.similarity)} /></Descriptions.Item>
					<Descriptions.Item label={intl.formatMessage({ id: 'dedup.desc.source' })}>{selectedContent.source}</Descriptions.Item>
					<Descriptions.Item label={intl.formatMessage({ id: 'dedup.desc.matchedText' })}><div className="p-2 bg-gray-50 rounded text-sm">{selectedContent.matched_text}</div></Descriptions.Item>
					<Descriptions.Item label={intl.formatMessage({ id: 'dedup.desc.advice' })}>{selectedContent.status === "high" ? intl.formatMessage({ id: 'dedup.advice.high' }) : selectedContent.status === "medium" ? intl.formatMessage({ id: 'dedup.advice.medium' }) : intl.formatMessage({ id: 'dedup.advice.low' })}</Descriptions.Item>
					</Descriptions>
				)}
			</Modal>
		</div>
	);
}
