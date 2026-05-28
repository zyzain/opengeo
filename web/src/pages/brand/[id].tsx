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
  Row,
  Col,
  Statistic,
  Breadcrumb,
} from 'antd';
import {
  EditOutlined,
  HistoryOutlined,
  ArrowLeftOutlined,
  DeleteOutlined,
  TeamOutlined,
  FileTextOutlined,
  BranchesOutlined,
  BookOutlined,
  HomeOutlined,
} from '@ant-design/icons';
import { useIntl } from 'react-intl';
import { useBrand, useBrandMetadata, useGlossary } from '../../hooks/useBrand';
import BrandForm from '../../components/brand/BrandForm';
import MetadataEditor from '../../components/brand/MetadataEditor';
import GlossaryManager from '../../components/brand/GlossaryManager';
import KnowledgeGraph from '../../components/brand/KnowledgeGraph';
import BrandSnapshot from '../../components/brand/BrandSnapshot';

const { TabPane } = Tabs;

const BrandDetailPage: React.FC = () => {
  const intl = useIntl();
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
      message.success(intl.formatMessage({ id: 'brand.message.updateSuccess' }));
      setEditModalVisible(false);
      refetchBrand();
    } catch (error) {
      message.error(intl.formatMessage({ id: 'common.message.updateFailed' }));
    }
  };

  const handleDeleteBrand = () => {
    Modal.confirm({
      title: intl.formatMessage({ id: 'common.confirmDelete' }),
      content: intl.formatMessage({ id: 'brand.confirmDelete' }, { name: brand?.name }),
      okText: intl.formatMessage({ id: 'common.action.delete' }),
      okType: 'danger',
      onOk: async () => {
        try {
          await fetch(`/api/v1/brand/${brandId}`, { method: 'DELETE' });
          message.success(intl.formatMessage({ id: 'brand.message.deleteSuccess' }));
          navigate('/brand');
        } catch (error) {
          message.error(intl.formatMessage({ id: 'common.message.deleteFailed' }));
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
      message.success(intl.formatMessage({ id: 'brand.message.metadataSaveSuccess' }));
    } catch (error) {
      message.error(intl.formatMessage({ id: 'brand.message.metadataSaveFailed' }));
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
      message.error(intl.formatMessage({ id: 'brand.message.addEntityFailed' }));
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
      message.error(intl.formatMessage({ id: 'brand.message.addRelationFailed' }));
    }
  };

  const handleDeleteEntity = async (entityId: number) => {
    try {
      await fetch(`/api/v1/brand/${brandId}/knowledge/entities/${entityId}`, {
        method: 'DELETE',
      });
      fetchKnowledgeData();
    } catch (error) {
      message.error(intl.formatMessage({ id: 'brand.message.deleteEntityFailed' }));
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
      message.success(intl.formatMessage({ id: 'brand.message.snapshotCreateSuccess' }));
    } catch (error) {
      message.error(intl.formatMessage({ id: 'brand.message.createSnapshotFailed' }));
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
        message.info(intl.formatMessage({ id: 'brand.message.snapshotCompareInDev' }));
      }
    } catch (error) {
      message.error(intl.formatMessage({ id: 'brand.message.compareFailed' }));
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
        <h2>{intl.formatMessage({ id: 'brand.message.brandNotFound' })}</h2>
        <Button onClick={() => navigate('/brand')}>{intl.formatMessage({ id: 'brand.action.backToList' })}</Button>
      </div>
    );
  }

  const statusMap: Record<number, { color: string; text: string }> = {
    1: { color: 'green', text: intl.formatMessage({ id: 'brand.status.active' }) },
    2: { color: 'orange', text: intl.formatMessage({ id: 'brand.status.archived' }) },
    3: { color: 'red', text: intl.formatMessage({ id: 'brand.status.disabled' }) },
  };
  const statusInfo = statusMap[brand.status] || { color: 'default', text: intl.formatMessage({ id: 'common.status.unknown' }) };

  return (
    <div style={{ padding: 24 }}>
      <Breadcrumb
        style={{ marginBottom: 16 }}
        items={[
          {
            title: <a onClick={() => navigate('/brand')}><HomeOutlined /> {intl.formatMessage({ id: 'brand.breadcrumb.management' })}</a>,
          },
          {
            title: brand.name,
          },
        ]}
      />

      <Card
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            {brand.logo_url ? (
              <img
                src={brand.logo_url}
                alt={brand.name}
                style={{ width: 48, height: 48, borderRadius: 8, objectFit: 'cover' }}
              />
            ) : (
              <div
                style={{
                  width: 48,
                  height: 48,
                  borderRadius: 8,
                  backgroundColor: '#1890ff',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#fff',
                  fontSize: 24,
                  fontWeight: 'bold',
                }}
              >
                {brand.name.charAt(0)}
              </div>
            )}
            <div>
              <h2 style={{ margin: 0, fontSize: 20 }}>{brand.name}</h2>
              <Tag color={statusInfo.color} style={{ marginTop: 4 }}>
                {statusInfo.text}
              </Tag>
            </div>
          </div>
        }
        extra={
          <Space>
            <Button icon={<EditOutlined />} onClick={() => setEditModalVisible(true)}>
              {intl.formatMessage({ id: 'common.action.edit' })}
            </Button>
            <Button icon={<HistoryOutlined />} onClick={fetchSnapshots}>
              {intl.formatMessage({ id: 'common.action.refresh' })}
            </Button>
            <Button danger icon={<DeleteOutlined />} onClick={handleDeleteBrand}>
              {intl.formatMessage({ id: 'common.action.delete' })}
            </Button>
          </Space>
        }
      >
        <Descriptions bordered column={2}>
          <Descriptions.Item label={intl.formatMessage({ id: 'brand.field.slug' })}>
            <Tag>{brand.slug}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label={intl.formatMessage({ id: 'brand.field.industry' })}>{brand.industry || '-'}</Descriptions.Item>
          <Descriptions.Item label={intl.formatMessage({ id: 'brand.field.website' })}>
            {brand.website ? (
              <a href={brand.website} target="_blank" rel="noopener noreferrer">
                {brand.website}
              </a>
            ) : '-'}
          </Descriptions.Item>
          <Descriptions.Item label={intl.formatMessage({ id: 'brand.field.status' })}>
            <Tag color={statusInfo.color}>{statusInfo.text}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label={intl.formatMessage({ id: 'brand.field.description' })} span={2}>
            {brand.description || '-'}
          </Descriptions.Item>
          <Descriptions.Item label={intl.formatMessage({ id: 'brand.field.createdAt' })}>
            {new Date(brand.created_at).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label={intl.formatMessage({ id: 'brand.field.updatedAt' })}>
            {new Date(brand.updated_at).toLocaleString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Row gutter={16} style={{ marginTop: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'brand.stat.termCount' })}
              value={glossary?.length || 0}
              prefix={<BookOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'brand.stat.knowledgeEntities' })}
              value={knowledgeEntities?.length || 0}
              prefix={<BranchesOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'brand.stat.snapshotCount' })}
              value={snapshots?.length || 0}
              prefix={<HistoryOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'common.field.status' })}
              value={statusInfo.text}
              valueStyle={{ color: statusInfo.color === 'green' ? '#3f8600' : '#cf1322' }}
              prefix={<TeamOutlined />}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ marginTop: 24 }}>
        <Tabs defaultActiveKey="metadata">
          <TabPane tab={intl.formatMessage({ id: 'brand.tab.metadata' })} key="metadata">
            <MetadataEditor
              brandId={brandId}
              metadata={metadata}
              onSave={handleSaveMetadata}
              loading={loading}
            />
          </TabPane>

          <TabPane tab={intl.formatMessage({ id: 'brand.tab.glossary' })} key="glossary">
            <GlossaryManager
              brandId={brandId}
              entries={glossary}
              loading={glossaryLoading}
              onRefresh={refetchGlossary}
            />
          </TabPane>

          <TabPane tab={intl.formatMessage({ id: 'brand.tab.knowledge' })} key="knowledge">
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

          <TabPane tab={intl.formatMessage({ id: 'brand.tab.snapshots' })} key="snapshots">
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
        title={intl.formatMessage({ id: 'brand.modal.editBrand' })}
        open={editModalVisible}
        onCancel={() => setEditModalVisible(false)}
        footer={null}
        width={800}
        destroyOnClose
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
