"use client";

import { useLogin, useRegister } from "@/hooks";
import { LockOutlined, MailOutlined, UserOutlined } from "@ant-design/icons";
import { Button, Card, Form, Input, Tabs, message } from "antd";
import { useState } from "react";
import { Link } from "react-router-dom";

export default function AuthPage() {
	const [activeTab, setActiveTab] = useState("login");
	const loginMutation = useLogin();
	const registerMutation = useRegister();

	const onLoginFinish = async (values: any) => {
		try {
			await loginMutation.mutateAsync(values);
			message.success("登录成功");
		} catch (error: any) {
			message.error(error.response?.data?.message || "登录失败");
		}
	};

	const onRegisterFinish = async (values: any) => {
		try {
			await registerMutation.mutateAsync(values);
			message.success("注册成功");
		} catch (error: any) {
			message.error(error.response?.data?.message || "注册失败");
		}
	};

	const items = [
		{
			key: "login",
			label: "登录",
			children: (
				<Form
					name="login"
					onFinish={onLoginFinish}
					autoComplete="off"
					size="large"
				>
					<Form.Item
						name="username"
						rules={[{ required: true, message: "请输入用户名" }]}
					>
						<Input prefix={<UserOutlined />} placeholder="用户名" />
					</Form.Item>

					<Form.Item
						name="password"
						rules={[{ required: true, message: "请输入密码" }]}
					>
						<Input.Password prefix={<LockOutlined />} placeholder="密码" />
					</Form.Item>

					<Form.Item>
						<Button
							type="primary"
							htmlType="submit"
							loading={loginMutation.isPending}
							block
						>
							登录
						</Button>
					</Form.Item>
				</Form>
			),
		},
		{
			key: "register",
			label: "注册",
			children: (
				<Form
					name="register"
					onFinish={onRegisterFinish}
					autoComplete="off"
					size="large"
				>
					<Form.Item
						name="username"
						rules={[
							{ required: true, message: "请输入用户名" },
							{ min: 3, message: "用户名至少3个字符" },
							{ max: 20, message: "用户名最多20个字符" },
						]}
					>
						<Input prefix={<UserOutlined />} placeholder="用户名" />
					</Form.Item>

					<Form.Item
						name="email"
						rules={[
							{ required: true, message: "请输入邮箱" },
							{ type: "email", message: "请输入有效的邮箱地址" },
						]}
					>
						<Input prefix={<MailOutlined />} placeholder="邮箱" />
					</Form.Item>

					<Form.Item
						name="password"
						rules={[
							{ required: true, message: "请输入密码" },
							{ min: 8, message: "密码至少8个字符" },
						]}
					>
						<Input.Password prefix={<LockOutlined />} placeholder="密码" />
					</Form.Item>

					<Form.Item
						name="confirmPassword"
						dependencies={["password"]}
						rules={[
							{ required: true, message: "请确认密码" },
							({ getFieldValue }) => ({
								validator(_, value) {
									if (!value || getFieldValue("password") === value) {
										return Promise.resolve();
									}
									return Promise.reject(new Error("两次输入的密码不一致"));
								},
							}),
						]}
					>
						<Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
					</Form.Item>

					<Form.Item>
						<Button
							type="primary"
							htmlType="submit"
							loading={registerMutation.isPending}
							block
						>
							注册
						</Button>
					</Form.Item>
				</Form>
			),
		},
	];

	return (
		<div className="min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center p-4">
			<Card
				className="w-full max-w-md shadow-2xl"
				title={
					<div className="text-center">
						<h1 className="text-2xl font-bold text-gray-800 mb-2">OpenGEO</h1>
						<p className="text-sm text-gray-500">智能发布平台</p>
					</div>
				}
			>
				<Tabs
					activeKey={activeTab}
					onChange={setActiveTab}
					items={items}
					centered
				/>
				<div className="text-center mt-4">
					<Link to="/" className="text-sm text-gray-500 hover:text-blue-500">
						返回首页
					</Link>
				</div>
			</Card>
		</div>
	);
}
