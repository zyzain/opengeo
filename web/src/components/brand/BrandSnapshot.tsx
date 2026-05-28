import React, { useState } from 'react';
import { Card, Button, Table, Modal, Form, Input, Tag, Timeline, Descriptions, message } from 'antd';
import { PlusOutlined, HistoryOutlined, DiffOutlined } from '@ant-design/icons';

const { TextArea } = Input;

interface Snapshot {
  id: number;
  brand_id: number;
  version: string;
  snapshot_data: string;
  change_log: string;
  created_by: number;
  created_at: string;
}

interface BrandSnapshotProps {
  brandId: number;
  snapshots: Snapshot[];
  loading?: boolean;
  onCreateSnapshot: (values: any) => Promise<void>;
  onCompareSnapshots: (id1: number, id2: number) => Promise<void>;
}

const BrandSnapshot: React.FC<BrandSnapshotProps> = ({
  brandId,
  snapshots,
  loading,
  onCreateSnapshot,
  onCompareSnapshots,
}) => {
  const [modalVisible, setModalVisible] = useState(false);
  const [compareModalVisible, setCompareModalVisible] = useState(false);
  const [selectedSnapshots, setSelectedSnapshots] = useState<number[]>([]);
  const [form] = Form.useForm();

  const columns = [
    {
      title: '版本',
      dataIndex: 'version',
      key: 'version',
      render: (text: string) => <Tag color="blue">{text}</Tag>,
    },
    {
      title: '变更说明',
      dataIndex: 'change_log',
      key: 'change_log',
      ellipsis: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Snapshot) => (
        <Button
          type="link"
          onClick={() => {
            if (selectedSnapshots.includes(record.id)) {
              setSelectedSnapshots(selectedSnapshots.filter(id => id !== record.id));
            } else if (selectedSnapshots.length < 2) {
              setSelectedSnapshots([...selectedSnapshots, record.id]);
            }
          }}
        >
          {selectedSnapshots.includes(record.id) ? '取消选择' : '选择对比'}
        </Button>
      ),
    },
  ];

  const handleCreate = async () => {
    try {
      const values = await form.validateFields();
      await onCreateSnapshot({ ...values, brand_id: brandId });
      setModalVisible(false);
      form.resetFields();
      message.success('快照创建成功');
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  const handleCompare = async () => {
    if (selectedSnapshots.length !== 2) {
      message.warning('请选择两个快照进行对比');
      return;
    }
    await onCompareSnapshots(selectedSnapshots[0], selectedSnapshots[1]);
    setCompareModalVisible(true);
  };

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h3>品牌快照</h3>
        <Space>
          <Button
            icon={<DiffOutlined />}
            onClick={handleCompare}
            disabled={selectedSnapshots.length !== 2}
          >
            对比快照
          </Button>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              form.resetFields();
              setModalVisible(true);
            }}
          >
            创建快照
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={snapshots}
        loading={loading}
        rowKey="id"
        pagination={{ pageSize: 10 }}
        rowSelection={{
          selectedRowKeys: selectedSnapshots,
          onChange: (keys) => setSelectedSnapshots(keys as number[]),
        }}
      />

      <Modal
        title="创建品牌快照"
        open={modalVisible}
        onOk={handleCreate}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="version"
            label="版本号"
            rules={[{ required: true, message: '请输入版本号' }]}
            extra="例如：v1.0.0"
          >
            <Input placeholder="请输入版本号" />
          </Form.Item>
          <Form.Item
            name="change_log"
            label="变更说明"
            rules={[{ required: true, message: '请输入变更说明' }]}
          >
            <TextArea rows={4} placeholder="请描述本次变更内容" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="快照对比"
        open={compareModalVisible}
        onCancel={() => setCompareModalVisible(false)}
        footer={null}
        width={800}
      >
        <div>
          <p>快照对比功能开发中...</p>
          <p>已选择快照：{selectedSnapshots.join(' vs ')}</p>
        </div>
      </Modal>
    </div>
  );
};

export default BrandSnapshot;
