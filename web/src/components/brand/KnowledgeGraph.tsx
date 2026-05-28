import React, { useState, useEffect } from 'react';
import { Card, Button, Space, Modal, Form, Input, Select, Tag, Empty, Spin, message } from 'antd';
import { PlusOutlined, NodeIndexOutlined } from '@ant-design/icons';

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
    brand: '品牌',
    product: '产品',
    person: '人物',
    org: '组织',
    event: '事件',
    concept: '概念',
    location: '地点',
    technology: '技术',
  };

  const handleAddEntity = async () => {
    try {
      const values = await entityForm.validateFields();
      await onAddEntity({ ...values, brand_id: brandId });
      setEntityModalVisible(false);
      entityForm.resetFields();
      message.success('实体添加成功');
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
      message.success('关系添加成功');
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
        <h3>知识图谱</h3>
        <Space>
          <Button
            icon={<PlusOutlined />}
            onClick={() => {
              entityForm.resetFields();
              setEntityModalVisible(true);
            }}
          >
            添加实体
          </Button>
          <Button
            icon={<NodeIndexOutlined />}
            onClick={() => {
              relationForm.resetFields();
              setRelationModalVisible(true);
            }}
            disabled={entities.length < 2}
          >
            添加关系
          </Button>
        </Space>
      </div>

      {entities.length === 0 ? (
        <Empty description="暂无知识实体" />
      ) : (
        <div>
          <Card title="实体列表" size="small" style={{ marginBottom: 16 }}>
            <div style={{ display: 'flex', flexWrap: 'wrap' }}>
              {entities.map(entity => renderEntityCard(entity))}
            </div>
          </Card>

          {relations.length > 0 && (
            <Card title="关系列表" size="small">
              {relations.map(relation => renderRelationLine(relation))}
            </Card>
          )}
        </div>
      )}

      <Modal
        title="添加知识实体"
        open={entityModalVisible}
        onOk={handleAddEntity}
        onCancel={() => setEntityModalVisible(false)}
        width={600}
      >
        <Form form={entityForm} layout="vertical">
          <Form.Item name="name" label="实体名称" rules={[{ required: true, message: '请输入实体名称' }]}>
            <Input placeholder="请输入实体名称" />
          </Form.Item>
          <Form.Item name="type" label="实体类型" rules={[{ required: true, message: '请选择实体类型' }]}>
            <Select placeholder="请选择实体类型">
              {Object.entries(entityTypeLabels).map(([value, label]) => (
                <Option key={value} value={value}>{label}</Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="description" label="实体描述">
            <Input.TextArea rows={3} placeholder="请输入实体描述" />
          </Form.Item>
          <Form.Item name="tags" label="标签">
            <Select mode="tags" placeholder="输入标签" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="添加知识关系"
        open={relationModalVisible}
        onOk={handleAddRelation}
        onCancel={() => setRelationModalVisible(false)}
        width={600}
      >
        <Form form={relationForm} layout="vertical">
          <Form.Item name="from_entity_id" label="源实体" rules={[{ required: true, message: '请选择源实体' }]}>
            <Select placeholder="请选择源实体">
              {entities.map(entity => (
                <Option key={entity.id} value={entity.id}>{entity.name}</Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="to_entity_id" label="目标实体" rules={[{ required: true, message: '请选择目标实体' }]}>
            <Select placeholder="请选择目标实体">
              {entities.map(entity => (
                <Option key={entity.id} value={entity.id}>{entity.name}</Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="type" label="关系类型" rules={[{ required: true, message: '请选择关系类型' }]}>
            <Select placeholder="请选择关系类型">
              <Option value="is_a">是一种</Option>
              <Option value="part_of">属于</Option>
              <Option value="related_to">相关</Option>
              <Option value="competes_with">竞争</Option>
              <Option value="mentions">提及</Option>
              <Option value="depends_on">依赖</Option>
              <Option value="owns">拥有</Option>
              <Option value="created_by">创建者</Option>
            </Select>
          </Form.Item>
          <Form.Item name="weight" label="关系权重">
            <Input type="number" min={0} max={1} step={0.1} placeholder="0-1" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="实体详情"
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
                <strong>标签：</strong>
                {selectedEntity.tags.map(tag => <Tag key={tag}>{tag}</Tag>)}
              </div>
            )}
            <div style={{ marginTop: 16 }}>
              <Button
                danger
                onClick={async () => {
                  await onDeleteEntity(selectedEntity.id);
                  setSelectedEntity(null);
                  message.success('实体删除成功');
                }}
              >
                删除实体
              </Button>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default KnowledgeGraph;
