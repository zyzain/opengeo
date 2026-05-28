import React, { useState } from 'react';
import { Modal, Upload, Button, Table, message, Alert, Space } from 'antd';
import { UploadOutlined, DownloadOutlined } from '@ant-design/icons';

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
  const [fileData, setFileData] = useState<any[]>([]);
  const [importing, setImporting] = useState(false);
  const [result, setResult] = useState<{ imported: number; skipped: number; errors: number } | null>(null);

  const columns = [
    { title: '术语', dataIndex: 'term', key: 'term' },
    { title: '定义', dataIndex: 'definition', key: 'definition', ellipsis: true },
    { title: '分类', dataIndex: 'category', key: 'category' },
    {
      title: '别名',
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
        message.success(`成功解析 ${entries.length} 条记录`);
      } catch (error) {
        message.error('文件解析失败，请检查文件格式');
      }
    };
    reader.readAsText(file);
    return false;
  };

  const handleImport = async () => {
    if (fileData.length === 0) {
      message.warning('请先上传文件');
      return;
    }

    setImporting(true);
    try {
      const importResult = await onImport(fileData);
      setResult(importResult);
      message.success(`导入完成：成功 ${importResult.imported} 条，跳过 ${importResult.skipped} 条`);
    } catch (error) {
      message.error('导入失败');
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
      title="批量导入术语"
      open={visible}
      onCancel={onClose}
      width={800}
      footer={[
        <Button key="cancel" onClick={onClose}>
          取消
        </Button>,
        <Button
          key="template"
          icon={<DownloadOutlined />}
          onClick={handleDownloadTemplate}
        >
          下载模板
        </Button>,
        <Button
          key="import"
          type="primary"
          loading={importing}
          onClick={handleImport}
          disabled={fileData.length === 0}
        >
          开始导入
        </Button>,
      ]}
    >
      {result && (
        <Alert
          type="success"
          message="导入结果"
          description={`成功导入 ${result.imported} 条，跳过 ${result.skipped} 条，失败 ${result.errors} 条`}
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
        <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
        <p className="ant-upload-hint">支持 CSV、TSV 格式文件</p>
      </Dragger>

      {fileData.length > 0 && (
        <div>
          <h4>预览数据（{fileData.length} 条）</h4>
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
