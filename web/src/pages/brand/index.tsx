import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Tag,
  Modal,
  Form,
  Input,
  Select,
  message,
  Card,
  Row,
  Col,
  Statistic,
  Popconfirm,
  Tooltip,
  Badge,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  SearchOutlined,
  ReloadOutlined,
  TeamOutlined,
  FileTextOutlined,
  BranchesOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useIntl } from 'react-intl';
import { useBrands } from '../../hooks/useBrand';
import BrandForm from '../../components/brand/BrandForm';

const { Option } = Select;

interface Brand {
  id: number;
  name: string;
  slug: string;
  description: string;
  industry: string;
  logo_url: string;
  website: string;
  status: number;
  created_at: string;
  updated_at: string;
}

const BrandListPage: React.FC = () => {
  const navigate = useNavigate();
  const intl = useIntl();
  const { brands, loading, refetch, createBrand, updateBrand, deleteBrand } = useBrands();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingBrand, setEditingBrand] = useState<Brand | null>(null);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [filterIndustry, setFilterIndustry] = useState<string | undefined>();
  const [filterStatus, setFilterStatus] = useState<number | undefined>();

  const statusMap: Record<number, { color: string; text: string }> = {
    1: { color: 'green', text: intl.formatMessage({ id: 'brand.status.active' }) },
    2: { color: 'orange', text: intl.formatMessage({ id: 'brand.status.archived' }) },
    3: { color: 'red', text: intl.formatMessage({ id: 'brand.status.disabled' }) },
  };

  const filteredBrands = brands.filter((brand: Brand) => {
    const matchesKeyword = !searchKeyword || 
      brand.name.toLowerCase().includes(searchKeyword.toLowerCase()) ||
      brand.slug.toLowerCase().includes(searchKeyword.toLowerCase());
    const matchesIndustry = !filterIndustry || brand.industry === filterIndustry;
    const matchesStatus = filterStatus === undefined || brand.status === filterStatus;
    return matchesKeyword && matchesIndustry && matchesStatus;
  });

  const stats = {
    total: brands.length,
    active: brands.filter((b: Brand) => b.status === 1).length,
    archived: brands.filter((b: Brand) => b.status === 2).length,
    disabled: brands.filter((b: Brand) => b.status === 3).length,
  };

  const columns = [
    {
      title: intl.formatMessage({ id: 'brand.column.name' }),
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: Brand) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          {record.logo_url ? (
            <img
              src={record.logo_url}
              alt={record.name}
              style={{ width: 32, height: 32, borderRadius: 4, objectFit: 'cover' }}
            />
          ) : (
            <div
              style={{
                width: 32,
                height: 32,
                borderRadius: 4,
                backgroundColor: '#1890ff',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: '#fff',
                fontWeight: 'bold',
              }}
            >
              {record.name.charAt(0)}
            </div>
          )}
          <a
            href={`/brand/${record.id}`}
            onClick={(e) => {
              e.preventDefault();
              navigate(`/brand/${record.id}`);
            }}
            style={{ fontWeight: 500, color: '#1890ff' }}
          >
            {text}
          </a>
        </div>
      ),
    },
    {
      title: intl.formatMessage({ id: 'brand.column.slug' }),
      dataIndex: 'slug',
      key: 'slug',
      render: (text: string) => (
        <Tag>{text}</Tag>
      ),
    },
    {
      title: intl.formatMessage({ id: 'brand.column.industry' }),
      dataIndex: 'industry',
      key: 'industry',
      render: (text: string) => text || '-',
    },
    {
      title: intl.formatMessage({ id: 'brand.column.description' }),
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      width: 200,
      render: (text: string) => text || '-',
    },
    {
      title: intl.formatMessage({ id: 'brand.column.status' }),
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => {
        const s = statusMap[status] || { color: 'default', text: intl.formatMessage({ id: 'common.status.unknown' }) };
        return <Tag color={s.color}>{s.text}</Tag>;
      },
    },
    {
      title: intl.formatMessage({ id: 'brand.column.createdAt' }),
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => new Date(text).toLocaleDateString(),
    },
    {
      title: intl.formatMessage({ id: 'common.column.action' }),
      key: 'action',
      render: (_: any, record: Brand) => (
        <Space size="middle">
          <Tooltip title={intl.formatMessage({ id: 'common.action.view' })}>
            <Button
              type="link"
              icon={<EyeOutlined />}
              onClick={() => navigate(`/brand/${record.id}`)}
            />
          </Tooltip>
          <Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}>
            <Button
              type="link"
              icon={<EditOutlined />}
              onClick={() => editBrand(record)}
            />
          </Tooltip>
          <Popconfirm
            title={intl.formatMessage({ id: 'brand.confirm.deleteTitle' })}
            description={intl.formatMessage({ id: 'brand.confirm.deleteDesc' }, { name: record.name })}
            onConfirm={() => deleteBrandHandler(record)}
            okText={intl.formatMessage({ id: 'brand.action.delete' })}
            cancelText={intl.formatMessage({ id: 'brand.action.cancel' })}
            okButtonProps={{ danger: true }}
          >
            <Tooltip title={intl.formatMessage({ id: 'brand.action.delete' })}>
              <Button type="link" danger icon={<DeleteOutlined />} />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const editBrand = (brand: Brand) => {
    setEditingBrand(brand);
    setModalVisible(true);
  };

  const deleteBrandHandler = async (brand: Brand) => {
    try {
      await deleteBrand(brand.id);
      message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
    } catch {
      message.error(intl.formatMessage({ id: 'common.message.deleteFailed' }));
    }
  };

  const handleCreate = () => {
    setEditingBrand(null);
    setModalVisible(true);
  };

  const handleSubmit = async (values: any) => {
    try {
      if (editingBrand) {
        await updateBrand(editingBrand.id, values);
        message.success(intl.formatMessage({ id: 'common.message.updateSuccess' }));
      } else {
        await createBrand(values);
        message.success(intl.formatMessage({ id: 'common.message.createSuccess' }));
      }
      setModalVisible(false);
    } catch {
      message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
    }
  };

  return (
    <div style={{ padding: 24 }}>
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 24, fontWeight: 600 }}>{intl.formatMessage({ id: 'brand.page.title' })}</h1>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={refetch}>
            {intl.formatMessage({ id: 'common.action.refresh' })}
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            {intl.formatMessage({ id: 'brand.action.create' })}
          </Button>
        </Space>
      </div>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'brand.stat.total' })}
              value={stats.total}
              prefix={<TeamOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'brand.stat.active' })}
              value={stats.active}
              valueStyle={{ color: '#3f8600' }}
              prefix={<FileTextOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'brand.stat.archived' })}
              value={stats.archived}
              valueStyle={{ color: '#fa8c16' }}
              prefix={<BranchesOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'brand.stat.disabled' })}
              value={stats.disabled}
              valueStyle={{ color: '#cf1322' }}
              prefix={<TeamOutlined />}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', gap: 16, flexWrap: 'wrap' }}>
          <Input
            placeholder={intl.formatMessage({ id: 'brand.placeholder.search' })}
            prefix={<SearchOutlined />}
            value={searchKeyword}
            onChange={(e) => setSearchKeyword(e.target.value)}
            style={{ width: 250 }}
            allowClear
          />
          <Select
            placeholder={intl.formatMessage({ id: 'brand.placeholder.selectIndustry' })}
            value={filterIndustry}
            onChange={setFilterIndustry}
            style={{ width: 150 }}
            allowClear
          >
            <Option value="科技">{intl.formatMessage({ id: 'industry.tech' })}</Option>
            <Option value="金融">{intl.formatMessage({ id: 'industry.finance' })}</Option>
            <Option value="医疗">{intl.formatMessage({ id: 'industry.medical' })}</Option>
            <Option value="教育">{intl.formatMessage({ id: 'industry.education' })}</Option>
            <Option value="电商">{intl.formatMessage({ id: 'industry.ecommerce' })}</Option>
            <Option value="其他">{intl.formatMessage({ id: 'industry.other' })}</Option>
          </Select>
          <Select
            placeholder={intl.formatMessage({ id: 'brand.placeholder.selectStatus' })}
            value={filterStatus}
            onChange={setFilterStatus}
            style={{ width: 150 }}
            allowClear
          >
            <Option value={1}>{intl.formatMessage({ id: 'brand.status.active' })}</Option>
            <Option value={2}>{intl.formatMessage({ id: 'brand.status.archived' })}</Option>
            <Option value={3}>{intl.formatMessage({ id: 'brand.status.disabled' })}</Option>
          </Select>
        </div>
      </Card>

      <Card>
        <Table
          columns={columns}
          dataSource={filteredBrands}
          loading={loading}
          rowKey="id"
          pagination={{
            pageSize: 20,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => intl.formatMessage({ id: 'common.pagination.total' }, { total }),
          }}
        />
      </Card>

      <Modal
        title={editingBrand ? intl.formatMessage({ id: 'brand.modal.edit' }) : intl.formatMessage({ id: 'brand.modal.create' })}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={800}
        destroyOnClose
      >
        <BrandForm
          initialValues={editingBrand}
          onSubmit={handleSubmit}
          onCancel={() => setModalVisible(false)}
          loading={loading}
        />
      </Modal>
    </div>
  );
};

export default BrandListPage;
