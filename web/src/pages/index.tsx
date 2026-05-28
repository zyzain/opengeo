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
import { useIntl } from "react-intl";

const { Title, Paragraph } = Typography;

export default function HomePage() {
	const intl = useIntl();
	const features = [
		{
			icon: <ThunderboltOutlined className="text-4xl text-blue-500" />,
			title: intl.formatMessage({ id: 'landing.features.aiOptimize' }),
			description:
				intl.formatMessage({ id: 'landing.features.aiOptimize.desc' }),
		},
		{
			icon: <GlobalOutlined className="text-4xl text-green-500" />,
			title: intl.formatMessage({ id: 'landing.features.multiChannel' }),
			description: intl.formatMessage({ id: 'landing.features.multiChannel.desc' }),
		},
		{
			icon: <SafetyOutlined className="text-4xl text-purple-500" />,
			title: intl.formatMessage({ id: 'landing.features.compliance' }),
			description: intl.formatMessage({ id: 'landing.features.compliance.desc' }),
		},
		{
			icon: <RocketOutlined className="text-4xl text-orange-500" />,
			title: intl.formatMessage({ id: 'landing.features.monitor' }),
			description: intl.formatMessage({ id: 'landing.features.monitor.desc' }),
		},
	];

	const stats = [
		{ title: intl.formatMessage({ id: 'landing.stats.platforms' }), value: 10, suffix: "+" },
		{ title: intl.formatMessage({ id: 'landing.stats.models' }), value: 4, suffix: intl.formatMessage({ id: 'landing.stats.modelsUnit' }) },
		{ title: intl.formatMessage({ id: 'landing.stats.optimization' }), value: 95, suffix: "%" },
		{ title: intl.formatMessage({ id: 'landing.stats.efficiency' }), value: 300, suffix: "%" },
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
							<Button type="text">{intl.formatMessage({ id: 'landing.login' })}</Button>
						</Link>
						<Link to="/auth/login">
							<Button type="primary">{intl.formatMessage({ id: 'landing.freeTrial' })}</Button>
						</Link>
					</Space>
				</div>
			</header>

			{/* 英雄区域 */}
			<section className="py-20 px-4">
				<div className="max-w-4xl mx-auto text-center">
					<Title level={1} className="mb-6">
						{intl.formatMessage({ id: 'landing.hero.title' })}
					</Title>
					<Paragraph className="text-lg text-gray-600 mb-8">
						{intl.formatMessage({ id: 'landing.hero.subtitle' })}
					</Paragraph>
					<Space size="large">
						<Link to="/auth/login">
							<Button type="primary" size="large" icon={<RocketOutlined />}>
								{intl.formatMessage({ id: 'landing.hero.start' })}
							</Button>
						</Link>
						<Button size="large" icon={<ArrowRightOutlined />}>
							{intl.formatMessage({ id: 'landing.hero.learnMore' })}
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
						<Title level={2}>{intl.formatMessage({ id: 'landing.features.title' })}</Title>
						<Paragraph className="text-gray-500">
							{intl.formatMessage({ id: 'landing.features.subtitle' })}
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
						<Title level={2}>{intl.formatMessage({ id: 'landing.arch.title' })}</Title>
						<Paragraph className="text-gray-500">
							{intl.formatMessage({ id: 'landing.arch.subtitle' })}
						</Paragraph>
					</div>
					<Row gutter={[32, 32]}>
						<Col xs={24} md={12}>
							<Card title={intl.formatMessage({ id: 'landing.arch.backend' })}>
								<ul className="space-y-3">
									<li className="flex items-center">
										<span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.http' })}</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.rpc' })}</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-purple-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.mq' })}</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-orange-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.db' })}</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.monitor' })}</span>
									</li>
								</ul>
							</Card>
						</Col>
						<Col xs={24} md={12}>
							<Card title={intl.formatMessage({ id: 'landing.arch.frontend' })}>
								<ul className="space-y-3">
									<li className="flex items-center">
										<span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.framework' })}</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.ui' })}</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-purple-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.state' })}</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-orange-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.style' })}</span>
									</li>
									<li className="flex items-center">
										<span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
										<span>{intl.formatMessage({ id: 'landing.arch.chart' })}</span>
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
					<Title level={2}>{intl.formatMessage({ id: 'landing.cta.title' })}</Title>
					<Paragraph className="text-gray-500 mb-8">
						{intl.formatMessage({ id: 'landing.cta.subtitle' })}
					</Paragraph>
					<Link to="/auth/login">
						<Button type="primary" size="large" icon={<RocketOutlined />}>
							{intl.formatMessage({ id: 'landing.cta.register' })}
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
							<p className="text-gray-400">{intl.formatMessage({ id: 'landing.hero.title' })}</p>
						</Col>
						<Col xs={12} md={4}>
							<h4 className="font-bold mb-4">{intl.formatMessage({ id: 'landing.footer.product' })}</h4>
							<ul className="space-y-2 text-gray-400">
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.features' })}
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.pricing' })}
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.docs' })}
									</a>
								</li>
							</ul>
						</Col>
						<Col xs={12} md={4}>
							<h4 className="font-bold mb-4">{intl.formatMessage({ id: 'landing.footer.resources' })}</h4>
							<ul className="space-y-2 text-gray-400">
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.blog' })}
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.cases' })}
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.help' })}
									</a>
								</li>
							</ul>
						</Col>
						<Col xs={12} md={4}>
							<h4 className="font-bold mb-4">{intl.formatMessage({ id: 'landing.footer.about' })}</h4>
							<ul className="space-y-2 text-gray-400">
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.company' })}
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.contact' })}
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.careers' })}
									</a>
								</li>
							</ul>
						</Col>
						<Col xs={12} md={4}>
							<h4 className="font-bold mb-4">{intl.formatMessage({ id: 'landing.footer.followUs' })}</h4>
							<ul className="space-y-2 text-gray-400">
								<li>
									<a href="#" className="hover:text-white">
										GitHub
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.wechat' })}
									</a>
								</li>
								<li>
									<a href="#" className="hover:text-white">
										{intl.formatMessage({ id: 'landing.footer.community' })}
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
