import React, { useState, useEffect } from 'react';
import { Card, Button, Space, Modal, Form, Input, Select, Tag, Empty, Spin, message } from 'antd';
import { PlusOutlined, NodeIndexOutlined } from '@ant-design/icons';
import { useIntl } from 'react-intl';

const { Option } = Select;

interface KnowledgeEntity {
  id: number;
  name: string;
  type: string;
  description: string;
  properties: Record<string, string>;
  tags: string[];
}

interface KnowledgeRelation {
  id: number;
  from_entity_id: number;
  to_entity_id: number;
  type: string;
  weight: number;
}

interface KnowledgeGraphProps {
  brandId: number;
  entities: KnowledgeEntity[];
  relations: KnowledgeRelation[];
  loading?: boolean;
  onAddEntity: (entity: any) => Promise<void>;
  onAddRelation: (relation: any) => Promise<void>;
  onDeleteEntity: (id: number) => Promise<void>;
}

const KnowledgeGraph: React.FC<KnowledgeGraphProps> = ({
  brandId,
  entities,
  relations,
  loading,
  onAddEntity,
  onAddRelation,
  onDeleteEntity,
}) => {
  const intl = useIntl();
  const [entityModalVisible, setEntityModalVisible] = useState(false);
  const [relationModalVisible, setRelationModalVisible] = useState(false);
  const [selectedEntity, setSelectedEntity] = useState<KnowledgeEntity | null>(null);
  const [entityForm] = Form.useForm();
  const [relationForm] = Form.useForm();

  const entityTypeColors: Record<string, string> = {
    brand: '#1890ff',
    product: '#52c41a',
    person: '#722ed1',
    org: '#13c2c2',
    event: '#fa8c16',
    concept: '#eb2f96',
    location: '#fadb14',
    technology: '#2f54eb',
  };

  const entityTypeLabels: Record<string, string> = {
    brand: intl.formatMessage({ id: 'knowledge.type.brand' }),
    product: intl.formatMessage({ id: 'knowledge.type.product' }),
    person: intl.formatMessage({ id: 'knowledge.type.person' }),
    org: intl.formatMessage({ id: 'knowledge.type.org' }),
    event: intl.formatMessage({ id: 'knowledge.type.event' }),
    concept: intl.formatMessage({ id: 'knowledge.type.concept' }),
    location: intl.formatMessage({ id: 'knowledge.type.location' }),
    technology: intl.formatMessage({ id: 'knowledge.type.technology' }),
  };

  const handleAddEntity = async () => {
    try {
      const values = await entityForm.validateFields();
      await onAddEntity({ ...values, brand_id: brandId });
      setEntityModalVisible(false);
      entityForm.resetFields();
      message.success(intl.formatMessage({ id: 'knowledge.message.entityAdded' }));
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  const handleAddRelation = async () => {
    try {
      const values = await relationForm.validateFields();
      await onAddRelation(values);
      setRelationModalVisible(false);
      relationForm.resetFields();
      message.success(intl.formatMessage({ id: 'knowledge.message.relationAdded' }));
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  const renderEntityCard = (entity: KnowledgeEntity) => (
    <Card
      key={entity.id}
      size="small"
      style={{
        width: 200,
        display: 'inline-block',
        margin: 8,
        borderColor: entityTypeColors[entity.type] || '#d9d9d9',
      }}
      hoverable
      onClick={() => setSelectedEntity(entity)}
    >
      <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
        <div
          style={{
            width: 12,
            height: 12,
            borderRadius: '50%',
            backgroundColor: entityTypeColors[entity.type] || '#d9d9d9',
          }}
        />
        <Tag color={entityTypeColors[entity.type]}>
          {entityTypeLabels[entity.type] || entity.type}
        </Tag>
      </div>
      <div style={{ fontWeight: 'bold', marginBottom: 4 }}>{entity.name}</div>
      {entity.description && (
        <div style={{ fontSize: 12, color: '#666', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
          {entity.description}
        </div>
      )}
      {entity.tags && entity.tags.length > 0 && (
        <div style={{ marginTop: 8 }}>
          {entity.tags.slice(0, 2).map(tag => (
            <Tag key={tag} style={{ fontSize: 10 }}>{tag}</Tag>
          ))}
          {entity.tags.length > 2 && <Tag style={{ fontSize: 10 }}>+{entity.tags.length - 2}</Tag>}
        </div>
      )}
    </Card>
  );

  const renderRelationLine = (relation: KnowledgeRelation) => {
    const fromEntity = entities.find(e => e.id === relation.from_entity_id);
    const toEntity = entities.find(e => e.id === relation.to_entity_id);
    if (!fromEntity || !toEntity) return null;

    return (
      <div key={relation.id} style={{ display: 'flex', alignItems: 'center', gap: 8, margin: '8px 0' }}>
        <Tag>{fromEntity.name}</Tag>
        <span style={{ color: '#999' }}>→</span>
        <Tag color="blue">{relation.type}</Tag>
        <span style={{ color: '#999' }}>→</span>
        <Tag>{toEntity.name}</Tag>
      </div>
    );
  };

  if (loading) {
    return <Spin size="large" style={{ display: 'flex', justifyContent: 'center', padding: 50 }} />;
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h3>{intl.formatMessage({ id: 'knowledge.title' })}</h3>
        <Space>
          <Button
            icon={<PlusOutlined />}
            onClick={() => {
              entityForm.resetFields();
              setEntityModalVisible(true);
            }}
          >
            {intl.formatMessage({ id: 'knowledge.action.addEntity' })}
          </Button>
          <Button
            icon={<NodeIndexOutlined />}
            onClick={() => {
              relationForm.resetFields();
              setRelationModalVisible(true);
            }}
            disabled={entities.length < 2}
          >
            {intl.formatMessage({ id: 'knowledge.action.addRelation' })}
          </Button>
        </Space>
      </div>

      {entities.length === 0 ? (
        <Empty description={intl.formatMessage({ id: 'knowledge.empty.noEntities' })} />
      ) : (
        <div>
          <Card title={intl.formatMessage({ id: 'knowledge.card.entityList' })} size="small" style={{ marginBottom: 16 }}>
            <div style={{ display: 'flex', flexWrap: 'wrap' }}>
              {entities.map(entity => renderEntityCard(entity))}
            </div>
          </Card>

          {relations.length > 0 && (
            <Card title={intl.formatMessage({ id: 'knowledge.card.relationList' })} size="small">
              {relations.map(relation => renderRelationLine(relation))}
            </Card>
          )}
        </div>
      )}

      <Modal
        title={intl.formatMessage({ id: 'knowledge.modal.addEntity' })}
        open={entityModalVisible}
        onOk={handleAddEntity}
        onCancel={() => setEntityModalVisible(false)}
        width={600}
      >
        <Form form={entityForm} layout="vertical">
          <Form.Item name="name" label={intl.formatMessage({ id: 'knowledge.form.entityName' })} rules={[{ required: true, message: intl.formatMessage({ id: 'knowledge.validation.enterEntityName' }) }]}>
            <Input placeholder={intl.formatMessage({ id: 'knowledge.placeholder.entityName' })} />
          </Form.Item>
          <Form.Item name="type" label={intl.formatMessage({ id: 'knowledge.form.entityType' })} rules={[{ required: true, message: intl.formatMessage({ id: 'knowledge.validation.selectEntityType' }) }]}>
            <Select placeholder={intl.formatMessage({ id: 'knowledge.placeholder.entityType' })}>
              {Object.entries(entityTypeLabels).map(([value, label]) => (
                <Option key={value} value={value}>{label}</Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="description" label={intl.formatMessage({ id: 'knowledge.form.entityDesc' })}>
            <Input.TextArea rows={3} placeholder={intl.formatMessage({ id: 'knowledge.placeholder.entityDesc' })} />
          </Form.Item>
          <Form.Item name="tags" label={intl.formatMessage({ id: 'knowledge.form.tags' })}>
            <Select mode="tags" placeholder={intl.formatMessage({ id: 'knowledge.placeholder.tags' })} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={intl.formatMessage({ id: 'knowledge.modal.addRelation' })}
        open={relationModalVisible}
        onOk={handleAddRelation}
        onCancel={() => setRelationModalVisible(false)}
        width={600}
      >
        <Form form={relationForm} layout="vertical">
          <Form.Item name="from_entity_id" label={intl.formatMessage({ id: 'knowledge.form.sourceEntity' })} rules={[{ required: true, message: intl.formatMessage({ id: 'knowledge.validation.selectSource' }) }]}>
            <Select placeholder={intl.formatMessage({ id: 'knowledge.placeholder.sourceEntity' })}>
              {entities.map(entity => (
                <Option key={entity.id} value={entity.id}>{entity.name}</Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="to_entity_id" label={intl.formatMessage({ id: 'knowledge.form.targetEntity' })} rules={[{ required: true, message: intl.formatMessage({ id: 'knowledge.validation.selectTarget' }) }]}>
            <Select placeholder={intl.formatMessage({ id: 'knowledge.placeholder.targetEntity' })}>
              {entities.map(entity => (
                <Option key={entity.id} value={entity.id}>{entity.name}</Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="type" label={intl.formatMessage({ id: 'knowledge.form.relationType' })} rules={[{ required: true, message: intl.formatMessage({ id: 'knowledge.validation.selectRelationType' }) }]}>
            <Select placeholder={intl.formatMessage({ id: 'knowledge.placeholder.relationType' })}>
              <Option value="is_a">{intl.formatMessage({ id: 'knowledge.relation.is_a' })}</Option>
              <Option value="part_of">{intl.formatMessage({ id: 'knowledge.relation.part_of' })}</Option>
              <Option value="related_to">{intl.formatMessage({ id: 'knowledge.relation.related_to' })}</Option>
              <Option value="competes_with">{intl.formatMessage({ id: 'knowledge.relation.competes_with' })}</Option>
              <Option value="mentions">{intl.formatMessage({ id: 'knowledge.relation.mentions' })}</Option>
              <Option value="depends_on">{intl.formatMessage({ id: 'knowledge.relation.depends_on' })}</Option>
              <Option value="owns">{intl.formatMessage({ id: 'knowledge.relation.owns' })}</Option>
              <Option value="created_by">{intl.formatMessage({ id: 'knowledge.relation.created_by' })}</Option>
            </Select>
          </Form.Item>
          <Form.Item name="weight" label={intl.formatMessage({ id: 'knowledge.form.weight' })}>
            <Input type="number" min={0} max={1} step={0.1} placeholder="0-1" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={intl.formatMessage({ id: 'knowledge.modal.entityDetail' })}
        open={!!selectedEntity}
        onCancel={() => setSelectedEntity(null)}
        footer={null}
        width={600}
      >
        {selectedEntity && (
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 16 }}>
              <div
                style={{
                  width: 16,
                  height: 16,
                  borderRadius: '50%',
                  backgroundColor: entityTypeColors[selectedEntity.type] || '#d9d9d9',
                }}
              />
              <Tag color={entityTypeColors[selectedEntity.type]}>
                {entityTypeLabels[selectedEntity.type] || selectedEntity.type}
              </Tag>
            </div>
            <h2>{selectedEntity.name}</h2>
            {selectedEntity.description && <p>{selectedEntity.description}</p>}
            {selectedEntity.tags && selectedEntity.tags.length > 0 && (
              <div>
                <strong>{intl.formatMessage({ id: 'knowledge.field.tags' })}</strong>
                {selectedEntity.tags.map(tag => <Tag key={tag}>{tag}</Tag>)}
              </div>
            )}
            <div style={{ marginTop: 16 }}>
              <Button
                danger
                onClick={async () => {
                  await onDeleteEntity(selectedEntity.id);
                  setSelectedEntity(null);
                  message.success(intl.formatMessage({ id: 'knowledge.message.entityDeleted' }));
                }}
              >
                {intl.formatMessage({ id: 'knowledge.action.deleteEntity' })}
              </Button>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default KnowledgeGraph;
