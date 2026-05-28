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
  Switch,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Statistic,
  Upload,
  Tooltip,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  UploadOutlined,
  DownloadOutlined,
  SearchOutlined,
  ReloadOutlined,
  ExclamationCircleOutlined,
  CheckCircleOutlined,
} from '@ant-design/icons';
import { useIntl } from 'react-intl';

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
  created_at: string;
  updated_at: string;
}

interface GlossaryManagerProps {
  brandId: number;
  entries: GlossaryEntry[];
  loading: boolean;
  onRefresh?: () => void;
}

const GlossaryManager: React.FC<GlossaryManagerProps> = ({
  brandId,
  entries,
  loading,
  onRefresh,
}) => {
  const intl = useIntl();
  const [modalVisible, setModalVisible] = useState(false);
  const [bulkImportModalVisible, setBulkImportModalVisible] = useState(false);
  const [editingEntry, setEditingEntry] = useState<GlossaryEntry | null>(null);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [filterCategory, setFilterCategory] = useState<string | undefined>();
  const [filterForbidden, setFilterForbidden] = useState<boolean | undefined>();
  const [form] = Form.useForm();

  const categoryMap: Record<string, { color: string; text: string }> = {
    product: { color: 'blue', text: intl.formatMessage({ id: 'glossary.category.product' }) },
    technology: { color: 'purple', text: intl.formatMessage({ id: 'glossary.category.technology' }) },
    concept: { color: 'cyan', text: intl.formatMessage({ id: 'glossary.category.concept' }) },
    person: { color: 'green', text: intl.formatMessage({ id: 'glossary.category.person' }) },
    place: { color: 'orange', text: intl.formatMessage({ id: 'glossary.category.place' }) },
  };

  const filteredEntries = entries.filter((entry) => {
    const matchesKeyword =
      !searchKeyword ||
      entry.term.toLowerCase().includes(searchKeyword.toLowerCase()) ||
      entry.definition.toLowerCase().includes(searchKeyword.toLowerCase());
    const matchesCategory = !filterCategory || entry.category === filterCategory;
    const matchesForbidden =
      filterForbidden === undefined || entry.is_forbidden === filterForbidden;
    return matchesKeyword && matchesCategory && matchesForbidden;
  });

  const stats = {
    total: entries.length,
    forbidden: entries.filter((e) => e.is_forbidden).length,
    preferred: entries.filter((e) => e.is_preferred).length,
    categories: new Set(entries.map((e) => e.category)).size,
  };

  const columns = [
    {
      title: intl.formatMessage({ id: 'glossary.column.term' }),
      dataIndex: 'term',
      key: 'term',
      render: (text: string, record: GlossaryEntry) => (
        <Space>
          <span style={{ fontWeight: record.is_preferred ? 'bold' : 'normal' }}>
            {text}
          </span>
          {record.is_forbidden && (
            <Tag color="red" icon={<ExclamationCircleOutlined />}>
              {intl.formatMessage({ id: 'glossary.tag.forbidden' })}
            </Tag>
          )}
          {record.is_preferred && (
            <Tag color="green" icon={<CheckCircleOutlined />}>
              {intl.formatMessage({ id: 'glossary.tag.preferred' })}
            </Tag>
          )}
        </Space>
      ),
    },
    {
      title: intl.formatMessage({ id: 'glossary.column.definition' }),
      dataIndex: 'definition',
      key: 'definition',
      ellipsis: true,
      width: 300,
    },
    {
      title: intl.formatMessage({ id: 'glossary.column.category' }),
      dataIndex: 'category',
      key: 'category',
      render: (category: string) => {
        const cat = categoryMap[category];
        return cat ? <Tag color={cat.color}>{cat.text}</Tag> : <Tag>{category}</Tag>;
      },
    },
    {
      title: intl.formatMessage({ id: 'glossary.column.aliases' }),
      dataIndex: 'aliases',
      key: 'aliases',
      render: (aliases: string[]) =>
        aliases && aliases.length > 0
          ? aliases.map((a) => <Tag key={a}>{a}</Tag>)
          : '-',
    },
    {
      title: intl.formatMessage({ id: 'common.column.action' }),
      key: 'action',
      render: (_: any, record: GlossaryEntry) => (
        <Space size="middle">
          <Tooltip title={intl.formatMessage({ id: 'common.action.edit' })}>
            <Button type="link" icon={<EditOutlined />} onClick={() => editEntry(record)} />
          </Tooltip>
          <Popconfirm
            title={intl.formatMessage({ id: 'common.confirmDelete' })}
            description={intl.formatMessage({ id: 'glossary.confirmDeleteTerm' }, { term: record.term })}
            onConfirm={() => deleteEntry(record)}
            okText={intl.formatMessage({ id: 'common.action.delete' })}
            cancelText={intl.formatMessage({ id: 'common.action.cancel' })}
            okButtonProps={{ danger: true }}
          >
            <Tooltip title={intl.formatMessage({ id: 'common.action.delete' })}>
              <Button type="link" danger icon={<DeleteOutlined />} />
            </Tooltip>
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
      await fetch(`/api/v1/brand/${brandId}/glossary/${entry.id}`, {
        method: 'DELETE',
      });
      message.success(intl.formatMessage({ id: 'common.message.deleteSuccess' }));
      onRefresh?.();
    } catch {
      message.error(intl.formatMessage({ id: 'common.message.deleteFailed' }));
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
        aliases: values.aliases
          ? values.aliases.split(',').map((s: string) => s.trim()).filter(Boolean)
          : [],
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

      message.success(editingEntry ? intl.formatMessage({ id: 'common.message.updateSuccess' }) : intl.formatMessage({ id: 'common.message.createSuccess' }));
      setModalVisible(false);
      onRefresh?.();
    } catch {
      message.error(intl.formatMessage({ id: 'common.message.operationFailed' }));
    }
  };

  const handleBulkImport = async (file: File) => {
    const reader = new FileReader();
    reader.onload = async (e) => {
      try {
        const text = e.target?.result as string;
        const lines = text.split('\n').filter((line) => line.trim());
        const headers = lines[0].split(',').map((h) => h.trim().toLowerCase());

        const entries = lines.slice(1).map((line) => {
          const values = line.split(',').map((v) => v.trim().replace(/^"|"$/g, ''));
          const entry: any = {};
          headers.forEach((header, i) => {
            if (header === 'aliases') {
              entry[header] = values[i]
                ? values[i].split(';').map((a) => a.trim())
                : [];
            } else if (header === 'is_forbidden') {
              entry[header] = values[i]?.toLowerCase() === 'true';
            } else {
              entry[header] = values[i] || '';
            }
          });
          return entry;
        });

        const response = await fetch(`/api/v1/brand/${brandId}/glossary/bulk-import`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ entries, overwrite_existing: false }),
        });

        if (response.ok) {
          const result = await response.json();
          message.success(
            intl.formatMessage({ id: 'glossary.importComplete' }, { imported: result.imported_count, skipped: result.skipped_count })
          );
          onRefresh?.();
        } else {
          message.error(intl.formatMessage({ id: 'glossary.importFailed' }));
        }
      } catch (error) {
        message.error(intl.formatMessage({ id: 'common.message.fileParseFailed' }));
      }
    };
    reader.readAsText(file);
    return false;
  };

  const handleExport = async () => {
    try {
      const response = await fetch(`/api/v1/brand/${brandId}/glossary/export?format=csv`);
      if (response.ok) {
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `glossary_${brandId}.csv`;
        a.click();
        URL.revokeObjectURL(url);
        message.success(intl.formatMessage({ id: 'common.message.exportSuccess' }));
      }
    } catch {
      message.error(intl.formatMessage({ id: 'common.message.exportFailed' }));
    }
  };

  const handleDownloadTemplate = () => {
    const template =
      'term,definition,category,aliases,context,is_forbidden\n' +
      'GEO,Generative Engine Optimization,technology,"生成式引擎优化;AI搜索优化",针对AI搜索引擎的内容优化策略,false\n' +
      'SEO,Search Engine Optimization,technology,"搜索引擎优化",针对传统搜索引擎的优化策略,false';

    const blob = new Blob([template], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'glossary_template.csv';
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div>
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic title={intl.formatMessage({ id: 'glossary.stat.totalTerms' })} value={stats.total} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'glossary.stat.forbidden' })}
              value={stats.forbidden}
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={intl.formatMessage({ id: 'glossary.stat.preferredTerms' })}
              value={stats.preferred}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title={intl.formatMessage({ id: 'glossary.stat.categories' })} value={stats.categories} />
          </Card>
        </Col>
      </Row>

      <Card>
        <div
          style={{
            marginBottom: 16,
            display: 'flex',
            justifyContent: 'space-between',
            flexWrap: 'wrap',
            gap: 16,
          }}
        >
          <Space wrap>
            <Input
              placeholder={intl.formatMessage({ id: 'glossary.placeholder.searchTerm' })}
              prefix={<SearchOutlined />}
              value={searchKeyword}
              onChange={(e) => setSearchKeyword(e.target.value)}
              style={{ width: 250 }}
              allowClear
            />
            <Select
              placeholder={intl.formatMessage({ id: 'glossary.placeholder.selectCategory' })}
              value={filterCategory}
              onChange={setFilterCategory}
              style={{ width: 150 }}
              allowClear
            >
              {Object.entries(categoryMap).map(([value, { text }]) => (
                <Option key={value} value={value}>
                  {text}
                </Option>
              ))}
            </Select>
            <Select
              placeholder={intl.formatMessage({ id: 'glossary.placeholder.filterForbidden' })}
              value={filterForbidden}
              onChange={setFilterForbidden}
              style={{ width: 150 }}
              allowClear
            >
              <Option value={true}>{intl.formatMessage({ id: 'glossary.option.forbidden' })}</Option>
              <Option value={false}>{intl.formatMessage({ id: 'glossary.option.normalTerm' })}</Option>
            </Select>
          </Space>
          <Space>
            <Button icon={<ReloadOutlined />} onClick={onRefresh}>
              {intl.formatMessage({ id: 'common.action.refresh' })}
            </Button>
            <Button icon={<DownloadOutlined />} onClick={handleDownloadTemplate}>
              {intl.formatMessage({ id: 'common.action.downloadTemplate' })}
            </Button>
            <Upload
              accept=".csv"
              showUploadList={false}
              beforeUpload={handleBulkImport}
            >
              <Button icon={<UploadOutlined />}>{intl.formatMessage({ id: 'common.action.bulkImport' })}</Button>
            </Upload>
            <Button icon={<DownloadOutlined />} onClick={handleExport}>
              {intl.formatMessage({ id: 'common.action.export' })}
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              {intl.formatMessage({ id: 'glossary.action.addTerm' })}
            </Button>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={filteredEntries}
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
        title={editingEntry ? intl.formatMessage({ id: 'glossary.modal.editTerm' }) : intl.formatMessage({ id: 'glossary.modal.addTerm' })}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="term"
            label={intl.formatMessage({ id: 'glossary.form.term' })}
            rules={[{ required: true, message: intl.formatMessage({ id: 'glossary.validation.enterTerm' }) }]}
          >
            <Input placeholder={intl.formatMessage({ id: 'glossary.placeholder.enterTerm' })} />
          </Form.Item>
          <Form.Item
            name="definition"
            label={intl.formatMessage({ id: 'glossary.form.definition' })}
            rules={[{ required: true, message: intl.formatMessage({ id: 'glossary.validation.enterDefinition' }) }]}
          >
            <TextArea rows={3} placeholder={intl.formatMessage({ id: 'glossary.placeholder.enterDefinition' })} />
          </Form.Item>
          <Form.Item name="category" label={intl.formatMessage({ id: 'glossary.form.category' })}>
            <Select placeholder={intl.formatMessage({ id: 'glossary.placeholder.selectCategory' })}>
              {Object.entries(categoryMap).map(([value, { text }]) => (
                <Option key={value} value={value}>
                  {text}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="aliases" label={intl.formatMessage({ id: 'glossary.form.aliases' })} help={intl.formatMessage({ id: 'glossary.help.aliasesHint' })}>
            <Input placeholder={intl.formatMessage({ id: 'glossary.placeholder.enterAliases' })} />
          </Form.Item>
          <Form.Item name="context" label={intl.formatMessage({ id: 'glossary.form.context' })}>
            <TextArea rows={2} placeholder={intl.formatMessage({ id: 'glossary.placeholder.enterContext' })} />
          </Form.Item>
          <Form.Item name="is_forbidden" label={intl.formatMessage({ id: 'glossary.form.forbidden' })} valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="is_preferred" label={intl.formatMessage({ id: 'glossary.form.preferred' })} valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default GlossaryManager;
