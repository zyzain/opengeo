import React, { useEffect } from 'react';
import { Form, Input, Select, Button, Card, message } from 'antd';
import { useIntl } from 'react-intl';

const { Option } = Select;
const { TextArea } = Input;

interface BrandFormProps {
  initialValues?: any;
  onSubmit: (values: any) => Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

const BrandForm: React.FC<BrandFormProps> = ({ initialValues, onSubmit, onCancel, loading }) => {
  const intl = useIntl();
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
      message.success(initialValues ? intl.formatMessage({ id: 'common.message.updateSuccess' }) : intl.formatMessage({ id: 'common.message.createSuccess' }));
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  return (
    <Card title={initialValues ? intl.formatMessage({ id: 'brand.form.editTitle' }) : intl.formatMessage({ id: 'brand.form.createTitle' })}>
      <Form form={form} layout="vertical">
        <Form.Item
          name="name"
          label={intl.formatMessage({ id: 'brand.form.name' })}
          rules={[{ required: true, message: intl.formatMessage({ id: 'brand.validation.enterName' }) }]}
        >
          <Input placeholder={intl.formatMessage({ id: 'brand.placeholder.enterName' })} />
        </Form.Item>

        <Form.Item
          name="slug"
          label={intl.formatMessage({ id: 'brand.form.slug' })}
          rules={[
            { required: true, message: intl.formatMessage({ id: 'brand.validation.enterSlug' }) },
            { pattern: /^[a-z0-9-]+$/, message: intl.formatMessage({ id: 'brand.validation.slugPattern' }) },
          ]}
          extra={intl.formatMessage({ id: 'brand.form.slugHelp' })}
        >
          <Input placeholder={intl.formatMessage({ id: 'brand.form.slugPlaceholder' })} disabled={!!initialValues} />
        </Form.Item>

        <Form.Item name="description" label={intl.formatMessage({ id: 'brand.form.description' })}>
          <TextArea rows={4} placeholder={intl.formatMessage({ id: 'brand.placeholder.enterDescription' })} />
        </Form.Item>

        <Form.Item name="industry" label={intl.formatMessage({ id: 'brand.form.industry' })}>
          <Select placeholder={intl.formatMessage({ id: 'brand.placeholder.selectIndustry' })}>
            <Option value="科技">{intl.formatMessage({ id: 'industry.tech' })}</Option>
            <Option value="金融">{intl.formatMessage({ id: 'industry.finance' })}</Option>
            <Option value="医疗">{intl.formatMessage({ id: 'industry.medical' })}</Option>
            <Option value="教育">{intl.formatMessage({ id: 'industry.education' })}</Option>
            <Option value="电商">{intl.formatMessage({ id: 'industry.ecommerce' })}</Option>
            <Option value="制造">{intl.formatMessage({ id: 'industry.manufacturing' })}</Option>
            <Option value="零售">{intl.formatMessage({ id: 'industry.retail' })}</Option>
            <Option value="其他">{intl.formatMessage({ id: 'industry.other' })}</Option>
          </Select>
        </Form.Item>

        <Form.Item name="website" label={intl.formatMessage({ id: 'brand.form.website' })}>
          <Input placeholder="https://example.com" />
        </Form.Item>

        <Form.Item name="logo_url" label={intl.formatMessage({ id: 'brand.form.logoUrl' })}>
          <Input placeholder={intl.formatMessage({ id: 'brand.placeholder.logoUrl' })} />
        </Form.Item>

        <Form.Item name="founded_year" label={intl.formatMessage({ id: 'brand.form.foundedYear' })}>
          <Input type="number" placeholder={intl.formatMessage({ id: 'brand.form.yearPlaceholder' })} />
        </Form.Item>

        <Form.Item name="headquarters" label={intl.formatMessage({ id: 'brand.form.headquarters' })}>
          <Input placeholder={intl.formatMessage({ id: 'brand.placeholder.headquartersExample' })} />
        </Form.Item>

        <Form.Item>
          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
            <Button onClick={onCancel}>{intl.formatMessage({ id: 'common.cancel' })}</Button>
            <Button type="primary" onClick={handleSubmit} loading={loading}>
              {initialValues ? intl.formatMessage({ id: 'common.update' }) : intl.formatMessage({ id: 'common.create' })}
            </Button>
          </div>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default BrandForm;
