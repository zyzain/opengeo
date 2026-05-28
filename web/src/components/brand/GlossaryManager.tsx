import React, { useState } from 'react';
import { Table, Button, Space, Tag, Modal, Form, Input, Select, Switch, message, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, UploadOutlined } from '@ant-design/icons';

const { Option } = Select;
const { TextArea } = Input;

interface GlossaryEntry {
  id: number;
  brand_id: number;
  term: string;
  definition: string;
  category: string;
  aliases: string[];
  context: string;
  is_forbidden: boolean;
  is_preferred: boolean;
}

interface GlossaryManagerProps {
  brandId: number;
  entries: GlossaryEntry[];
  loading: boolean;
}

const GlossaryManager: React.FC<GlossaryManagerProps> = ({ brandId, entries, loading }) => {
  const [modalVisible, setModalVisible] = useState(false);
  const [editingEntry, setEditingEntry] = useState<GlossaryEntry | null>(null);
  const [form] = Form.useForm();

  const columns = [
    {
      title: '术语',
      dataIndex: 'term',
      key: 'term',
      render: (text: string, record: GlossaryEntry) => (
        <Space>
          <span style={{ fontWeight: record.is_preferred ? 'bold' : 'normal' }}>{text}</span>
          {record.is_forbidden && <Tag color="red">禁用词</Tag>}
          {record.is_preferred && <Tag color="green">首选</Tag>}
        </Space>
      ),
    },
    {
      title: '定义',
      dataIndex: 'definition',
      key: 'definition',
      ellipsis: true,
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      render: (category: string) => {
        const categoryMap: Record<string, string> = {
          product: '产品',
          technology: '技术',
          concept: '概念',
          person: '人物',
          place: '地点',
        };
        return <Tag>{categoryMap[category] || category}</Tag>;
      },
    },
    {
      title: '别名',
      dataIndex: 'aliases',
      key: 'aliases',
      render: (aliases: string[]) => aliases?.map(a => <Tag key={a}>{a}</Tag>),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: GlossaryEntry) => (
        <Space size="middle">
          <Button type="link" icon={<EditOutlined />} onClick={() => editEntry(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除？" onConfirm={() => deleteEntry(record)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const editEntry = (entry: GlossaryEntry) => {
    setEditingEntry(entry);
    form.setFieldsValue({
      ...entry,
      aliases: entry.aliases?.join(', ') || '',
    });
    setModalVisible(true);
  };

  const deleteEntry = async (entry: GlossaryEntry) => {
    try {
      await fetch(`/api/v1/brand/${brandId}/glossary/${entry.id}`, { method: 'DELETE' });
      message.success('删除成功');
    } catch {
      message.error('删除失败');
    }
  };

  const handleCreate = () => {
    setEditingEntry(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const data = {
        ...values,
        aliases: values.aliases ? values.aliases.split(',').map((s: string) => s.trim()) : [],
      };

      const url = editingEntry
        ? `/api/v1/brand/${brandId}/glossary/${editingEntry.id}`
        : `/api/v1/brand/${brandId}/glossary`;
      const method = editingEntry ? 'PUT' : 'POST';

      await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      message.success(editingEntry ? '更新成功' : '创建成功');
      setModalVisible(false);
    } catch {
      message.error('操作失败');
    }
  };

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h3>品牌术语表</h3>
        <Space>
          <Button icon={<UploadOutlined />}>批量导入</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            添加术语
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={entries}
        loading={loading}
        rowKey="id"
        pagination={{ pageSize: 20 }}
      />

      <Modal
        title={editingEntry ? '编辑术语' : '添加术语'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="term" label="术语" rules={[{ required: true, message: '请输入术语' }]}>
            <Input placeholder="请输入术语" />
          </Form.Item>
          <Form.Item name="definition" label="定义" rules={[{ required: true, message: '请输入定义' }]}>
            <TextArea rows={3} placeholder="请输入术语定义" />
          </Form.Item>
          <Form.Item name="category" label="分类">
            <Select placeholder="请选择分类">
              <Option value="product">产品</Option>
              <Option value="technology">技术</Option>
              <Option value="concept">概念</Option>
              <Option value="person">人物</Option>
              <Option value="place">地点</Option>
            </Select>
          </Form.Item>
          <Form.Item name="aliases" label="别名" help="多个别名用逗号分隔">
            <Input placeholder="请输入别名，多个用逗号分隔" />
          </Form.Item>
          <Form.Item name="context" label="使用上下文">
            <TextArea rows={2} placeholder="请输入使用上下文示例" />
          </Form.Item>
          <Form.Item name="is_forbidden" label="禁用词" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="is_preferred" label="首选术语" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default GlossaryManager;
