import React, { useState } from 'react';
import { Card, Select, Spin, Empty, message } from 'antd';
import { useIntl } from 'react-intl';
import { useBrands, useBrandMetadata } from '../../hooks/useBrand';
import MetadataEditor from '../../components/brand/MetadataEditor';

const { Option } = Select;

const BrandMetadataPage: React.FC = () => {
  const intl = useIntl();
  const { brands, loading: brandsLoading } = useBrands();
  const [selectedBrandId, setSelectedBrandId] = useState<number | undefined>();
  const { metadata, loading: metadataLoading, refetch } = useBrandMetadata(selectedBrandId || 0);
  const [saving, setSaving] = useState(false);

  const handleSave = async (values: any) => {
    if (!selectedBrandId) return;
    setSaving(true);
    try {
      await fetch(`/api/v1/brand/${selectedBrandId}/metadata`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });
      message.success(intl.formatMessage({ id: 'brand.message.metadataSaveSuccess' }));
      refetch();
    } catch (error) {
      message.error(intl.formatMessage({ id: 'brand.message.metadataSaveFailed' }));
    } finally {
      setSaving(false);
    }
  };

  return (
    <div style={{ padding: 24 }}>
      <Card title={intl.formatMessage({ id: 'brand.metadata.title' })}>
        <div style={{ marginBottom: 24 }}>
          <span style={{ marginRight: 16 }}>{intl.formatMessage({ id: 'brand.metadata.selectBrand' })}</span>
          <Select
            placeholder={intl.formatMessage({ id: 'brand.metadata.selectBrandPlaceholder' })}
            style={{ width: 300 }}
            value={selectedBrandId}
            onChange={setSelectedBrandId}
            loading={brandsLoading}
          >
            {brands.map((brand: any) => (
              <Option key={brand.id} value={brand.id}>
                {brand.name}
              </Option>
            ))}
          </Select>
        </div>

        {!selectedBrandId ? (
          <Empty description={intl.formatMessage({ id: 'brand.metadata.emptySelect' })} />
        ) : metadataLoading ? (
          <div style={{ textAlign: 'center', padding: 50 }}>
            <Spin size="large" />
          </div>
        ) : (
          <MetadataEditor
            brandId={selectedBrandId}
            metadata={metadata}
            onSave={handleSave}
            loading={saving}
          />
        )}
      </Card>
    </div>
  );
};

export default BrandMetadataPage;
