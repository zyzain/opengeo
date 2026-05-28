"use client";

import {
	useAICitations,
	useAccounts,
	useContents,
	usePublishTasks,
} from "@/hooks";
import {
	CheckCircleOutlined,
	ClockCircleOutlined,
	ExclamationCircleOutlined,
	FileTextOutlined,
	RiseOutlined,
	SendOutlined,
	TeamOutlined,
} from "@ant-design/icons";
import {
	Avatar,
	Card,
	Col,
	List,
	Progress,
	Row,
	Statistic,
	Table,
	Tag,
} from "antd";
import { useIntl } from "react-intl";

export default function DashboardPage() {
	const intl = useIntl();
	const { data: contentsData } = useContents({ page: 1, page_size: 100 });
	const { data: tasksData } = usePublishTasks({ page: 1, page_size: 100 });
	const { data: accountsData } = useAccounts({ page: 1, page_size: 5 });
	const { data: citationsData } = useAICitations({ page: 1, page_size: 1000 });

	// 统计数据
	const successRate = (() => {
		const items = tasksData?.items || [];
		if (items.length === 0) return 0;
		const successCount = items.filter((t: any) => t.status === 2).length;
		return +((successCount / items.length) * 100).toFixed(1);
	})();

	const stats = {
		totalContents: contentsData?.total || 0,
		totalTasks: tasksData?.total || 0,
		totalAccounts: accountsData?.total || 0,
		successRate,
	};

	// 最近发布任务
	const recentTasks = tasksData?.items?.slice(0, 5) || [];

	// 任务状态标签
	const getTaskStatusTag = (status: number) => {
		const statusMap: Record<number, { color: string; text: string }> = {
			0: { color: "processing", text: intl.formatMessage({ id: 'publish.status.pending' }) },
			1: { color: "processing", text: intl.formatMessage({ id: 'publish.status.publishing' }) },
			2: { color: "success", text: intl.formatMessage({ id: 'publish.status.published' }) },
			3: { color: "error", text: intl.formatMessage({ id: 'publish.status.failed' }) },
			4: { color: "default", text: intl.formatMessage({ id: 'publish.status.cancelled' }) },
		};
		const config = statusMap[status] || { color: "default", text: intl.formatMessage({ id: 'common.status.unknown' }) };
		return <Tag color={config.color}>{config.text}</Tag>;
	};

	// 内容类型统计
	const typeColorMap: Record<string, string> = {
		文章: "#1890ff",
		视频: "#52c41a",
		图片: "#faad14",
	};
	const typeLabelMap: Record<string, string> = {
		文章: intl.formatMessage({ id: 'content.type.article' }),
		视频: intl.formatMessage({ id: 'content.type.video' }),
		图片: intl.formatMessage({ id: 'content.type.image' }),
	};
	const contentTypeStats = (() => {
		const items = contentsData?.items || [];
		const countMap: Record<string, number> = {};
		items.forEach((item: any) => {
			const t = item.content_type || intl.formatMessage({ id: 'common.status.unknown.label' });
			countMap[t] = (countMap[t] || 0) + 1;
		});
		return Object.entries(countMap).map(([type, count]) => ({
			type: typeLabelMap[type] || type,
			count,
			color: typeColorMap[type] || "#999",
		}));
	})();

	// 最近活动
	const recentActivities = (() => {
		const activities: {
			id: number;
			action: string;
			target: string;
			time: string;
		}[] = [];
		const recentContents = contentsData?.items?.slice(0, 3) || [];
		const recentTasksList = tasksData?.items?.slice(0, 3) || [];

		recentContents.forEach((item: any, idx: number) => {
			activities.push({
				id: idx + 1,
				action: intl.formatMessage({ id: 'dashboard.activity.published' }),
				target: `《${item.title || intl.formatMessage({ id: 'dashboard.unnamed' })}》`,
				time: item.created_at ? new Date(item.created_at).toLocaleString() : "",
			});
		});
		recentTasksList.forEach((item: any, idx: number) => {
			const statusText = item.status === 2 ? intl.formatMessage({ id: 'dashboard.activity.completed' }) : intl.formatMessage({ id: 'dashboard.activity.updated' });
			activities.push({
				id: 100 + idx,
				action: `${statusText}${intl.formatMessage({ id: 'dashboard.activity.task' })}`,
				target: `${intl.formatMessage({ id: 'dashboard.task.title' })}#${item.id}`,
				time: item.created_at ? new Date(item.created_at).toLocaleString() : "",
			});
		});
		return activities.slice(0, 5);
	})();

	// AI引用统计
	const citationStats = (() => {
		const items = citationsData?.items || [];
		const citedItems = items.filter((item: any) => item.is_cited);
		const totalCitations = citedItems.length;
		const uniqueContents = new Set(
			citedItems.map((item: any) => item.content_id),
		).size;
		const uniqueModels = new Set(items.map((item: any) => item.ai_model)).size;
		const positions = citedItems
			.map((item: any) => item.citation_position)
			.filter((p: number) => p > 0);
		const avgPosition =
			positions.length > 0
				? +(
						positions.reduce((a: number, b: number) => a + b, 0) /
						positions.length
					).toFixed(1)
				: 0;
		return { totalCitations, uniqueContents, uniqueModels, avgPosition };
	})();

	return (
		<div className="page-container">
			<div className="page-header">
				<h1 className="text-2xl font-bold text-gray-800">{intl.formatMessage({ id: 'dashboard.page.title' })}</h1>
				<p className="text-gray-500 mt-1">{intl.formatMessage({ id: 'dashboard.page.subtitle' })}</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-6">
				<Col xs={24} sm={12} lg={6}>
					<Card hoverable>
						<Statistic
							title={intl.formatMessage({ id: 'dashboard.stat.totalContents' })}
							value={stats.totalContents}
							prefix={<FileTextOutlined />}
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
				<Col xs={24} sm={12} lg={6}>
					<Card hoverable>
						<Statistic
							title={intl.formatMessage({ id: 'dashboard.stat.totalTasks' })}
							value={stats.totalTasks}
							prefix={<SendOutlined />}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={24} sm={12} lg={6}>
					<Card hoverable>
						<Statistic
							title={intl.formatMessage({ id: 'dashboard.stat.totalAccounts' })}
							value={stats.totalAccounts}
							prefix={<TeamOutlined />}
							valueStyle={{ color: "#722ed1" }}
						/>
					</Card>
				</Col>
				<Col xs={24} sm={12} lg={6}>
					<Card hoverable>
						<Statistic
							title={intl.formatMessage({ id: 'dashboard.stat.successRate' })}
							value={stats.successRate}
							suffix="%"
							prefix={<RiseOutlined />}
							valueStyle={{ color: "#faad14" }}
						/>
					</Card>
				</Col>
			</Row>

			<Row gutter={[16, 16]}>
				{/* 最近发布任务 */}
				<Col xs={24} lg={16}>
					<Card title={intl.formatMessage({ id: 'dashboard.recentTasks' })} className="h-full">
						<Table
							dataSource={recentTasks}
							rowKey="id"
							pagination={false}
							size="small"
							columns={[
								{
									title: intl.formatMessage({ id: 'publish.column.taskId' }),
									dataIndex: "id",
									key: "id",
									width: 80,
								},
								{
									title: intl.formatMessage({ id: 'publish.column.contentId' }),
									dataIndex: "content_id",
									key: "content_id",
									width: 80,
								},
								{
									title: intl.formatMessage({ id: 'common.column.status' }),
									dataIndex: "status",
									key: "status",
									render: (status: number) => getTaskStatusTag(status),
								},
								{
									title: intl.formatMessage({ id: 'common.column.createdAt' }),
									dataIndex: "created_at",
									key: "created_at",
									render: (text: string) => new Date(text).toLocaleString(),
								},
							]}
						/>
					</Card>
				</Col>

				{/* 内容类型统计 */}
				<Col xs={24} lg={8}>
					<Card title={intl.formatMessage({ id: 'dashboard.contentTypeDist' })} className="h-full">
						<div className="space-y-4">
							{contentTypeStats.map((item) => (
								<div key={item.type} className="flex items-center">
									<span className="w-16 text-gray-600">{item.type}</span>
									<Progress
										percent={Math.round(
											(item.count / stats.totalContents) * 100,
										)}
										strokeColor={item.color}
										className="flex-1 mx-4"
									/>
									<span className="w-12 text-right text-gray-500">
										{item.count}
									</span>
								</div>
							))}
						</div>
					</Card>
				</Col>
			</Row>

			<Row gutter={[16, 16]} className="mt-4">
				{/* AI引用统计 */}
				<Col xs={24} lg={12}>
					<Card title={intl.formatMessage({ id: 'dashboard.aiCitations' })} className="h-full">
						<div className="grid grid-cols-2 gap-4">
							<div className="text-center p-4 bg-blue-50 rounded-lg">
								<div className="text-3xl font-bold text-blue-600">
									{citationStats.totalCitations}
								</div>
								<div className="text-sm text-gray-500 mt-1">{intl.formatMessage({ id: 'dashboard.totalCitations' })}</div>
							</div>
							<div className="text-center p-4 bg-green-50 rounded-lg">
								<div className="text-3xl font-bold text-green-600">
									{citationStats.uniqueContents}
								</div>
								<div className="text-sm text-gray-500 mt-1">{intl.formatMessage({ id: 'dashboard.citedContents' })}</div>
							</div>
							<div className="text-center p-4 bg-purple-50 rounded-lg">
								<div className="text-3xl font-bold text-purple-600">
									{citationStats.uniqueModels}
								</div>
								<div className="text-sm text-gray-500 mt-1">{intl.formatMessage({ id: 'dashboard.coveredModels' })}</div>
							</div>
							<div className="text-center p-4 bg-orange-50 rounded-lg">
								<div className="text-3xl font-bold text-orange-600">
									{citationStats.avgPosition}
								</div>
								<div className="text-sm text-gray-500 mt-1">{intl.formatMessage({ id: 'dashboard.avgPosition' })}</div>
							</div>
						</div>
					</Card>
				</Col>

				{/* 最近活动 */}
				<Col xs={24} lg={12}>
					<Card title={intl.formatMessage({ id: 'dashboard.recentActivity' })} className="h-full">
						<List
							dataSource={recentActivities}
							renderItem={(item) => (
								<List.Item>
									<List.Item.Meta
										avatar={
											<Avatar
												icon={<CheckCircleOutlined />}
												style={{ backgroundColor: "#52c41a" }}
												size="small"
											/>
										}
										title={
											<span>
												{item.action}
												<span className="text-blue-500 ml-1">
													{item.target}
												</span>
											</span>
										}
										description={item.time}
									/>
								</List.Item>
							)}
						/>
					</Card>
				</Col>
			</Row>
		</div>
	);
}
