"use client";

import {
	ArrowRightOutlined,
	GlobalOutlined,
	RocketOutlined,
	SafetyOutlined,
	ThunderboltOutlined,
} from "@ant-design/icons";
import { Button, Card, Col, Row, Space, Statistic, Typography } from "antd";
import { Link } from "react-router-dom";

const { Title, Paragraph } = Typography;

export default function HomePage() {
	const features = [
		{
			icon: <ThunderboltOutlined className="text-4xl text-blue-500" />,
			title: "AI内容优化",
			description:
				"基于DeepSeek、Kimi、豆包等多模型的GEO语义增强，提升AI搜索可见性",
		},
		{
			icon: <GlobalOutlined className="text-4xl text-green-500" />,
			title: "多渠道分发",
			description: "支持微信、微博、抖音、小红书等主流平台一键发布",
		},
		{
			icon: <SafetyOutlined className="text-4xl text-purple-500" />,
			title: "智能合规检测",
			description: "敏感词过滤、广告法合规校验、AIGC标识自动添加",
		},
		{
			icon: <RocketOutlined className="text-4xl text-orange-500" />,
			title: "效果监测",
			description: "AI引用追踪、信源权重评分、ROI归因分析",
		},
	];

	const stats = [
		{ title: "支持平台", value: 10, suffix: "+" },
		{ title: "AI模型", value: 4, suffix: "个" },
		{ title: "内容优化", value: 95, suffix: "%" },
		{ title: "发布效率", value: 300, suffix: "%" },
	];

	return (
		<div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
			{/* 导航栏 */}
			<header className="bg-white shadow-sm sticky top-0 z-50">
				<div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
					<div className="flex items-center space-x-2">
						<ThunderboltOutlined className="text-2xl text-blue-500" />
						<span className="text-xl font-bold text-gray-800">OpenGEO</span>
					</div>
					<Space>
						<Link to="/auth/login">
							<Button type="text">登录</Button>
						</Link>
						<Link to="/auth/login">
							<Button type="primary">免费试用</Button>
						</Link>
					</Space>
				</div>
			</header>

			{/* 英雄区域 */}
			<section className="py-20 px-4">
				<div className="max-w-4xl mx-auto text-center">
					<Title level={1} className="mb-6">
						AI时代的<span className="text-blue-500">GEO智能发布平台</span>
					</Title>
					<Paragraph className="text-lg text-gray-600 mb-8">
						基于CloudWeGo高性能微服务架构，帮助企业在AI搜索时代提升品牌可见性。
						支持多模型内容优化、多渠道智能分发、全链路效果监测。
					</Paragraph>
					<Space size="large">
						<Link to="/auth/login">
							<Button type="primary" size="large" icon={<RocketOutlined />}>
								开始使用
							</Button>
						</Link>
						<Button size="large" icon={<ArrowRightOutlined />}>
							了解更多
						</Button>
					</Space>
				</div>
			</section>

			{/* 统计数据 */}
			<section className="py-12 bg-blue-500">
				<div className="max-w-5xl mx-auto px-4">
					<Row gutter={[32, 32]} justify="center">
						{stats.map((stat, index) => (
							<Col key={index} xs={12} sm={6}>
								<div className="text-center text-white">
									<div className="text-4xl font-bold mb-2">
										{stat.value}
										<span className="text-lg">{stat.suffix}</span>
									</div>
									<div className="text-blue-100">{stat.title}</div>
								</div>
							</Col>
						))}
					</Row>
				</div>
			</section>

			{/* 功能特性 */}
			<section className="py-20 px-4">
				<div className="max-w-6xl mx-auto">
					<div className="text-center mb-16">
						<Title level={2}>核心功能</Title>
						<Paragraph className="text-gray-500">
							全链路GEO优化解决方案
						</Paragraph>
					</div>
					<Row gutter={[32, 32]}>
						{features.map((feature, index) => (
							<Col key={index} xs={24} sm={12} lg={6}>
								<Card
									hoverable
									className="h-full text-center"
									cover={<div className="py-8 bg-gray-50">{feature.icon}</div>}
								>
									<Card.Meta
										title={feature.title}
										description={feature.description}
									/>
								</Card>
							</Col>
						))}
					</Row>
				</div>
			</section>

			{/* 技术架构 */}
			<section className="py-20 px-4 bg-gray-50">
				<div className="max-w-6xl mx-auto">
					<div className="text-center mb-16">
						<Title level={2}>技术架构</Title>
						<Paragraph className="text-gray-500">
							基于CloudWeGo高性能微服务架构
						</Paragraph>
					</div>
					<Row gutter={[32, 32]}>
						<Col xs={24} md={12}>
							<Card title="后端技术栈">
								<ul className="space-y-3">
									<li className="flex items-center">
										<span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
										<span>HTTP框架: Hertz (CloudWeGo)</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
										<span>RPC框架: Kitex (CloudWeGo)</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-purple-500 rounded-full mr-2"></span>
										<span>消息队列: Kafka</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-orange-500 rounded-full mr-2"></span>
										<span>数据库: MySQL + Redis + Milvus</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
										<span>监控: Prometheus + Grafana + Jaeger</span>
									</li>
								</ul>
							</Card>
						</Col>
						<Col xs={24} md={12}>
							<Card title="前端技术栈">
								<ul className="space-y-3">
									<li className="flex items-center">
										<span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
										<span>框架: Next.js 15+ (App Router)</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
										<span>UI: Ant Design 5.x + ProComponents</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-purple-500 rounded-full mr-2"></span>
										<span>状态: TanStack Query + Zustand</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-orange-500 rounded-full mr-2"></span>
										<span>样式: Tailwind CSS 4</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
										<span>图表: ECharts 5</span>
									</li>
								</ul>
							</Card>
						</Col>
					</Row>
				</div>
			</section>

			{/* CTA */}
			<section className="py-20 px-4">
				<div className="max-w-4xl mx-auto text-center">
					<Title level={2}>开始优化您的AI搜索可见性</Title>
					<Paragraph className="text-gray-500 mb-8">
						立即注册，免费体验OpenGEO智能发布平台
					</Paragraph>
					<Link to="/auth/login">
						<Button type="primary" size="large" icon={<RocketOutlined />}>
							免费注册
						</Button>
					</Link>
				</div>
			</section>

			{/* 页脚 */}
			<footer className="bg-gray-800 text-white py-12 px-4">
				<div className="max-w-6xl mx-auto">
					<Row gutter={[32, 32]}>
						<Col xs={24} md={8}>
							<div className="flex items-center space-x-2 mb-4">
								<ThunderboltOutlined className="text-2xl text-blue-400" />
								<span className="text-xl font-bold">OpenGEO</span>
							</div>
							<p className="text-gray-400">AI时代的GEO智能发布平台</p>
						</Col>
						<Col xs={12} md={4}>
							<h4 className="font-bold mb-4">产品</h4>
							<ul className="space-y-2 text-gray-400">
								<li>
									<a href="#" className="hover:text-white">
										功能特性
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										定价方案
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										使用文档
									</a>
								</li>
							</ul>
						</Col>
						<Col xs={12} md={4}>
							<h4 className="font-bold mb-4">资源</h4>
							<ul className="space-y-2 text-gray-400">
								<li>
									<a href="#" className="hover:text-white">
										博客
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										案例
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										帮助中心
									</a>
								</li>
							</ul>
						</Col>
						<Col xs={12} md={4}>
							<h4 className="font-bold mb-4">关于我们</h4>
							<ul className="space-y-2 text-gray-400">
								<li>
									<a href="#" className="hover:text-white">
										公司介绍
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										联系我们
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										加入我们
									</a>
								</li>
							</ul>
						</Col>
						<Col xs={12} md={4}>
							<h4 className="font-bold mb-4">关注我们</h4>
							<ul className="space-y-2 text-gray-400">
								<li>
									<a href="#" className="hover:text-white">
										GitHub
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										微信公众号
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										技术社区
									</a>
								</li>
							</ul>
						</Col>
					</Row>
					<div className="border-t border-gray-700 mt-8 pt-8 text-center text-gray-400">
						<p>© 2024 OpenGEO. All rights reserved.</p>
					</div>
				</div>
			</footer>
		</div>
	);
}
