"use client";

import { useLogin, useRegister } from "@/hooks";
import { LockOutlined, MailOutlined, UserOutlined } from "@ant-design/icons";
import { Button, Card, Form, Input, Tabs, message } from "antd";
import { useState } from "react";
import { Link } from "react-router-dom";
import { useIntl } from "react-intl";

export default function AuthPage() {
	const intl = useIntl();
	const [activeTab, setActiveTab] = useState("login");
	const loginMutation = useLogin();
	const registerMutation = useRegister();

	const onLoginFinish = async (values: any) => {
		try {
			await loginMutation.mutateAsync(values);
			message.success(intl.formatMessage({ id: 'login.success' }));
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'login.failed' }));
		}
	};

	const onRegisterFinish = async (values: any) => {
		try {
			await registerMutation.mutateAsync(values);
			message.success(intl.formatMessage({ id: 'login.registerSuccess' }));
		} catch (error: any) {
			message.error(error.response?.data?.message || intl.formatMessage({ id: 'login.registerFailed' }));
		}
	};

	const items = [
		{
			key: "login",
			label: intl.formatMessage({ id: 'login.title' }),
			children: (
				<Form
					name="login"
					onFinish={onLoginFinish}
					autoComplete="off"
					size="large"
				>
					<Form.Item
						name="username"
						rules={[{ required: true, message: intl.formatMessage({ id: 'login.validation.enterUsername' }) }]}
					>
						<Input prefix={<UserOutlined />} placeholder={intl.formatMessage({ id: 'login.username' })} />
					</Form.Item>

					<Form.Item
						name="password"
						rules={[{ required: true, message: intl.formatMessage({ id: 'login.validation.enterPassword' }) }]}
					>
						<Input.Password prefix={<LockOutlined />} placeholder={intl.formatMessage({ id: 'login.password' })} />
					</Form.Item>

					<Form.Item>
						<Button
							type="primary"
							htmlType="submit"
							loading={loginMutation.isPending}
							block
						>
							{intl.formatMessage({ id: 'login.submit' })}
						</Button>
					</Form.Item>
				</Form>
			),
		},
		{
			key: "register",
			label: intl.formatMessage({ id: 'login.register' }),
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
							{ required: true, message: intl.formatMessage({ id: 'login.validation.enterUsername' }) },
							{ min: 3, message: intl.formatMessage({ id: 'login.validation.usernameMin' }) },
							{ max: 20, message: intl.formatMessage({ id: 'login.validation.usernameMax' }) },
						]}
					>
						<Input prefix={<UserOutlined />} placeholder={intl.formatMessage({ id: 'login.username' })} />
					</Form.Item>

					<Form.Item
						name="email"
						rules={[
							{ required: true, message: intl.formatMessage({ id: 'login.validation.enterEmail' }) },
							{ type: "email", message: intl.formatMessage({ id: 'login.validation.validEmail' }) },
						]}
					>
						<Input prefix={<MailOutlined />} placeholder={intl.formatMessage({ id: 'login.email' })} />
					</Form.Item>

					<Form.Item
						name="password"
						rules={[
							{ required: true, message: intl.formatMessage({ id: 'login.validation.enterPassword' }) },
							{ min: 8, message: intl.formatMessage({ id: 'login.validation.passwordMin' }) },
						]}
					>
						<Input.Password prefix={<LockOutlined />} placeholder={intl.formatMessage({ id: 'login.password' })} />
					</Form.Item>

					<Form.Item
						name="confirmPassword"
						dependencies={["password"]}
						rules={[
							{ required: true, message: intl.formatMessage({ id: 'login.validation.confirmPassword' }) },
							({ getFieldValue }) => ({
								validator(_, value) {
									if (!value || getFieldValue("password") === value) {
										return Promise.resolve();
									}
									return Promise.reject(new Error(intl.formatMessage({ id: 'login.validation.passwordMismatch' })));
								},
							}),
						]}
					>
						<Input.Password prefix={<LockOutlined />} placeholder={intl.formatMessage({ id: 'login.confirmPassword' })} />
					</Form.Item>

					<Form.Item>
						<Button
							type="primary"
							htmlType="submit"
							loading={registerMutation.isPending}
							block
						>
							{intl.formatMessage({ id: 'login.registerSubmit' })}
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
						<p className="text-sm text-gray-500">{intl.formatMessage({ id: 'app.subtitle' })}</p>
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
						{intl.formatMessage({ id: 'login.backToHome' })}
					</Link>
				</div>
			</Card>
		</div>
	);
}
