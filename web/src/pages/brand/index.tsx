import React, { useState } from 'react';
import { Table, Button, Space, Tag, Modal, Form, Input, Select, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { useBrands } from '../../hooks/useBrand';

const { Option } = Select;

interface Brand {
  id: number;
  name: string;
  slug: string;
  description: string;
  industry: string;
  status: number;
  created_at: string;
}

const BrandListPage: React.FC = () => {
  const { brands, loading, refetch } = useBrands();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingBrand, setEditingBrand] = useState<Brand | null>(null);
  const [form] = Form.useForm();

  const columns = [
    {
      title: '品牌名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: Brand) => (
        <a href={`/brand/${record.id}`}>{text}</a>
      ),
    },
    {
      title: '品牌标识',
      dataIndex: 'slug',
      key: 'slug',
    },
    {
      title: '行业',
      dataIndex: 'industry',
      key: 'industry',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => {
        const statusMap: Record<number, { color: string; text: string }> = {
          1: { color: 'green', text: '活跃' },
          2: { color: 'orange', text: '已归档' },
          3: { color: 'red', text: '已禁用' },
        };
        const s = statusMap[status] || { color: 'default', text: '未知' };
        return <Tag color={s.color}>{s.text}</Tag>;
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => new Date(text).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Brand) => (
        <Space size="middle">
          <Button type="link" icon={<EyeOutlined />} onClick={() => viewBrand(record)}>
            查看
          </Button>
          <Button type="link" icon={<EditOutlined />} onClick={() => editBrand(record)}>
            编辑
          </Button>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => deleteBrand(record)}>
            删除
          </Button>
        </Space>
      ),
    },
  ];

  const viewBrand = (brand: Brand) => {
    window.location.href = `/brand/${brand.id}`;
  };

  const editBrand = (brand: Brand) => {
    setEditingBrand(brand);
    form.setFieldsValue(brand);
    setModalVisible(true);
  };

  const deleteBrand = async (brand: Brand) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除品牌 "${brand.name}" 吗？`,
      onOk: async () => {
        try {
          await fetch(`/api/v1/brand/${brand.id}`, { method: 'DELETE' });
          message.success('删除成功');
          refetch();
        } catch {
          message.error('删除失败');
        }
      },
    });
  };

  const handleCreate = () => {
    setEditingBrand(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const url = editingBrand ? `/api/v1/brand/${editingBrand.id}` : '/api/v1/brands';
      const method = editingBrand ? 'PUT' : 'POST';

      await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });

      message.success(editingBrand ? '更新成功' : '创建成功');
      setModalVisible(false);
      refetch();
    } catch {
      message.error('操作失败');
    }
  };

  return (
    <div style={{ padding: 24 }}>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h1>品牌管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          创建品牌
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={brands}
        loading={loading}
        rowKey="id"
        pagination={{ pageSize: 20 }}
      />

      <Modal
        title={editingBrand ? '编辑品牌' : '创建品牌'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="品牌名称" rules={[{ required: true, message: '请输入品牌名称' }]}>
            <Input placeholder="请输入品牌名称" />
          </Form.Item>
          <Form.Item name="slug" label="品牌标识" rules={[{ required: true, message: '请输入品牌标识' }]}>
            <Input placeholder="请输入品牌标识（URL友好）" disabled={!!editingBrand} />
          </Form.Item>
          <Form.Item name="description" label="品牌描述">
            <Input.TextArea rows={3} placeholder="请输入品牌描述" />
          </Form.Item>
          <Form.Item name="industry" label="所属行业">
            <Select placeholder="请选择行业">
              <Option value="科技">科技</Option>
              <Option value="金融">金融</Option>
              <Option value="医疗">医疗</Option>
              <Option value="教育">教育</Option>
              <Option value="电商">电商</Option>
              <Option value="其他">其他</Option>
            </Select>
          </Form.Item>
          <Form.Item name="website" label="品牌官网">
            <Input placeholder="请输入品牌官网" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default BrandListPage;
