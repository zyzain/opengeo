import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Tabs, Descriptions, Tag, Button, Space, Spin } from 'antd';
import { EditOutlined, HistoryOutlined } from '@ant-design/icons';
import { useBrand, useBrandMetadata, useGlossary } from '../../hooks/useBrand';
import GlossaryManager from '../../components/brand/GlossaryManager';

const { TabPane } = Tabs;

const BrandDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const brandId = Number(id);
  const { brand, loading: brandLoading } = useBrand(brandId);
  const { metadata, loading: metadataLoading } = useBrandMetadata(brandId);
  const { entries: glossary, loading: glossaryLoading } = useGlossary(brandId);

  if (brandLoading || metadataLoading) {
    return <Spin size="large" style={{ display: 'flex', justifyContent: 'center', marginTop: 100 }} />;
  }

  if (!brand) {
    return <div>品牌不存在</div>;
  }

  return (
    <div style={{ padding: 24 }}>
      <Card
        title={brand.name}
        extra={
          <Space>
            <Button icon={<EditOutlined />}>编辑</Button>
            <Button icon={<HistoryOutlined />}>快照</Button>
          </Space>
        }
      >
        <Descriptions bordered column={2}>
          <Descriptions.Item label="品牌标识">{brand.slug}</Descriptions.Item>
          <Descriptions.Item label="行业">{brand.industry}</Descriptions.Item>
          <Descriptions.Item label="官网">{brand.website}</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={brand.status === 1 ? 'green' : 'red'}>
              {brand.status === 1 ? '活跃' : '已禁用'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="描述" span={2}>{brand.description}</Descriptions.Item>
          <Descriptions.Item label="创建时间">{new Date(brand.created_at).toLocaleString()}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{new Date(brand.updated_at).toLocaleString()}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card style={{ marginTop: 24 }}>
        <Tabs defaultActiveKey="metadata">
          <TabPane tab="品牌元数据" key="metadata">
            {metadata && (
              <div>
                <h3>VI 规范</h3>
                <Descriptions bordered column={2}>
                  <Descriptions.Item label="主色">
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <div style={{ width: 24, height: 24, backgroundColor: metadata.vi_profile.primary_color, borderRadius: 4 }} />
                      {metadata.vi_profile.primary_color}
                    </div>
                  </Descriptions.Item>
                  <Descriptions.Item label="副色">
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <div style={{ width: 24, height: 24, backgroundColor: metadata.vi_profile.secondary_color, borderRadius: 4 }} />
                      {metadata.vi_profile.secondary_color}
                    </div>
                  </Descriptions.Item>
                  <Descriptions.Item label="字体">{metadata.vi_profile.font_family}</Descriptions.Item>
                  <Descriptions.Item label="口号">{metadata.vi_profile.slogan}</Descriptions.Item>
                  <Descriptions.Item label="关键词" span={2}>
                    {metadata.vi_profile.brand_keywords.map(k => <Tag key={k}>{k}</Tag>)}
                  </Descriptions.Item>
                </Descriptions>

                <h3 style={{ marginTop: 24 }}>语调规范</h3>
                <Descriptions bordered column={2}>
                  <Descriptions.Item label="正式度">{metadata.tone_profile.formality}</Descriptions.Item>
                  <Descriptions.Item label="个性">{metadata.tone_profile.personality}</Descriptions.Item>
                  <Descriptions.Item label="偏好短语" span={2}>
                    {metadata.tone_profile.preferred_phrases.map(p => <Tag key={p}>{p}</Tag>)}
                  </Descriptions.Item>
                </Descriptions>

                <h3 style={{ marginTop: 24 }}>受众画像</h3>
                {metadata.audience_profiles.map((audience, index) => (
                  <Card key={index} size="small" style={{ marginBottom: 16 }}>
                    <Descriptions bordered column={2}>
                      <Descriptions.Item label="受众名称">{audience.name}</Descriptions.Item>
                      <Descriptions.Item label="年龄范围">{audience.age_range}</Descriptions.Item>
                      <Descriptions.Item label="兴趣" span={2}>
                        {audience.interests.map(i => <Tag key={i}>{i}</Tag>)}
                      </Descriptions.Item>
                    </Descriptions>
                  </Card>
                ))}
              </div>
            )}
          </TabPane>

          <TabPane tab="术语表" key="glossary">
            <GlossaryManager brandId={brandId} entries={glossary} loading={glossaryLoading} />
          </TabPane>

          <TabPane tab="知识图谱" key="knowledge">
            <div>知识图谱功能开发中...</div>
          </TabPane>

          <TabPane tab="品牌快照" key="snapshots">
            <div>品牌快照功能开发中...</div>
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default BrandDetailPage;
