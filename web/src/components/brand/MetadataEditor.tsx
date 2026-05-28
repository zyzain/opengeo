import React, { useState } from 'react';
import { Card, Form, Input, Select, Button, ColorPicker, Space, Tag, message } from 'antd';
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons';
import { useIntl } from 'react-intl';

const { Option } = Select;
const { TextArea } = Input;

interface MetadataEditorProps {
  brandId: number;
  metadata?: any;
  onSave: (values: any) => Promise<void>;
  loading?: boolean;
}

const MetadataEditor: React.FC<MetadataEditorProps> = ({ brandId, metadata, onSave, loading }) => {
  const intl = useIntl();
  const [form] = Form.useForm();
  const [viKeywords, setViKeywords] = useState<string[]>(metadata?.vi_profile?.brand_keywords || []);
  const [preferredPhrases, setPreferredPhrases] = useState<string[]>(metadata?.tone_profile?.preferred_phrases || []);
  const [avoidWords, setAvoidWords] = useState<string[]>(metadata?.tone_profile?.avoid_words || []);

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      const metadataValues = {
        vi_profile: {
          ...values.vi_profile,
          brand_keywords: viKeywords,
        },
        tone_profile: {
          ...values.tone_profile,
          preferred_phrases: preferredPhrases,
          avoid_words: avoidWords,
        },
        audience_profiles: values.audience_profiles || [],
        brand_values: values.brand_values || [],
        unique_selling_points: values.unique_selling_points || [],
      };
      await onSave(metadataValues);
      message.success(intl.formatMessage({ id: 'metadata.message.saveSuccess' }));
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  return (
    <div>
      <Form
        form={form}
        layout="vertical"
        initialValues={{
          vi_profile: metadata?.vi_profile || {},
          tone_profile: metadata?.tone_profile || {},
          audience_profiles: metadata?.audience_profiles || [],
          brand_values: metadata?.brand_values || [],
          unique_selling_points: metadata?.unique_selling_points || [],
        }}
      >
        <Card title={intl.formatMessage({ id: 'metadata.section.viSpec' })} style={{ marginBottom: 16 }}>
          <Form.Item name={['vi_profile', 'primary_color']} label={intl.formatMessage({ id: 'metadata.form.primaryColor' })}>
            <Input placeholder="#1890ff" />
          </Form.Item>
          <Form.Item name={['vi_profile', 'secondary_color']} label={intl.formatMessage({ id: 'metadata.form.secondaryColor' })}>
            <Input placeholder="#52c41a" />
          </Form.Item>
          <Form.Item name={['vi_profile', 'font_family']} label={intl.formatMessage({ id: 'metadata.form.fontFamily' })}>
            <Input placeholder="PingFang SC" />
          </Form.Item>
          <Form.Item name={['vi_profile', 'slogan']} label={intl.formatMessage({ id: 'metadata.form.slogan' })}>
            <Input placeholder={intl.formatMessage({ id: 'metadata.placeholder.slogan' })} />
          </Form.Item>
          <Form.Item label={intl.formatMessage({ id: 'metadata.form.brandKeywords' })}>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, marginBottom: 8 }}>
              {viKeywords.map((keyword, index) => (
                <Tag
                  key={index}
                  closable
                  onClose={() => setViKeywords(viKeywords.filter((_, i) => i !== index))}
                >
                  {keyword}
                </Tag>
              ))}
            </div>
            <Input
              placeholder={intl.formatMessage({ id: 'metadata.placeholder.keywords' })}
              onPressEnter={(e) => {
                const value = (e.target as HTMLInputElement).value.trim();
                if (value && !viKeywords.includes(value)) {
                  setViKeywords([...viKeywords, value]);
                  (e.target as HTMLInputElement).value = '';
                }
              }}
            />
          </Form.Item>
        </Card>

        <Card title={intl.formatMessage({ id: 'metadata.section.toneSpec' })} style={{ marginBottom: 16 }}>
          <Form.Item name={['tone_profile', 'formality']} label={intl.formatMessage({ id: 'metadata.form.formality' })}>
            <Select placeholder={intl.formatMessage({ id: 'metadata.placeholder.formality' })}>
              <Option value="formal">{intl.formatMessage({ id: 'metadata.option.formal' })}</Option>
              <Option value="casual">{intl.formatMessage({ id: 'metadata.option.casual' })}</Option>
              <Option value="technical">{intl.formatMessage({ id: 'metadata.option.technical' })}</Option>
            </Select>
          </Form.Item>
          <Form.Item name={['tone_profile', 'personality']} label={intl.formatMessage({ id: 'metadata.form.personality' })}>
            <Select placeholder={intl.formatMessage({ id: 'metadata.placeholder.personality' })}>
              <Option value="friendly">{intl.formatMessage({ id: 'metadata.option.friendly' })}</Option>
              <Option value="professional">{intl.formatMessage({ id: 'metadata.option.professional' })}</Option>
              <Option value="playful">{intl.formatMessage({ id: 'metadata.option.playful' })}</Option>
              <Option value="authoritative">{intl.formatMessage({ id: 'metadata.option.authoritative' })}</Option>
            </Select>
          </Form.Item>
          <Form.Item name={['tone_profile', 'style_guide']} label={intl.formatMessage({ id: 'metadata.form.styleGuide' })}>
            <TextArea rows={3} placeholder={intl.formatMessage({ id: 'metadata.placeholder.styleGuide' })} />
          </Form.Item>
          <Form.Item label={intl.formatMessage({ id: 'metadata.form.preferredPhrases' })}>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, marginBottom: 8 }}>
              {preferredPhrases.map((phrase, index) => (
                <Tag
                  key={index}
                  closable
                  onClose={() => setPreferredPhrases(preferredPhrases.filter((_, i) => i !== index))}
                >
                  {phrase}
                </Tag>
              ))}
            </div>
            <Input
              placeholder={intl.formatMessage({ id: 'metadata.placeholder.preferredPhrases' })}
              onPressEnter={(e) => {
                const value = (e.target as HTMLInputElement).value.trim();
                if (value && !preferredPhrases.includes(value)) {
                  setPreferredPhrases([...preferredPhrases, value]);
                  (e.target as HTMLInputElement).value = '';
                }
              }}
            />
          </Form.Item>
          <Form.Item label={intl.formatMessage({ id: 'metadata.form.avoidWords' })}>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, marginBottom: 8 }}>
              {avoidWords.map((word, index) => (
                <Tag
                  key={index}
                  color="red"
                  closable
                  onClose={() => setAvoidWords(avoidWords.filter((_, i) => i !== index))}
                >
                  {word}
                </Tag>
              ))}
            </div>
            <Input
              placeholder={intl.formatMessage({ id: 'metadata.placeholder.avoidWords' })}
              onPressEnter={(e) => {
                const value = (e.target as HTMLInputElement).value.trim();
                if (value && !avoidWords.includes(value)) {
                  setAvoidWords([...avoidWords, value]);
                  (e.target as HTMLInputElement).value = '';
                }
              }}
            />
          </Form.Item>
        </Card>

        <Card title={intl.formatMessage({ id: 'metadata.section.audience' })} style={{ marginBottom: 16 }}>
          <Form.List name="audience_profiles">
            {(fields, { add, remove }) => (
              <>
                {fields.map(({ key, name, ...restField }) => (
                  <Card
                    key={key}
                    size="small"
                    style={{ marginBottom: 16 }}
                    extra={
                      <Button
                        type="text"
                        danger
                        icon={<MinusCircleOutlined />}
                        onClick={() => remove(name)}
                      >
                        {intl.formatMessage({ id: 'common.action.delete' })}
                      </Button>
                    }
                  >
                    <Form.Item
                      {...restField}
                      name={[name, 'name']}
                      label={intl.formatMessage({ id: 'metadata.form.audienceName' })}
                      rules={[{ required: true, message: intl.formatMessage({ id: 'metadata.validation.enterAudienceName' }) }]}
                    >
                      <Input placeholder={intl.formatMessage({ id: 'metadata.placeholder.audienceNameExample' })} />
                    </Form.Item>
                    <Form.Item {...restField} name={[name, 'age_range']} label={intl.formatMessage({ id: 'metadata.form.ageRange' })}>
                      <Input placeholder={intl.formatMessage({ id: 'metadata.form.agePlaceholder' })} />
                    </Form.Item>
                    <Form.Item {...restField} name={[name, 'interests']} label={intl.formatMessage({ id: 'metadata.form.interests' })}>
                      <Select mode="tags" placeholder={intl.formatMessage({ id: 'metadata.placeholder.interests' })} />
                    </Form.Item>
                    <Form.Item {...restField} name={[name, 'pain_points']} label={intl.formatMessage({ id: 'metadata.form.painPoints' })}>
                      <Select mode="tags" placeholder={intl.formatMessage({ id: 'metadata.placeholder.painPoints' })} />
                    </Form.Item>
                  </Card>
                ))}
                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                  {intl.formatMessage({ id: 'metadata.action.addAudience' })}
                </Button>
              </>
            )}
          </Form.List>
        </Card>

        <Card title={intl.formatMessage({ id: 'metadata.section.brandValues' })} style={{ marginBottom: 16 }}>
          <Form.Item name="brand_values">
            <Select mode="tags" placeholder={intl.formatMessage({ id: 'metadata.placeholder.brandValues' })} />
          </Form.Item>
        </Card>

        <Card title={intl.formatMessage({ id: 'metadata.section.usp' })} style={{ marginBottom: 16 }}>
          <Form.Item name="unique_selling_points">
            <Select mode="tags" placeholder={intl.formatMessage({ id: 'metadata.placeholder.usp' })} />
          </Form.Item>
        </Card>

        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
          <Button type="primary" onClick={handleSave} loading={loading}>
            {intl.formatMessage({ id: 'metadata.action.save' })}
          </Button>
        </div>
      </Form>
    </div>
  );
};

export default MetadataEditor;
