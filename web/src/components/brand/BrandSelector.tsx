import React, { useState, useEffect } from 'react';
import { Select, Spin, Tag } from 'antd';
import { useBrands } from '../../hooks/useBrand';

const { Option } = Select;

interface BrandSelectorProps {
  value?: number;
  onChange?: (value: number) => void;
  placeholder?: string;
  disabled?: boolean;
  style?: React.CSSProperties;
}

const BrandSelector: React.FC<BrandSelectorProps> = ({
  value,
  onChange,
  placeholder = '请选择品牌',
  disabled = false,
  style,
}) => {
  const { brands, loading } = useBrands();

  return (
    <Select
      value={value}
      onChange={onChange}
      placeholder={placeholder}
      disabled={disabled}
      loading={loading}
      style={style}
      showSearch
      optionFilterProp="children"
      filterOption={(input, option) =>
        (option?.children as unknown as string)?.toLowerCase().includes(input.toLowerCase())
      }
    >
      {brands.map(brand => (
        <Option key={brand.id} value={brand.id}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            {brand.logo_url && (
              <img
                src={brand.logo_url}
                alt={brand.name}
                style={{ width: 20, height: 20, borderRadius: 4 }}
              />
            )}
            <span>{brand.name}</span>
            <Tag color={brand.status === 1 ? 'green' : 'red'} style={{ marginLeft: 'auto' }}>
              {brand.industry}
            </Tag>
          </div>
        </Option>
      ))}
    </Select>
  );
};

export default BrandSelector;
