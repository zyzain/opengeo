import React, { useState, useEffect } from 'react';
import { Card, Select, Spin, Empty, message } from 'antd';
import { useIntl } from 'react-intl';
import { useBrands } from '../../hooks/useBrand';
import KnowledgeGraph from '../../components/brand/KnowledgeGraph';

const { Option } = Select;

const BrandKnowledgePage: React.FC = () => {
  const intl = useIntl();
  const { brands, loading: brandsLoading } = useBrands();
  const [selectedBrandId, setSelectedBrandId] = useState<number | undefined>();
  const [entities, setEntities] = useState<any[]>([]);
  const [relations, setRelations] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (selectedBrandId) {
      fetchKnowledgeData();
    }
  }, [selectedBrandId]);

  const fetchKnowledgeData = async () => {
    if (!selectedBrandId) return;
    setLoading(true);
    try {
      const [entitiesRes, relationsRes] = await Promise.all([
        fetch(`/api/v1/brand/${selectedBrandId}/knowledge/entities`),
        fetch(`/api/v1/brand/${selectedBrandId}/knowledge/relations`),
      ]);

      if (entitiesRes.ok) {
        const entitiesData = await entitiesRes.json();
        setEntities(entitiesData.entities || []);
      }

      if (relationsRes.ok) {
        const relationsData = await relationsRes.json();
        setRelations(relationsData.relations || []);
      }
    } catch (error) {
      console.error('Failed to fetch knowledge data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddEntity = async (entity: any) => {
    try {
      await fetch(`/api/v1/brand/${selectedBrandId}/knowledge/entities`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(entity),
      });
      message.success(intl.formatMessage({ id: 'knowledge.message.entityAdded' }));
      fetchKnowledgeData();
    } catch (error) {
      message.error(intl.formatMessage({ id: 'brand.message.addEntityFailed' }));
    }
  };

  const handleAddRelation = async (relation: any) => {
    try {
      await fetch(`/api/v1/brand/${selectedBrandId}/knowledge/relations`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(relation),
      });
      message.success(intl.formatMessage({ id: 'knowledge.message.relationAdded' }));
      fetchKnowledgeData();
    } catch (error) {
      message.error(intl.formatMessage({ id: 'brand.message.addRelationFailed' }));
    }
  };

  const handleDeleteEntity = async (entityId: number) => {
    try {
      await fetch(`/api/v1/brand/${selectedBrandId}/knowledge/entities/${entityId}`, {
        method: 'DELETE',
      });
      message.success(intl.formatMessage({ id: 'knowledge.message.entityDeleted' }));
      fetchKnowledgeData();
    } catch (error) {
      message.error(intl.formatMessage({ id: 'brand.message.deleteEntityFailed' }));
    }
  };

  return (
    <div style={{ padding: 24 }}>
      <Card title={intl.formatMessage({ id: 'brand.knowledge.title' })}>
        <div style={{ marginBottom: 24 }}>
          <span style={{ marginRight: 16 }}>{intl.formatMessage({ id: 'brand.knowledge.selectBrand' })}</span>
          <Select
            placeholder={intl.formatMessage({ id: 'brand.knowledge.selectBrandPlaceholder' })}
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
          <Empty description={intl.formatMessage({ id: 'brand.knowledge.emptySelect' })} />
        ) : (
          <KnowledgeGraph
            brandId={selectedBrandId}
            entities={entities}
            relations={relations}
            loading={loading}
            onAddEntity={handleAddEntity}
            onAddRelation={handleAddRelation}
            onDeleteEntity={handleDeleteEntity}
          />
        )}
      </Card>
    </div>
  );
};

export default BrandKnowledgePage;
