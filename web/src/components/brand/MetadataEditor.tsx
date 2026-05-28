import React, { useState } from 'react';
import { Card, Form, Input, Select, Button, ColorPicker, Space, Tag, message } from 'antd';
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons';

const { Option } = Select;
const { TextArea } = Input;

interface MetadataEditorProps {
  brandId: number;
  metadata?: any;
  onSave: (values: any) => Promise<void>;
  loading?: boolean;
}

const MetadataEditor: React.FC<MetadataEditorProps> = ({ brandId, metadata, onSave, loading }) => {
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
      message.success('元数据保存成功');
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
        <Card title="VI 规范" style={{ marginBottom: 16 }}>
          <Form.Item name={['vi_profile', 'primary_color']} label="主色">
            <Input placeholder="#1890ff" />
          </Form.Item>
          <Form.Item name={['vi_profile', 'secondary_color']} label="副色">
            <Input placeholder="#52c41a" />
          </Form.Item>
          <Form.Item name={['vi_profile', 'font_family']} label="字体">
            <Input placeholder="PingFang SC" />
          </Form.Item>
          <Form.Item name={['vi_profile', 'slogan']} label="品牌口号">
            <Input placeholder="请输入品牌口号" />
          </Form.Item>
          <Form.Item label="品牌关键词">
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
              placeholder="输入关键词后按回车"
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

        <Card title="语调规范" style={{ marginBottom: 16 }}>
          <Form.Item name={['tone_profile', 'formality']} label="正式度">
            <Select placeholder="请选择正式度">
              <Option value="formal">正式</Option>
              <Option value="casual">随意</Option>
              <Option value="technical">技术</Option>
            </Select>
          </Form.Item>
          <Form.Item name={['tone_profile', 'personality']} label="个性">
            <Select placeholder="请选择个性">
              <Option value="friendly">友好</Option>
              <Option value="professional">专业</Option>
              <Option value="playful">活泼</Option>
              <Option value="authoritative">权威</Option>
            </Select>
          </Form.Item>
          <Form.Item name={['tone_profile', 'style_guide']} label="风格指南">
            <TextArea rows={3} placeholder="请输入写作风格指南" />
          </Form.Item>
          <Form.Item label="偏好短语">
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
              placeholder="输入偏好短语后按回车"
              onPressEnter={(e) => {
                const value = (e.target as HTMLInputElement).value.trim();
                if (value && !preferredPhrases.includes(value)) {
                  setPreferredPhrases([...preferredPhrases, value]);
                  (e.target as HTMLInputElement).value = '';
                }
              }}
            />
          </Form.Item>
          <Form.Item label="禁用词">
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
              placeholder="输入禁用词后按回车"
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

        <Card title="受众画像" style={{ marginBottom: 16 }}>
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
                        删除
                      </Button>
                    }
                  >
                    <Form.Item
                      {...restField}
                      name={[name, 'name']}
                      label="受众名称"
                      rules={[{ required: true, message: '请输入受众名称' }]}
                    >
                      <Input placeholder="例如：企业用户" />
                    </Form.Item>
                    <Form.Item {...restField} name={[name, 'age_range']} label="年龄范围">
                      <Input placeholder="例如：25-45" />
                    </Form.Item>
                    <Form.Item {...restField} name={[name, 'interests']} label="兴趣">
                      <Select mode="tags" placeholder="输入兴趣标签" />
                    </Form.Item>
                    <Form.Item {...restField} name={[name, 'pain_points']} label="痛点">
                      <Select mode="tags" placeholder="输入痛点" />
                    </Form.Item>
                  </Card>
                ))}
                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                  添加受众画像
                </Button>
              </>
            )}
          </Form.List>
        </Card>

        <Card title="品牌价值观" style={{ marginBottom: 16 }}>
          <Form.Item name="brand_values">
            <Select mode="tags" placeholder="输入品牌价值观" />
          </Form.Item>
        </Card>

        <Card title="独特卖点" style={{ marginBottom: 16 }}>
          <Form.Item name="unique_selling_points">
            <Select mode="tags" placeholder="输入独特卖点" />
          </Form.Item>
        </Card>

        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
          <Button type="primary" onClick={handleSave} loading={loading}>
            保存元数据
          </Button>
        </div>
      </Form>
    </div>
  );
};

export default MetadataEditor;
