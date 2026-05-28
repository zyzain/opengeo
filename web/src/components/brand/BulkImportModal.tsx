import React, { useState } from 'react';
import { Modal, Upload, Button, Table, message, Alert, Space } from 'antd';
import { UploadOutlined, DownloadOutlined } from '@ant-design/icons';
import { useIntl } from 'react-intl';

const { Dragger } = Upload;

interface BulkImportModalProps {
  visible: boolean;
  onClose: () => void;
  onImport: (entries: any[]) => Promise<{ imported: number; skipped: number; errors: number }>;
  brandId: number;
}

const BulkImportModal: React.FC<BulkImportModalProps> = ({
  visible,
  onClose,
  onImport,
  brandId,
}) => {
  const intl = useIntl();
  const [fileData, setFileData] = useState<any[]>([]);
  const [importing, setImporting] = useState(false);
  const [result, setResult] = useState<{ imported: number; skipped: number; errors: number } | null>(null);

  const columns = [
    { title: intl.formatMessage({ id: 'glossary.column.term' }), dataIndex: 'term', key: 'term' },
    { title: intl.formatMessage({ id: 'glossary.column.definition' }), dataIndex: 'definition', key: 'definition', ellipsis: true },
    { title: intl.formatMessage({ id: 'glossary.column.category' }), dataIndex: 'category', key: 'category' },
    {
      title: intl.formatMessage({ id: 'glossary.column.aliases' }),
      dataIndex: 'aliases',
      key: 'aliases',
      render: (aliases: string[]) => aliases?.join(', ') || '-',
    },
  ];

  const handleUpload = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const text = e.target?.result as string;
        const lines = text.split('\n').filter(line => line.trim());
        const headers = lines[0].split(',').map(h => h.trim().toLowerCase());

        const entries = lines.slice(1).map((line, index) => {
          const values = line.split(',').map(v => v.trim().replace(/^"|"$/g, ''));
          const entry: any = {};
          headers.forEach((header, i) => {
            if (header === 'aliases') {
              entry[header] = values[i] ? values[i].split(';').map(a => a.trim()) : [];
            } else {
              entry[header] = values[i] || '';
            }
          });
          return entry;
        });

        setFileData(entries);
        message.success(intl.formatMessage({ id: 'bulkImport.parseSuccess' }, { count: entries.length }));
      } catch (error) {
        message.error(intl.formatMessage({ id: 'bulkImport.parseFailed' }));
      }
    };
    reader.readAsText(file);
    return false;
  };

  const handleImport = async () => {
    if (fileData.length === 0) {
      message.warning(intl.formatMessage({ id: 'bulkImport.uploadFirst' }));
      return;
    }

    setImporting(true);
    try {
      const importResult = await onImport(fileData);
      setResult(importResult);
      message.success(intl.formatMessage({ id: 'bulkImport.importComplete' }, { imported: importResult.imported, skipped: importResult.skipped }));
    } catch (error) {
      message.error(intl.formatMessage({ id: 'bulkImport.importFailed' }));
    } finally {
      setImporting(false);
    }
  };

  const handleDownloadTemplate = () => {
    const template = 'term,definition,category,aliases,context,is_forbidden\n' +
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
    <Modal
      title={intl.formatMessage({ id: 'bulkImport.modalTitle' })}
      open={visible}
      onCancel={onClose}
      width={800}
      footer={[
        <Button key="cancel" onClick={onClose}>
          {intl.formatMessage({ id: 'common.action.cancel' })}
        </Button>,
        <Button
          key="template"
          icon={<DownloadOutlined />}
          onClick={handleDownloadTemplate}
        >
          {intl.formatMessage({ id: 'common.action.downloadTemplate' })}
        </Button>,
        <Button
          key="import"
          type="primary"
          loading={importing}
          onClick={handleImport}
          disabled={fileData.length === 0}
        >
          {intl.formatMessage({ id: 'bulkImport.action.startImport' })}
        </Button>,
      ]}
    >
      {result && (
        <Alert
          type="success"
          message={intl.formatMessage({ id: 'bulkImport.result.title' })}
          description={intl.formatMessage({ id: 'bulkImport.result.detail' }, { imported: result.imported, skipped: result.skipped, errors: result.errors })}
          style={{ marginBottom: 16 }}
          closable
          afterClose={() => setResult(null)}
        />
      )}

      <Dragger
        accept=".csv,.tsv"
        beforeUpload={handleUpload}
        showUploadList={false}
        style={{ marginBottom: 16 }}
      >
        <p className="ant-upload-drag-icon">
          <UploadOutlined />
        </p>
        <p className="ant-upload-text">{intl.formatMessage({ id: 'bulkImport.uploadText' })}</p>
        <p className="ant-upload-hint">{intl.formatMessage({ id: 'bulkImport.uploadHint' })}</p>
      </Dragger>

      {fileData.length > 0 && (
        <div>
          <h4>{intl.formatMessage({ id: 'bulkImport.preview' }, { count: fileData.length })}</h4>
          <Table
            columns={columns}
            dataSource={fileData.map((item, index) => ({ ...item, key: index }))}
            size="small"
            pagination={{ pageSize: 5 }}
            scroll={{ y: 200 }}
          />
        </div>
      )}
    </Modal>
  );
};

export default BulkImportModal;
