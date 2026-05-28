import React, { useState } from 'react';
import { Card, Button, Table, Modal, Form, Input, Tag, Timeline, Descriptions, message, Space } from 'antd';
import { PlusOutlined, HistoryOutlined, DiffOutlined } from '@ant-design/icons';
import { useIntl } from 'react-intl';

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
  const intl = useIntl();
  const [modalVisible, setModalVisible] = useState(false);
  const [compareModalVisible, setCompareModalVisible] = useState(false);
  const [selectedSnapshots, setSelectedSnapshots] = useState<number[]>([]);
  const [form] = Form.useForm();

  const columns = [
    {
      title: intl.formatMessage({ id: 'snapshot.column.version' }),
      dataIndex: 'version',
      key: 'version',
      render: (text: string) => <Tag color="blue">{text}</Tag>,
    },
    {
      title: intl.formatMessage({ id: 'snapshot.column.changeLog' }),
      dataIndex: 'change_log',
      key: 'change_log',
      ellipsis: true,
    },
    {
      title: intl.formatMessage({ id: 'snapshot.column.createdAt' }),
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: intl.formatMessage({ id: 'snapshot.column.action' }),
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
          {selectedSnapshots.includes(record.id) ? intl.formatMessage({ id: 'snapshot.action.deselect' }) : intl.formatMessage({ id: 'snapshot.action.selectCompare' })}
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
      message.success(intl.formatMessage({ id: 'snapshot.message.createSuccess' }));
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  const handleCompare = async () => {
    if (selectedSnapshots.length !== 2) {
      message.warning(intl.formatMessage({ id: 'snapshot.message.selectTwo' }));
      return;
    }
    await onCompareSnapshots(selectedSnapshots[0], selectedSnapshots[1]);
    setCompareModalVisible(true);
  };

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h3>{intl.formatMessage({ id: 'snapshot.title' })}</h3>
        <Space>
          <Button
            icon={<DiffOutlined />}
            onClick={handleCompare}
            disabled={selectedSnapshots.length !== 2}
          >
            {intl.formatMessage({ id: 'snapshot.action.compare' })}
          </Button>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              form.resetFields();
              setModalVisible(true);
            }}
          >
            {intl.formatMessage({ id: 'snapshot.action.create' })}
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
        title={intl.formatMessage({ id: 'snapshot.modal.createTitle' })}
        open={modalVisible}
        onOk={handleCreate}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="version"
            label={intl.formatMessage({ id: 'snapshot.form.version' })}
            rules={[{ required: true, message: intl.formatMessage({ id: 'snapshot.validation.enterVersion' }) }]}
            extra={intl.formatMessage({ id: 'snapshot.form.versionExample' })}
          >
            <Input placeholder={intl.formatMessage({ id: 'snapshot.placeholder.enterVersion' })} />
          </Form.Item>
          <Form.Item
            name="change_log"
            label={intl.formatMessage({ id: 'snapshot.form.changeLog' })}
            rules={[{ required: true, message: intl.formatMessage({ id: 'snapshot.validation.enterChangeLog' }) }]}
          >
            <TextArea rows={4} placeholder={intl.formatMessage({ id: 'snapshot.placeholder.changeLog' })} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={intl.formatMessage({ id: 'snapshot.modal.compareTitle' })}
        open={compareModalVisible}
        onCancel={() => setCompareModalVisible(false)}
        footer={null}
        width={800}
      >
        <div>
          <p>{intl.formatMessage({ id: 'snapshot.message.compareInDev' })}</p>
          <p>{intl.formatMessage({ id: 'snapshot.message.selectedSnapshots' })}{selectedSnapshots.join(' vs ')}</p>
        </div>
      </Modal>
    </div>
  );
};

export default BrandSnapshot;
