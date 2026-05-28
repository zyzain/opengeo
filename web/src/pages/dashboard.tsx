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

export default function DashboardPage() {
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
			0: { color: "processing", text: "待发布" },
			1: { color: "processing", text: "发布中" },
			2: { color: "success", text: "已发布" },
			3: { color: "error", text: "失败" },
			4: { color: "default", text: "已取消" },
		};
		const config = statusMap[status] || { color: "default", text: "未知" };
		return <Tag color={config.color}>{config.text}</Tag>;
	};

	// 内容类型统计
	const typeColorMap: Record<string, string> = {
		文章: "#1890ff",
		视频: "#52c41a",
		图片: "#faad14",
	};
	const contentTypeStats = (() => {
		const items = contentsData?.items || [];
		const countMap: Record<string, number> = {};
		items.forEach((item: any) => {
			const t = item.content_type || "未知";
			countMap[t] = (countMap[t] || 0) + 1;
		});
		return Object.entries(countMap).map(([type, count]) => ({
			type,
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
				action: "发布了新内容",
				target: `《${item.title || "未命名"}》`,
				time: item.created_at ? new Date(item.created_at).toLocaleString() : "",
			});
		});
		recentTasksList.forEach((item: any, idx: number) => {
			const statusText = item.status === 2 ? "已完成" : "更新了";
			activities.push({
				id: 100 + idx,
				action: `${statusText}发布任务`,
				target: `任务#${item.id}`,
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
				<h1 className="text-2xl font-bold text-gray-800">仪表盘</h1>
				<p className="text-gray-500 mt-1">欢迎回来，这里是您的数据概览</p>
			</div>

			{/* 统计卡片 */}
			<Row gutter={[16, 16]} className="mb-6">
				<Col xs={24} sm={12} lg={6}>
					<Card hoverable>
						<Statistic
							title="内容总数"
							value={stats.totalContents}
							prefix={<FileTextOutlined />}
							valueStyle={{ color: "#1890ff" }}
						/>
					</Card>
				</Col>
				<Col xs={24} sm={12} lg={6}>
					<Card hoverable>
						<Statistic
							title="发布任务"
							value={stats.totalTasks}
							prefix={<SendOutlined />}
							valueStyle={{ color: "#52c41a" }}
						/>
					</Card>
				</Col>
				<Col xs={24} sm={12} lg={6}>
					<Card hoverable>
						<Statistic
							title="账号数量"
							value={stats.totalAccounts}
							prefix={<TeamOutlined />}
							valueStyle={{ color: "#722ed1" }}
						/>
					</Card>
				</Col>
				<Col xs={24} sm={12} lg={6}>
					<Card hoverable>
						<Statistic
							title="发布成功率"
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
					<Card title="最近发布任务" className="h-full">
						<Table
							dataSource={recentTasks}
							rowKey="id"
							pagination={false}
							size="small"
							columns={[
								{
									title: "任务ID",
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
									title: "状态",
									dataIndex: "status",
									key: "status",
									render: (status: number) => getTaskStatusTag(status),
								},
								{
									title: "创建时间",
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
					<Card title="内容类型分布" className="h-full">
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
					<Card title="AI引用统计" className="h-full">
						<div className="grid grid-cols-2 gap-4">
							<div className="text-center p-4 bg-blue-50 rounded-lg">
								<div className="text-3xl font-bold text-blue-600">
									{citationStats.totalCitations}
								</div>
								<div className="text-sm text-gray-500 mt-1">总引用次数</div>
							</div>
							<div className="text-center p-4 bg-green-50 rounded-lg">
								<div className="text-3xl font-bold text-green-600">
									{citationStats.uniqueContents}
								</div>
								<div className="text-sm text-gray-500 mt-1">被引用内容</div>
							</div>
							<div className="text-center p-4 bg-purple-50 rounded-lg">
								<div className="text-3xl font-bold text-purple-600">
									{citationStats.uniqueModels}
								</div>
								<div className="text-sm text-gray-500 mt-1">覆盖AI模型</div>
							</div>
							<div className="text-center p-4 bg-orange-50 rounded-lg">
								<div className="text-3xl font-bold text-orange-600">
									{citationStats.avgPosition}
								</div>
								<div className="text-sm text-gray-500 mt-1">平均引用位置</div>
							</div>
						</div>
					</Card>
				</Col>

				{/* 最近活动 */}
				<Col xs={24} lg={12}>
					<Card title="最近活动" className="h-full">
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
