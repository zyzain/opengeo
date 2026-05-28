import React, { useState } from 'react';
import { Card, Select, Spin, Empty, message } from 'antd';
import { useIntl } from 'react-intl';
import { useBrands, useGlossary } from '../../hooks/useBrand';
import GlossaryManager from '../../components/brand/GlossaryManager';

const { Option } = Select;

const BrandGlossaryPage: React.FC = () => {
  const intl = useIntl();
  const { brands, loading: brandsLoading } = useBrands();
  const [selectedBrandId, setSelectedBrandId] = useState<number | undefined>();
  const { entries, loading: glossaryLoading, refetch } = useGlossary(selectedBrandId || 0);

  return (
    <div style={{ padding: 24 }}>
      <Card title={intl.formatMessage({ id: 'brand.glossary.title' })}>
        <div style={{ marginBottom: 24 }}>
          <span style={{ marginRight: 16 }}>{intl.formatMessage({ id: 'brand.glossary.selectBrand' })}</span>
          <Select
            placeholder={intl.formatMessage({ id: 'brand.glossary.selectBrandPlaceholder' })}
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
          <Empty description={intl.formatMessage({ id: 'brand.glossary.emptySelect' })} />
        ) : (
          <GlossaryManager
            brandId={selectedBrandId}
            entries={entries}
            loading={glossaryLoading}
            onRefresh={refetch}
          />
        )}
      </Card>
    </div>
  );
};

export default BrandGlossaryPage;
