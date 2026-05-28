import { useState, useEffect, useCallback } from 'react';

interface KnowledgeEntity {
  id: number;
  brand_id: number;
  name: string;
  type: string;
  description: string;
  properties: Record<string, string>;
  tags: string[];
  created_at: string;
  updated_at: string;
}

interface KnowledgeRelation {
  id: number;
  from_entity_id: number;
  to_entity_id: number;
  type: string;
  weight: number;
  description: string;
  created_at: string;
}

export function useKnowledgeEntities(brandId: number) {
  const [entities, setEntities] = useState<KnowledgeEntity[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchEntities = useCallback(async () => {
    if (!brandId) return;
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/brand/${brandId}/knowledge/entities`);
      if (!response.ok) throw new Error('Failed to fetch entities');
      const data = await response.json();
      setEntities(data.entities || []);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setLoading(false);
    }
  }, [brandId]);

  useEffect(() => {
    fetchEntities();
  }, [fetchEntities]);

  const addEntity = async (entity: Omit<KnowledgeEntity, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await fetch(`/api/v1/brand/${brandId}/knowledge/entities`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(entity),
    });
    if (!response.ok) throw new Error('Failed to add entity');
    await fetchEntities();
    return response.json();
  };

  const updateEntity = async (id: number, updates: Partial<KnowledgeEntity>) => {
    const response = await fetch(`/api/v1/brand/${brandId}/knowledge/entities/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(updates),
    });
    if (!response.ok) throw new Error('Failed to update entity');
    await fetchEntities();
    return response.json();
  };

  const deleteEntity = async (id: number) => {
    const response = await fetch(`/api/v1/brand/${brandId}/knowledge/entities/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete entity');
    await fetchEntities();
  };

  return {
    entities,
    loading,
    error,
    refetch: fetchEntities,
    addEntity,
    updateEntity,
    deleteEntity,
  };
}

export function useKnowledgeRelations(brandId: number) {
  const [relations, setRelations] = useState<KnowledgeRelation[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchRelations = useCallback(async () => {
    if (!brandId) return;
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/v1/brand/${brandId}/knowledge/relations`);
      if (!response.ok) throw new Error('Failed to fetch relations');
      const data = await response.json();
      setRelations(data.relations || []);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setLoading(false);
    }
  }, [brandId]);

  useEffect(() => {
    fetchRelations();
  }, [fetchRelations]);

  const addRelation = async (relation: Omit<KnowledgeRelation, 'id' | 'created_at'>) => {
    const response = await fetch(`/api/v1/brand/${brandId}/knowledge/relations`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(relation),
    });
    if (!response.ok) throw new Error('Failed to add relation');
    await fetchRelations();
    return response.json();
  };

  const deleteRelation = async (id: number) => {
    const response = await fetch(`/api/v1/brand/${brandId}/knowledge/relations/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete relation');
    await fetchRelations();
  };

  return {
    relations,
    loading,
    error,
    refetch: fetchRelations,
    addRelation,
    deleteRelation,
  };
}

export function useKnowledgeGraph(brandId: number) {
  const { entities, loading: entitiesLoading, refetch: refetchEntities } = useKnowledgeEntities(brandId);
  const { relations, loading: relationsLoading, refetch: refetchRelations } = useKnowledgeRelations(brandId);

  const loading = entitiesLoading || relationsLoading;

  const refetch = () => {
    refetchEntities();
    refetchRelations();
  };

  return {
    entities,
    relations,
    loading,
    refetch,
  };
}
