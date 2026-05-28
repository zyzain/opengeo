import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card,
  Tabs,
  Descriptions,
  Tag,
  Button,
  Space,
  Spin,
  message,
  Modal,
} from 'antd';
import {
  EditOutlined,
  HistoryOutlined,
  ArrowLeftOutlined,
  DeleteOutlined,
} from '@ant-design/icons';
import { useBrand, useBrandMetadata, useGlossary } from '../../hooks/useBrand';
import BrandForm from '../../components/brand/BrandForm';
import MetadataEditor from '../../components/brand/MetadataEditor';
import GlossaryManager from '../../components/brand/GlossaryManager';
import KnowledgeGraph from '../../components/brand/KnowledgeGraph';
import BrandSnapshot from '../../components/brand/BrandSnapshot';

const { TabPane } = Tabs;

const BrandDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const brandId = Number(id);

  const { brand, loading: brandLoading, refetch: refetchBrand } = useBrand(brandId);
  const { metadata, loading: metadataLoading, refetch: refetchMetadata } = useBrandMetadata(brandId);
  const { entries: glossary, loading: glossaryLoading, refetch: refetchGlossary } = useGlossary(brandId);

  const [editModalVisible, setEditModalVisible] = useState(false);
  const [knowledgeEntities, setKnowledgeEntities] = useState<any[]>([]);
  const [knowledgeRelations, setKnowledgeRelations] = useState<any[]>([]);
  const [snapshots, setSnapshots] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchKnowledgeData();
    fetchSnapshots();
  }, [brandId]);

  const fetchKnowledgeData = async () => {
    try {
      const response = await fetch(`/api/v1/brand/${brandId}/knowledge/entities`);
      if (response.ok) {
        const data = await response.json();
        setKnowledgeEntities(data.entities || []);
      }
    } catch (error) {
      console.error('Failed to fetch knowledge entities:', error);
    }
  };

  const fetchSnapshots = async () => {
    try {
      const response = await fetch(`/api/v1/brand/${brandId}/snapshots`);
      if (response.ok) {
        const data = await response.json();
        setSnapshots(data.snapshots || []);
      }
    } catch (error) {
      console.error('Failed to fetch snapshots:', error);
    }
  };

  const handleUpdateBrand = async (values: any) => {
    try {
      await fetch(`/api/v1/brand/${brandId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });
      message.success('品牌更新成功');
      setEditModalVisible(false);
      refetchBrand();
    } catch (error) {
      message.error('更新失败');
    }
  };

  const handleDeleteBrand = () => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除品牌 "${brand?.name}" 吗？此操作不可恢复。`,
      okText: '删除',
      okType: 'danger',
      onOk: async () => {
        try {
          await fetch(`/api/v1/brand/${brandId}`, { method: 'DELETE' });
          message.success('品牌删除成功');
          navigate('/brand');
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const handleSaveMetadata = async (values: any) => {
    try {
      await fetch(`/api/v1/brand/${brandId}/metadata`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });
      refetchMetadata();
    } catch (error) {
      message.error('元数据保存失败');
    }
  };

  const handleAddEntity = async (entity: any) => {
    try {
      await fetch(`/api/v1/brand/${brandId}/knowledge/entities`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(entity),
      });
      fetchKnowledgeData();
    } catch (error) {
      message.error('添加实体失败');
    }
  };

  const handleAddRelation = async (relation: any) => {
    try {
      await fetch(`/api/v1/brand/${brandId}/knowledge/relations`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(relation),
      });
      fetchKnowledgeData();
    } catch (error) {
      message.error('添加关系失败');
    }
  };

  const handleDeleteEntity = async (entityId: number) => {
    try {
      await fetch(`/api/v1/brand/${brandId}/knowledge/entities/${entityId}`, {
        method: 'DELETE',
      });
      fetchKnowledgeData();
    } catch (error) {
      message.error('删除实体失败');
    }
  };

  const handleCreateSnapshot = async (values: any) => {
    try {
      await fetch(`/api/v1/brand/${brandId}/snapshots`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });
      fetchSnapshots();
    } catch (error) {
      message.error('创建快照失败');
    }
  };

  const handleCompareSnapshots = async (id1: number, id2: number) => {
    try {
      const response = await fetch(`/api/v1/brand/${brandId}/snapshots/compare`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ snapshot_id_1: id1, snapshot_id_2: id2 }),
      });
      if (response.ok) {
        const data = await response.json();
        message.info('快照对比功能开发中...');
      }
    } catch (error) {
      message.error('对比失败');
    }
  };

  if (brandLoading || metadataLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!brand) {
    return (
      <div style={{ padding: 24, textAlign: 'center' }}>
        <h2>品牌不存在</h2>
        <Button onClick={() => navigate('/brand')}>返回品牌列表</Button>
      </div>
    );
  }

  return (
    <div style={{ padding: 24 }}>
      <div style={{ marginBottom: 16 }}>
        <Button
          icon={<ArrowLeftOutlined />}
          onClick={() => navigate('/brand')}
          style={{ marginBottom: 16 }}
        >
          返回品牌列表
        </Button>
      </div>

      <Card
        title={
          <Space>
            {brand.logo_url && (
              <img src={brand.logo_url} alt={brand.name} style={{ height: 32 }} />
            )}
            <span>{brand.name}</span>
          </Space>
        }
        extra={
          <Space>
            <Button icon={<EditOutlined />} onClick={() => setEditModalVisible(true)}>
              编辑
            </Button>
            <Button icon={<HistoryOutlined />} onClick={() => fetchSnapshots()}>
              刷新快照
            </Button>
            <Button danger icon={<DeleteOutlined />} onClick={handleDeleteBrand}>
              删除
            </Button>
          </Space>
        }
      >
        <Descriptions bordered column={2}>
          <Descriptions.Item label="品牌标识">{brand.slug}</Descriptions.Item>
          <Descriptions.Item label="行业">{brand.industry}</Descriptions.Item>
          <Descriptions.Item label="官网">
            {brand.website ? (
              <a href={brand.website} target="_blank" rel="noopener noreferrer">
                {brand.website}
              </a>
            ) : '-'}
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={brand.status === 1 ? 'green' : 'red'}>
              {brand.status === 1 ? '活跃' : '已禁用'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="描述" span={2}>
            {brand.description || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">
            {new Date(brand.created_at).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="更新时间">
            {new Date(brand.updated_at).toLocaleString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card style={{ marginTop: 24 }}>
        <Tabs defaultActiveKey="metadata">
          <TabPane tab="品牌元数据" key="metadata">
            <MetadataEditor
              brandId={brandId}
              metadata={metadata}
              onSave={handleSaveMetadata}
              loading={loading}
            />
          </TabPane>

          <TabPane tab="术语表" key="glossary">
            <GlossaryManager
              brandId={brandId}
              entries={glossary}
              loading={glossaryLoading}
              onRefresh={refetchGlossary}
            />
          </TabPane>

          <TabPane tab="知识图谱" key="knowledge">
            <KnowledgeGraph
              brandId={brandId}
              entities={knowledgeEntities}
              relations={knowledgeRelations}
              loading={loading}
              onAddEntity={handleAddEntity}
              onAddRelation={handleAddRelation}
              onDeleteEntity={handleDeleteEntity}
            />
          </TabPane>

          <TabPane tab="品牌快照" key="snapshots">
            <BrandSnapshot
              brandId={brandId}
              snapshots={snapshots}
              loading={loading}
              onCreateSnapshot={handleCreateSnapshot}
              onCompareSnapshots={handleCompareSnapshots}
            />
          </TabPane>
        </Tabs>
      </Card>

      <Modal
        title="编辑品牌"
        open={editModalVisible}
        onCancel={() => setEditModalVisible(false)}
        footer={null}
        width={800}
      >
        <BrandForm
          initialValues={brand}
          onSubmit={handleUpdateBrand}
          onCancel={() => setEditModalVisible(false)}
          loading={loading}
        />
      </Modal>
    </div>
  );
};

export default BrandDetailPage;
