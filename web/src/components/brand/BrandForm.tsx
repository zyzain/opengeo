import React, { useEffect } from 'react';
import { Form, Input, Select, Button, Card, message } from 'antd';

const { Option } = Select;
const { TextArea } = Input;

interface BrandFormProps {
  initialValues?: any;
  onSubmit: (values: any) => Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

const BrandForm: React.FC<BrandFormProps> = ({ initialValues, onSubmit, onCancel, loading }) => {
  const [form] = Form.useForm();

  useEffect(() => {
    if (initialValues) {
      form.setFieldsValue(initialValues);
    } else {
      form.resetFields();
    }
  }, [initialValues, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      await onSubmit(values);
      message.success(initialValues ? '更新成功' : '创建成功');
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  return (
    <Card title={initialValues ? '编辑品牌' : '创建品牌'}>
      <Form form={form} layout="vertical">
        <Form.Item
          name="name"
          label="品牌名称"
          rules={[{ required: true, message: '请输入品牌名称' }]}
        >
          <Input placeholder="请输入品牌名称" />
        </Form.Item>

        <Form.Item
          name="slug"
          label="品牌标识"
          rules={[
            { required: true, message: '请输入品牌标识' },
            { pattern: /^[a-z0-9-]+$/, message: '只能包含小写字母、数字和连字符' },
          ]}
          extra="URL友好的标识，创建后不可修改"
        >
          <Input placeholder="例如：my-brand" disabled={!!initialValues} />
        </Form.Item>

        <Form.Item name="description" label="品牌描述">
          <TextArea rows={4} placeholder="请输入品牌描述" />
        </Form.Item>

        <Form.Item name="industry" label="所属行业">
          <Select placeholder="请选择行业">
            <Option value="科技">科技</Option>
            <Option value="金融">金融</Option>
            <Option value="医疗">医疗</Option>
            <Option value="教育">教育</Option>
            <Option value="电商">电商</Option>
            <Option value="制造">制造</Option>
            <Option value="零售">零售</Option>
            <Option value="其他">其他</Option>
          </Select>
        </Form.Item>

        <Form.Item name="website" label="品牌官网">
          <Input placeholder="https://example.com" />
        </Form.Item>

        <Form.Item name="logo_url" label="Logo URL">
          <Input placeholder="请输入Logo图片地址" />
        </Form.Item>

        <Form.Item name="founded_year" label="成立年份">
          <Input type="number" placeholder="例如：2024" />
        </Form.Item>

        <Form.Item name="headquarters" label="总部所在地">
          <Input placeholder="例如：北京" />
        </Form.Item>

        <Form.Item>
          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
            <Button onClick={onCancel}>取消</Button>
            <Button type="primary" onClick={handleSubmit} loading={loading}>
              {initialValues ? '更新' : '创建'}
            </Button>
          </div>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default BrandForm;
